package subscription_service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		return -1, fmt.Errorf("%w: price не может быть отрицательным", services.ErrValidation)
	}
	if data.EndDate != nil && data.EndDate.Before(data.StartDate) {
		return -1, fmt.Errorf("%w: end_date не может быть раньше start_date", services.ErrValidation)
	}
	return s.repo.Create(ctx, data)
}

func (s *subscriptionService) GetByID(ctx context.Context, id int) (models.Subscription, error) {
	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrNotFound) {
			return models.Subscription{}, services.ErrNotFound
		}
		return models.Subscription{}, fmt.Errorf("service GetByID(): %w", err)
	}
	return res, nil
}

func (s *subscriptionService) Update(ctx context.Context, data models.Subscription) error {
	if data.Price < 0 {
		return fmt.Errorf("%w: price не может быть отрицательным", services.ErrValidation)
	}
	if data.EndDate != nil && data.EndDate.Before(data.StartDate) {
		return fmt.Errorf("%w: end_date не может быть раньше start_date", services.ErrValidation)
	}
	if err := s.repo.Update(ctx, data); err != nil {
		if errors.Is(err, infra.ErrNotFound) {
			return services.ErrNotFound
		}
		return fmt.Errorf("service Update(): %w", err)
	}
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, infra.ErrNotFound) {
			return services.ErrNotFound
		}
		return fmt.Errorf("service Delete(): %w", err)
	}
	return nil
}

func (s *subscriptionService) List(ctx context.Context, filter services.ListFilter) ([]models.Subscription, error) {
	var uid, sname *string
	if filter.HasUserID {
		uid = &filter.UserID
	}
	if filter.HasServiceName {
		sname = &filter.ServiceName
	}

	return s.repo.List(ctx, infra.ListFilter{
		UserID:      uid,
		ServiceName: sname,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
	})
}

func (s *subscriptionService) TotalCost(ctx context.Context, periodStart, periodEnd time.Time, filter services.ListFilter) (int, error) {
	normalizeDate := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	}

	ps := normalizeDate(periodStart)
	pe := normalizeDate(periodEnd)

	if pe.Before(ps) {
		return 0, fmt.Errorf("%w: end_date не может быть раньше start_date", services.ErrValidation)
	}

	filter.Limit, filter.Offset = 0, 0

	data, err := s.List(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("service TotalCost(): %w", err)
	}

	sum := 0
	for _, item := range data {
		start := normalizeDate(item.StartDate)
		if start.Before(ps) {
			start = ps
		}

		end := pe
		if item.EndDate != nil {
			e := normalizeDate(*item.EndDate)
			if e.Before(end) {
				end = e
			}
		}
		if end.Before(start) {
			continue
		}

		months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month()) + 1
		if months > 0 {
			sum += months * item.Price
		}

	}
	return sum, nil
}
