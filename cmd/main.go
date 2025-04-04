package main

import (
	"context"
	"fmt"

	"github.com/fatihsen-dev/kanban-backend/config"
	"github.com/fatihsen-dev/kanban-backend/internal/application/item"
	"github.com/fatihsen-dev/kanban-backend/internal/infra/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	appConfig := config.Read()
	ctx := context.Background()

	pg, err := postgres.NewPG(ctx, appConfig.DatabaseUrl)
	if err != nil {
		panic(fmt.Sprint("Database connection error: %w", err))
	}

	app := fiber.New()

	app.Use("/ws/projects/:id", websocket.New(func(c *websocket.Conn) {
		fmt.Println("Connected to websocket")
	}))

	createItemHandler := item.NewCreateItemHandler(pg)

	app.Post("/items", createItemHandler.Handle)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", appConfig.Port)); err != nil {
		panic(fmt.Sprint("http server starting error: %w", err))
	}
}
