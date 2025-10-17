package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *SubscriptionHandler) Store(ctx *fiber.Ctx) error {
	const op = "handler.Store"
	sub, err := h.parser.ParseStoreRequest(ctx)
	if err != nil {
		h.logger.Error("failed to parse store request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	err = h.usecase.Store(ctx.Context(), sub)
	if err != nil {
		h.logger.Error("failed to store subscription", "operation", op, "error", err)
		return errorResponse(ctx, fiber.StatusInternalServerError, "Failed to create subscription")
	}

	response := h.mapper.ToStoreResponse(sub)

	h.logger.Info("subscription created successfully", 
		"operation", op, 
		"subscription_id", sub.Id,
		"user_id", sub.UserID,
	)

	return ctx.Status(fiber.StatusCreated).JSON(response)
}
