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
func (pgr *PostgresRepo) GetRides(ctx context.Context) (*[]entity.Ride, error) {

	sql := "select * from rides"
	rows, err := pgr.db.Pool.Query(ctx, sql)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rides, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Ride])

	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(rides) == 0 {
		return nil, errors.New("no rides")
	}
	log.Info(rides)

	return &rides, nil
}

func (pgr *PostgresRepo) GetRidesById(ctx context.Context, userId string) (*[]entity.Ride, error) {

	sql := fmt.Sprintf("select * from rides where user_id = '%s' and status = 'DONE'", userId)
	rows, err := pgr.db.Pool.Query(ctx, sql)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rides, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Ride])

	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(rides) == 0 {
		return nil, errors.New("no rides")
	}
	log.Info(rides)

	return &rides, nil
}

func (pgr *PostgresRepo) GetActiveRide(ctx context.Context) (*entity.Ride, error) {

	sql := "select * from rides where status != 'DONE'"
	rows, err := pgr.db.Pool.Query(ctx, sql)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ride, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Ride])

	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return &ride, nil
}

func (pgr *PostgresRepo) GetUserById(ctx context.Context, sql string) (string, error) {

	res := ""

	if err := pgr.db.Pool.QueryRow(ctx, sql).Scan(&res); err == nil {
		return res, nil
	} else {
		if err == pgx.ErrNoRows {
			log.Info(err)
			return "", err
		}
		return "", errors.New("something else")
	}

}

func (pgr *PostgresRepo) AddUser(ctx context.Context, sql string) error {

	if res, err := pgr.db.Pool.Exec(ctx, sql); err != nil {
		log.Error(res)
		return err
	}
	return nil
}

func (pgr *PostgresRepo) AddRide(ctx context.Context, t entity.Ride) error {

	sql := fmt.Sprintf("insert into rides values ('%s','%s','%s','%s')", t.RideId, t.ScooterId, t.UserId, t.Status)
	log.Info(sql)
	if res, err := pgr.db.Pool.Exec(ctx, sql); err != nil {
		log.Error(res)
		return err
	}
	return nil
}

func (pgr *PostgresRepo) UpdateRide(ctx context.Context, sql string) error {

	log.Info(sql)
	if res, err := pgr.db.Pool.Exec(ctx, sql); err != nil {
		log.Error(res)
		log.Error(err)
		return err
	}
	return nil

}
