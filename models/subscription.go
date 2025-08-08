package models

import "time"

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}
