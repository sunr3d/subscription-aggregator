package infra

import (
	"context"

	"github.com/sunr3d/subscription-aggregator/models"
)

type ListFilter struct {
	// TODO: Фильтр для List метода
}

type Database interface {
	Create(ctx context.Context, data models.Subscription) (string, error) // Create (C)
	GetByID(ctx context.Context, id int) (models.Subscription, error) // Read (R)
	Update(ctx context.Context, data models.Subscription) error // Update (U)
	Delete(ctx context.Context, id int) error // Delete (D)
	List(ctx context.Context, filter ListFilter) ([]models.Subscription, error) // List (L)
}
