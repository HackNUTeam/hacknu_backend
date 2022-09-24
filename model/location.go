package model

type LocationData struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Altitude           float64 `json:"altitude"`
	Identifier         string  `json:"identifier"`
	Timestamp          int64   `json:"timestamp"`
	FloorLabel         string  `json:"floorLabel"`
	HorizontalAccuracy int     `json:"horizontalAccuracy"`
	VerticalAccuracy   int     `json:"verticalAccuracy"`
	Activity           string  `json:"activity"`
}
