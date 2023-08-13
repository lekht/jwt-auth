package postgres

import (
	"context"
	"fmt"
	"jwt-auth/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PostgreDB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Postgres) (*PostgreDB, error) {
	var connstr string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgPort, cfg.PgDB)
	dbpool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db pool")
	}

	return &PostgreDB{
		Pool: dbpool,
	}, nil
}

func (pg *PostgreDB) Close() {
	pg.Pool.Close()
}
