package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/ucok-man/pixelrental/internal/contract"
	"github.com/ucok-man/pixelrental/internal/repo"
)

// func (app *Application) gameCreateHandler(ctx echo.Context) error {
// 	var input contract.ReqGameCreate

// 	if err := ctx.Bind(&input); err != nil {
// 		return app.ErrBadRequest(ctx, err)
// 	}

// 	if err := ctx.Validate(&input); err != nil {
// 		return app.ErrFailedValidation(ctx, err)
// 	}

// 	var game entity.Game
// 	if err := copier.Copy(&game, &input); err != nil {
// 		return app.ErrInternalServer(ctx, err)
// 	}

// 	err := app.repo.Game.Create(&game)
// 	if err != nil {
// 		return app.ErrInternalServer(ctx, err)
// 	}

// 	var response contract.ResGameCreate
// 	if err := copier.Copy(&response.Game, &game); err != nil {
// 		return app.ErrInternalServer(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, &response)
// }

func (app *Application) gameGetByIdHandler(ctx echo.Context) error {
	gameid, err := app.getParamId(ctx)
	if err != nil {
		return app.ErrNotFound(ctx)
	}

	game, err := app.repo.Game.GetByID(gameid)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrNotFound(ctx)
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	var response contract.ResGameGetByID
	if err := copier.Copy(&response.Game, game); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

func (app *Application) gameGetAllHandler(ctx echo.Context) error {
	var input contract.ReqGameGetAll
	queryparam := ctx.QueryParams()

	input.Title = app.getQueryString(queryparam, "title", "")
	input.Genres = app.getQuerySlice(queryparam, "genres", []string{})
	input.Filters.Sort = app.getQueryString(queryparam, "sort", "game_id")
	var err error
	if input.Filters.Page, err = app.getQueryInt(queryparam, "page", 1); err != nil {
		return app.ErrFailedValidation(ctx, fmt.Errorf("page: must be positive integer"))
	}
	if input.Filters.PageSize, err = app.getQueryInt(queryparam, "page_size", 5); err != nil {
		return app.ErrFailedValidation(ctx, fmt.Errorf("page_size: must be positive integer"))
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	games, metadata, err := app.repo.Game.GetAll(&input)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResGameGetAll
	if err := copier.Copy(&response.Games, games); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	response.Metadata = *metadata

	return ctx.JSON(http.StatusOK, &response)
}
