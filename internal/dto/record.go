package dto

// RecordValueInput represents a record value in the request.
// RecordValueInput 表示请求中的记录值
type RecordValueInput struct {
	// Record value, e.g. IP address for A record, domain for CNAME.
	// 记录值，如A记录的IP地址，CNAME记录的域名
	// @example 192.168.1.1
	Value string `json:"value" binding:"required" example:"192.168.1.1"`
	// MX record priority. Only required for MX records, value range 0-65535.
	// MX记录优先级。仅MX记录需要，取值范围0-65535
	MXPriority *int `json:"mxPriority" example:"10"`
	// Whether this is the default value. Required for A/AAAA/CNAME records.
	// 是否为默认值。A/AAAA/CNAME记录必须有一个默认值
	IsDefault bool `json:"isDefault" example:"true"`
	// Client IP CIDR for smart DNS resolution. Should not be set for default values.
	// 客户端IP网段，用于智能DNS解析。默认值不应设置此字段
	// @example 192.168.1.0/24
	ClientCIDR string `json:"clientCidr" example:"192.168.1.0/24"`
}

// CreateRecordRequest represents the request body for creating a DNS record.
// CreateRecordRequest 表示创建DNS解析记录的请求体
type CreateRecordRequest struct {
	// Record type ID (1=A, 2=AAAA, 3=CNAME, 4=MX, etc.).
	// 记录类型ID（1=A, 2=AAAA, 3=CNAME, 4=MX等）
	RecordTypeID uint `json:"recordTypeId" binding:"required" example:"1"`
	// Host record, e.g. www, @, *.
	// 主机记录，如 www, @, *
	Host string `json:"host" binding:"required" example:"www"`
	// Time to live in seconds. Default is 300.
	// 生存时间（秒），默认为300
	TTL int `json:"ttl" example:"300"`
	// Remark or description for this record.
	// 记录的备注或说明
	Remark string `json:"remark" example:"Web server"`
	// Record values. At least one value is required.
	// 记录值列表。至少需要一个值
	Values []RecordValueInput `json:"values" binding:"required,min=1"`
}

// UpdateRecordRequest represents the request body for updating a DNS record.
// UpdateRecordRequest 表示更新DNS解析记录的请求体
type UpdateRecordRequest struct {
	// Host record, e.g. www, @, *.
	// 主机记录，如 www, @, *
	Host string `json:"host" binding:"required" example:"www"`
	// Time to live in seconds. Default is 300.
	// 生存时间（秒），默认为300
	TTL int `json:"ttl" example:"300"`
	// Remark or description for this record.
	// 记录的备注或说明
	Remark string `json:"remark" example:"Web server"`
	// Whether the record is disabled.
	// 记录是否禁用
	Disabled *bool `json:"disabled" example:"false"`
	// Record values. At least one value is required.
	// 记录值列表。至少需要一个值
	Values []RecordValueInput `json:"values" binding:"required,min=1"`
}

// RecordValueResponse represents a record value in the response.
// RecordValueResponse 表示响应中的记录值
type RecordValueResponse struct {
	// Record value ID.
	// 记录值ID
	ID uint `json:"id" example:"1"`
	// Record value.
	// 记录值
	Value string `json:"value" example:"192.168.1.1"`
	// MX record priority. Only present for MX records.
	// MX记录优先级。仅MX记录有此字段
	MXPriority *int `json:"mxPriority,omitempty" example:"10"`
	// Whether this is the default value.
	// 是否为默认值
	IsDefault bool `json:"isDefault" example:"true"`
	// Client IP CIDR for smart DNS resolution.
	// 客户端IP网段，用于智能DNS解析
	ClientCIDR string `json:"clientCidr" example:"192.168.1.0/24"`
}

// RecordResponse represents the response body for record operations.
// RecordResponse 表示记录操作的响应体
type RecordResponse struct {
	// Record ID.
	// 记录ID
	ID uint `json:"id" example:"1"`
	// Domain ID that this record belongs to.
	// 此记录所属的域名ID
	DomainID uint `json:"domainId" example:"1"`
	// Domain name.
	// 域名
	DomainName string `json:"domainName" example:"example.com"`
	// Record type ID.
	// 记录类型ID
	RecordTypeID uint `json:"recordTypeId" example:"1"`
	// Record type name (A, AAAA, CNAME, MX, etc.).
	// 记录类型名称（A, AAAA, CNAME, MX等）
	RecordType string `json:"recordType" example:"A"`
	// Host record.
	// 主机记录
	Host string `json:"host" example:"www"`
	// Time to live in seconds.
	// 生存时间（秒）
	TTL int `json:"ttl" example:"300"`
	// Remark or description.
	// 备注或说明
	Remark string `json:"remark" example:"Web server"`
	// Whether the record is disabled.
	// 记录是否禁用
	Disabled bool `json:"disabled" example:"false"`
	// Record values.
	// 记录值列表
	Values []RecordValueResponse `json:"values"`
	// Creation time.
	// 创建时间
	CreatedAt string `json:"createdAt" example:"2024-01-01 12:00:00"`
	// Last update time.
	// 最后更新时间
	UpdatedAt string `json:"updatedAt" example:"2024-01-01 12:00:00"`
}

// RecordListResponse represents the response for listing records.
// RecordListResponse 表示记录列表的响应
type RecordListResponse struct {
	// Total number of records.
	// 记录总数
	Total int `json:"total" example:"10"`
	// Record list.
	// 记录列表
	Records []RecordResponse `json:"records"`
}
