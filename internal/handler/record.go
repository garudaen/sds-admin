package handler

import (
	"fmt"
	"net/http"

	"sds-admin/internal/dto"
	"sds-admin/internal/logger"
	"sds-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// RecordHandler handles DNS record-related HTTP requests.
// RecordHandler 处理DNS解析记录相关的HTTP请求
type RecordHandler struct {
	recordService *service.RecordService
}

// NewRecordHandler creates a new RecordHandler instance.
// NewRecordHandler 创建新的RecordHandler实例
func NewRecordHandler(recordService *service.RecordService) *RecordHandler {
	return &RecordHandler{recordService: recordService}
}

// CreateRecord godoc
// @Summary      Create a new DNS record
// @Description  Create a new DNS resolution record under the specified domain
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id       path      int                      true  "Domain ID"
// @Param        request  body      dto.CreateRecordRequest  true  "Record creation request"
// @Success      201  {object}  dto.Response{data=dto.RecordResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records [post]
func (h *RecordHandler) CreateRecord(c *gin.Context) {
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

	var req dto.CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	record, err := h.recordService.CreateRecord(domainID, &req)
	if err != nil {
		logger.Errorf("Failed to create record: %v", err)
		switch err {
		case service.ErrDomainNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
		case service.ErrRecordTypeNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record type not found",
			})
		case service.ErrRecordExists:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Record already exists",
			})
		case service.ErrInvalidIPv4Address,
			service.ErrInvalidIPv6Address,
			service.ErrMXPriorityRequired,
			service.ErrDefaultValueRequired,
			service.ErrInvalidCIDR,
			service.ErrDuplicateDefaultValue,
			service.ErrDefaultValueWithCIDR,
			service.ErrCNAMEConflict:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create record",
				Detail:  err.Error(),
			})
		}
		return
	}

	logger.Infof("Record created successfully: domain_id=%d, host=%s", domainID, record.Host)
	c.JSON(http.StatusCreated, dto.Response{
		Code:    http.StatusCreated,
		Message: "Record created successfully",
		Data:    record,
	})
}

// GetRecord godoc
// @Summary      Get a record by ID
// @Description  Retrieve a single DNS record by its ID
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "Domain ID"
// @Param        recordId  path      int  true  "Record ID"
// @Success      200  {object}  dto.Response{data=dto.RecordResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records/{recordId} [get]
func (h *RecordHandler) GetRecord(c *gin.Context) {
	recordId := c.Param("recordId")
	var id uint
	if _, err := fmt.Sscanf(recordId, "%d", &id); err != nil {
		logger.Warnf("Invalid record ID: %s", recordId)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record ID",
		})
		return
	}

	record, err := h.recordService.GetRecordByID(id)
	if err != nil {
		logger.Errorf("Failed to get record: %v", err)
		if err == service.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get record",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    record,
	})
}

// ListRecords godoc
// @Summary      List DNS records for a domain
// @Description  Retrieve a list of DNS records for the specified domain
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Domain ID"
// @Success      200  {object}  dto.Response{data=dto.RecordListResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records [get]
func (h *RecordHandler) ListRecords(c *gin.Context) {
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

	records, err := h.recordService.ListRecords(domainID)
	if err != nil {
		logger.Errorf("Failed to list records: %v", err)
		if err == service.ErrDomainNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Domain not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list records",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    records,
	})
}

// UpdateRecord godoc
// @Summary      Update a DNS record
// @Description  Update an existing DNS record with new information
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id        path      int                     true  "Domain ID"
// @Param        recordId  path      int                     true  "Record ID"
// @Param        request   body      dto.UpdateRecordRequest true  "Record update request"
// @Success      200  {object}  dto.Response{data=dto.RecordResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records/{recordId} [put]
func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	recordId := c.Param("recordId")
	var id uint
	if _, err := fmt.Sscanf(recordId, "%d", &id); err != nil {
		logger.Warnf("Invalid record ID: %s", recordId)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record ID",
		})
		return
	}

	var req dto.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	record, err := h.recordService.UpdateRecord(id, &req)
	if err != nil {
		logger.Errorf("Failed to update record: %v", err)
		switch err {
		case service.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record not found",
			})
		case service.ErrInvalidIPv4Address,
			service.ErrInvalidIPv6Address,
			service.ErrMXPriorityRequired,
			service.ErrDefaultValueRequired,
			service.ErrInvalidCIDR,
			service.ErrDuplicateDefaultValue,
			service.ErrDefaultValueWithCIDR,
			service.ErrCNAMEConflict:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update record",
				Detail:  err.Error(),
			})
		}
		return
	}

	logger.Infof("Record updated successfully: id=%d, host=%s", record.ID, record.Host)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Record updated successfully",
		Data:    record,
	})
}

// DeleteRecord godoc
// @Summary      Delete a DNS record
// @Description  Delete a DNS record by its ID (soft delete)
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "Domain ID"
// @Param        recordId  path      int  true  "Record ID"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records/{recordId} [delete]
func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	recordId := c.Param("recordId")
	var id uint
	if _, err := fmt.Sscanf(recordId, "%d", &id); err != nil {
		logger.Warnf("Invalid record ID: %s", recordId)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record ID",
		})
		return
	}

	if err := h.recordService.DeleteRecord(id); err != nil {
		logger.Errorf("Failed to delete record: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete record",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Record deleted successfully: ID %d", id)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Record deleted successfully",
	})
}

// DisableRecord godoc
// @Summary      Disable a DNS record
// @Description  Disable a DNS record by setting its disabled flag to true
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "Domain ID"
// @Param        recordId  path      int  true  "Record ID"
// @Success      200  {object}  dto.Response{data=dto.RecordResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records/{recordId}/disable [post]
func (h *RecordHandler) DisableRecord(c *gin.Context) {
	recordId := c.Param("recordId")
	var id uint
	if _, err := fmt.Sscanf(recordId, "%d", &id); err != nil {
		logger.Warnf("Invalid record ID: %s", recordId)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record ID",
		})
		return
	}

	record, err := h.recordService.DisableRecord(id)
	if err != nil {
		logger.Errorf("Failed to disable record: %v", err)
		if err == service.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to disable record",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Record disabled successfully: id=%d, host=%s", record.ID, record.Host)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Record disabled successfully",
		Data:    record,
	})
}

// EnableRecord godoc
// @Summary      Enable a DNS record
// @Description  Enable a DNS record by setting its disabled flag to false
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "Domain ID"
// @Param        recordId  path      int  true  "Record ID"
// @Success      200  {object}  dto.Response{data=dto.RecordResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /domains/{id}/records/{recordId}/enable [post]
func (h *RecordHandler) EnableRecord(c *gin.Context) {
	recordId := c.Param("recordId")
	var id uint
	if _, err := fmt.Sscanf(recordId, "%d", &id); err != nil {
		logger.Warnf("Invalid record ID: %s", recordId)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record ID",
		})
		return
	}

	record, err := h.recordService.EnableRecord(id)
	if err != nil {
		logger.Errorf("Failed to enable record: %v", err)
		if err == service.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to enable record",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Record enabled successfully: id=%d, host=%s", record.ID, record.Host)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Record enabled successfully",
		Data:    record,
	})
}
