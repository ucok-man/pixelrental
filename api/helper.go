package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	// "github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/ucok-man/pixelrental/internal/entity"
)

type envelope map[string]any

func (app *Application) getCurrentUser(ctx echo.Context) *entity.User {
	obj := ctx.Get(app.ctxkey.user)
	user, ok := obj.(*entity.User)
	if !ok {
		panic("[app.getCurrentUser]: user should be *entity.User")
	}
	return user
}

func (app *Application) getParamId(ctx echo.Context) (int, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return int(0), fmt.Errorf("invalid id parameter")
	}

	return int(id), nil
}

func (app *Application) getQueryString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *Application) getQuerySlice(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *Application) getQueryInt(qs url.Values, key string, defaultValue int) (int, error) {
	s := qs.Get(key)

	if s == "" {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		// v.AddError(key, "must be an integer value")
		return defaultValue, err
	}

	return i, nil
}

func (app *Application) background(fn func()) {
	app.wg.Add(1)

	go func() {

		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Errorf("%s", err), "failed processing background task", nil)
			}
		}()

		fn()
	}()
}
