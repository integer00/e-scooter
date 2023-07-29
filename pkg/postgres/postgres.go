package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int32
	connAttempts int
	connTimeout  time.Duration
	Pool         *pgxpool.Pool
}

func New(url string) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, errors.New("failed to parse config")
	}

	poolConfig.MaxConns = pg.maxPoolSize

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		defer pg.Pool.Close()

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	if ok := pg.Pool.Ping(context.Background()); ok != nil {
		log.Fatal(ok)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
