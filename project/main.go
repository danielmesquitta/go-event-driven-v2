package main

import (
	"context"

	"tickets/internal/app"
)

func main() {
	err := app.New().Run(context.Background())
	if err != nil {
		panic(err)
	}
}
