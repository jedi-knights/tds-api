package main

import (
	"github.com/brpaz/echozap"
	_ "github.com/jedi-knights/tds-api/docs"
	"github.com/jedi-knights/tds-api/pkg/api"
	"github.com/jedi-knights/tds-api/pkg/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"strconv"
	"strings"
)

//type (
//	QueryGenderParam struct {
//		Gender string `query:"gender" swagger:"default(female),enum(male|female|both),desc(Specify a gender)"`
//	}
//
//	QueryDivisionParam struct {
//		Division string `query:"division" swagger:"default(all),enum(all|di|dii|diii|naia|njcaa),desc(Specify a division you are interested in)"`
//	}
//
//	QueryConferenceParam struct {
//		Conference string `query:"conference" swagger:"desc(Specify a conference you are interested in)"`
//	}
//
//	QueryNameParam struct {
//		Name string `query:"name" swagger:"desc(The name of the entity you are looking for)"`
//	}
//
//	QueryNameLikeParam struct {
//		NameLike string `query:"nameLike" swagger:"desc(A partial name of the entity you are looking for)"`
//	}
//
//	QueryIdParam struct {
//		Id int `query:"id" swagger:"desc(Specify a target id)"`
//	}
//
//	PathDivisionParam struct {
//		Division string `param:"division" swagger:"default(di),enum(di|dii|diii|naia|njcaa),desc(Specify a division you are interested in),required"`
//	}
//)

func extractCode(msg string) (int, error) {
	parts := strings.Split(msg, ", ")
	code := 0 // Default value in case "code" is not found

	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]

			if key == "code" {
				// Attempt to convert the "code" value to an integer
				parsedCode, err := strconv.Atoi(value)
				if err == nil {
					code = parsedCode
					break
				}
			}
		}
	}

	return code, nil
}

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

	e := echo.New()

	e.Use(echozap.ZapLogger(logger))

	e.Use(middleware.CORS())

	// centralized error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if apiErr, ok := err.(api.Error); ok {
			_ = c.JSON(apiErr.Status, map[string]any{"error": apiErr.Msg})
			return
		}

		msg := err.Error()
		code, err := extractCode(msg)
		if err != nil {
			_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		}

		switch code {
		case http.StatusNotFound:
			_ = c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
			return
		case http.StatusInternalServerError:
			_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		case http.StatusBadRequest:
			_ = c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
			return
		case http.StatusForbidden:
			_ = c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden"})
			return
		case http.StatusUnauthorized:
			_ = c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		case http.StatusConflict:
			_ = c.JSON(http.StatusConflict, map[string]string{"error": "Conflict"})
			return
		case http.StatusUnprocessableEntity:
			_ = c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
			return
		default:
			_ = c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			return
		}
	}

	v1 := e.Group("/api/v1")
	{
		v1.GET("/health", handlers.HandleHealthCheck)
		//SetSummary("Health Check").
		//SetDescription("Returns the health of the API").
		//AddResponse(http.StatusOK, "Health Check", models.HealthCheckResponse{}, nil).
		//AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		v1.GET("/version", handlers.HandleGetVersion)
		//SetSummary("Get version").
		//SetDescription("Returns the version of the API").
		//AddResponse(http.StatusOK, "Version", models.VersionResponse{}, nil).
		//AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		v1.GET("/conferences/:division", handlers.HandleGetConferences)
		//SetSummary("List conferences").
		//SetDescription("Returns all the conferences").
		//AddParamQuery(QueryDivisionParam{}, "division", "", false).
		//AddResponse(http.StatusOK, "Conferences", []models.Conference{}, nil).
		//AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		//v1.GET("/teams/conferences", handlers.HandleGetTeamsConferences).
		//	SetSummary("List teams conferences").
		//	SetDescription("Returns all teams conferences").
		//	AddResponse(http.StatusOK, "Conferences", []models.Conference{}, nil).
		//	AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		//v1.GET("/teams/conferences", handlers.HandleGetTeamsConferencesByDivision).
		//	SetSummary("List teams conferences by division").
		//	SetDescription("Returns all teams conferences by division").
		//	AddParamPath(PathDivisionParam{}, "division", "").
		//	AddResponse(http.StatusOK, "Conferences", []models.Conference{}, nil).
		//	AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)

		v1.GET("/teams", handlers.HandleGetTeams)
		//SetSummary("List teams").
		//SetDescription("Returns all teams").
		//AddParamQuery(QueryGenderParam{}, "gender", "", false).
		//AddParamQuery(QueryDivisionParam{}, "division", "", false).
		//AddParamQuery(QueryConferenceParam{}, "conference", "", false).
		//AddParamQuery(QueryNameParam{}, "name", "", false).
		//AddParamQuery(QueryNameLikeParam{}, "nameLike", "", false).
		//AddParamQuery(QueryIdParam{}, "id", "", false).
		//AddResponse(http.StatusOK, "Teams", []models.Team{}, nil).
		//AddResponse(http.StatusInternalServerError, "Internal Server Error", api.Error{}, nil)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
