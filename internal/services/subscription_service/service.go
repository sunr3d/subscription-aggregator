package subscription_service

import (
	"context"
	"errors"
	"fmt"

	infraErr "github.com/sunr3d/subscription-aggregator/internal/infra"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/infra"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/services"
	"github.com/sunr3d/subscription-aggregator/models"
)

var _ services.SubscriptionService = (*subscriptionService)(nil)

type subscriptionService struct {
	repo infra.Database
}

func New(repo infra.Database) services.SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(ctx context.Context, data models.Subscription) (int, error) {
	if data.Price < 0 {
		return -1, fmt.Errorf("price не может быть отрицательным")
	}
	if data.EndDate != nil && data.EndDate.Before(data.StartDate) {
		return -1, fmt.Errorf("end_date не может быть раньше start_date")
	}
	return s.repo.Create(ctx, data)
}

func (s *subscriptionService) GetByID(ctx context.Context, id int) (models.Subscription, error) {
	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infraErr.ErrNotFound) {
			return models.Subscription{}, ErrNotFound
		}
		return models.Subscription{}, fmt.Errorf("postgres GetByID(): %w", err)
	}
	return res, nil
}

func (s *subscriptionService) Update(ctx context.Context, data models.Subscription) error {
	if data.Price < 0 {
		return fmt.Errorf("price не может быть отрицательным")
	}
	if data.EndDate != nil && data.EndDate.Before(data.StartDate) {
		return fmt.Errorf("end_date не может быть раньше start_date")
	}
	if err := s.repo.Update(ctx, data); err != nil {
		if errors.Is(err, infraErr.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("postgres Update(): %w", err)
	}
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, infraErr.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("postgres Delete(): %w", err)
	}
	return nil
}

func (s *subscriptionService) List(ctx context.Context, filter services.ListFilter) ([]models.Subscription, error) {
	return s.repo.List(ctx, infra.ListFilter{
		UserID:      filter.UserID,
		ServiceName: filter.ServiceName,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
	})
}
