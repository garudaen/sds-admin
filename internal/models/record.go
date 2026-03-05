package models

import (
	"time"

	"gorm.io/gorm"
)

// Record represents a DNS resolution record.
// Record 表示DNS解析记录
type Record struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DomainID     uint           `gorm:"not null;index" json:"domain_id"`
	Domain       *Domain        `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	RecordTypeID uint           `gorm:"not null;index" json:"record_type_id"`
	RecordType   *RecordType    `gorm:"foreignKey:RecordTypeID" json:"record_type,omitempty"`
	Host         string         `gorm:"type:varchar(255);not null;index;comment:'Host record, e.g. www, @, *'" json:"host"`
	TTL          int            `gorm:"default:300;comment:'Time to live in seconds'" json:"ttl"`
	Remark       string         `gorm:"type:text" json:"remark"`
	Disabled     bool           `gorm:"default:false;comment:'Whether the record is disabled'" json:"disabled"`
	Values       []RecordValue  `gorm:"foreignKey:RecordID" json:"values,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Record model.
// TableName 指定Record模型的表名
func (Record) TableName() string {
	return "records"
}

// RecordValue represents a value of a DNS record with optional client CIDR matching.
// RecordValue 表示DNS记录的值，支持可选的客户端CIDR匹配
type RecordValue struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RecordID   uint           `gorm:"not null;index" json:"record_id"`
	Record     *Record        `gorm:"foreignKey:RecordID" json:"record,omitempty"`
	Value      string         `gorm:"type:varchar(255);not null;comment:'Record value'" json:"value"`
	MXPriority *int           `gorm:"comment:'MX record priority, only for MX type'" json:"mx_priority,omitempty"`
	IsDefault  bool           `gorm:"default:false;comment:'Whether this is the default value for A/AAAA/CNAME types'" json:"is_default"`
	ClientCIDR string         `gorm:"type:varchar(50);comment:'Client IP CIDR for smart resolution, e.g. 192.168.1.0/24'" json:"client_cidr"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for RecordValue model.
// TableName 指定RecordValue模型的表名
func (RecordValue) TableName() string {
	return "record_values"
}

// RecordWithValues is a struct for returning record with its values in API responses.
// RecordWithValues 用于API响应中返回记录及其值
type RecordWithValues struct {
	ID           uint        `json:"id"`
	DomainID     uint        `json:"domain_id"`
	DomainName   string      `json:"domain_name"`
	RecordTypeID uint        `json:"record_type_id"`
	RecordType   string      `json:"record_type"`
	Host         string      `json:"host"`
	TTL          int         `json:"ttl"`
	Remark       string      `json:"remark"`
	Disabled     bool        `json:"disabled"`
	Values       []ValueInfo `json:"values"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// ValueInfo represents a record value with its metadata.
// ValueInfo 表示记录值及其元数据
type ValueInfo struct {
	ID         uint   `json:"id"`
	Value      string `json:"value"`
	MXPriority *int   `json:"mx_priority,omitempty"`
	IsDefault  bool   `json:"is_default"`
	ClientCIDR string `json:"client_cidr,omitempty"`
}
