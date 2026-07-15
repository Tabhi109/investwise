package main

import (
	"log"

	"github.com/Tabhi109/investwise/internal/app"
)

func main() {
	application := app.NewApplication()

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
