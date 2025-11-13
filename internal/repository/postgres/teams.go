package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-review/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrTeamAlredyExists = errors.New("team already exists")
)

func (s *Storage) GetTeamMembers(ctx context.Context, name string) ([]*models.User, error) {
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

	var users []*models.User
	for rows.Next() {
		var user models.User
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

func (s *Storage) TeamExists(ctx context.Context, name string) (bool, error) {
	const op = "postgres.TeamExists"

	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM teams WHERE name = $1`, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (s *Storage) CreateTeam(ctx context.Context, tx pgx.Tx, team *models.Team, members []*models.User) error {
	const op = "postgres.CreateTeam"

	_, err := tx.Exec(ctx, `INSERT INTO teams (id, name) VALUES ($1, $2)`, team.Id, team.Name)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrTeamAlredyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddOrUpdateTeamMembers(ctx context.Context, tx pgx.Tx, teamId string, members []*models.User) error {
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
