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
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Status   bool     `json:"status"`
	Roles    []string `json:"roles"`
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

type GetRolesResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"descriptions"`
	Permissions []string `json:"permissions"`
}

type GetRoleRequest struct {
	ID uint `json:"id" binding:"required,min=1"`
}

type GetRoleResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions"`
}

type UpdateRoleRequest struct {
	ID          uint   `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions,omitempty"`
}

type GetPermissionsResponse struct {
	Name        string   `json:"permissions"`
	Description string   `json:"descriptions"`
	Sections    []string `json:"sections"`
}

type GetDelPermissionRequest struct {
	ID uint `json:"id" binding:"required,min=1"`
}

type GetPermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sections    []uint `json:"sections"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sections    []uint `json:"sections"`
}

type UpdatePermissionRequest struct {
	ID          uint   `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sections    []uint `json:"sections,omitempty"`
}

type GetSectionsResponse struct {
	Module   string `json:"module"`
	Endpoint string `json:"endpoint"`
}

type GetDelSectionRequest struct {
	ID uint `json:"id" binding:"required,min=1"`
}

type CreateSectionRequest struct {
	Module   string `json:"module" binding:"required"`
	Endpoint string `json:"endpoint"`
}

type UpdateSectionRequest struct {
	ID       uint   `json:"id" binding:"required,min=1"`
	Module   string `json:"module"`
	Endpoint string `json:"endpoint"`
}
