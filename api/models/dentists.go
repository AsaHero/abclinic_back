package models

type GetDentistsListResponse struct {
	ID        int64  `json:"id"`
	CloneName string `json:"clone_name"`
	Name      string `json:"name"`
	Info      string `json:"info"`
	Img       string `json:"img"`
	Side      string `json:"side"`
	Priority  int16  `json:"priority"`
	Language  string `json:"language"`
}

type UpdateDentistRequest struct {
	Name string `json:"name"`
	Info string `json:"info"`
	Img  string `json:"img"`
}
