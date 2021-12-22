package errors

import (
	"cathub.me/go-web-examples/web/formatter"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var (
	Unauthorized        = Problem{Guide: "https://cathub.me/go-web-examples/faq/Unauthorized", Status: http.StatusUnauthorized, Title: "Unauthorized"}
	TokenExpired        = Problem{Guide: "https://cathub.me/go-web-examples/faq/TokenExpired", Status: http.StatusUnauthorized, Title: "TokenExpired"}
	AccountNotActivated = Problem{Guide: "https://cathub.me/go-web-examples/faq/AccountNotActivated", Status: http.StatusUnauthorized, Title: "AccountNotActivated"}

	BadRequest = Problem{Guide: "https://cathub.me/go-web-examples/faq/BadRequest", Status: http.StatusBadRequest, Title: "BadRequest"}
)

var (
	UserNotFound          = Problem{Guide: "https://cathub.me/go-web-examples/faq/UserNotFound", Status: http.StatusNotFound, Title: "UserNotFound"}
	EmailNotFound         = Problem{Guide: "https://cathub.me/go-web-examples/faq/EmailNotFound", Status: http.StatusNotFound, Title: "EmailNotFound"}
	UsernameAlreadyExists = Problem{Guide: "https://cathub.me/go-web-examples/faq/UsernameAlreadyExists", Status: http.StatusInternalServerError, Title: "UsernameAlreadyExists"}
	EmailAlreadyExists    = Problem{Guide: "https://cathub.me/go-web-examples/faq/EmailAlreadyExists", Status: http.StatusInternalServerError, Title: "EmailAlreadyExists"}
	InvalidActivateCode   = Problem{Guide: "https://cathub.me/go-web-examples/faq/InvalidActivateCode", Status: http.StatusInternalServerError, Title: "InvalidActivateCode"}
	UserAlreadyActivated  = Problem{Guide: "https://cathub.me/go-web-examples/faq/UserAlreadyActivated", Status: http.StatusInternalServerError, Title: "UserAlreadyActivated"}
)

func ErrorHandler(error error) Problem {
	switch e := error.(type) {
	case Problem:
		return e
	case validator.ValidationErrors:
		descriptive := formatter.Descriptive(e)
		problem := Problem{
			Status:     http.StatusBadRequest,
			Title:      "InvalidParameters",
			Parameters: descriptive,
		}
		return problem
	default:
		problem := Problem{
			Status:      http.StatusInternalServerError,
			Title:       "InternalServerError",
			Description: error.Error(),
		}
		return problem
	}
}
