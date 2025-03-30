package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MDmitryM/async-order-system/services/api/kafka"
	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type createOrderInput struct {
	UserID        int32  `json:"user_id" validate:"required" example:"1"`
	Total         int32  `json:"total" validate:"required" example:"15000"`
	PaymentMethod string `json:"payment_method" validate:"required" example:"SBP"`
	ProductID     int32  `json:"product_id" validate:"required" example:"1499"`
}

// trying to avoid pgx.Timestamptz
type orderDetailsResponce struct {
	ID            int32     `json:"id"`
	UserID        int32     `json:"user_id"`
	Total         int32     `json:"total"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	PaymentMethod string    `json:"payment_method"`
	ProductID     int32     `json:"product_id"`
}

// @Summary Create order
// @Description Create order, produces kafka message
// @Tags orderAPI
// @Accept json
// @Produce json
// @Param order body createOrderInput true "Order"
// @Success 200 {object} orderDetailsResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /order/create/ [post]
func (h *Handler) OrderCreate(ctx *fiber.Ctx) error {
	var input createOrderInput
	if err := ctx.BodyParser(&input); err != nil {
		SendJsonError(ctx, http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(&input); err != nil {
		return SendJsonError(ctx, http.StatusBadRequest, err.Error())
	}

	createdOrder, err := h.repo.CreateOrder(ctx.Context(), repository.CreateOrderParams{
		UserID:        input.UserID,
		Total:         input.Total,
		Status:        "created",
		PaymentMethod: input.PaymentMethod,
		ProductID:     input.ProductID,
	})
	if err != nil {
		return SendJsonError(ctx, http.StatusInternalServerError, err.Error())
	}

	kafkaOrderModel := kafka.OrderMessage{
		ID:            createdOrder.ID,
		UserID:        createdOrder.UserID,
		Total:         createdOrder.Total,
		Status:        createdOrder.Status,
		PaymentMethod: createdOrder.PaymentMethod,
		ProductID:     createdOrder.ProductID,
	}

	if err := h.producer.SendOrder(ctx.Context(), kafkaOrderModel); err != nil {
		return SendJsonError(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(orderDetailsResponce{
		ID:            createdOrder.ID,
		UserID:        createdOrder.UserID,
		Total:         createdOrder.Total,
		Status:        createdOrder.Status,
		CreatedAt:     createdOrder.CreatedAt.Time,
		PaymentMethod: createdOrder.PaymentMethod,
		ProductID:     createdOrder.ProductID,
	})
}

// @Summary Order details
// @Description Returns order info by orderID
// @Tags orderAPI
// @Produce json
// @Param orderID query string true "Order ID"
// @Success 200 {object} orderDetailsResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /order/details/ [get]
func (h *Handler) OrderDetails(ctx *fiber.Ctx) error {
	orderID := ctx.Query("orderID")
	if orderID == "" {
		return SendJsonError(ctx, http.StatusBadRequest, ErrNotEnoughParamsInQuery.Error())
	}

	intOrderID, err := strconv.Atoi(orderID)
	if err != nil {
		return SendJsonError(ctx, http.StatusBadRequest, ErrQueryParamConversionError.Error())
	}

	order, err := h.repo.GetOrderByID(ctx.Context(), int32(intOrderID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return SendJsonError(ctx, http.StatusNotFound, ErrRecordNotFound.Error())
		}
		return SendJsonError(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(orderDetailsResponce{
		ID:            order.ID,
		UserID:        order.UserID,
		Total:         order.Total,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt.Time,
		PaymentMethod: order.PaymentMethod,
		ProductID:     order.ProductID,
	})
}

type deleteOrderResponce struct {
	Status string `json:"status"`
}

// @Summary Delete order
// @Description Deletes order by orderID
// @Tags orderAPI
// @Produce json
// @Param orderID query string true "Order ID"
// @Success 200 {object} deleteOrderResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /order/delete/ [delete]
func (h *Handler) OrderDelete(ctx *fiber.Ctx) error {
	orderID := ctx.Query("orderID")
	if orderID == "" {
		return SendJsonError(ctx, http.StatusBadRequest, ErrNotEnoughParamsInQuery.Error())
	}

	intOrderID, err := strconv.Atoi(orderID)
	if err != nil {
		return SendJsonError(ctx, http.StatusBadRequest, ErrQueryParamConversionError.Error())
	}

	rowsAffected, err := h.repo.DeleteOrderByID(ctx.Context(), int32(intOrderID))
	if err != nil {
		return SendJsonError(ctx, http.StatusInternalServerError, err.Error())
	}
	if rowsAffected == 0 {
		return SendJsonError(ctx, http.StatusNotFound, ErrRecordNotFound.Error())
	}

	return ctx.Status(http.StatusOK).JSON(deleteOrderResponce{"ok"})
}

// @Summary List orders
// @Description Returns list of orders
// @Tags orderAPI
// @Produce json
// @Param page query string true "Page"
// @Param pageSize query string true "Page size"
// @Success 200 {object} []orderDetailsResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /order/list/ [get]
func (h *Handler) OrderList(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		return SendJsonError(ctx, http.StatusBadRequest, err.Error())
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "5"))
	if err != nil || pageSize < 1 {
		return SendJsonError(ctx, http.StatusBadRequest, err.Error())
	}

	offset := (page - 1) * pageSize

	orders, err := h.repo.ListOrders(ctx.Context(), repository.ListOrdersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return SendJsonError(ctx, http.StatusInternalServerError, err.Error())
	}
	if orders == nil {
		return ctx.Status(http.StatusOK).JSON([]orderDetailsResponce{})
	}

	ordersResponce := make([]orderDetailsResponce, len(orders))
	for i, order := range orders {
		ordersResponce[i] = orderDetailsResponce{
			ID:            order.ID,
			UserID:        order.UserID,
			Total:         order.Total,
			Status:        order.Status,
			CreatedAt:     order.CreatedAt.Time,
			PaymentMethod: order.PaymentMethod,
			ProductID:     order.ProductID,
		}
	}

	return ctx.Status(http.StatusOK).JSON(ordersResponce)
}
