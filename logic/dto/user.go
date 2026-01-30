package dto

type CreateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateProfileReq struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type UpdateStatusReq struct {
	Status int8 `json:"status" binding:"required"`
}
