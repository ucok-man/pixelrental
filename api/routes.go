package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucok-man/pixelrental/internal/serializer"
	"github.com/ucok-man/pixelrental/internal/validator"
)

func (app *Application) routes() http.Handler {
	router := echo.New()
	router.HTTPErrorHandler = app.httpErrorHandler
	router.JSONSerializer = serializer.JSONSerializer{}
	router.Validator = validator.New()

	root := router.Group("/v1")
	root.Use(app.withRecover())
	root.Use(app.withLogger())
	root.Use(app.withRateLimit)

	root.GET("/healthcheck", app.healthcheckHandler)

	games := root.Group("/games")
	{
		// games.POST("", app.gameCreateHandler)
		games.GET("", app.gameGetAllHandler)
		games.GET("/:id", app.gameGetByIdHandler)
	}

	users := root.Group("/users")
	{
		users.POST("/register", app.userRegisterHandler)
		users.POST("/reactivated", app.userResendActivationTokenHandler)
		users.PUT("/activated", app.userActivatedHandler)
		users.POST("/login", app.userLoginHandler)

		users.POST("/deposit", app.userDepositHandler, app.withLogin, app.withActivated)
		users.GET("/me", app.userProfileHandler, app.withLogin, app.withUpdateDeposit)
	}

	carts := root.Group("/carts")
	carts.Use(app.withLogin, app.withActivated)
	{
		carts.GET("/estimate", app.cartEstimateHandler)
		carts.GET("", app.cartGetAllHandler)
		carts.POST("", app.cartCreateHandler)
		carts.DELETE("/:id", app.cartDeleteHandler)
		carts.PUT("/:id", app.cartUpdateHandler)
	}

	orders := root.Group("/orders")
	orders.Use(app.withLogin, app.withActivated)
	{
		orders.POST("", app.orderPayHandler, app.withUpdateDeposit)
		orders.GET("", app.orderGetAllHandler, app.withCheckPayment)
	}

	return router
}
