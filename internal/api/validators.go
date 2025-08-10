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
