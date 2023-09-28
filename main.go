package main

import (
	"github.com/brpaz/echozap"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/handlers"
	"github.com/jedi-knights/tds-api/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"net/http"
)

type (
	QueryGenderParam struct {
		Gender string `query:"gender" swagger:"default(male),enum(male,female,both),desc(Specify a target gender)"`
	}

	QueryDivisionParam struct {
		Division string `query:"division" swagger:"default(di),enum(di,dii,diii,naia,njcaa),desc(Specify a target division)"`
	}
)

// @title TopDrawerSoccer API
// @version 1.0
// @description This is a simple API providing access to data from TopDrawerSoccer.
// @contact.name Omar Crosby
// @contact.email omar.crosby@gmail.com
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger := api.GetLogger()

	logger.Info("Starting server...")

	r := echoswagger.New(echo.New(), "/swagger", nil)

	r.Echo().Use(echozap.ZapLogger(logger))

	// centralized error handler
	r.Echo().HTTPErrorHandler = func(err error, c echo.Context) {
		if apiErr, ok := err.(api.Error); ok {
			_ = c.JSON(apiErr.Status, map[string]any{"error": apiErr.Msg})
			return
		}

		_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	v1 := r.Group("v1", "/api/v1")
	{
		v1.GET("/version", handlers.HandleGetVersion).
			SetSummary("Get version").
			SetDescription("Returns the version of the API").
			// AddResponse(http.StatusOK, "Version", , nil).
			AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		v1.GET("/conference/names", handlers.HandleGetConferenceNames).
			SetSummary("List conference names").
			SetDescription("Returns all the conference names").
			AddResponse(http.StatusOK, "Conference Names", []string{}, nil).
			AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		v1.GET("/conferences", handlers.HandleGetConferences).
			SetSummary("List conferences").
			SetDescription("Returns all the conferences").
			AddParamQuery(QueryGenderParam{}, "Gender", "", false).
			AddParamQuery(QueryDivisionParam{}, "Division", "", false).
			AddResponse(http.StatusOK, "Conferences", []models.Conference{}, nil).
			AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)
	}

	r.Echo().Logger.Fatal(r.Echo().Start(":8080"))
}
