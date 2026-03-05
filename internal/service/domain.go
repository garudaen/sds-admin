package service

import (
	"errors"
	"fmt"

	"sds-admin/internal/dto"
	"sds-admin/internal/models"

	"gorm.io/gorm"
)

var (
	ErrDomainNotFound     = errors.New("domain not found")
	ErrDomainExists       = errors.New("domain already exists")
	ErrRecordTypeNotFound = errors.New("record type not found")
	ErrRecordTypeExists   = errors.New("record type already exists")
)

// DomainService handles domain-related business logic.
// DomainService 处理域名相关的业务逻辑
type DomainService struct {
	db *gorm.DB
}

// NewDomainService creates a new DomainService instance.
// NewDomainService 创建新的DomainService实例
func NewDomainService(db *gorm.DB) *DomainService {
	return &DomainService{db: db}
}

// CreateDomain creates a new domain.
// CreateDomain 创建新域名
func (s *DomainService) CreateDomain(req *dto.CreateDomainRequest) (*dto.DomainResponse, error) {
	var existingDomain models.Domain
	result := s.db.Where("domain_name = ?", req.DomainName).First(&existingDomain)
	if result.Error == nil {
		return nil, ErrDomainExists
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check domain existence: %w", result.Error)
	}

	domain := &models.Domain{
		DomainName:  req.DomainName,
		Recursive:   req.Recursive,
		Description: req.Description,
	}

	if err := s.db.Create(domain).Error; err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return s.toDomainResponse(domain), nil
}

// GetDomainByID retrieves a domain by ID.
// GetDomainByID 根据ID获取域名
func (s *DomainService) GetDomainByID(id uint) (*dto.DomainResponse, error) {
	var domain models.Domain
	if err := s.db.First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return s.toDomainResponse(&domain), nil
}

// ListDomains retrieves all domains.
// ListDomains 获取所有域名
func (s *DomainService) ListDomains() ([]dto.DomainResponse, error) {
	var domains []models.Domain
	if err := s.db.Find(&domains).Error; err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	responses := make([]dto.DomainResponse, len(domains))
	for i, domain := range domains {
		responses[i] = *s.toDomainResponse(&domain)
	}

	return responses, nil
}

// UpdateDomain updates an existing domain.
// UpdateDomain 更新现有域名
func (s *DomainService) UpdateDomain(id uint, req *dto.CreateDomainRequest) (*dto.DomainResponse, error) {
	var domain models.Domain
	if err := s.db.First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	domain.DomainName = req.DomainName
	domain.Recursive = req.Recursive
	domain.Description = req.Description

	if err := s.db.Save(&domain).Error; err != nil {
		return nil, fmt.Errorf("failed to update domain: %w", err)
	}

	return s.toDomainResponse(&domain), nil
}

// DeleteDomain deletes a domain by ID.
// DeleteDomain 根据ID删除域名
func (s *DomainService) DeleteDomain(id uint) error {
	if err := s.db.Delete(&models.Domain{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}

// DisableDomain disables a domain by ID.
// DisableDomain 根据ID禁用域名
func (s *DomainService) DisableDomain(id uint) (*dto.DomainResponse, error) {
	var domain models.Domain
	if err := s.db.First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	domain.Status = "inactive"
	if err := s.db.Save(&domain).Error; err != nil {
		return nil, fmt.Errorf("failed to disable domain: %w", err)
	}

	return s.toDomainResponse(&domain), nil
}

// EnableDomain enables a domain by ID.
// EnableDomain 根据ID启用域名
func (s *DomainService) EnableDomain(id uint) (*dto.DomainResponse, error) {
	var domain models.Domain
	if err := s.db.First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	domain.Status = "active"
	if err := s.db.Save(&domain).Error; err != nil {
		return nil, fmt.Errorf("failed to enable domain: %w", err)
	}

	return s.toDomainResponse(&domain), nil
}

// UpdateDomainConfig updates domain config (description and recursive).
// UpdateDomainConfig 更新域名配置（描述和是否递归）
func (s *DomainService) UpdateDomainConfig(id uint, req *dto.UpdateDomainConfigRequest) (*dto.DomainResponse, error) {
	var domain models.Domain
	if err := s.db.First(&domain, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDomainNotFound
		}
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	// 只更新提供的字段
	if req.Description != nil {
		domain.Description = *req.Description
	}
	if req.Recursive != nil {
		domain.Recursive = *req.Recursive
	}

	if err := s.db.Save(&domain).Error; err != nil {
		return nil, fmt.Errorf("failed to update domain config: %w", err)
	}

	return s.toDomainResponse(&domain), nil
}

// toDomainResponse converts a domain model to response DTO.
// toDomainResponse 将域名模型转换为响应DTO
func (s *DomainService) toDomainResponse(domain *models.Domain) *dto.DomainResponse {
	return &dto.DomainResponse{
		ID:          domain.ID,
		DomainName:  domain.DomainName,
		Recursive:   domain.Recursive,
		Description: domain.Description,
		Status:      domain.Status,
		CreatedAt:   domain.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   domain.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
