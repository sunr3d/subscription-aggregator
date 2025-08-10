package services

import (
	"context"

	"github.com/sunr3d/subscription-aggregator/models"
)

type ListFilter struct {
	UserID      string
	HasUserID bool
	ServiceName string
	HasServiceName bool
	Limit       int
	Offset      int
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=SubscriptionService --output=../../../mocks --filename=mock_subscription_service.go --with-expecter
type SubscriptionService interface {
	Create(ctx context.Context, data models.Subscription) (int, error)
	GetByID(ctx context.Context, id int) (models.Subscription, error)
	Update(ctx context.Context, data models.Subscription) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter ListFilter) ([]models.Subscription, error)
}
