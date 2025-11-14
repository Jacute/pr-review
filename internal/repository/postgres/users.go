package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"

	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound = errors.New("user not found")
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

	var users []*models.Member
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
	const op = "postgres.AddTeamMembers"

	for _, member := range members {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (id, username, team_id, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE
			SET username = EXCLUDED.username,
				team_id = EXCLUDED.team_id,
				is_active = EXCLUDED.is_active
		`, member.Id, member.Username, teamId, member.IsActive)

		if err != nil {
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
