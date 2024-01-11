package contract

type ReqOrderPay struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=ALFAMART INDOMARET OVO SALDO"`
	GameID        *int   `json:"game_id"`
}
