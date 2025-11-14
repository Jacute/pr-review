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

func (s *Storage) CreateTeam(ctx context.Context, tx pgx.Tx, team *models.Team) error {
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
