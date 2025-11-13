package postgres

import (
	"context"
	"fmt"
	"pr-review/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func (s *Storage) Stop() {
	s.db.Close()
}

func New(config *config.DatabaseConfig) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		panic("Failed to create connection pool: " + err.Error())
	}

	err = db.Ping(ctx)
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	return &Storage{db}
}

func (s *Storage) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}
