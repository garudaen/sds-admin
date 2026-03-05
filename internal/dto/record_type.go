package dto

// RecordTypeResponse represents the response for record type operations.
// RecordTypeResponse 表示记录类型操作的响应
type RecordTypeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateRecordTypeRequest represents the request for creating a record type.
// CreateRecordTypeRequest 表示创建记录类型的请求
type CreateRecordTypeRequest struct {
	Name        string `json:"name" binding:"required,max=20"`
	Description string `json:"description" binding:"max=255"`
}

// UpdateRecordTypeRequest represents the request for updating a record type.
// UpdateRecordTypeRequest 表示更新记录类型的请求
type UpdateRecordTypeRequest struct {
	Name        string `json:"name" binding:"max=20"`
	Description string `json:"description" binding:"max=255"`
}
