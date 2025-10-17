package handler

import (
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/dto"
	"github.com/gofiber/fiber/v2"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(dto.ErrorResponse{Error: msg})
}