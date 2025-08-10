package services

import (
	"context"

	"github.com/sunr3d/subscription-aggregator/models"
)

type ListFilter struct {
	UserID      *string
	ServiceName *string
	Limit       int
	Offset      int
}

type SubscriptionService interface {
	Create(ctx context.Context, data models.Subscription) (int, error)
	GetByID(ctx context.Context, id int) (models.Subscription, error)
	Update(ctx context.Context, data models.Subscription) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter ListFilter) ([]models.Subscription, error)
}
