package contract

type ReqGameCreate struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Year        int32    `json:"year" validate:"required"`
	Genres      []string `json:"genres" validate:"required"`
	Price       float64  `json:"price" validate:"required,min=0"`
	Stock       int32    `json:"stock" validate:"required,min=0"`
}

type ReqGameGetAll struct {
	Title  string
	Genres []string
	Filters
}
