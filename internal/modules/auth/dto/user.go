package dto

type UpdateUserRequest struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	City     string `json:"city"`
}

type InviteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=editor viewer"`
}
