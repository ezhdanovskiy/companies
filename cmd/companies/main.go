package main

import (
	"log"

	"github.com/ezhdanovskiy/companies/internal/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
