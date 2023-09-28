package api

import "net/http"

type Error struct {
	Status int
	Msg    string
}

func (e Error) Error() string {
	return e.Msg
}

func NotFoundError(r string) Error {
	return Error{
		Status: http.StatusNotFound,
		Msg:    r + " not found",
	}
}

func InternalServerError(r string) Error {
	return Error{
		Status: http.StatusInternalServerError,
		Msg:    "internal server error: " + r,
	}
}

func BadRequestError(r string) Error {
	return Error{
		Status: http.StatusBadRequest,
		Msg:    "bad request: " + r,
	}
}

func ForbiddenError(r string) Error {
	return Error{
		Status: http.StatusForbidden,
		Msg:    "forbidden: " + r,
	}
}

func UnauthorizedError(r string) Error {
	return Error{
		Status: http.StatusUnauthorized,
		Msg:    "unauthorized: " + r,
	}
}

func ConflictError(r string) Error {
	return Error{
		Status: http.StatusConflict,
		Msg:    "conflict: " + r,
	}
}

func UnprocessableEntityError(r string) Error {
	return Error{
		Status: http.StatusUnprocessableEntity,
		Msg:    "unprocessable entity: " + r,
	}
}
