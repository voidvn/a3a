package dto

type CreateConnectionRequest struct {
	Name        string                 `json:"name" binding:"required"`
	ServiceName string                 `json:"serviceName" binding:"required"`
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
	IsActive    *bool                  `json:"isActive"`
}

type UpdateConnectionRequest struct {
	Name        string                 `json:"name"`
	Credentials map[string]interface{} `json:"credentials"`
	IsActive    *bool                  `json:"isActive"`
}

type TestConnectionRequest struct {
	// Add any test-specific fields here
}
