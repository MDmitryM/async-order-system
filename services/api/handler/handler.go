package handler

import (
	_ "github.com/MDmitryM/async-order-system/services/api/docs"
	"github.com/MDmitryM/async-order-system/services/api/kafka"
	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Handler struct {
	repo     repository.Repository
	producer *kafka.Producer
}

func NewHandler(repo repository.Repository, prod *kafka.Producer) *Handler {
	return &Handler{
		repo:     repo,
		producer: prod,
	}
}

func (h *Handler) InitRouts(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	order := app.Group("/order") //host/order

	orderCreate := order.Group("/create") //order/create
	orderCreate.Post("/", h.OrderCreate)  //order/create/?orderID

	orderDetails := order.Group("/details") //order/detils
	orderDetails.Get("/", h.OrderDetails)   //order/details/?orderID

	orderDelete := order.Group("/delete")  //order/delete
	orderDelete.Delete("/", h.OrderDelete) //order/delete/?orderID

	orderList := order.Group("/list") //order/list
	orderList.Get("/", h.OrderList)   //order/list
}
