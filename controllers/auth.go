package controllers

import (
	"net/http"

	"github.com/ksungcaya/todo-echo/models"
	"github.com/ksungcaya/todo-echo/repositories"
	"github.com/ksungcaya/todo-echo/requests"
	"github.com/labstack/echo/v4"
)

// AuthController todo
type AuthController struct {
	ur repositories.UserRepository
}

// userResponse is a private struct for response
type userResponse struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
}

// NewAuth creates AuthController instance
func NewAuth(ur repositories.UserRepository) *AuthController {
	return &AuthController{ur}
}

// Register handles register route
// POST /auth/register
func (ac *AuthController) Register(ctx echo.Context) error {
	rr := new(requests.RegisterRequest)
	if code, err := rr.Validate(ctx); err != nil {
		return ctx.JSON(code, requests.NewResponseError(err))
	}
	if user := ac.ur.ByEmail(rr.Email); user != nil {
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			requests.NewValidationError("email", "The email already exist"),
		)
	}
	user := rr.UserModel()
	if err := ac.ur.Create(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, requests.NewResponseError(err))
	}

	return ctx.JSON(http.StatusOK, newUserResponse(user))
}

// newUserResponse is a private function for creating *userResponse
func newUserResponse(u *models.User) ResponseData {
	r := new(userResponse)
	r.Username = u.Username
	r.Email = u.Email
	r.Name = u.Name

	return NewResponseData(r)
}
