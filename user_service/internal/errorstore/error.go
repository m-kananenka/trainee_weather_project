package errorstore

type AppError struct {
	Error   string `json:"Error"`
	Message string `json:"Message"`
	Code    int    `json:"Code"`
}

func EntityNotFound(err error) *AppError {
	return &AppError{
		Error:   "Entity not found",
		Message: err.Error(),
		Code:    404,
	}
}

func BadRequest(err error) *AppError {
	return &AppError{
		Error:   "Bad request",
		Message: err.Error(),
		Code:    400,
	}
}

func StatusInternalServerError(err error) *AppError {
	return &AppError{
		Error:   "Failed!",
		Message: err.Error(),
		Code:    500,
	}
}

func MethodNotAllowed(err error) *AppError {
	return &AppError{
		Error:   "Method Not Allowed ",
		Message: err.Error(),
		Code:    405,
	}
}

func Unauthorized(err error) *AppError {
	return &AppError{
		Error:   "Unauthorized",
		Message: err.Error(),
		Code:    401,
	}
}
