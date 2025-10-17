package mapper

import (
	"strconv"
	"time"

	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/dto"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/entity"
)

type SubscriptionMapper struct{}

func New() *SubscriptionMapper {
	return &SubscriptionMapper{}
}

func (m *SubscriptionMapper) ToStoreResponse(sub *entity.Subscription) dto.StoreSubscriptionHandlerResponse {
	return dto.StoreSubscriptionHandlerResponse{
		Id: strconv.FormatInt(sub.Id, 10),
	}
}

func (m *SubscriptionMapper) ToGetResponse(sub *entity.Subscription) dto.GetSubscriptionHandlerResponse {
	response := dto.GetSubscriptionHandlerResponse{
		ServiceName: sub.ServiceName,
		Price:       strconv.FormatUint(sub.Price, 10),
		UserId:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format(time.RFC3339),
	}

	if !sub.EndDate.IsZero() {
		response.EndDate = sub.EndDate.Format(time.RFC3339)
	}

	return response
}

func (m *SubscriptionMapper) ToListResponse(
	subscriptions []*entity.Subscription,
	total int,
	page int,
	pageSize int,
) dto.ListSubscriptionsHandlerResponse {
	response := dto.ListSubscriptionsHandlerResponse{
		Subscriptions: make([]dto.SubscriptionItem, len(subscriptions)),
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    (total + pageSize - 1) / pageSize,
	}

	for i, sub := range subscriptions {
		item := dto.SubscriptionItem{
			ID:          strconv.FormatInt(sub.Id, 10),
			ServiceName: sub.ServiceName,
			Price:       strconv.FormatUint(sub.Price, 10),
			UserID:      sub.UserID.String(),
			StartDate:   sub.StartDate.Format(time.RFC3339),
		}

		if !sub.EndDate.IsZero() {
			item.EndDate = sub.EndDate.Format(time.RFC3339)
		}

		response.Subscriptions[i] = item
	}

	return response
}

func (m *SubscriptionMapper) ToUpdateResponse(sub *entity.Subscription) dto.GetSubscriptionHandlerResponse {
    return m.ToGetResponse(sub)
}

func (m *SubscriptionMapper) ToTotalCostResponse(total uint64, req *dto.TotalCostHandlerRequest) dto.TotalCostHandlerResponse {
    response := dto.TotalCostHandlerResponse{
        TotalCost: total,
        Period: dto.Period{
            StartDate: *req.StartDate,
            EndDate:   *req.EndDate,
        },
        Filters: dto.TotalCostFilters{
            UserID:      req.UserID,
            ServiceName: req.ServiceName,
        },
    }
    
    return response
}
