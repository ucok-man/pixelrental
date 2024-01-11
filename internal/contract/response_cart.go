package contract

import "time"

type ResCartGetAll struct {
	Carts []struct {
		CartID   int     `json:"cart_id"`
		Quantity int     `json:"quantity"`
		SubTotal float64 `json:"sub_total"`
		Game     struct {
			GameID      int      `json:"game_id"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Year        int      `json:"year"`
			Genres      []string `json:"genres"`
			Price       float64  `json:"price"`
			Stock       int      `json:"stock"`
		} `json:"game"`
	} `json:"carts"`
}

type ResCartCreate struct {
	Cart struct {
		CartID   int     `json:"cart_id"`
		GameID   int     `json:"game_id"`
		Quantity int     `json:"quantity"`
		SubTotal float64 `json:"sub_total"`
	} `json:"cart"`
	Message string `json:"message"`
}

type ResCartDelete struct {
	Cart struct {
		CartID   int     `json:"cart_id"`
		GameID   int     `json:"game_id"`
		Quantity int     `json:"quantity"`
		SubTotal float64 `json:"sub_total"`
	} `json:"cart"`
	Message string `json:"message"`
}

type ResCartUpdate struct {
	Cart struct {
		CartID   int     `json:"cart_id"`
		GameID   int     `json:"game_id"`
		Quantity int     `json:"quantity"`
		SubTotal float64 `json:"sub_total"`
	} `json:"cart"`
	Message string `json:"message"`
}

type ResCartEstimate struct {
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
	PriceTotal float64   `json:"price_total"`
	CreatedAt  time.Time `json:"created_at"`
}
