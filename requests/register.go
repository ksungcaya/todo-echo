package requests

import (
	"net/http"

	"github.com/ksungcaya/todo-echo/models"
	"github.com/labstack/echo/v4"
	"gopkg.in/thedevsaddam/govalidator.v1"
)

// RegisterRequest is the struct for registration request
type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

// make sure to implement Request interface
var _ Request = &RegisterRequest{}

// Validate will validate the request with the given context
func (rr *RegisterRequest) Validate(ctx echo.Context) (int, error) {
	if err := ctx.Bind(rr); err != nil {
		return http.StatusBadRequest, err
	}
	if err := ValidateRequest(rr); err != nil {
		return http.StatusUnprocessableEntity, err
	}
	return http.StatusOK, nil
}

// UserModel creates a *models.User using request data
func (rr *RegisterRequest) UserModel() *models.User {
	return &models.User{
		Username: rr.Username,
		Email:    rr.Email,
		Name:     rr.Name,
		Password: rr.Password,
	}
}

// rules is a privated function called on request validation
func (rr *RegisterRequest) rules() govalidator.MapData {
	return govalidator.MapData{
		"username": []string{"required", "between:3,5"},
		"email":    []string{"required", "min:4", "max:20", "email"},
		"name":     []string{"required", "min:4", "max:20"},
		"password": []string{"required", "min:6"},
	}
}
