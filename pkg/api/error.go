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
