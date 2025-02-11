package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func OrderCreate(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/create/?orderID")
}

func OrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/details/?orderID")
}

func OrderDelete(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/delete/?orderID")
}
