package main

import (
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	logger := api.GetLogger()

	logger.Info("Starting server...")

	e := echo.New()

	e.GET("/version", handlers.HandleGetVersion)

	e.GET("/conference/names", handlers.HandleGetConferenceNames)
	e.GET("/conferences", handlers.HandleGetConferences)

	// centralized error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if apiErr, ok := err.(api.Error); ok {
			_ = c.JSON(apiErr.Status, map[string]any{"error": apiErr.Msg})
			return
		}

		_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	e.Logger.Fatal(e.Start(":8080"))
}
