package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var ErrNotEnoughParamsInQuery = errors.New("not enough params in query")
var ErrQueryParamConversionError = errors.New("not enough params in query")
var ErrRecordNotFound = errors.New("record not found")

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendJsonError(ctx *fiber.Ctx, status int, msg string) error {
	logrus.Error(msg)
	return ctx.Status(status).JSON(ErrorResponse{msg})
}
