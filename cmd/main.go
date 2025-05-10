package main

import (
	"log"
	"schedule-app/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	a.RunApp()
}
