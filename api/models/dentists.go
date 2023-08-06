package models

type GetDentistsListResponse struct {
	ID        int64  `json:"id"`
	CloneName string `json:"clone_name"`
	Img       string `json:"img"`
	Priority  int16  `json:"priority"`
	Side      string `json:"side"`
	Name      string `json:"name"`
	Info      string `json:"info"`
}

type UpdateDentistRequest struct {
	Name string `json:"name"`
	Info string `json:"info"`
	Img  string `json:"img"`
}
