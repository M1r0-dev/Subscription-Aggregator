package usecase

import (
	"context"

	"github.com/M1r0-dev/Subscription-Aggregator/internal/entity"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/repo/persistence"
)

type SubscriptionUsecase interface {
	Store(ctx context.Context, sub *entity.Subscription) error
	Get(ctx context.Context, id int) (*entity.Subscription, error)
	Update(cxt context.Context, sub *entity.Subscription) error
	Delete(cxt context.Context, id int) error
	List(cxt context.Context, opts ...persistence.ListOption) (*[]entity.Subscription, error)
}
