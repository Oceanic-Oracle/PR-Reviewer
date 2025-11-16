package main

import (
	"pr/internal/app"
	"pr/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app := app.NewBootstrap(cfg)
	app.Run()
}
