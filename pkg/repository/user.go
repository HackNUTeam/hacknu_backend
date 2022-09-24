package repository

import (
	"context"
	"database/sql"
	"hacknu/model"
	"time"
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

	stmt := `insert into locations (latitude, longitude, altitude, timestamp, floorLabel, horizontalAccuracy, verticalAccuracy, activity, userID)
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
