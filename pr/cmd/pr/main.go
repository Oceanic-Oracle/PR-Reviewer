package main

import "pr/internal/app"

func main() {
	app := app.NewBootstrap()
	app.Run()
}
