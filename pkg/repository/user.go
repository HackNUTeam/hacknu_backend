package repository

import (
	"context"
	"database/sql"
	"errors"
	"hacknu/model"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

const dbTimeout = time.Second * 3

type UserDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

func (u *UserDB) CreateReading(location *model.LocationData) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into positions (latitude, longitude, altitude, floorLabel, h_accuracy, v_accuracy, activity, user_id, _created_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	log.Print(location.UserID)
	err := u.db.QueryRowContext(ctx, stmt,
		location.Latitude,
		location.Longitude,
		location.Altitude,
		location.FloorLabel,
		location.HorizontalAccuracy,
		location.VerticalAccuracy,
		location.Activity,
		location.UserID,
		location.Timestamp,
	)
	log.Print(err)

	if err != nil {
		return err.Err()
	}

	return nil
}
func (u *UserDB) CreateUser(name string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	var id int64
	defer cancel()
	stmt := `insert into users (name) values ($1) ON CONFLICT (name) DO NOTHING`

	row := u.db.QueryRowContext(ctx, stmt, name)
	if row.Err() != nil {
		return -1, row.Err()
	}
	stmt = `select id from users where name=$1;`
	qrow := u.db.QueryRowContext(ctx, stmt, name)

	err := qrow.Scan(&id)
	log.Printf("User with id %d and name %s received", id, name)
	if err != nil {
		log.Println("Error while getting user with name ", name)
		return -1, err
	}
	return id, nil
}
func (u *UserDB) GetHistoryLocation(name string, timestamp int64) ([]*model.LocationData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var stmt string

	res := make([]*model.LocationData, 0, 1)
	var row *sql.Rows
	var err error
	var id int64

	stmt = `select id from users where name=$1;`
	qrow := u.db.QueryRowContext(ctx, stmt, name)

	err = qrow.Scan(&id)
	if err != nil {
		log.Println("Error while getting user with name ", name)
		return nil, err
	}
	log.Println("user_id:", id)
	if timestamp == -1 {
		stmt = `select latitude, longitude, altitude, _created_at, floorLabel, h_accuracy, v_accuracy, activity from positions where user_id=$1;`
		row, err = u.db.QueryContext(ctx, stmt, id)
		if err != nil {
			return nil, err
		}
		log.Println("no timestamp")
	} else {
		stmt = `select latitude, longitude, altitude, _created_at, floorLabel, h_accuracy, v_accuracy, activity from positions where _created_at >= $1 AND user_id = $2`
		row, err = u.db.QueryContext(ctx, stmt, timestamp, id)
		if err != nil {
			return nil, err
		}
	}
	if row.Err() != nil {
		if errors.Is(pgx.ErrNoRows, row.Err()) {
			return nil, model.ErrNoDataForSuchUser
		}
		return nil, row.Err()
	}

	for row.Next() {
		locationData := &model.LocationData{}
		err := row.Scan(
			&locationData.Latitude,
			&locationData.Longitude,
			&locationData.Altitude,
			&locationData.Timestamp,
			&locationData.FloorLabel,
			&locationData.HorizontalAccuracy,
			&locationData.VerticalAccuracy,
			&locationData.Activity,
		)
		if err != nil {
			if errors.Is(pgx.ErrNoRows, err) {
				return nil, model.ErrNoDataForSuchUser
			}
			return nil, err
		}
		res = append(res, locationData)
	}

	return res, nil
}
