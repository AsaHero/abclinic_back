package entity

import "time"

type PriceList struct {
	GUID       string
	Title      string
	Service    string
	Price      float64
	Created_at time.Time
}
