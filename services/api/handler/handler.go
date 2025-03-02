package handler

import (
	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) InitRouts(app *fiber.App) {
	order := app.Group("/order") //host/order

	orderCreate := order.Group("/create") //order/create
	orderCreate.Post("/", h.OrderCreate)  //order/create/?orderID

	orderDetails := order.Group("/details") //order/detils
	orderDetails.Get("/", h.OrderDetails)   //order/details/?orderID

	orderDelete := order.Group("/delete")  //order/delete
	orderDelete.Delete("/", h.OrderDelete) //order/delete/?orderID
}
