package repo

import (
	"context"

	"github.com/M1r0-dev/Subscription-Aggregator/internal/entity"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/repo/persistence"
)

type SubscriptionRepo interface {
	Store(cxt context.Context, sub *entity.Subscription) error
	Get(cxt context.Context, id int) (*entity.Subscription, error)
	Update(cxt context.Context, sub *entity.Subscription) error
	Delete(cxt context.Context, id int) error
	List(cxt context.Context, opts ...persistence.ListOptions) (*[]entity.Subscription, error)
}
