package requests

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/thedevsaddam/govalidator.v1"
)

// Request is a contract for HTTP requests
type Request interface {
	rules() govalidator.MapData
	Validate(ctx echo.Context) (int, error)
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
