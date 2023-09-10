package models

type GetRolesResponse struct {
	Roles []string `json:"roles"`
}

type GetUserInfoResponse struct {
	GUID      string `json:"guid"`
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `jsom:"username"`
}

type CreateUserRequest struct {
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `jsom:"username"`
	Password  string `json:"password"`
}

type UpdateUserRequest struct {
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `jsom:"username"`
	Password  string `json:"password"`
}

type GetAllUsersResponse struct {
	GUID      string `json:"guid"`
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `jsom:"username"`
}

type GetUserResponse struct {
	GUID      string `json:"guid"`
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `jsom:"username"`
}
