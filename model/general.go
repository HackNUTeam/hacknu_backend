package model

const (
	statusSuccess = 1000
)

// DefaultPage
type DefaultPage struct {
	TotalPages  int64       `json:"totalPages" bson:"totalPages"`
	CurrentPage int         `json:"currentPage" bson:"currentPage"`
	Data        interface{} `json:"data" bson:"data"`
}
type PongStruct struct{}

// DefaultResponse used in every response
type DefaultResponse struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

// SuccessResponse returns success
func SuccessResponse() *DefaultResponse {
	dr := &DefaultResponse{StatusCode: statusSuccess}
	return dr
}
