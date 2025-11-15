package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("one or more users with this usernames already exists")
)

func (s *Storage) GetTeamMembers(ctx context.Context, name string) ([]*models.Member, error) {
	const op = "postgres.GetTeamMembers"

	rows, err := s.db.Query(ctx, `SELECT id, username, is_active
	FROM users
	WHERE team_id = (
		SELECT id FROM teams WHERE name = $1
	)`, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]*models.Member, 0)
	for rows.Next() {
		var user models.Member
		if err := rows.Scan(&user.Id, &user.Username, &user.IsActive); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *Storage) AddOrUpdateTeamMembers(ctx context.Context, tx pgx.Tx, teamId string, members []*models.Member) error {
	const op = "postgres.AddOrUpdateTeamMembers"

	for _, member := range members {
		var count int
		err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE id = $1", member.Id).Scan(&count)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if count == 0 {
			_, err = tx.Exec(ctx, `
				INSERT INTO users (id, username, team_id, is_active)
				VALUES ($1, $2, $3, $4)
			`, member.Id, member.Username, teamId, member.IsActive)
		} else {
			_, err = tx.Exec(ctx, `
				UPDATE users SET username = $2, team_id = $3, is_active = $4
				WHERE id = $1
			`, member.Id, member.Username, teamId, member.IsActive)
		}
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, ErrUserExists)
			}
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) UserSetIsActive(ctx context.Context, reqDTO *dto.SetIsActiveRequest) error {
	const op = "postgres.UserSetIsActive"

	if reqDTO.IsActive == nil {
		return fmt.Errorf("%s: isActive is nil", op)
	}
	isActive := *reqDTO.IsActive

	cmd, err := s.db.Exec(ctx, `UPDATE users SET is_active = $1 WHERE id = $2`, isActive, reqDTO.UserId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *Storage) GetUserById(ctx context.Context, id string) (*models.User, error) {
	const op = "postgres.GetUserById"

	var user models.User
	err := s.db.QueryRow(ctx, `SELECT u.id, u.username, t.name, u.is_active FROM users u JOIN teams t ON u.team_id = t.id WHERE u.id = $1`, id).Scan(&user.Id, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
