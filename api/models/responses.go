package models

type ResponseOK struct {
	Message struct{} `json:"message"`
}

type ResponseSuccess struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ResponseError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

type Empty struct{}
