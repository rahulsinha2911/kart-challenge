package model

type GenericErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Ecode   int    `json:"ecode"`
	Edesc   string `json:"edesc"`
}

// //////
type GenericSuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
