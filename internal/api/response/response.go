package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}
