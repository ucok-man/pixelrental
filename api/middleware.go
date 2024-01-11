package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ucok-man/pixelrental/internal/entity"
	"github.com/ucok-man/pixelrental/internal/jwt"
	"github.com/ucok-man/pixelrental/internal/logging"
	"github.com/ucok-man/pixelrental/internal/repo"
	"github.com/xendit/xendit-go/v4/invoice"
	"golang.org/x/time/rate"
)

func (app *Application) withRecover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			app.logger.Error(err, "PANIC RECOVER", logging.Meta{
				"stack": string(stack),
			})
			return err
		},
	})
}

func (app *Application) withLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP:     true,
		LogStatus:       true,
		LogMethod:       true,
		LogURI:          true,
		LogLatency:      true,
		LogResponseSize: true,
		LogError:        true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			switch v.Status {
			case http.StatusInternalServerError:
				httperr, _ := v.Error.(*echo.HTTPError)
				app.logger.Error(httperr.Internal, http.StatusText(500), logging.Meta{
					"code":          v.Status,
					"method":        v.Method,
					"url":           v.URI,
					"ip_addr":       v.RemoteIP,
					"response_time": v.Latency,
					"response_size": v.ResponseSize,
					"stack":         v.Error,
				})
			default:
				app.logger.Info(http.StatusText(v.Status), logging.Meta{
					"code":          v.Status,
					"method":        v.Method,
					"url":           v.URI,
					"ip_addr":       v.RemoteIP,
					"response_time": v.Latency,
					"response_size": v.ResponseSize,
				})
			}

			return nil
		},
	})
}

func (app *Application) withRateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(ctx echo.Context) error {
		if app.config.Limiter.Enabled {
			ip := ctx.RealIP()

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.Limiter.Rps), app.config.Limiter.Burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				return app.ErrRateLimitExceeded(ctx)
			}

			mu.Unlock()
		}

		return next(ctx)
	}
}

func (app *Application) withLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorizationHeader := ctx.Request().Header.Get("Authorization")
		if authorizationHeader == "" {
			return app.ErrInvalidCredentials(ctx)
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return app.ErrInvalidCredentials(ctx)
		}
		tokenstr := headerParts[1]

		var claim jwt.JWTClaim
		err := jwt.DecodeToken(tokenstr, &claim, app.config.Jwt.Secret)
		if err != nil {
			return app.ErrInvalidCredentials(ctx)
		}

		user, err := app.repo.User.GetByID(claim.UserID)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrRecordNotFound):
				return app.ErrInvalidCredentials(ctx)
			default:
				return app.ErrInternalServer(ctx, err)
			}
		}

		// set context
		ctx.Set(app.ctxkey.user, user)
		return next(ctx)
	}
}

func (app *Application) withActivated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := app.getCurrentUser(ctx)

		if !user.Activated {
			return app.ErrInactiveAccount(ctx)
		}
		return next(ctx)
	}
}

func (app *Application) withUpdateDeposit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cu := app.getCurrentUser(ctx)

		invoicetoken, err := app.repo.Token.GetTokenInvoice(entity.ScopeInvoiceTopUp, cu.UserID)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrRecordNotFound):
				return next(ctx)
			default:
				return app.ErrInternalServer(ctx, err)
			}
		}

		inv, errxendit := app.httpclient.Payment.GetInvoice(string(invoicetoken.TokenHash))
		if errxendit != nil {
			switch {
			case errxendit.Status() == "404":
				return next(ctx)
			default:
				app.ErrInternalServer(ctx, errxendit)
			}
		}

		// UPDATE DEPOSIT
		cu.Deposit = inv.GetAmount()

		if err := app.repo.User.Update(cu); err != nil {
			return app.ErrInternalServer(ctx, err)
		}

		// DELETE TOKEN
		if err := app.repo.Token.DeleteTokenAll(entity.ScopeInvoiceTopUp, cu.UserID); err != nil {
			return app.ErrInternalServer(ctx, err)
		}

		ctx.Set(app.ctxkey.user, cu)
		return next(ctx)
	}
}

func (app *Application) withCheckPayment(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cu := app.getCurrentUser(ctx)

		invoicetoken, err := app.repo.Token.GetTokenInvoice(entity.ScopeInvoiceOrder, cu.UserID)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrRecordNotFound):
				return next(ctx)
			default:
				return app.ErrInternalServer(ctx, err)
			}
		}

		inv, errxendit := app.httpclient.Payment.GetInvoice(string(invoicetoken.TokenHash))
		if errxendit != nil {
			switch {
			case errxendit.Status() == "404":
				return next(ctx)
			default:
				app.ErrInternalServer(ctx, errxendit)
			}
		}

		if inv.GetStatus() == invoice.INVOICESTATUS_PAID {
			orderid, err := strconv.Atoi(inv.GetExternalId())
			if err != nil {
				return app.ErrInternalServer(ctx, err)
			}

			payment, err := app.repo.Payment.GetByorderID(orderid)
			if err != nil {
				return app.ErrInternalServer(ctx, err)
			}

			// update payment
			payment.Status = inv.GetStatus().String()
			if err := app.repo.Payment.UpdateStatus(payment); err != nil {
				return app.ErrInternalServer(ctx, err)
			}

			// update stock game
			orderdetails, err := app.repo.OrderDetail.GetAll(orderid)
			if err != nil {
				return app.ErrInternalServer(ctx, err)
			}
			for _, od := range orderdetails {
				od.Game.Stock = od.Game.Stock - int32(od.Quantity)
				if err := app.repo.Game.Update(od.Game); err != nil {
					return app.ErrInternalServer(ctx, err)
				}
			}

			// update order
			order, err := app.repo.Order.GetByID(orderid)
			if err != nil {
				return app.ErrInternalServer(ctx, err)
			}
			order.Status = entity.ORDER_STATUS_SENDING
			if err := app.repo.Order.UpdateStatus(order); err != nil {
				return app.ErrInternalServer(ctx, err)
			}

			// DELETE TOKEN
			if err := app.repo.Token.DeleteTokenAll(entity.ScopeInvoiceOrder, cu.UserID); err != nil {
				return app.ErrInternalServer(ctx, err)
			}

			ctx.Set(app.ctxkey.user, cu)
		}
		return next(ctx)
	}
}
