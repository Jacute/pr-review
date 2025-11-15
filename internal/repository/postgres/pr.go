package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-review/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPRNotFound      = errors.New("pull request not found")
	ErrPRAlreadyExists = errors.New("pull request already exists")
)

// UnassignPRsFromUser удаляет ассайни на юзера со всех открытых PR
func (s *Storage) UnassignPRsFromUser(ctx context.Context, tx pgx.Tx, userId string) ([]string, error) {
	const op = "postgres.UnassignPRsFromUser"

	rows, err := tx.Query(ctx, `
		DELETE FROM pull_requests_users pru
		WHERE pru.user_id = $1 AND
		(
			SELECT s.name FROM pull_requests pr
			JOIN statuses s ON pr.status_id = s.id
			WHERE pr.id = pru.pr_id
		) = 'OPEN'
		RETURNING pr_id
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var prIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		prIDs = append(prIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return prIDs, nil
}

// GetMembers возвращает идентификаторы пользователей из команды автора PR'а, кроме самого автора
func (s *Storage) GetMembers(ctx context.Context, tx pgx.Tx, prId string) ([]*models.Member, error) {
	const op = "postgres.GetMembers"

	var authorId string
	err := tx.QueryRow(ctx, `SELECT author_id FROM pull_requests WHERE id = $1`, prId).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var teamId string
	err = tx.QueryRow(ctx, `SELECT team_id FROM users WHERE id = $1`, authorId).Scan(&teamId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var members []*models.Member
	rows, err := tx.Query(ctx, `SELECT id, username, is_active FROM users WHERE team_id = $1 AND id != $2`, teamId, authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var member models.Member
		if err := rows.Scan(&member.Id, &member.Username, &member.IsActive); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		members = append(members, &member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return members, nil
}

// AssignPRToUser ассайнит первого пользователя из userIds на PR
func (s *Storage) AssignPRToUser(ctx context.Context, tx pgx.Tx, prId string, members []*models.Member) (string, error) {
	const op = "postgres.AssignPRToUser"

	for _, member := range members {
		cmd, err := tx.Exec(ctx, `
			INSERT INTO pull_requests_users (pr_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, prId, member.Id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue // пропускаем, если участник уже зассайнен
			}
			return "", fmt.Errorf("%s: %w", op, err)
		}
		if cmd.RowsAffected() == 1 {
			return member.Id, nil
		}
	}

	return "", nil
}

// SetNeedMoreReviewers устанавливает PR'у флаг need_more_reviewers = true
func (s *Storage) SetNeedMoreReviewers(ctx context.Context, tx pgx.Tx, prId string) error {
	const op = "postgres.SetNeedMoreReviewers"

	cmd, err := tx.Exec(ctx, `UPDATE pull_requests SET need_more_reviewers = true WHERE id = $1`, prId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrPRNotFound)
	}

	return nil
}

func (s *Storage) GetPRsByUserId(ctx context.Context, id string) ([]*models.PullRequest, error) {
	const op = "postgres.GetPRsByUserId"

	rows, err := s.db.Query(ctx, `
		SELECT pr.id, pr.title, pr.author_id, s.name, pr.need_more_reviewers, pr.merged_at
		FROM pull_requests_users pru
		JOIN pull_requests pr ON pru.pr_id = pr.id
		JOIN statuses s ON pr.status_id = s.id
		WHERE user_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var prs []*models.PullRequest
	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(
			&pr.Id, &pr.Title, &pr.AuthorId, &pr.Status, &pr.NeedMoreReviewers, &pr.MergedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		prs = append(prs, &pr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, pr := range prs {
		rows, err := s.db.Query(ctx, `
			SELECT user_id FROM pull_requests_users WHERE pr_id = $1
		`, pr.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
		reviewers := make([]string, 0)
		for rows.Next() {
			var reviewerId string
			if err := rows.Scan(&reviewerId); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			reviewers = append(reviewers, reviewerId)
		}
		pr.Reviewers = reviewers
	}

	return prs, nil
}

func (s *Storage) CreatePR(ctx context.Context, tx pgx.Tx, pr *models.PullRequestShort) error {
	const op = "postgres.CreatePR"

	_, err := tx.Exec(ctx, `
		INSERT INTO pull_requests (id, title, author_id, status_id, need_more_reviewers)
		VALUES ($1, $2, $3, (SELECT id FROM statuses WHERE name = $4), $5)
	`, pr.Id, pr.Title, pr.AuthorId, pr.Status, pr.NeedMoreReviewers)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrPRAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdatePR(ctx context.Context, tx pgx.Tx, pr *models.PullRequestShort) error {
	const op = "postgres.UpdatePR"

	values := make(map[string]interface{})
	if pr.Title != "" {
		values["title"] = pr.Title
	}
	if pr.AuthorId != "" {
		values["author_id"] = pr.AuthorId
	}
	if pr.Status != "" {
		values["status_id"] = sq.Expr("(SELECT id FROM statuses WHERE name = ?)", pr.Status)
	}
	values["need_more_reviewers"] = pr.NeedMoreReviewers

	builder := sq.Update("pull_requests").SetMap(values).Where(sq.Eq{"id": pr.Id}).PlaceholderFormat(sq.Dollar)
	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cmd, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrPRNotFound)
	}

	return nil
}

func (s *Storage) MergePR(ctx context.Context, tx pgx.Tx, id string) error {
	const op = "postgres.MergePR"

	cmd, err := tx.Exec(ctx, "UPDATE pull_requests SET status_id = (SELECT id FROM statuses WHERE name = 'MERGED'), merged_at = NOW() WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrPRNotFound)
	}

	return nil
}

func (s *Storage) GetPRById(ctx context.Context, tx pgx.Tx, id string) (*models.PullRequest, error) {
	const op = "postgres.GetPRById"

	var pr models.PullRequest

	err := tx.QueryRow(ctx, `
		SELECT pr.id, pr.title, pr.author_id, s.name, pr.need_more_reviewers
		FROM pull_requests pr
		JOIN statuses s ON pr.status_id = s.id
		WHERE pr.id = $1
	`, id).Scan(&pr.Id, &pr.Title, &pr.AuthorId, &pr.Status, &pr.NeedMoreReviewers)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.Query(ctx, `
		SELECT user_id FROM pull_requests_users WHERE pr_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewerId string
		if err := rows.Scan(&reviewerId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviewers = append(reviewers, reviewerId)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	pr.Reviewers = reviewers

	err = tx.QueryRow(ctx, `
		SELECT merged_at FROM pull_requests WHERE id = $1
	`, id).Scan(&pr.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &pr, nil
}

func (s *Storage) UnassignPRFromUser(ctx context.Context, tx pgx.Tx, prId string, userId string) error {
	const op = "postgres.UnassignPRFromUser"

	cmd, err := tx.Exec(ctx, `
		DELETE FROM pull_requests_users
		WHERE pr_id = $1 AND user_id = $2
	`, prId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrPRNotFound)
	}

	return nil
}

func (s *Storage) UserIsReviewerOfPR(ctx context.Context, tx pgx.Tx, prId string, userId string) (bool, error) {
	const op = "postgres.UserIsReviewerOfPR"

	var count int
	err := tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM pull_requests_users
		WHERE pr_id = $1 AND user_id = $2
	`, prId, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}
