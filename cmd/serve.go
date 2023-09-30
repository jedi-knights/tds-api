/*
Copyright © 2023 Omar Crosby <omar.crosby@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	_ "github.com/jedi-knights/tds-api/docs"
	"github.com/jedi-knights/tds-api/routes"
	v1routes "github.com/jedi-knights/tds-api/routes/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()

		e.Logger.SetLevel(log.DEBUG)

		e.Logger.Info("Starting server")

		e.Use(middleware.CORS())

		e.Pre(middleware.RemoveTrailingSlash())
		e.Pre(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "time=${time_rfc3339_nano} remote_ip=${remote_ip} host=${host} method=${method} uri=${uri} ${user_agent} status=${status} error=${error} latency=${latency_human}\n",
		}))

		// centralized error handler
		e.HTTPErrorHandler = func(err error, c echo.Context) {
			if apiErr, ok := err.(routes.Error); ok {
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

		configureV1Routes(v1)

		e.GET("/swagger/*", echoSwagger.WrapHandler)

		e.Logger.Fatal(e.Start(":8080"))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configureV1Routes(v1 *echo.Group) {
	v1.GET("/health", v1routes.HandleHealthCheck)
	v1.GET("/version", v1routes.HandleVersion)
	v1.GET("/conferences", v1routes.HandleGetConferences)
	v1.GET("/conferences/:division", v1routes.HandleGetConferencesByDivision)
	v1.GET("/teams", v1routes.HandleGetTeams)
}

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
