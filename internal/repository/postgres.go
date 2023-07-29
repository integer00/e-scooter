package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/integer00/e-scooter/internal/entity"
	"github.com/integer00/e-scooter/pkg/postgres"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type PostgresRepo struct {
	db *postgres.Postgres
}

func NewPostgresRepo(pg *postgres.Postgres) *PostgresRepo {
	return &PostgresRepo{
		db: pg,
	}
}

//FindUserById(ctx Context, s string) (string,error)
//FindRideById(ctx Context, s string) (string,error)

//AddUserById(s string) error
//AddRideById(s string) error
//GetUsers(ctx Contex) []User
//GetRides(ctx Contex) []Ride

func (pgr *PostgresRepo) GetUsers(ctx context.Context) ([]entity.User, error) {

	rows, err := pgr.db.Pool.Query(ctx, "select * from users")
	if err != nil {
		return nil, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (pgr *PostgresRepo) GetRides(ctx context.Context) ([]entity.Ride, error) {

	rows, err := pgr.db.Pool.Query(ctx, "select * from rides")
	if err != nil {
		return nil, err
	}
	rides, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Ride])
	if err != nil {
		return nil, err
	}
	return rides, nil
}

func (pgr *PostgresRepo) FindUserById(ctx context.Context, s string) (string, error) {

	res := ""

	if err := pgr.db.Pool.QueryRow(ctx, s).Scan(&res); err == nil {
		return res, nil
	} else {
		if err == pgx.ErrNoRows {
			log.Info(err)
			return "", err
		}
		return "", errors.New("something else")
	}

}

func (pgr *PostgresRepo) AddUserById(ctx context.Context, s string) error {

	if res, err := pgr.db.Pool.Exec(ctx, s); err != nil {
		log.Error(res)
		return err
	}
	return nil
}

func (pgr *PostgresRepo) AddRide(ctx context.Context, t entity.Ride) error {
	// RideID      string
	// Scooter     Scooter
	// User        User
	// Date        string
	// Time        string
	// Status      string
	// FareCharged string
	// Distance    string
	// StartTime   string
	// StopTime    string

	sql := fmt.Sprintf("insert into rides (ride_id, scooter_id, user_id, status) values ('%s','%s','%s','%s')", t.RideID, t.Scooter.Id, t.User.Name, t.Status)
	log.Info(sql)
	if res, err := pgr.db.Pool.Exec(ctx, sql); err != nil {
		log.Error(res)
		return err
	}
	return nil
}
