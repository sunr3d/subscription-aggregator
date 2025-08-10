package api

import (
	"fmt"
	"strings"
	"time"
)

func validateCreateSubscription(req createSubscriptionReq) error {
	if strings.TrimSpace(req.ServiceName) == "" {
		return fmt.Errorf("service_name обязателен")
	}
	if req.Price < 0 {
		return fmt.Errorf("price не может быть отрицательным")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return fmt.Errorf("user_id обязателен")
	}

	start, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return fmt.Errorf("start_date должен быть в формате MM-YYYY")
	}
	
	if strings.TrimSpace(req.EndDate) != "" {
		end, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return fmt.Errorf("end_date должен быть в формате MM-YYYY")
		}
		if end.Before(start) {
			return fmt.Errorf("end_date не может быть раньше start_date")
		}
	}
	return nil
}
