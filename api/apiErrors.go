package api	

import (
	"fmt"
)

var (
	ErrGeneric = func() string {
		return Error{Code: 0, Message: "An error occured."}.Err()
	}()

	ErrMissingParam = func() string {
		return Error{Code: 10, Message: "Required Parameter is missing."}.Err()
	}()

	ErrServer = func() string {
		return Error{Code: 500, Message: "Server Error."}.Err()
	}()
)

type Error struct {
	Code 		int 		`json:"code,attr"`
	Message string 	`json:"message,attr"`
}

func (e Error) Err() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var (
	permissionErr = ErrorResponse{&Error{Code: 403, Message: "Permission Denied"}}
	serverErr = ErrorResponse{&Error{Code: 500, Message: "Server Error"}}
)

type ErrorResponse struct {
	Error *Error `json:"error"`
}


