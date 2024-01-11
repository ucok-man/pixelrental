package contract

type ResUserRegister struct {
	User struct {
		UserID    int     `json:"user_id"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Address   string  `json:"address" validate:"required"`
		Phone     string  `json:"phone" validate:"required"`
		Email     string  `json:"email"`
		Deposit   float64 `json:"deposit"`
		Activated bool    `json:"activated"`
	} `json:"user"`
}

type ResResendActivationToken struct {
	Message string `json:"message"`
}

type ResUserActivated struct {
	User struct {
		UserID    int     `json:"user_id"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Address   string  `json:"address" validate:"required"`
		Phone     string  `json:"phone" validate:"required"`
		Email     string  `json:"email"`
		Deposit   float64 `json:"deposit"`
		Activated bool    `json:"activated"`
	} `json:"user"`
}

type ResUserLogin struct {
	AuthenticationToken struct {
		Token  string `json:"token"`
		Expiry string `json:"expiry"`
	} `json:"auhentication_token"`
}

type ResUserDeposit struct {
	Message string `json:"message"`
}

type ResUserProfile struct {
	User struct {
		UserID    int     `json:"user_id"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Address   string  `json:"address" validate:"required"`
		Phone     string  `json:"phone" validate:"required"`
		Email     string  `json:"email"`
		Deposit   float64 `json:"deposit"`
		Activated bool    `json:"activated"`
	}
}
