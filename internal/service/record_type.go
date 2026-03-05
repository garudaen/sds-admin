package service

import (
	"sds-admin/internal/dto"
	"sds-admin/internal/models"

	"gorm.io/gorm"
)

// RecordTypeService handles record type business logic.
// RecordTypeService 处理记录类型的业务逻辑
type RecordTypeService struct {
	db *gorm.DB
}

// NewRecordTypeService creates a new RecordTypeService instance.
// NewRecordTypeService 创建新的RecordTypeService实例
func NewRecordTypeService(db *gorm.DB) *RecordTypeService {
	return &RecordTypeService{db: db}
}

// ListRecordTypes returns a list of all record types.
// ListRecordTypes 返回所有记录类型的列表
func (s *RecordTypeService) ListRecordTypes() ([]dto.RecordTypeResponse, error) {
	var recordTypes []models.RecordType
	if err := s.db.Find(&recordTypes).Error; err != nil {
		return nil, err
	}

	response := make([]dto.RecordTypeResponse, len(recordTypes))
	for i, recordType := range recordTypes {
		response[i] = dto.RecordTypeResponse{
			ID:          recordType.ID,
			Name:        recordType.Name,
			Description: recordType.Description,
		}
	}

	return response, nil
}

// GetRecordTypeByID returns a record type by ID.
// GetRecordTypeByID 根据ID返回记录类型
func (s *RecordTypeService) GetRecordTypeByID(id uint) (*dto.RecordTypeResponse, error) {
	var recordType models.RecordType
	if err := s.db.First(&recordType, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordTypeNotFound
		}
		return nil, err
	}

	return &dto.RecordTypeResponse{
		ID:          recordType.ID,
		Name:        recordType.Name,
		Description: recordType.Description,
	}, nil
}

// CreateRecordType creates a new record type.
// CreateRecordType 创建新的记录类型
func (s *RecordTypeService) CreateRecordType(req *dto.CreateRecordTypeRequest) (*dto.RecordTypeResponse, error) {
	// Check if record type already exists
	var existingRecordType models.RecordType
	if err := s.db.Where("name = ?", req.Name).First(&existingRecordType).Error; err == nil {
		return nil, ErrRecordTypeExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	recordType := models.RecordType{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&recordType).Error; err != nil {
		return nil, err
	}

	return &dto.RecordTypeResponse{
		ID:          recordType.ID,
		Name:        recordType.Name,
		Description: recordType.Description,
	}, nil
}

// UpdateRecordType updates an existing record type.
// UpdateRecordType 更新现有的记录类型
func (s *RecordTypeService) UpdateRecordType(id uint, req *dto.UpdateRecordTypeRequest) (*dto.RecordTypeResponse, error) {
	var recordType models.RecordType
	if err := s.db.First(&recordType, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordTypeNotFound
		}
		return nil, err
	}

	// Check if name is being changed and if it already exists
	if req.Name != "" && req.Name != recordType.Name {
		var existingRecordType models.RecordType
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingRecordType).Error; err == nil {
			return nil, ErrRecordTypeExists
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		recordType.Name = req.Name
	}

	if req.Description != "" {
		recordType.Description = req.Description
	}

	if err := s.db.Save(&recordType).Error; err != nil {
		return nil, err
	}

	return &dto.RecordTypeResponse{
		ID:          recordType.ID,
		Name:        recordType.Name,
		Description: recordType.Description,
	}, nil
}

// DeleteRecordType deletes a record type by ID (soft delete).
// DeleteRecordType 根据ID删除记录类型（软删除）
func (s *RecordTypeService) DeleteRecordType(id uint) error {
	var recordType models.RecordType
	if err := s.db.First(&recordType, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrRecordTypeNotFound
		}
		return err
	}

	return s.db.Delete(&recordType).Error
}
