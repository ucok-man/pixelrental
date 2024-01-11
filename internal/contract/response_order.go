package contract

import "time"

type ResOrderPay struct {
	Order struct {
		OrderID    int       `json:"order_id"`
		Status     string    `json:"status"`
		TotalPrice float64   `json:"total_price"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		User       struct {
			UserID    int    `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Address   string `json:"address" validate:"required"`
			Phone     string `json:"phone" validate:"required"`
			Email     string `json:"email"`
		} `json:"user"`
		OrderDetails []struct {
			Quantity int     `json:"quantity"`
			SubTotal float64 `json:"sub_total"`
			Game     struct {
				GameID      int      `json:"game_id"`
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Year        int      `json:"year"`
				Genres      []string `json:"genres"`
				Price       float64  `json:"price"`
			} `json:"game"`
		} `json:"order_details"`
		Payment struct {
			Status        string `json:"status"`
			PaymentMethod string `json:"payment_method"`
			InvoiceUrl    string `json:"invoice_url,omitempty"`
		} `json:"payment"`
	} `json:"order"`
}

type ResOrderGetAll struct {
	Orders []struct {
		OrderID    int       `json:"order_id"`
		Status     string    `json:"status"`
		TotalPrice float64   `json:"total_price"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		User       struct {
			UserID    int    `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Address   string `json:"address" validate:"required"`
			Phone     string `json:"phone" validate:"required"`
			Email     string `json:"email"`
		} `json:"user"`
		OrderDetails []struct {
			Quantity int     `json:"quantity"`
			SubTotal float64 `json:"sub_total"`
			Game     struct {
				GameID      int      `json:"game_id"`
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Year        int      `json:"year"`
				Genres      []string `json:"genres"`
				Price       float64  `json:"price"`
			} `json:"game"`
		} `json:"order_details"`
		Payment struct {
			Status        string `json:"status"`
			PaymentMethod string `json:"payment_method"`
			InvoiceUrl    string `json:"invoice_url,omitempty"`
		} `json:"payment"`
	} `json:"orders"`
}
