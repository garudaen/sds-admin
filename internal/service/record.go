package service

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"sds-admin/internal/dto"
	"sds-admin/internal/models"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrRecordExists          = errors.New("record already exists")
	ErrInvalidIPv4Address    = errors.New("invalid IPv4 address for A record")
	ErrInvalidIPv6Address    = errors.New("invalid IPv6 address for AAAA record")
	ErrCNAMEConflict         = errors.New("CNAME record cannot coexist with other record types for the same host")
	ErrMXPriorityRequired    = errors.New("MX priority is required for MX record")
	ErrDefaultValueRequired  = errors.New("A/AAAA/CNAME record must have at least one default value")
	ErrInvalidCIDR           = errors.New("invalid CIDR format")
	ErrDuplicateDefaultValue = errors.New("only one default value is allowed for the same record")
	ErrDefaultValueWithCIDR  = errors.New("default value should not have client CIDR")
)

// RecordService handles record-related business logic.
// RecordService 处理解析记录相关的业务逻辑
type RecordService struct {
	db *gorm.DB
}

// NewRecordService creates a new RecordService instance.
// NewRecordService 创建新的RecordService实例
func NewRecordService(db *gorm.DB) *RecordService {
	return &RecordService{db: db}
}

// CreateRecord creates a new DNS record.
// CreateRecord 创建新的DNS解析记录
func (s *RecordService) CreateRecord(domainID uint, req *dto.CreateRecordRequest) (*dto.RecordResponse, error) {
	// Validate domain exists.
	// 验证域名是否存在
	var domain models.Domain
	if err := s.db.First(&domain, domainID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// Validate record type exists.
	// 验证记录类型是否存在
	var recordType models.RecordType
	if err := s.db.First(&recordType, req.RecordTypeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordTypeNotFound
		}
		return nil, fmt.Errorf("failed to get record type: %w", err)
	}

	// Validate values based on record type.
	// 根据记录类型验证值
	validatedValues, err := s.validateValues(recordType.Name, req.Values)
	if err != nil {
		return nil, err
	}

	// Check CNAME conflict.
	// 检查CNAME冲突
	if err := s.checkCNAMEConflict(domainID, req.Host, recordType.Name, 0); err != nil {
		return nil, err
	}

	// Check if record already exists.
	// 检查记录是否已存在
	var existingRecord models.Record
	err = s.db.Where("domain_id = ? AND host = ? AND record_type_id = ?",
		domainID, req.Host, req.RecordTypeID).First(&existingRecord).Error
	if err == nil {
		return nil, ErrRecordExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check record existence: %w", err)
	}

	// Create record.
	// 创建记录
	ttl := req.TTL
	if ttl <= 0 {
		ttl = 300
	}

	record := &models.Record{
		DomainID:     domainID,
		RecordTypeID: req.RecordTypeID,
		Host:         req.Host,
		TTL:          ttl,
		Remark:       req.Remark,
	}

	// Start transaction.
	// 开始事务
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return fmt.Errorf("failed to create record: %w", err)
		}

		// Create record values.
		// 创建记录值
		for _, v := range validatedValues {
			recordValue := models.RecordValue{
				RecordID:   record.ID,
				Value:      v.Value,
				MXPriority: v.MXPriority,
				IsDefault:  v.IsDefault,
				ClientCIDR: v.ClientCIDR,
			}
			if err := tx.Create(&recordValue).Error; err != nil {
				return fmt.Errorf("failed to create record value: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetRecordByID(record.ID)
}

// GetRecordByID retrieves a record by ID.
// GetRecordByID 根据ID获取解析记录
func (s *RecordService) GetRecordByID(id uint) (*dto.RecordResponse, error) {
	var record models.Record
	if err := s.db.Preload("Domain").Preload("RecordType").Preload("Values").First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	return s.toRecordResponse(&record), nil
}

// ListRecords retrieves all records for a domain by domain ID.
// ListRecords 根据域名ID获取所有解析记录
func (s *RecordService) ListRecords(domainID uint) (*dto.RecordListResponse, error) {
	// Validate domain exists.
	// 验证域名是否存在
	var domain models.Domain
	if err := s.db.First(&domain, domainID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	var records []models.Record
	if err := s.db.Preload("Domain").Preload("RecordType").Preload("Values").
		Where("domain_id = ?", domainID).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to list records: %w", err)
	}

	responses := make([]dto.RecordResponse, len(records))
	for i, record := range records {
		responses[i] = *s.toRecordResponse(&record)
	}

	return &dto.RecordListResponse{
		Total:   len(responses),
		Records: responses,
	}, nil
}

// UpdateRecord updates an existing DNS record.
// UpdateRecord 更新现有的DNS解析记录
func (s *RecordService) UpdateRecord(id uint, req *dto.UpdateRecordRequest) (*dto.RecordResponse, error) {
	var record models.Record
	if err := s.db.Preload("RecordType").First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	// Validate values based on record type.
	// 根据记录类型验证值
	validatedValues, err := s.validateValues(record.RecordType.Name, req.Values)
	if err != nil {
		return nil, err
	}

	// Check CNAME conflict if host changed.
	// 如果主机记录改变，检查CNAME冲突
	if record.Host != req.Host {
		if err := s.checkCNAMEConflict(record.DomainID, req.Host, record.RecordType.Name, id); err != nil {
			return nil, err
		}
	}

	// Update record.
	// 更新记录
	ttl := req.TTL
	if ttl <= 0 {
		ttl = 300
	}

	record.Host = req.Host
	record.TTL = ttl
	record.Remark = req.Remark
	if req.Disabled != nil {
		record.Disabled = *req.Disabled
	}

	// Start transaction.
	// 开始事务
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&record).Error; err != nil {
			return fmt.Errorf("failed to update record: %w", err)
		}

		// Delete old values.
		// 删除旧值
		if err := tx.Where("record_id = ?", record.ID).Delete(&models.RecordValue{}).Error; err != nil {
			return fmt.Errorf("failed to delete old values: %w", err)
		}

		// Create new values.
		// 创建新值
		for _, v := range validatedValues {
			recordValue := models.RecordValue{
				RecordID:   record.ID,
				Value:      v.Value,
				IsDefault:  v.IsDefault,
				ClientCIDR: v.ClientCIDR,
			}
			if err := tx.Create(&recordValue).Error; err != nil {
				return fmt.Errorf("failed to create record value: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetRecordByID(record.ID)
}

// DeleteRecord deletes a record by ID.
// DeleteRecord 根据ID删除解析记录
func (s *RecordService) DeleteRecord(id uint) error {
	// Start transaction.
	// 开始事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete record values.
		// 删除记录值
		if err := tx.Where("record_id = ?", id).Delete(&models.RecordValue{}).Error; err != nil {
			return fmt.Errorf("failed to delete record values: %w", err)
		}

		// Delete record.
		// 删除记录
		if err := tx.Delete(&models.Record{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete record: %w", err)
		}

		return nil
	})
}

// DisableRecord disables a record by ID.
// DisableRecord 根据ID禁用解析记录
func (s *RecordService) DisableRecord(id uint) (*dto.RecordResponse, error) {
	var record models.Record
	if err := s.db.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	record.Disabled = true
	if err := s.db.Save(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to disable record: %w", err)
	}

	return s.GetRecordByID(record.ID)
}

// EnableRecord enables a record by ID.
// EnableRecord 根据ID启用解析记录
func (s *RecordService) EnableRecord(id uint) (*dto.RecordResponse, error) {
	var record models.Record
	if err := s.db.First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	record.Disabled = false
	if err := s.db.Save(&record).Error; err != nil {
		return nil, fmt.Errorf("failed to enable record: %w", err)
	}

	return s.GetRecordByID(record.ID)
}

// validateValues validates and processes record values based on record type.
// validateValues 根据记录类型验证和处理记录值
func (s *RecordService) validateValues(recordTypeName string, values []dto.RecordValueInput) ([]dto.RecordValueInput, error) {
	if len(values) == 0 {
		return nil, errors.New("at least one value is required")
	}

	result := make([]dto.RecordValueInput, 0, len(values))
	hasDefault := false
	hasMXPriority := false

	for _, v := range values {
		validatedValue := v

		// Validate based on record type.
		// 根据记录类型验证
		switch recordTypeName {
		case "A":
			// Validate IPv4 address.
			// 验证IPv4地址
			if net.ParseIP(v.Value) == nil || strings.Contains(v.Value, ":") {
				return nil, ErrInvalidIPv4Address
			}
		case "AAAA":
			// Validate IPv6 address.
			// 验证IPv6地址
			if net.ParseIP(v.Value) == nil || !strings.Contains(v.Value, ":") {
				return nil, ErrInvalidIPv6Address
			}
		case "CNAME":
			// Ensure CNAME ends with dot.
			// 确保CNAME以点结尾
			if !strings.HasSuffix(v.Value, ".") {
				validatedValue.Value = v.Value + "."
			}
		case "MX":
			// MX record requires priority.
			// MX记录需要优先级
			if v.MXPriority == nil {
				return nil, ErrMXPriorityRequired
			}
			hasMXPriority = true
		}

		// Validate CIDR if provided.
		// 如果提供了CIDR则验证
		if v.ClientCIDR != "" {
			// Default value should not have CIDR.
			// 默认值不应该有CIDR
			if v.IsDefault {
				return nil, ErrDefaultValueWithCIDR
			}
			// Validate CIDR format.
			// 验证CIDR格式
			if _, _, err := net.ParseCIDR(v.ClientCIDR); err != nil {
				return nil, ErrInvalidCIDR
			}
		}

		// Check for A/AAAA/CNAME: must have at least one default value.
		// 检查A/AAAA/CNAME：必须有至少一个默认值
		if recordTypeName == "A" || recordTypeName == "AAAA" || recordTypeName == "CNAME" {
			if v.IsDefault {
				if hasDefault {
					return nil, ErrDuplicateDefaultValue
				}
				hasDefault = true
			}
		}

		result = append(result, validatedValue)
	}

	// A/AAAA/CNAME must have at least one default value.
	// A/AAAA/CNAME必须有至少一个默认值
	if (recordTypeName == "A" || recordTypeName == "AAAA" || recordTypeName == "CNAME") && !hasDefault {
		return nil, ErrDefaultValueRequired
	}

	// MX record must have priority for each value.
	// MX记录每个值都必须有优先级
	if recordTypeName == "MX" && !hasMXPriority {
		return nil, ErrMXPriorityRequired
	}

	return result, nil
}

// checkCNAMEConflict checks if a CNAME record conflicts with other records.
// checkCNAMEConflict 检查CNAME记录是否与其他记录冲突
func (s *RecordService) checkCNAMEConflict(domainID uint, host, recordTypeName string, excludeID uint) error {
	var existingRecords []models.Record
	query := s.db.Where("domain_id = ? AND host = ?", domainID, host)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Find(&existingRecords).Error; err != nil {
		return fmt.Errorf("failed to check CNAME conflict: %w", err)
	}

	// If creating CNAME, check if other records exist.
	// 如果创建CNAME，检查是否存在其他记录
	if recordTypeName == "CNAME" && len(existingRecords) > 0 {
		return ErrCNAMEConflict
	}

	// If creating other type, check if CNAME exists.
	// 如果创建其他类型，检查是否存在CNAME
	if recordTypeName != "CNAME" {
		for _, r := range existingRecords {
			var rt models.RecordType
			if err := s.db.First(&rt, r.RecordTypeID).Error; err != nil {
				continue
			}
			if rt.Name == "CNAME" {
				return ErrCNAMEConflict
			}
		}
	}

	return nil
}

// toRecordResponse converts a record model to response DTO.
// toRecordResponse 将记录模型转换为响应DTO
func (s *RecordService) toRecordResponse(record *models.Record) *dto.RecordResponse {
	domainName := ""
	if record.Domain != nil {
		domainName = record.Domain.DomainName
	}

	recordTypeName := ""
	if record.RecordType != nil {
		recordTypeName = record.RecordType.Name
	}

	values := make([]dto.RecordValueResponse, len(record.Values))
	for i, v := range record.Values {
		values[i] = dto.RecordValueResponse{
			ID:         v.ID,
			Value:      v.Value,
			MXPriority: v.MXPriority,
			IsDefault:  v.IsDefault,
			ClientCIDR: v.ClientCIDR,
		}
	}

	return &dto.RecordResponse{
		ID:           record.ID,
		DomainID:     record.DomainID,
		DomainName:   domainName,
		RecordTypeID: record.RecordTypeID,
		RecordType:   recordTypeName,
		Host:         record.Host,
		TTL:          record.TTL,
		Remark:       record.Remark,
		Disabled:     record.Disabled,
		Values:       values,
		CreatedAt:    record.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    record.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
