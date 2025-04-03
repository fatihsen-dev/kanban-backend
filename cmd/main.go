package main

import (
	"context"
	"fmt"

	"github.com/fatihsen-dev/kanban-backend/app/item"
	"github.com/fatihsen-dev/kanban-backend/config"
	"github.com/fatihsen-dev/kanban-backend/infra/postgres"
	"github.com/gofiber/fiber/v2"
)

func main() {
	appConfig := config.Read()
	ctx := context.Background()

	pg, err := postgres.NewPG(ctx, appConfig.DatabaseUrl)
	if err != nil {
		panic(fmt.Sprint("Database connection error: %w", err))
	}

	app := fiber.New()

	createItemHandler := item.NewCreateItemHandler(pg)

	app.Post("/items", createItemHandler.Handle)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", appConfig.Port)); err != nil {
		panic(fmt.Sprint("http server starting error: %w", err))
	}
}
