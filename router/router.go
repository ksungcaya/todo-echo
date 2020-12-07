package router

import (
	"github.com/ksungcaya/todo-echo/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Router definition
type Router struct {
	*echo.Echo
}

// New creates a new Router instance.
func New() *Router {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	return &Router{e}
}

// SetAuthRoutes define auth routes
func (r *Router) SetAuthRoutes(ac *controllers.AuthController) {

}
