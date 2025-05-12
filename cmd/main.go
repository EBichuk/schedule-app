package main

import (
	"log"
	"schedule-app/internal/app"
)

func main() {
	a, err := app.Init()
	if err != nil {
		log.Fatal(err)
	}

	a.RunApp()
}
