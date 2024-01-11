package api

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/ucok-man/pixelrental/internal/contract"
	"github.com/ucok-man/pixelrental/internal/entity"
	"github.com/ucok-man/pixelrental/internal/repo"
)

func (app *Application) orderGetAllHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	orders, err := app.repo.Order.GetAllByUserID(cu.UserID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	var response contract.ResOrderGetAll
	if err := copier.Copy(&response.Orders, orders); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, response)
}

func (app *Application) orderPayHandler(ctx echo.Context) error {
	cu := app.getCurrentUser(ctx)

	var input contract.ReqOrderPay

	if err := ctx.Bind(&input); err != nil {
		return app.ErrBadRequest(ctx, err)
	}

	if err := ctx.Validate(&input); err != nil {
		return app.ErrFailedValidation(ctx, err)
	}

	carts, err := app.repo.Cart.GetAll(cu.UserID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// check chart length or input.gameid if user buy direct
	if len(carts) == 0 && input.GameID == nil {
		return app.ErrFailedValidation(ctx, fmt.Errorf("currently you have no carts available"))
	}

	if input.GameID != nil {
		switch input.PaymentMethod {
		case "SALDO":
			return app.orderDirectWithSaldo(ctx, cu, &input)
		default:
			return app.orderDirectWithProvider(ctx, cu, &input)
		}
	} else {
		switch input.PaymentMethod {
		case "SALDO":
			return app.orderCartWithSaldo(ctx, cu, &input, carts)
		default:
			return app.orderCartWithProvider(ctx, cu, &input, carts)
		}
	}
}

func (app *Application) orderCartWithSaldo(ctx echo.Context, user *entity.User, input *contract.ReqOrderPay, charts []entity.Cart) error {
	priceTotal := repo.CalculatePriceTotal(charts)
	if user.Deposit < priceTotal {
		return app.ErrFailedValidation(ctx, fmt.Errorf("your saldo is not enough to make payment"))
	}

	// insert order
	var order = &entity.Order{
		UserID:     user.UserID,
		Status:     entity.ORDER_STATUS_SENDING,
		TotalPrice: priceTotal,
	}
	if err := app.repo.Order.Create(order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to orderdetail
	var orderdetails []*entity.OrderDetail
	if err := copier.Copy(&orderdetails, charts); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	for _, od := range orderdetails {
		od.OrderID = order.OrderID
	}

	if err := app.repo.OrderDetail.CreateBatch(orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to payment
	var payment = &entity.Payment{
		OrderID:       order.OrderID,
		Status:        entity.PAYMENT_STATUS_PAID,
		PaymentMethod: input.PaymentMethod,
	}
	if err := app.repo.Payment.Create(payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// update gamestock
	for _, cart := range charts {
		cart.Game.Stock = cart.Game.Stock - int32(cart.Quantity)
		if err := app.repo.Game.Update(cart.Game); err != nil {
			return app.ErrInternalServer(ctx, err)
		}
	}

	// delete chart
	if err := app.repo.Cart.DeleteAll(user.UserID); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// substract deposit
	user.Deposit = user.Deposit - priceTotal
	if err := app.repo.User.Update(user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// send response
	var response contract.ResOrderPay
	if err := copier.Copy(&response.Order, order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.User, user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.Payment, payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, &response)
}

func (app *Application) orderCartWithProvider(ctx echo.Context, user *entity.User, input *contract.ReqOrderPay, charts []entity.Cart) error {
	priceTotal := repo.CalculatePriceTotal(charts)
	// insert order
	var order = &entity.Order{
		UserID:     user.UserID,
		Status:     entity.ORDER_STATUS_WAITING,
		TotalPrice: priceTotal,
	}
	if err := app.repo.Order.Create(order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	inv, errxendit := app.httpclient.Payment.CreateInvoice(order.OrderID, priceTotal, &input.PaymentMethod)
	if errxendit != nil {
		return app.ErrInternalServer(ctx, errxendit)
	}

	// DONT FORGET THIS!
	err := app.repo.Token.Insert(&entity.Token{
		TokenHash: []byte(*inv.Id),
		UserID:    user.UserID,
		Expiry:    inv.ExpiryDate,
		Scope:     entity.ScopeInvoiceOrder,
	})
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to orderdetail
	var orderdetails []*entity.OrderDetail
	if err := copier.Copy(&orderdetails, charts); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	for _, od := range orderdetails {
		od.OrderID = order.OrderID
	}

	if err := app.repo.OrderDetail.CreateBatch(orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to payment
	var payment = &entity.Payment{
		OrderID:       order.OrderID,
		Status:        inv.GetStatus().String(),
		PaymentMethod: input.PaymentMethod,
		InvoiceID:     inv.GetId(),
		InvoiceUrl:    inv.GetInvoiceUrl(),
	}
	if err := app.repo.Payment.Create(payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// // update gamestock
	// for _, cart := range charts {
	// 	cart.Game.Stock = cart.Game.Stock - int32(cart.Quantity)
	// 	if err := app.repo.Game.Update(cart.Game); err != nil {
	// 		return app.ErrInternalServer(ctx, err)
	// 	}
	// }

	// delete chart
	if err := app.repo.Cart.DeleteAll(user.UserID); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// send response
	var response contract.ResOrderPay
	if err := copier.Copy(&response.Order, order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.User, user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.Payment, payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// send email
	app.background(func() {
		data := map[string]any{
			"InvoiceURL": inv.InvoiceUrl,
		}

		if err := app.mailer.Send(user.Email, "user_order_invoice.html", data); err != nil {
			app.logger.Error(err, "error sending email to user", nil)
		}
	})

	return ctx.JSON(http.StatusCreated, &response)
}

func (app *Application) orderDirectWithSaldo(ctx echo.Context, user *entity.User, input *contract.ReqOrderPay) error {
	game, err := app.repo.Game.GetByID(*input.GameID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	priceTotal := game.Price

	if user.Deposit < priceTotal {
		return app.ErrFailedValidation(ctx, fmt.Errorf("your saldo is not enough to make payment"))
	}

	// insert order
	var order = &entity.Order{
		UserID:     user.UserID,
		Status:     entity.ORDER_STATUS_SENDING,
		TotalPrice: priceTotal,
	}
	if err := app.repo.Order.Create(order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to orderdetail
	var orderdetails []*entity.OrderDetail
	orderdetails = append(orderdetails, &entity.OrderDetail{
		OrderID:  order.OrderID,
		GameID:   game.GameID,
		Quantity: 1,
		SubTotal: game.Price,
	})

	if err := app.repo.OrderDetail.CreateBatch(orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	orderdetails, err = app.repo.OrderDetail.GetAll(order.OrderID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to payment
	var payment = &entity.Payment{
		OrderID:       order.OrderID,
		Status:        entity.PAYMENT_STATUS_PAID,
		PaymentMethod: input.PaymentMethod,
	}
	if err := app.repo.Payment.Create(payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// update gamestock
	game.Stock = game.Stock - 1
	if err := app.repo.Game.Update(game); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// substract deposit
	user.Deposit = user.Deposit - priceTotal
	if err := app.repo.User.Update(user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// send response
	var response contract.ResOrderPay
	if err := copier.Copy(&response.Order, order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.User, user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.Payment, payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, &response)
}

func (app *Application) orderDirectWithProvider(ctx echo.Context, user *entity.User, input *contract.ReqOrderPay) error {
	game, err := app.repo.Game.GetByID(*input.GameID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	priceTotal := game.Price

	// insert order
	var order = &entity.Order{
		UserID:     user.UserID,
		Status:     entity.ORDER_STATUS_WAITING,
		TotalPrice: priceTotal,
	}
	if err := app.repo.Order.Create(order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	inv, errxendit := app.httpclient.Payment.CreateInvoice(order.OrderID, priceTotal, &input.PaymentMethod)
	if errxendit != nil {
		return app.ErrInternalServer(ctx, errxendit)
	}

	// DONT FORGET THIS!
	err = app.repo.Token.Insert(&entity.Token{
		TokenHash: []byte(*inv.Id),
		UserID:    user.UserID,
		Expiry:    inv.ExpiryDate,
		Scope:     entity.ScopeInvoiceOrder,
	})
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to orderdetail
	var orderdetails []*entity.OrderDetail
	orderdetails = append(orderdetails, &entity.OrderDetail{
		OrderID:  order.OrderID,
		GameID:   game.GameID,
		Quantity: 1,
		SubTotal: game.Price,
	})

	if err := app.repo.OrderDetail.CreateBatch(orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	orderdetails, err = app.repo.OrderDetail.GetAll(order.OrderID)
	if err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// insert to payment
	var payment = &entity.Payment{
		OrderID:       order.OrderID,
		Status:        inv.GetStatus().String(),
		PaymentMethod: input.PaymentMethod,
		InvoiceID:     inv.GetId(),
		InvoiceUrl:    inv.GetInvoiceUrl(),
	}
	if err := app.repo.Payment.Create(payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// // update gamestock
	// game.Stock = game.Stock - 1
	// if err := app.repo.Game.Update(game); err != nil {
	// 	return app.ErrInternalServer(ctx, err)
	// }

	// send response
	var response contract.ResOrderPay
	if err := copier.Copy(&response.Order, order); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.OrderDetails, orderdetails); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.User, user); err != nil {
		return app.ErrInternalServer(ctx, err)
	}
	if err := copier.Copy(&response.Order.Payment, payment); err != nil {
		return app.ErrInternalServer(ctx, err)
	}

	// send email
	app.background(func() {
		data := map[string]any{
			"InvoiceURL": inv.InvoiceUrl,
		}

		if err := app.mailer.Send(user.Email, "user_order_invoice.html", data); err != nil {
			app.logger.Error(err, "error sending email to user", nil)
		}
	})

	return ctx.JSON(http.StatusCreated, &response)
}
