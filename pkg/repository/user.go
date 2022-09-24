package repository

import (
	"context"
	"database/sql"
	"errors"
	"hacknu/model"
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

	stmt := `insert into positions (latitude, longitude, altitude, _created_at, floorLabel, h_accuracy, v_accuracy, activity, userID)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	err := u.db.QueryRowContext(ctx, stmt,
		location.Latitude,
		location.Longitude,
		location.Altitude,
		location.Timestamp,
		location.FloorLabel,
		location.HorizontalAccuracy,
		location.VerticalAccuracy,
		location.Activity,
		location.UserID,
	)

	if err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (u *UserDB) GetHistoryLocation(user *model.GetLocationRequest) (*model.LocationData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var stmt string
	locationData := &model.LocationData{}
	var row *sql.Row
	if user.Timestamp == -1 {
		stmt = `select (latitude, longitude, altitude, _created_at, floorLabel, h_accuracy, v_accuracy, activity) from positions where user_id=$1;`
		row = u.db.QueryRowContext(ctx, stmt, user.UserID)

	} else {
		stmt = `select (latitude, longitude, altitude, _created_at, floorLabel, h_accuracy, v_accuracy, activity) from positions where _created_at >= $1 AND user_id = $2`
		row = u.db.QueryRowContext(ctx, stmt, user.Timestamp, user.UserID)
	}
	if row.Err() != nil {
		if errors.Is(pgx.ErrNoRows, row.Err()) {
			return nil, model.ErrNoDataForSuchUser
		}
		return nil, row.Err()
	}
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
	return locationData, nil
}
