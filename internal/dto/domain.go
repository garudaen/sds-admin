package dto

// CreateDomainRequest represents the request body for creating a domain.
// CreateDomainRequest 表示创建域名的请求体
type CreateDomainRequest struct {
	DomainName  string `json:"domainName" binding:"required"`
	Recursive   bool   `json:"recursive"`
	Description string `json:"description"`
}

// UpdateDomainConfigRequest represents the request body for updating domain config.
// UpdateDomainConfigRequest 表示更新域名配置的请求体
type UpdateDomainConfigRequest struct {
	Description *string `json:"description"`
	Recursive   *bool   `json:"recursive"`
}

// DomainResponse represents the response body for domain operations.
// DomainResponse 表示域名操作的响应体
type DomainResponse struct {
	ID          uint   `json:"id"`
	DomainName  string `json:"domainName"`
	Recursive   bool   `json:"recursive"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Response represents a standard API response structure.
// Response 表示标准的API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response structure.
// ErrorResponse 表示错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}
