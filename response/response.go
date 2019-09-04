package response

type ResponseError struct {
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return e.Message
}
