package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-review/internal/models"

	"github.com/jackc/pgx/v5"
)

// UnassignPRsFromUser удаляет ассайни на юзера со всех PR
func (s *Storage) UnassignPRsFromUser(ctx context.Context, tx pgx.Tx, userId string) ([]string, error) {
	const op = "postgres.UnassignPRsFromUser"

	rows, err := tx.Query(ctx, `
		DELETE FROM pull_requests_users
		WHERE user_id = $1
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

// GetTeammates возвращает идентификаторы пользователей из команды пользователя oldPRUserId, исключая при этом oldPRUserId и authorId PR'а
func (s *Storage) GetTeammates(ctx context.Context, tx pgx.Tx, prId string, oldPRUserId string) ([]string, error) {
	const op = "postgres.AssignPRToUser"

	var teamId string
	err := tx.QueryRow(ctx, `SELECT team_id FROM users WHERE id = $1`, oldPRUserId).Scan(&teamId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var authorId string
	err = tx.QueryRow(ctx, `SELECT author_id FROM pull_requests WHERE id = $1`, prId).Scan(&authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var teammates []string
	rows, err := tx.Query(ctx, `SELECT id FROM users WHERE team_id = $1 AND id != $2 AND id != $3`, teamId, oldPRUserId, authorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		teammates = append(teammates, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return teammates, nil
}

// AssignPRToUser ассайнит первого пользователя из userIds на PR
func (s *Storage) AssignPRToUser(ctx context.Context, tx pgx.Tx, prId string, userIds []string) (string, error) {
	const op = "postgres.AssignPRToUser"

	var id string
	for _, userId := range userIds {
		err := tx.QueryRow(ctx, `
			INSERT INTO pull_requests (pr_id, user_id)
			VALUES ($1, $2) RETURNING id
		`, prId, userId).Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		break
	}

	return id, nil
}

// SetNeedMoreReviewers устанавливает PR'у флаг need_more_reviewers = true
func (s *Storage) SetNeedMoreReviewers(ctx context.Context, tx pgx.Tx, prId string) error {
	const op = "postgres.SetNeedMoreReviewers"

	_, err := tx.Exec(ctx, `UPDATE pull_requests SET need_more_reviewers = true WHERE id = $1`, prId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetPRsByUserId(ctx context.Context, id string) ([]*models.PullRequest, error) {
	const op = "postgres.GetPRsByUserId"

	rows, err := s.db.Query(ctx, `
		SELECT pr.id, pr.title, pr.author_id, s.name, pr.need_more_reviewers
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
		var pr *models.PullRequest
		if err := rows.Scan(
			&pr.Id, &pr.Title, &pr.AuthorId, &pr.Status, &pr.NeedMoreReviewers,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		prs = append(prs, pr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return prs, nil
}
