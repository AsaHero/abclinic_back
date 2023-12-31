package models

type Services struct {
	GUID        string  `json:"guid"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	UrgentPrice float64 `json:"urgent_price"`
}

type ServicesGroup struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type GUIDResponse struct {
	GUID string `json:"guid"`
}

type CreateServiceRequest struct {
	GroupID     string  `json:"group_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	UrgentPrice float64 `json:"urgent_price"`
}

type UpdateServiceRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	UrgentPrice float64 `json:"urgent_price"`
}

type CreateServiceGroupRequest struct {
	Name string `json:"name"`
}
