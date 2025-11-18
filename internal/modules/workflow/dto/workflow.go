package dto

type CreateWorkflowRequest struct {
	Name       string `json:"name" binding:"required"`
	JSON       string `json:"json" binding:"required"`
	MaxTimeout int    `json:"maxTimeout"`
	RetryCount int    `json:"retryCount"`
	RetryDelay int    `json:"retryDelay"`
}

type UpdateWorkflowRequest struct {
	Name       string `json:"name"`
	JSON       string `json:"json"`
	Active     *bool  `json:"active"`
	MaxTimeout int    `json:"maxTimeout"`
	RetryCount int    `json:"retryCount"`
	RetryDelay int    `json:"retryDelay"`
}

type TestWorkflowRequest struct {
	TestData map[string]interface{} `json:"testData"`
}
