package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sunr3d/subscription-aggregator/internal/interfaces/services"
)

func validateCreateSubscription(req createSubscriptionReq) error {
	if strings.TrimSpace(req.ServiceName) == "" {
		return fmt.Errorf("service_name обязателен")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return fmt.Errorf("user_id обязателен")
	}

	if _, err := time.Parse("01-2006", req.StartDate); err != nil {
		return fmt.Errorf("start_date должен быть в формате MM-YYYY")
	}

	if strings.TrimSpace(req.EndDate) != "" {
		if _, err := time.Parse("01-2006", req.EndDate); err != nil {
			return fmt.Errorf("end_date должен быть в формате MM-YYYY")
		}
	}
	return nil
}

func validateUpdateSubscription(req updateSubscriptionReq) error {
	if req.ServiceName == nil && req.Price == nil && req.UserID == nil && req.StartDate == nil && req.EndDate == nil {
		return fmt.Errorf("необходимо указать хотя бы одно поле для обновления")
	}

	if req.ServiceName != nil && strings.TrimSpace(*req.ServiceName) == "" {
		return fmt.Errorf("service_name не может быть пустым")
	}

	if req.UserID != nil && strings.TrimSpace(*req.UserID) == "" {
		return fmt.Errorf("user_id не может быть пустым")
	}

	if req.StartDate != nil {
		if _, err := time.Parse("01-2006", *req.StartDate); err != nil {
			return fmt.Errorf("start_date должен быть в формате MM-YYYY")
		}
	}

	if req.EndDate != nil {
		if _, err := time.Parse("01-2006", *req.EndDate); err != nil {
			return fmt.Errorf("end_date должен быть в формате MM-YYYY")
		}
	}

	return nil
}

func validateListSubscription(query url.Values, filter *services.ListFilter) error {
	filter.Limit = 50
	filter.Offset = 0

	if userID := strings.TrimSpace(query.Get("user_id")); userID != "" {
		filter.UserID, filter.HasUserID = userID, true
	}
	if serviceName := strings.TrimSpace(query.Get("service_name")); serviceName != "" {
		filter.ServiceName, filter.HasServiceName = serviceName, true
	}

	if limitStr := strings.TrimSpace(query.Get("limit")); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			return fmt.Errorf("limit должен быть числом от 1 до 100")
		}
		filter.Limit = limit
	}

	if offsetStr := strings.TrimSpace(query.Get("offset")); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return fmt.Errorf("offset должен быть числом >= 0")
		}
		filter.Offset = offset
	}
	return nil
}

func validateTotalCost(query url.Values) error {
	startDate := strings.TrimSpace(query.Get("period_start"))
	endDate := strings.TrimSpace(query.Get("period_end"))

	if startDate == "" || endDate == "" {
		return fmt.Errorf("period_start и period_end не могут быть пустыми")
	}

	if _, err := time.Parse("01-2006", startDate); err != nil {
		return fmt.Errorf("period_start должен быть в формате MM-YYYY")
	}

	if _, err := time.Parse("01-2006", endDate); err != nil {
		return fmt.Errorf("period_end должен быть в формате MM-YYYY")
	}

	return nil
}
