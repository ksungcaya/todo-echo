package requests

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/thedevsaddam/govalidator.v1"
)

// Request is a contract for HTTP requests
type Request interface {
	rules() govalidator.MapData
	Validate(ctx echo.Context) (int, error)
}

// BindRequest binds the context parameters to request struct
func BindRequest(request Request, ctx echo.Context) (int, error) {
	if err := ctx.Bind(request); err != nil {
		return http.StatusBadRequest, errors.New("Invalid request payload")
	}

	return http.StatusOK, nil
}

// ValidateRequest validates the request with set rules
func ValidateRequest(request Request) error {
	opts := govalidator.Options{
		Data:  request,
		Rules: request.rules(),
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	if len(e) == 0 {
		return nil
	}

	return ValidationErrors{Errors: e}
}
