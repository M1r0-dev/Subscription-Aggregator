package subscriptionservice

import (
	"context"

	"github.com/M1r0-dev/Subscription-Aggregator/internal/entity"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/repo"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/repo/persistence"
)

type SubscriptionUsecase struct {
	repo repo.SubscriptionRepo
}

func New(repo repo.SubscriptionRepo) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		repo: repo,
	}
}

func (u *SubscriptionUsecase) Store(ctx context.Context, sub *entity.Subscription) error {
	return u.repo.Store(ctx, sub)
}

func (u *SubscriptionUsecase) Get(ctx context.Context, id int) (*entity.Subscription, error) {
	return u.repo.Get(ctx, id)
}

func (u *SubscriptionUsecase) Update(ctx context.Context, sub *entity.Subscription) error {
	return u.repo.Update(ctx, sub)
}

func (u *SubscriptionUsecase) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

func (u *SubscriptionUsecase) List(ctx context.Context, opts ...persistence.ListOption) (*[]entity.Subscription, error) {
	return u.repo.List(ctx, opts...)
}
