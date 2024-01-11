package contract

type ReqUserRegister struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type ReqUserReactivate struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqUserActivated struct {
	ActivationToken string `json:"activation_token" validate:"required,min=26,max=26"`
}

type ReqUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ReqUserDeposit struct{
	Amount float64 `json:"amount" validate:"required,min=50000"`
}
