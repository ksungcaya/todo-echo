package controllers

// Controllers is
type Controllers struct {
	Auth AuthController
}

// New is..
func New() *Controllers {
	return &Controllers{
		Auth: AuthController{},
	}
}
