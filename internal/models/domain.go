package models

import (
	"time"

	"gorm.io/gorm"
)

// Domain represents a top-level domain model.
// Domain 表示一级域名模型
type Domain struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DomainName  string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"domain_name"`
	Recursive   bool           `gorm:"default:false" json:"recursive"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Domain model.
// TableName 指定Domain模型的表名
func (Domain) TableName() string {
	return "domains"
}

// BeforeCreate is a hook that runs before creating a new record.
// BeforeCreate 是在创建新记录前运行的钩子
func (d *Domain) BeforeCreate(tx *gorm.DB) error {
	d.Status = "active"
	return nil
}
