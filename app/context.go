package app

import (
	v1routes "github.com/jedi-knights/tds-api/routes/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"time"
)

// Context is a struct that holds the context of the application
type Context struct {
	StartTime time.Time
	Echo      *echo.Echo
}

// NewContext returns a new Context
func NewContext() *Context {
	ctx := &Context{
		StartTime: time.Now(),
		Echo:      echo.New(),
	}

	ctx.Echo.Use(middleware.Logger())
	ctx.Echo.Use(middleware.CORS())

	ctx.Echo.Logger.SetLevel(log.DEBUG)

	v1 := ctx.Echo.Group("/api/v1")

	configureV1Routes(v1)

	ctx.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	return ctx
}

func configureV1Routes(v1 *echo.Group) {
	v1.GET("/health", v1routes.HandleHealthCheck)
	v1.GET("/version", v1routes.HandleVersion)
}
