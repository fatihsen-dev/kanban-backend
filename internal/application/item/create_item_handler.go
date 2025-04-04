package item

import (
	"fmt"

	"github.com/fatihsen-dev/kanban-backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type CreateItemHandler struct {
	repository Repository
}

func NewCreateItemHandler(repository Repository) *CreateItemHandler {
	return &CreateItemHandler{
		repository: repository,
	}
}

func (h *CreateItemHandler) Handle(c *fiber.Ctx) error {
	item := domain.Item{}

	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"status":  fiber.StatusBadRequest,
		})
	}

	err := h.repository.InsertItem(c.UserContext(), &item)
	if err != nil {
		fmt.Print(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "item create error",
			"status":  fiber.StatusInternalServerError,
		})
	}

	return c.Status(201).JSON(item)
}
