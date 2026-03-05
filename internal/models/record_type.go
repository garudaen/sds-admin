package models

import (
	"time"

	"gorm.io/gorm"
)

// RecordType represents a DNS record type (e.g., A, AAAA, CNAME).
// RecordType 表示DNS记录类型（如A、AAAA、CNAME等）
type RecordType struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for RecordType model.
// TableName 指定RecordType模型的表名
func (RecordType) TableName() string {
	return "record_types"
}
