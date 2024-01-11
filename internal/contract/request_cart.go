package contract

type ReqCartCreate struct {
	GameID   int `json:"game_id" validate:"required,min=1"`
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type ReqCartUpdate struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
