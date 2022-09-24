package model

import "errors"

type LocationData struct {
	Latitude           float64 `json:"lat"`
	Longitude          float64 `json:"lng"`
	Altitude           float64 `json:"altitude"`
	Identifier         string  `json:"identifier"`
	Timestamp          int64   `json:"timestamp"`
	FloorLabel         string  `json:"floor"`
	HorizontalAccuracy int     `json:"horizontalAccuracy"`
	VerticalAccuracy   int     `json:"verticalAccuracy"`
	Activity           string  `json:"activity"`
	UserID             int64   `json:"userID"`
}

type GetLocationRequest struct {
	Timestamp int64 `json:"timestamp" binding:"required"`
	UserID    int64 `json:"user_id"`
}

var (
	ErrNoDataForSuchUser = errors.New("NO_DATA")
)
