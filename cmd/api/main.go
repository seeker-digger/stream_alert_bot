package main

import "main.go/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
