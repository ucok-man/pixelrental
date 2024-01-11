package main

import (
	"log"

	"github.com/ucok-man/pixelrental/api"
	_ "github.com/ucok-man/pixelrental/docs"
)

// @title PixelRental API
// @version 1.0
// @description Documentation for PixelRental API
// @termsOfService http://swagger.io/terms/

// @contact.name ucok-man
// @contact.email ucokkocu411@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
func main() {
	if err := api.New().Serve(); err != nil {
		log.Fatal(err)
	}
}
