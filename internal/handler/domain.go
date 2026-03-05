package handler

import (
	"fmt"
	"net/http"

	"sds-admin/internal/dto"
	"sds-admin/internal/logger"
	"sds-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// DomainHandler handles domain-related HTTP requests.
// DomainHandler 处理域名相关的HTTP请求
type DomainHandler struct {
	domainService *service.DomainService
}

// NewDomainHandler creates a new DomainHandler instance.
// NewDomainHandler 创建新的DomainHandler实例
func NewDomainHandler(domainService *service.DomainService) *DomainHandler {
	return &DomainHandler{domainService: domainService}
}

// CreateDomain godoc
// @Summary      Create a new domain
// @Description  Create a new domain with the provided information
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateDomainRequest  true  "Domain creation request"
// @Success      201  {object}  dto.Response{data=dto.DomainResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains [post]
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req dto.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	domain, err := h.domainService.CreateDomain(&req)
	if err != nil {
		logger.Errorf("Failed to create domain: %v", err)
		if err == service.ErrDomainExists {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Domain already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create domain",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Domain created successfully: %s", domain.DomainName)
	c.JSON(http.StatusCreated, dto.Response{
		Code:    http.StatusCreated,
		Message: "Domain created successfully",
		Data:    domain,
	})
}

// GetDomain godoc
// @Summary      Get a domain by ID
// @Description  Retrieve a single domain by its ID
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Domain ID"
// @Success      200  {object}  dto.Response{data=dto.DomainResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id} [get]
func (h *DomainHandler) GetDomain(c *gin.Context) {
	id := c.Param("id")
	var domainID uint
	if _, err := fmt.Sscanf(id, "%d", &domainID); err != nil {
		logger.Warnf("Invalid domain ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid domain ID",
		})
		return
	}

	domain, err := h.domainService.GetDomainByID(domainID)
	if err != nil {
		logger.Errorf("Failed to get domain: %v", err)
		if err == service.ErrDomainNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get domain",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    domain,
	})
}

// ListDomains godoc
// @Summary      List all domains
// @Description  Retrieve a list of all domains
// @Tags         domains
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.Response{data=[]dto.DomainResponse}
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains [get]
func (h *DomainHandler) ListDomains(c *gin.Context) {
	domains, err := h.domainService.ListDomains()
	if err != nil {
		logger.Errorf("Failed to list domains: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list domains",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    domains,
	})
}

// UpdateDomain godoc
// @Summary      Update a domain
// @Description  Update an existing domain with new information
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Domain ID"
// @Param        request  body      dto.CreateDomainRequest  true  "Domain update request"
// @Success      200  {object}  dto.Response{data=dto.DomainResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id} [put]
func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	id := c.Param("id")
	var domainID uint
	if _, err := fmt.Sscanf(id, "%d", &domainID); err != nil {
		logger.Warnf("Invalid domain ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid domain ID",
		})
		return
	}

	var req dto.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	domain, err := h.domainService.UpdateDomain(domainID, &req)
	if err != nil {
		logger.Errorf("Failed to update domain: %v", err)
		if err == service.ErrDomainNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update domain",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Domain updated successfully: %s", domain.DomainName)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Domain updated successfully",
		Data:    domain,
	})
}

// DeleteDomain godoc
// @Summary      Delete a domain
// @Description  Delete a domain by its ID (soft delete)
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Domain ID"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id} [delete]
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	id := c.Param("id")
	var domainID uint
	if _, err := fmt.Sscanf(id, "%d", &domainID); err != nil {
		logger.Warnf("Invalid domain ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid domain ID",
		})
		return
	}

	if err := h.domainService.DeleteDomain(domainID); err != nil {
		logger.Errorf("Failed to delete domain: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete domain",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Domain deleted successfully: ID %d", domainID)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Domain deleted successfully",
	})
}

// DisableDomain godoc
// @Summary      Disable a domain
// @Description  Disable a domain by setting its status to inactive
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Domain ID"
// @Success      200  {object}  dto.Response{data=dto.DomainResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/disable [post]
func (h *DomainHandler) DisableDomain(c *gin.Context) {
	id := c.Param("id")
	var domainID uint
	if _, err := fmt.Sscanf(id, "%d", &domainID); err != nil {
		logger.Warnf("Invalid domain ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid domain ID",
		})
		return
	}

	domain, err := h.domainService.DisableDomain(domainID)
	if err != nil {
		logger.Errorf("Failed to disable domain: %v", err)
		if err == service.ErrDomainNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to disable domain",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Domain disabled successfully: %s", domain.DomainName)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Domain disabled successfully",
		Data:    domain,
	})
}

// EnableDomain godoc
// @Summary      Enable a domain
// @Description  Enable a domain by setting its status to active
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Domain ID"
// @Success      200  {object}  dto.Response{data=dto.DomainResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/enable [post]
func (h *DomainHandler) EnableDomain(c *gin.Context) {
	id := c.Param("id")
	var domainID uint
	if _, err := fmt.Sscanf(id, "%d", &domainID); err != nil {
		logger.Warnf("Invalid domain ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid domain ID",
		})
		return
	}

	domain, err := h.domainService.EnableDomain(domainID)
	if err != nil {
		logger.Errorf("Failed to enable domain: %v", err)
		if err == service.ErrDomainNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to enable domain",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Domain enabled successfully: %s", domain.DomainName)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Domain enabled successfully",
		Data:    domain,
	})
}

// UpdateDomainConfig godoc
// @Summary      Update domain config
// @Description  Update domain config (description and recursive)
// @Tags         domains
// @Accept       json
// @Produce      json
// @Param        id       path      int                            true  "Domain ID"
// @Param        request  body      dto.UpdateDomainConfigRequest  true  "Domain config update request"
// @Success      200  {object}  dto.Response{data=dto.DomainResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/config [patch]
func (h *DomainHandler) UpdateDomainConfig(c *gin.Context) {
	id := c.Param("id")
	var domainID uint
	if _, err := fmt.Sscanf(id, "%d", &domainID); err != nil {
		logger.Warnf("Invalid domain ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid domain ID",
		})
		return
	}

	var req dto.UpdateDomainConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	domain, err := h.domainService.UpdateDomainConfig(domainID, &req)
	if err != nil {
		logger.Errorf("Failed to update domain config: %v", err)
		if err == service.ErrDomainNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update domain config",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Domain config updated successfully: %s", domain.DomainName)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Domain config updated successfully",
		Data:    domain,
	})
}
