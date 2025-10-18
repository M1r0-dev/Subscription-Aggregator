package handler

import (
	"github.com/M1r0-dev/Subscription-Aggregator/internal/repo/persistence"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Store creates a new subscription
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.StoreSubscriptionHandlerRequest true "Subscription data"
// @Success 201 {object} dto.StoreSubscriptionHandlerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
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


// Get retrieves a subscription by ID
// @Summary Get subscription
// @Description Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} dto.GetSubscriptionHandlerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(ctx *fiber.Ctx) error {
	const op = "handler.Get"

	id, err := h.parser.ParseGetRequest(ctx)
	if err != nil {
		h.logger.Error("failed to parse get request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	sub, err := h.usecase.Get(ctx.Context(), id)
	if err != nil {
		h.logger.Error("failed to get subscription", "operation", op, "id", id, "error", err)
		return errorResponse(ctx, fiber.StatusNotFound, "Subscription not found")
	}

	response := h.mapper.ToGetResponse(sub)

	h.logger.Info("subscription retrieved successfully",
		"operation", op,
		"subscription_id", sub.Id,
	)

	return ctx.Status(fiber.StatusOK).JSON(response)
}


// Update updates a subscription
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body dto.UpdateSubscriptionHandlerRequest true "Subscription data"
// @Success 200 {object} dto.GetSubscriptionHandlerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(ctx *fiber.Ctx) error {
	const op = "handler.Update"

	id, err := h.parser.ParseGetRequest(ctx)
	if err != nil {
		h.logger.Error("failed to parse update request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	existingSub, err := h.usecase.Get(ctx.Context(), id)
	if err != nil {
		h.logger.Error("subscription not found for update", "operation", op, "id", id, "error", err)
		return errorResponse(ctx, fiber.StatusNotFound, "Subscription not found")
	}

	err = h.parser.ParseUpdateRequest(ctx, existingSub)
	if err != nil {
		h.logger.Error("failed to parse update request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	err = h.usecase.Update(ctx.Context(), existingSub)
	if err != nil {
		h.logger.Error("failed to update subscription", "operation", op, "id", id, "error", err)
		return errorResponse(ctx, fiber.StatusInternalServerError, "Failed to update subscription")
	}

	response := h.mapper.ToUpdateResponse(existingSub)

	h.logger.Info("subscription updated successfully",
		"operation", op,
		"subscription_id", existingSub.Id,
		"user_id", existingSub.UserID,
	)

	return ctx.Status(fiber.StatusOK).JSON(response)
}


// Delete deletes a subscription
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(ctx *fiber.Ctx) error {
	const op = "handler.Delete"

	id, err := h.parser.ParseDeleteRequest(ctx)
	if err != nil {
		h.logger.Error("failed to parse delete request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	err = h.usecase.Delete(ctx.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete subscription", "operation", op, "id", id, "error", err)
		return errorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete subscription")
	}

	h.logger.Info("subscription deleted successfully",
		"operation", op,
		"subscription_id", id,
	)

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}


// List retrieves subscriptions with filtering and pagination
// @Summary List subscriptions
// @Description Get list of subscriptions with filtering and pagination
// @Tags subscriptions
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Page size" default(10) minimum(1) maximum(100)
// @Param user_id query string false "User ID filter (UUID)"
// @Param service_name query string false "Service name filter"
// @Param sort_by query string false "Sort field" default(start_date)
// @Param sort_order query string false "Sort order" default(desc) Enums(asc, desc)
// @Success 200 {object} dto.ListSubscriptionsHandlerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(ctx *fiber.Ctx) error {
	const op = "handler.List"

	req, err := h.parser.ParseListRequest(ctx)
	if err != nil {
		h.logger.Error("failed to parse list request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	opts := []persistence.ListOption{}
	
	if req.UserID != nil {
		userID, err := uuid.Parse(*req.UserID)
		if err != nil {
			h.logger.Error("failed to parse list request", "operation", op, "error", err)
			return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
		}
		opts = append(opts, persistence.WithUserID(userID))
	}
	if req.ServiceName != nil {
		opts = append(opts, persistence.WithServiceName(*req.ServiceName))
	}

	subscriptions, err := h.usecase.List(ctx.Context(), opts...)
	if err != nil {
		h.logger.Error("failed to list subscriptions", "operation", op, "error", err)
		return errorResponse(ctx, fiber.StatusInternalServerError, "Failed to get subscriptions list")
	}

	total := len(subscriptions)
	response := h.mapper.ToListResponse(subscriptions, total, req.Page, req.PageSize)

	h.logger.Info("subscriptions listed successfully",
		"operation", op,
		"count", len(subscriptions),
		"page", req.Page,
	)

	return ctx.Status(fiber.StatusOK).JSON(response)
}


// GetTotalCost calculates total cost of subscriptions for a period
// @Summary Get total cost
// @Description Calculate total cost of subscriptions for a specific period with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID filter (UUID)"
// @Param service_name query string false "Service name filter"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.TotalCostHandlerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) GetTotalCost(ctx *fiber.Ctx) error {
	const op = "handler.GetTotalCost"

	req, err := h.parser.ParseTotalCostRequest(ctx)
	if err != nil {
		h.logger.Error("failed to parse total cost request", "operation", op, "error", err)
		return errorResponse(ctx, err.(*fiber.Error).Code, err.Error())
	}

	total, err := h.usecase.GetTotalCost(ctx.Context(), req.UserID, req.ServiceName, *req.StartDate, *req.EndDate)
	if err != nil {
		h.logger.Error("failed to calculate total cost", "operation", op, "error", err)
		return errorResponse(ctx, fiber.StatusInternalServerError, "Failed to calculate total cost")
	}

	response := h.mapper.ToTotalCostResponse(total, req)

	h.logger.Info("total cost calculated successfully",
		"operation", op,
		"user_id", req.UserID,
		"service_name", req.ServiceName,
		"start_date", *req.StartDate,
		"end_date", *req.EndDate,
		"total_cost", total,
	)

	return ctx.Status(fiber.StatusOK).JSON(response)
}
