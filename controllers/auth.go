package controllers

import (
	"github.com/ksungcaya/todo-echo/repositories"
	"github.com/labstack/echo/v4"
)

// AuthController todo
type AuthController struct {
	ur repositories.UserRepository
}

// NewAuth creates AuthController instance
func NewAuth(ur repositories.UserRepository) *AuthController {
	return &AuthController{ur}
}

// Register handles register route
// POST /auth/register
func (ac *AuthController) Register(ctx echo.Context) error {
	return nil
}
