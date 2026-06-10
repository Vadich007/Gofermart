package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Vadich007/Gofermart/internal/app"
	"github.com/Vadich007/Gofermart/internal/config"
)

func main() {
	cfg := config.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
