package handler

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func InitRouts(app *fiber.App) {
	order := app.Group("/order") //host/order

	orderCreate := order.Group("/create") //order/create
	orderCreate.Post("/", OrderCreate)    //order/create/?orderID

	orderDetails := order.Group("/details") //order/detils
	orderDetails.Get("/", OrderDetails)     //order/details/?orderID

	orderDelete := order.Group("/delete") //order/delete
	orderDelete.Delete("/", OrderDelete)  //order/delete/?orderID
}
