package contract

type ResGameGetByID struct {
	Game struct {
		GameID      int      `json:"game_id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Year        int      `json:"year"`
		Genres      []string `json:"genres"`
		Price       float64  `json:"price"`
		Stock       int      `json:"stock"`
	} `json:"game"`
}

type ResGameCreate struct {
	Message string
	Game    struct {
		GameID      int      `json:"game_id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Year        int      `json:"year"`
		Genres      []string `json:"genres"`
		Price       float64  `json:"price"`
		Stock       int      `json:"stock"`
	} `json:"game"`
}

type ResGameGetAll struct {
	Games []struct {
		GameID      int      `json:"game_id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Year        int      `json:"year"`
		Genres      []string `json:"genres"`
		Price       float64  `json:"price"`
		Stock       int      `json:"stock"`
	} `json:"games"`

	Metadata Metadata `json:"metadata"`
}
