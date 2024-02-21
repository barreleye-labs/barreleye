package restful

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Error      *Error      `json:"error"`
}

func ResponseBadRequest(message string) *Response {
	err := &Error{
		Message: message,
	}
	return &Response{
		StatusCode: 400,
		Error:      err,
	}
}

func ResponseNotFound(message string) *Response {
	err := &Error{
		Message: message,
	}
	return &Response{
		StatusCode: 404,
		Error:      err,
	}
}

func ResponseServerError(message string) *Response {
	err := &Error{
		Message: message,
	}
	return &Response{
		StatusCode: 500,
		Error:      err,
	}
}

func ResponseOk(data interface{}) *Response {
	return &Response{
		StatusCode: 200,
		Data:       data,
	}
}
