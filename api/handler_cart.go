package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/ucok-man/pixelrental/internal/contract"
	"github.com/ucok-man/pixelrental/internal/entity"
	"github.com/ucok-man/pixelrental/internal/repo"
)

// carts godoc
// @Tags carts
// @Summary Get all carts
// @Description Get all available carts record
// @Accept  json
// @Produce json
// @Success 200 {object} contract.ResGameGetAll
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /carts [get]
func (app *Application) cartGetAllHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	carts, err := app.repo.Cart.GetAll(cu.UserID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResCartGetAll
	if err := copier.Copy(&response.Carts, carts); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

// carts godoc
// @Tags carts
// @Summary Create carts
// @Description Create new carts record
// @Accept  json
// @Produce json
// @Param payload body contract.ReqCartCreate true "Create Cart"
// @Success 201 {object} contract.ResCartCreate
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /carts [post]
func (app *Application) cartCreateHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	var input contract.ReqCartCreate
	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	existingcart, err := app.repo.Cart.GetByGameID(cu.UserID, input.GameID)
	if err != nil && !errors.Is(err, repo.ErrRecordNotFound) {
		return app.ErrInternalServer(ctx, err)
	}
	if existingcart != nil {
		return app.ErrFailedValidation(ctx, fmt.Errorf("game_id: game with id %d already added to your cart", input.GameID))
	}

	game, err := app.repo.Game.GetByID(input.GameID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrFailedValidation(ctx, fmt.Errorf("game_id : did not exists"))
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	if input.Quantity > int(game.Stock) {
		return app.ErrFailedValidation(ctx, fmt.Errorf("sorry, your requested quantity exceeds the current stock availability"))
	}

	var cart entity.Cart
	if err := copier.Copy(&cart, input); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	cart.UserID = cu.UserID
	cart.SubTotal = float64(cart.Quantity) * game.Price

	if err := app.repo.Cart.Create(&cart); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResCartCreate
	if err := copier.Copy(&response.Cart, cart); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	response.Message = "success creating cart"

	return ctx.JSON(http.StatusCreated, &response)
}

// carts godoc
// @Tags carts
// @Summary Delete carts
// @Description Delete carts record
// @Accept  json
// @Produce json
// @Param id path int true "cart id"
// @Success 200 {object} contract.ResCartDelete
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 404 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /carts/:id [delete]
func (app *Application) cartDeleteHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)
	cartid, err := app.getParamId(ctx)
	if err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	cart, err := app.repo.Cart.GetByID(cu.UserID, cartid)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrNotFound(ctx)
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	err = app.repo.Cart.DeleteOne(cartid)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResCartDelete
	if err := copier.Copy(&response.Cart, cart); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	response.Message = "success removing cart"

	return ctx.JSON(http.StatusOK, &response)
}

// carts godoc
// @Tags carts
// @Summary Update carts
// @Description Update carts record
// @Accept  json
// @Produce json
// @Param id path int true "cart id"
// @Param payload body contract.ReqCartUpdate true "Update Cart"
// @Success 200 {object} contract.ResCartUpdate
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 404 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /carts/:id [put]
func (app *Application) cartUpdateHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	cartid, err := app.getParamId(ctx)
	if err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	var input contract.ReqCartUpdate
	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	cart, err := app.repo.Cart.GetByID(cu.UserID, cartid)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrNotFound(ctx)
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	game, err := app.repo.Game.GetByID(cart.GameID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrFailedValidation(ctx, fmt.Errorf("game_id : did not exists"))
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	cart.Quantity = input.Quantity
	cart.SubTotal = game.Price * float64(input.Quantity)

	if err := app.repo.Cart.Update(cart); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResCartUpdate
	if err := copier.Copy(&response.Cart, cart); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	response.Message = "success updating cart"

	return ctx.JSON(http.StatusOK, response)
}

// carts godoc
// @Tags carts
// @Summary Get calculated carts
// @Description Get estimation calculated carts record
// @Accept  json
// @Produce json
// @Success 200 {object} contract.ResCartEstimate
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /carts/estimate [get]
func (app *Application) cartEstimateHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	carts, err := app.repo.Cart.GetAll(cu.UserID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	if len(carts) == 0 {
		return app.ErrFailedValidation(ctx, fmt.Errorf("currently you have no carts available"))
	}

	var response contract.ResCartEstimate
	if err := copier.CopyWithOption(&response.OrderDetails, carts, copier.Option{
		DeepCopy:      true,
		CaseSensitive: false,
	}); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	response.PriceTotal = repo.CalculatePriceTotal(carts)
	response.CreatedAt = time.Now()

	return ctx.JSON(http.StatusOK, response)
}
