package main

import (
	"log"
	"schedule-app/internal/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	a.RunApp()
}
