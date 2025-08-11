package subscription_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sunr3d/subscription-aggregator/internal/interfaces/infra"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/services"
	"github.com/sunr3d/subscription-aggregator/internal/services/subscription_service"
	"github.com/sunr3d/subscription-aggregator/mocks"
	"github.com/sunr3d/subscription-aggregator/models"
)

func ym(y int, m time.Month) time.Time {
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
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
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     nil,
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
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       -100,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     nil,
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
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     &endDate,
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
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     nil,
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
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "u-1",
		StartDate:   ym(2025, time.July),
		EndDate:     nil,
	}

	repo.EXPECT().Update(ctx, in).Return(errors.New("ошибка БД"))

	err := svc.Update(ctx, in)
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}

// DELETE Tests
func TestService_Delete_OK(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	repo.EXPECT().Delete(ctx, 1).Return(nil)

	err := svc.Delete(ctx, 1)
	require.NoError(t, err)
}

func TestService_Delete_ErrNotFound(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	repo.EXPECT().Delete(ctx, 1).Return(infra.ErrNotFound)

	err := svc.Delete(ctx, 1)
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrNotFound))
}

func TestService_Delete_ErrDatabase(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	repo.EXPECT().Delete(ctx, 1).Return(errors.New("ошибка БД"))

	err := svc.Delete(ctx, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}

// LIST Tests
func TestService_List_OK_AllFilters(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	filter := services.ListFilter{
		UserID:         "u-1",
		HasUserID:      true,
		ServiceName:    "Yandex Plus",
		HasServiceName: true,
		Limit:          10,
		Offset:         20,
	}

	want := []models.Subscription{
		{
			ID:          1,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      "u-1",
			StartDate:   ym(2025, time.July),
			EndDate:     nil,
		},
	}

	repo.EXPECT().List(ctx, mock.MatchedBy(func(ifl infra.ListFilter) bool {
		return ifl.UserID != nil && *ifl.UserID == "u-1" &&
			ifl.ServiceName != nil && *ifl.ServiceName == "Yandex Plus" &&
			ifl.Limit == 10 && ifl.Offset == 20
	})).Return(want, nil)

	subs, err := svc.List(ctx, filter)
	require.NoError(t, err)
	require.Equal(t, want, subs)
}

func TestService_List_OK_UserOnly(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	filter := services.ListFilter{
		UserID:    "u-1",
		HasUserID: true,
		Limit:     25,
		Offset:    5,
	}

	repo.EXPECT().List(ctx, mock.MatchedBy(func(ifl infra.ListFilter) bool {
		return ifl.UserID != nil && *ifl.UserID == "u-1" &&
			ifl.ServiceName == nil && ifl.Limit == 25 && ifl.Offset == 5
	})).Return([]models.Subscription{}, nil)

	_, err := svc.List(ctx, filter)
	require.NoError(t, err)
}

func TestService_List_OK_NoFilters(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	filter := services.ListFilter{
		Limit:  10,
		Offset: 20,
	}

	repo.EXPECT().List(ctx, mock.MatchedBy(func(ifl infra.ListFilter) bool {
		return ifl.UserID == nil && ifl.ServiceName == nil && ifl.Limit == 10 && ifl.Offset == 20
	})).Return([]models.Subscription{}, nil)

	_, err := svc.List(ctx, filter)
	require.NoError(t, err)
}

func TestService_List_ErrDatabase(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	filter := services.ListFilter{
		Limit:  10,
		Offset: 0,
	}

	repo.EXPECT().List(ctx, mock.MatchedBy(func(ifl infra.ListFilter) bool {
		return ifl.UserID == nil && ifl.ServiceName == nil && ifl.Limit == 10 && ifl.Offset == 0
	})).Return(nil, errors.New("ошибка БД"))

	subs, err := svc.List(ctx, filter)
	require.Nil(t, subs)
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}

// TotalCost Tests
func TestService_TotalCost_OK_1(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	// Период: с января по март 2025, подписка: с января по марта 2025, цена: 400 - 400*3 = 1200
	periodStart, periodEnd := ym(2025, time.January), ym(2025, time.March)
	data := []models.Subscription{
		{
			ID:          1,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      "u-1",
			StartDate:   periodStart,
			EndDate:     &periodEnd,
		},
	}
	repo.EXPECT().List(ctx, mock.AnythingOfType("infra.ListFilter")).Return(data, nil)

	sum, err := svc.TotalCost(ctx, periodStart, periodEnd, services.ListFilter{})
	require.NoError(t, err)
	require.Equal(t, 1200, sum)
}

func TestService_TotalCost_OK_2(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	// Период: с января по март 2025, подписка: с февраля по апрель 2025, цена: 400, пересечение: февраль-март = 400*2 = 800
	periodStart, periodEnd := ym(2025, time.January), ym(2025, time.March)
	subStart, subEnd := ym(2025, time.February), ym(2025, time.April)
	data := []models.Subscription{
		{
			ID:          1,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      "u-1",
			StartDate:   subStart,
			EndDate:     &subEnd,
		},
	}

	repo.EXPECT().List(ctx, mock.AnythingOfType("infra.ListFilter")).Return(data, nil)

	sum, err := svc.TotalCost(ctx, periodStart, periodEnd, services.ListFilter{})
	require.NoError(t, err)
	require.Equal(t, 800, sum)
}

func TestService_TotalCost_OK_3(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	// Период: с января по март 2025, подписка: с декабря 2024 без конца, цена: 400*3 = 1200
	periodStart, periodEnd := ym(2025, time.January), ym(2025, time.March)
	subStart := ym(2024, time.December)
	data := []models.Subscription{
		{
			ID:          1,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      "u-1",
			StartDate:   subStart,
			EndDate:     nil,
		},
	}

	repo.EXPECT().List(ctx, mock.AnythingOfType("infra.ListFilter")).Return(data, nil)

	sum, err := svc.TotalCost(ctx, periodStart, periodEnd, services.ListFilter{})
	require.NoError(t, err)
	require.Equal(t, 1200, sum)
}

func TestService_TotalCost_OK_4(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	// Период: с января по март 2025, подписка: с октября 2024 по декабрь 2024, цена: 400 = 0 (период не пересекается)
	periodStart, periodEnd := ym(2025, time.January), ym(2025, time.March)
	subStart, subEnd := ym(2024, time.October), ym(2024, time.December)
	data := []models.Subscription{
		{
			ID:          1,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      "u-1",
			StartDate:   subStart,
			EndDate:     &subEnd,
		},
	}

	repo.EXPECT().List(ctx, mock.AnythingOfType("infra.ListFilter")).Return(data, nil)

	sum, err := svc.TotalCost(ctx, periodStart, periodEnd, services.ListFilter{})
	require.NoError(t, err)
	require.Equal(t, 0, sum)
}

func TestService_TotalCost_ListMapping(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	filter := services.ListFilter{
		UserID: "u-1", HasUserID: true,
		ServiceName: "Yandex Plus", HasServiceName: true,
	}

	periodStart, periodEnd := ym(2025, time.January), ym(2025, time.March)

	repo.EXPECT().List(ctx, mock.MatchedBy(func(ifl infra.ListFilter) bool {
		return ifl.UserID != nil && *ifl.UserID == "u-1" &&
			ifl.ServiceName != nil && *ifl.ServiceName == "Yandex Plus"
	})).Return([]models.Subscription{}, nil)

	sum, err := svc.TotalCost(ctx, periodStart, periodEnd, filter)
	require.NoError(t, err)
	require.Equal(t, 0, sum)
}

func TestService_TotalCost_ErrValidation_Period(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	// Период начало март, конец январь = ошибка валидации
	periodStart, periodEnd := ym(2025, time.March), ym(2025, time.January)

	_, err := svc.TotalCost(ctx, periodStart, periodEnd, services.ListFilter{})
	require.Error(t, err)
	require.True(t, errors.Is(err, services.ErrValidation))
}

func TestService_TotalCost_ErrDatabase(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewDatabase(t)
	svc := subscription_service.New(repo)

	periodStart, periodEnd := ym(2025, time.January), ym(2025, time.March)
	data := []models.Subscription{
		{
			ID:          1,
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      "u-1",
			StartDate:   periodStart,
			EndDate:     &periodEnd,
		},
	}

	repo.EXPECT().List(ctx, mock.AnythingOfType("infra.ListFilter")).Return(data, errors.New("ошибка БД"))

	_, err := svc.TotalCost(ctx, periodStart, periodEnd, services.ListFilter{})
	require.Error(t, err)
	require.ErrorContains(t, err, "ошибка БД")
}
