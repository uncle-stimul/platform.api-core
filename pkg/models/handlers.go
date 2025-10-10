package models

type DefaultResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type ModuleInfoResponse struct {
	Status  string `json:"status"`
	Module  string `json:"module"`
	Version string `json:"version"`
	Message string `json:"msg"`
}

type UsersResponse struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Status   bool     `json:"status"`
}

type GetUserRequest struct {
	ID uint `json:"id" binding:"required,min=1"`
}

type GetUserResponse struct {
	ID       uint   `json:"ID"`
	Username string `json:"Username"`
	Status   bool   `json:"Status"`
	Roles    []uint `json:"Roles"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Status   bool   `json:"status"`
	Roles    []uint `json:"roles"`
}

type UpdateUserRequest struct {
	ID       uint   `json:"id" binding:"required,min=1"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Status   bool   `json:"status"`
	Roles    []uint `json:"roles,omitempty"`
}

type DeleteUserRequest struct {
	ID uint `json:"id" binding:"required,min=1"`
}
