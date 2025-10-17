package parser

import (
	"strconv"
	"time"

	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/dto"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/entity"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscriptionParser struct {
	logger logger.Interface
}

func New(logger logger.Interface) *SubscriptionParser {
	return &SubscriptionParser{
		logger: logger,
	}
}

func (p *SubscriptionParser) ParseStoreRequest(ctx *fiber.Ctx) (*entity.Subscription, error) {
	var req dto.StoreSubscriptionHandlerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.ServiceName == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Service name is required")
	}
	if req.Price == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Price is required")
	}
	if req.UserId == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "User ID is required")
	}
	if req.StartDate == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Start date is required")
	}

	price, err := strconv.ParseUint(req.Price, 10, 64)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid price format")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format, must be UUID")
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start date format. Use RFC3339 format")
	}

	var endDate time.Time
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end date format. Use RFC3339 format")
		}
		endDate = parsedEndDate
	} else {
		endDate = time.Time{}
	}

	return &entity.Subscription{
		ServiceName: req.ServiceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}


func (p *SubscriptionParser) ParseUpdateRequest(ctx *fiber.Ctx, existingSub *entity.Subscription) error {
	var req dto.UpdateSubscriptionHandlerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ServiceName != "" {
		existingSub.ServiceName = req.ServiceName
	}

	if req.Price != "" {
		price, err := strconv.ParseUint(req.Price, 10, 64)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid price format")
		}
		existingSub.Price = price
	}

	if req.UserId != "" {
		userID, err := uuid.Parse(req.UserId)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format, must be UUID")
		}
		existingSub.UserID = userID 
	}

	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid start date format")
		}
		existingSub.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid end date format")
		}
		existingSub.EndDate = endDate
	}

	return nil
}

func (p *SubscriptionParser) ParseListRequest(ctx *fiber.Ctx) (*dto.ListSubscriptionsHandlerRequest, error) {
	var req dto.ListSubscriptionsHandlerRequest
	if err := ctx.QueryParser(&req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.SortBy == "" {
		req.SortBy = "start_date"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	if req.Page < 1 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Page must be greater than 0")
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Page size must be between 1 and 100")
	}

	if req.UserID != nil {
		if _, err := uuid.Parse(*req.UserID); err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format in filter, must be UUID")
		}
	}

	return &req, nil
}

func (p *SubscriptionParser) ParseGetRequest(ctx *fiber.Ctx) (int, error) {
	idStr := ctx.Params("id")
	if idStr == "" {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Subscription ID is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Invalid subscription ID format")
	}

	return id, nil
}

func (p *SubscriptionParser) ParseDeleteRequest(ctx *fiber.Ctx) (int, error) {
	return p.ParseGetRequest(ctx)
}

func (p *SubscriptionParser) ParseTotalCostRequest(ctx *fiber.Ctx) (*dto.TotalCostHandlerRequest, error) {
    var req dto.TotalCostHandlerRequest
    if err := ctx.QueryParser(&req); err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
    }

    if req.StartDate == nil || *req.StartDate == "" {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Start date is required")
    }
    if req.EndDate == nil || *req.EndDate == "" {
        return nil, fiber.NewError(fiber.StatusBadRequest, "End date is required")
    }

    if _, err := time.Parse("2006-01-02", *req.StartDate); err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD")
    }
    if _, err := time.Parse("2006-01-02", *req.EndDate); err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD")
    }

    if req.UserID != nil && *req.UserID != "" {
        if _, err := uuid.Parse(*req.UserID); err != nil {
            return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format in filter, must be UUID")
        }
    }

    return &req, nil
}