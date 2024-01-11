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
	"github.com/ucok-man/pixelrental/internal/jwt"
	"github.com/ucok-man/pixelrental/internal/repo"
)

// users godoc
// @Tags users
// @Summary Create user
// @Description Create new user record
// @Accept  json
// @Produce json
// @Param payload body contract.ReqUserRegister true "Create User"
// @Success 202 {object} contract.ResUserRegister
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /users/register [post]
func (app *Application) userRegisterHandler(ctx echo.Context) error {
	var input contract.ReqUserRegister

	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	var user entity.User
	if err := copier.Copy(&user, &input); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	if err := user.SetPassword(input.Password); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	err := app.repo.User.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDuplicateRecord):
			return app.ErrFailedValidation(ctx, fmt.Errorf("email: already exists"))
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	token, err := app.repo.Token.GenerateAndInsert(user.UserID, 3*24*time.Hour, entity.ScopeActivation)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	app.background(func() {
		data := map[string]any{
			"UserID":          user.UserID,
			"ActivationToken": token.Plaintext,
		}

		if err := app.mailer.Send(user.Email, "user_welcome.html", data); err != nil {
			app.logger.Error(err, "error sending email to user", nil)
		}
	})

	var response contract.ResUserRegister
	if err := copier.Copy(&response.User, user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusAccepted, response)
}

// users godoc
// @Tags users
// @Summary Resend activation token
// @Description Resending new activation token user
// @Accept  json
// @Produce json
// @Param payload body contract.ReqUserReactivate true "Email User"
// @Success 200 {object} contract.ResResendActivationToken
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /users/reactivated [post]
func (app *Application) userResendActivationTokenHandler(ctx echo.Context) error {
	var input contract.ReqUserReactivate

	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	user, err := app.repo.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrFailedValidation(ctx, fmt.Errorf("email: no matching email address found"))
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	if user.Activated {
		return app.ErrFailedValidation(ctx, fmt.Errorf("email: user has already been activated"))
	}

	token, err := app.repo.Token.GenerateAndInsert(user.UserID, 3*24*time.Hour, entity.ScopeActivation)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	app.background(func() {
		data := map[string]any{
			"ActivationToken": token.Plaintext,
		}

		if err := app.mailer.Send(user.Email, "resend_token_activation.html", data); err != nil {
			app.logger.Error(err, "error sending email to user", nil)
		}
	})

	var response = &contract.ResResendActivationToken{
		Message: "an email will be sent to you containing activation instructions",
	}

	return ctx.JSON(http.StatusOK, response)
}

// users godoc
// @Tags users
// @Summary Activate user record
// @Description Activate registered user record
// @Accept  json
// @Produce json
// @Param payload body contract.ReqUserActivated true "Activation Token"
// @Success 200 {object} contract.ResUserActivated
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /users/activated [post]
func (app *Application) userActivatedHandler(ctx echo.Context) error {
	var input contract.ReqUserActivated

	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	user, err := app.repo.User.GetToken(entity.ScopeActivation, input.ActivationToken)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrFailedValidation(ctx, fmt.Errorf("token: invalid or expired activation token"))
		default:
			return app.ErrInternalServer(ctx, err)

		}
	}
	user.Activated = true

	err = app.repo.User.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrEditConflict):
			// return app.ErrEditConflict(ctx, err)
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	err = app.repo.Token.DeleteTokenAll(entity.ScopeActivation, user.UserID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResUserActivated
	if err := copier.Copy(&response.User, user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// users godoc
// @Tags users
// @Summary Login user
// @Description Login user record
// @Accept  json
// @Produce json
// @Param payload body contract.ReqUserLogin true "Login User"
// @Success 200 {object} contract.ResUserLogin
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /users/login [post]
func (app *Application) userLoginHandler(ctx echo.Context) error {
	var input contract.ReqUserLogin

	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	user, err := app.repo.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrRecordNotFound):
			return app.ErrInvalidCredentials(ctx)
		default:
			return app.ErrInternalServer(ctx, err)
		}
	}

	if err := user.MatchesPassword(input.Password); err != nil {
		return app.ErrInvalidCredentials(ctx)
	}

	expiration := time.Now().Add(24 * time.Hour)
	claims := jwt.NewJWTClaim(user.UserID, expiration)
	token, err := jwt.GenerateToken(&claims, app.config.Jwt.Secret)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response = &contract.ResUserLogin{}
	response.AuthenticationToken.Token = token
	response.AuthenticationToken.Expiry = expiration.String()

	return ctx.JSONPretty(http.StatusOK, response, "\t")
}

// users godoc
// @Tags users
// @Summary Deposit saldo user
// @Description Top up saldo user record
// @Accept  json
// @Produce json
// @Param payload body contract.ReqUserDeposit true "Deposit Saldo"
// @Success 200 {object} contract.ResUserDeposit
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 400 {object} object{error=object{message=string}}
// @Failure 422 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /users/deposit [post]
func (app *Application) userDepositHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	var input contract.ReqUserDeposit
	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	invoice, err := app.httpclient.Payment.CreateInvoice(cu.UserID, input.Amount, nil)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// DONT FORGET THIS!
	err = app.repo.Token.Insert(&entity.Token{
		TokenHash: []byte(*invoice.Id),
		UserID:    cu.UserID,
		Expiry:    invoice.ExpiryDate,
		Scope:     entity.ScopeInvoiceTopUp,
	})
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	app.background(func() {
		data := map[string]any{
			"InvoiceURL": invoice.InvoiceUrl,
		}

		if err := app.mailer.Send(cu.Email, "user_deposit_invoice.html", data); err != nil {
			app.logger.Error(err, "error sending email to user", nil)
		}
	})

	var response contract.ResUserDeposit
	response.Message = "an email will be sent to you containing payment instructions"

	return ctx.JSON(http.StatusAccepted, &response)
}

// users godoc
// @Tags users
// @Summary Get profile user
// @Description Get info of login user
// @Accept  json
// @Produce json
// @Success 200 {object} contract.ResUserDeposit
// @Failure 429 {object} object{error=object{message=string}}
// @Failure 403 {object} object{error=object{message=string}}
// @Failure 401 {object} object{error=object{message=string}}
// @Failure 500 {object} object{error=object{message=string}}
// @Router /users/me [get]
func (app *Application) userProfileHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	var response contract.ResUserProfile
	if err := copier.Copy(&response.User, cu); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, &response)
}
