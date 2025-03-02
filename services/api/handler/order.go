package handler

import (
	"net/http"
	"strconv"

	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type createOrderInput struct {
	UserID        int32  `json:"user_id" validate:"required"`
	Total         int32  `json:"total" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	ProductID     int32  `json:"product_id" validate:"required"`
}

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

	return ctx.Status(http.StatusOK).JSON(createdOrder)
}

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

	return ctx.Status(http.StatusOK).JSON(order)
}

type deleteOrderResponce struct {
	Status string `json:"status"`
}

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
		return ctx.Status(http.StatusOK).JSON([]repository.Order{})
	}

	return ctx.Status(http.StatusOK).JSON(orders)
}
