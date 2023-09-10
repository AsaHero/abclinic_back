package entity

import "time"

type Services struct {
	GUID      string
	GroupID   string
	Name      string
	Price     []float64
	CreatedAt time.Time
	UpdateAt  time.Time
}

type ServiceGroups struct {
	GUID      string
	Name      string
	CreatedAt time.Time
}
