package infra

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

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=Database --output=../../../mocks --filename=mock_database.go --with-expecter
type Database interface {
	Create(ctx context.Context, data models.Subscription) (int, error)          // Create (C)
	GetByID(ctx context.Context, id int) (models.Subscription, error)           // Read (R)
	Update(ctx context.Context, data models.Subscription) error                 // Update (U)
	Delete(ctx context.Context, id int) error                                   // Delete (D)
	List(ctx context.Context, filter ListFilter) ([]models.Subscription, error) // List (L)
}
