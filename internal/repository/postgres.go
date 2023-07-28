package repository

import (
	"context"

	"github.com/integer00/e-scooter/pkg/postgres"
)

type PostgresRepo struct {
	db *postgres.Postgres
}

func NewPostgresRepo(pg *postgres.Postgres) *PostgresRepo {
	return &PostgresRepo{
		db: pg,
	}
}

func (pgr *PostgresRepo) DoSelect(ctx context.Context, s string) error {
	return nil
}
