package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error struct {
		Message any `json:"message"`
	} `json:"error"`
}

func (app *Application) httpErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}
	httperr, ok := err.(*echo.HTTPError)
	if !ok {
		httperr = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	if httperr.Internal != nil {
		if herr, ok := httperr.Internal.(*echo.HTTPError); ok {
			httperr = herr
		}
	}

	switch httperr.Code {
	case http.StatusNotFound:
		httperr = app.ErrNotFound(ctx)
	case http.StatusMethodNotAllowed:
		httperr = app.ErrMethodNotAllowed(ctx)
	default:
	}

	var response = &ErrorResponse{}
	response.Error.Message = fmt.Sprintf("%v", httperr.Message)

	err = ctx.JSON(httperr.Code, response)
	if err != nil {
		app.logger.Error(err, "ECHO INTERNAL", nil)
		ctx.Response().WriteHeader(http.StatusInternalServerError)
	}

}

func (app *Application) ErrInternalServer(ctx echo.Context, err error) error {
	message := "the server encountered a problem and could not process your request"
	return echo.NewHTTPError(http.StatusInternalServerError, message).SetInternal(err)
}

func (app *Application) ErrNotFound(ctx echo.Context, customeMsg ...string) *echo.HTTPError {
	if len(customeMsg) == 0 {
		customeMsg = append(customeMsg, "the requested resource could not be found")
	}
	message := customeMsg[0]
	return echo.NewHTTPError(http.StatusNotFound, message)
}

func (app *Application) ErrMethodNotAllowed(ctx echo.Context) *echo.HTTPError {
	message := fmt.Sprintf("the %s method is not supported for this resource", ctx.Request().Method)
	return echo.NewHTTPError(http.StatusMethodNotAllowed, message)
}

func (app *Application) ErrBadRequest(ctx echo.Context, err error) error {
	httperr, ok := err.(*echo.HTTPError)
	if ok {
		err = fmt.Errorf(fmt.Sprintf("%v", httperr.Message))
	}
	return echo.NewHTTPError(http.StatusBadRequest, err)
}

func (app *Application) ErrFailedValidation(ctx echo.Context, err error) error {
	httperr, ok := err.(*echo.HTTPError)
	if ok {
		err = fmt.Errorf(fmt.Sprintf("%v", httperr.Message))
	}
	return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
}

// func (app *Application) ErrEditConflict(ctx echo.Context, err error) error {
// 	message := "unable to update the record due to an edit conflict, please try again"
// 	return echo.NewHTTPError(http.StatusConflict, message)
// }

func (app *Application) ErrRateLimitExceeded(ctx echo.Context) error {
	message := "rate limit exceeded"
	return echo.NewHTTPError(http.StatusTooManyRequests, message)
}

func (app *Application) ErrInvalidCredentials(ctx echo.Context) error {
	message := "invalid authentication credentials"
	return echo.NewHTTPError(http.StatusUnauthorized, message)
}

func (app *Application) ErrInvalidAuthenticationToken(ctx echo.Context) error {
	ctx.Request().Header.Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	return echo.NewHTTPError(http.StatusUnauthorized, message)
}

// func (app *Application) ErrAuthenticationRequired(ctx echo.Context) error {
// 	message := "you must be authenticated to access this resource"
// 	return echo.NewHTTPError(http.StatusBadRequest, message)
// }

func (app *Application) ErrInactiveAccount(ctx echo.Context) error {
	message := "your user account must be activated to access this resource"
	return echo.NewHTTPError(http.StatusForbidden, message)
}

// func (app *Application) ErrNotPermitted(ctx echo.Context) error {
// 	message := "your user account doesn't have the necessary permissions to access this resource"
// 	return echo.NewHTTPError(http.StatusForbidden, message)
// }
