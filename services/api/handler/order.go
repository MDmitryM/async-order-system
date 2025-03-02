package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// type createOrderInput struct {
// 	UserID        int32  `json:"user_id" validate:"required"`
// 	Total         int32  `json:"total" validate:"required"`
// 	PaymentMethod string `json:"payment_method" validate:"required"`
// }

func (h *Handler) OrderCreate(ctx *fiber.Ctx) error {
	// var input createOrderInput
	// if err := ctx.BodyParser(&input)

	return ctx.Status(http.StatusOK).JSON("order/create/?orderID")
}

func (h *Handler) OrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/details/?orderID")
}

func (h *Handler) OrderDelete(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON("order/delete/?orderID")
}
