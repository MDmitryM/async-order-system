package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) OrderCreate(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/create/?orderID")
}

func (h *Handler) OrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/details/?orderID")
}

func (h *Handler) OrderDelete(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/delete/?orderID")
}
