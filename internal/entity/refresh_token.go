package entity

import "time"

type RefreshToken struct {
	GUID         string
	RefreshToken string
	ExpiryDate   time.Time
	CreatedAt    time.Time
}
