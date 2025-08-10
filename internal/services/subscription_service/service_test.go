package subscription_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sunr3d/subscription-aggregator/internal/interfaces/infra"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/services"
	"github.com/sunr3d/subscription-aggregator/internal/services/subscription_service"
	"github.com/sunr3d/subscription-aggregator/mocks"
	"github.com/sunr3d/subscription-aggregator/models"
)

func ym(y int, m time.Month) time.Time {
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

// CREATE Tests
func TestService_Create_OK(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     nil,
	}

	repo.EXPECT().Create(ctx, in).Return(1, nil)

	id, err := svc.Create(ctx, in)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestService_Create_OK_EndDateAfterStartDate(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	endDate := ym(2025, time.August)
	in := models.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     &endDate,
	}

	repo.EXPECT().Create(ctx, in).Return(1, nil)

	id, err := svc.Create(ctx, in)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestService_Create_OK_PriceZero(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ServiceName: "Yandex Plus",
		Price:       0,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
	}

	repo.EXPECT().Create(ctx, in).Return(1, nil)

	id, err := svc.Create(ctx, in)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestService_Create_ErrValidation_PriceNegative(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ServiceName: "Yandex Plus",
		Price:       -100,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
	}

	_, err := svc.Create(ctx, in)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrValidation))
}

func TestService_Create_ErrValidation_EndDateBeforeStartDate(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	endDate := ym(2025, time.June)
	in := models.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     &endDate,
	}

	_, err := svc.Create(ctx, in)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrValidation))
}

func TestService_Create_ErrDatabase(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
	}

	repo.EXPECT().Create(ctx, in).Return(-1, errors.New("ошибка БД"))

	_, err := svc.Create(ctx, in)
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}

// READ Tests
func TestService_GetByID_OK(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	repo.EXPECT().GetByID(ctx, 1).Return(models.Subscription{
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     nil,
	}, nil)

	sub, err := svc.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, 1, sub.ID)
	require.Equal(t, "Yandex Plus", sub.ServiceName)
	require.Equal(t, 400, sub.Price)
	require.Equal(t, "u-1", sub.UserID)
	require.Equal(t, ym(2025, time.July), sub.StartDate)
	require.Nil(t, sub.EndDate)
}

func TestService_GetByID_ErrNotFound(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	repo.EXPECT().GetByID(ctx, 1).Return(models.Subscription{}, infra.ErrNotFound)

	_, err := svc.GetByID(ctx, 1)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrNotFound))
}

func TestService_GetByID_ErrDatabase(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	repo.EXPECT().GetByID(ctx, 1).Return(models.Subscription{}, errors.New("ошибка БД"))

	_, err := svc.GetByID(ctx, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}

// UPDATE Tests
func TestService_Update_OK(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ID: 1,
		ServiceName: "Yandex Plus",
		Price: 400,
		UserID: "u-1",
		StartDate: ym(2025, time.July),
		EndDate: nil,
	}

	repo.EXPECT().Update(ctx, in).Return(nil)

	err := svc.Update(ctx, in)
	require.NoError(t, err)
}

func TestService_Update_ErrValidation_PriceNegative(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ID: 1,
		ServiceName: "Yandex Plus",
		Price: -100,
		UserID: "u-1",
		StartDate: ym(2025, time.July),
		EndDate: nil,
	}

	err := svc.Update(ctx, in)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrValidation))
}

func TestService_Update_ErrValidation_EndDateBeforeStartDate(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	endDate := ym(2025, time.June)
	in := models.Subscription{
		ID: 1,
		ServiceName: "Yandex Plus",
		Price: 400,
		UserID: "u-1",
		StartDate: ym(2025, time.July),
		EndDate: &endDate,
	}

	err := svc.Update(ctx, in)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrValidation))
}

func TestService_Update_ErrNotFound(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ID: 1,
		ServiceName: "Yandex Plus",
		Price: 400,
		UserID: "u-1",
		StartDate: ym(2025, time.July),
		EndDate: nil,
	}

	repo.EXPECT().Update(ctx, in).Return(infra.ErrNotFound)

	err := svc.Update(ctx, in)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrNotFound))
}

func TestService_Update_ErrDatabase(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	in := models.Subscription{
		ID: 1,
		ServiceName: "Yandex Plus",
		Price: 400,
		UserID: "u-1",
		StartDate: ym(2025, time.July),
		EndDate: nil,
	}

	repo.EXPECT().Update(ctx, in).Return(errors.New("ошибка БД"))

	err := svc.Update(ctx, in)
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}

// DELETE Tests
// TODO: Добавить тесты для Delete

// LIST Tests
// TODO: Добавить тесты для List