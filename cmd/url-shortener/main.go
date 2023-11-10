package main

import (
	"fmt"

	"github.com/alexapps/url-shortener/internal/config"
)

func main() {
	// read config
	cfg := config.MustLoad()

	// TODO: remove fmt
	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: storage: sqllite

	// TODO: router: chi, "chi render"

	// TODO: run server
}
