package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Application) healthcheckHandler(ctx echo.Context) error {
	response := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.Environment,
			"version":     version,
		},
	}

	return ctx.JSON(http.StatusOK, &response)
}
