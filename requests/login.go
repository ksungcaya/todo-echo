package requests

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/thedevsaddam/govalidator.v1"
)

// LoginRequest is the struct for registration request
type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// make sure to implement Request interface
var _ Request = &LoginRequest{}

// Validate will validate the request with the given context
func (lr *LoginRequest) Validate(ctx echo.Context) (int, error) {
	if code, err := BindRequest(lr, ctx); err != nil {
		return code, err
	}
	if err := ValidateRequest(lr); err != nil {
		return http.StatusUnprocessableEntity, err
	}
	return http.StatusOK, nil
}

// rules is a privated function called on request validation
func (lr *LoginRequest) rules() govalidator.MapData {
	return govalidator.MapData{
		"username": []string{"required", "min:3"},
		"password": []string{"required", "min:3"},
	}
}
