package handler

import (
	"fmt"
	"net/http"

	"sds-admin/internal/dto"
	"sds-admin/internal/logger"
	"sds-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// RecordTypeHandler handles record type-related HTTP requests.
// RecordTypeHandler 处理记录类型相关的HTTP请求
type RecordTypeHandler struct {
	recordTypeService *service.RecordTypeService
}

// NewRecordTypeHandler creates a new RecordTypeHandler instance.
// NewRecordTypeHandler 创建新的RecordTypeHandler实例
func NewRecordTypeHandler(recordTypeService *service.RecordTypeService) *RecordTypeHandler {
	return &RecordTypeHandler{recordTypeService: recordTypeService}
}

// ListRecordTypes godoc
// @Summary      List all record types
// @Description  Retrieve a list of all record types
// @Tags         record-types
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.Response{data=[]dto.RecordTypeResponse}
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /record-types [get]
func (h *RecordTypeHandler) ListRecordTypes(c *gin.Context) {
	recordTypes, err := h.recordTypeService.ListRecordTypes()
	if err != nil {
		logger.Errorf("Failed to list record types: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list record types",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    recordTypes,
	})
}

// GetRecordType godoc
// @Summary      Get a record type by ID
// @Description  Retrieve a single record type by its ID
// @Tags         record-types
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Record Type ID"
// @Success      200  {object}  dto.Response{data=dto.RecordTypeResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /record-types/{id} [get]
func (h *RecordTypeHandler) GetRecordType(c *gin.Context) {
	id := c.Param("id")
	var recordTypeID uint
	if _, err := fmt.Sscanf(id, "%d", &recordTypeID); err != nil {
		logger.Warnf("Invalid record type ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record type ID",
		})
		return
	}

	recordType, err := h.recordTypeService.GetRecordTypeByID(recordTypeID)
	if err != nil {
		logger.Errorf("Failed to get record type: %v", err)
		if err == service.ErrRecordTypeNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record type not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get record type",
			Detail:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    recordType,
	})
}

// CreateRecordType godoc
// @Summary      Create a new record type
// @Description  Create a new record type with the provided information
// @Tags         record-types
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateRecordTypeRequest  true  "Record type creation request"
// @Success      201  {object}  dto.Response{data=dto.RecordTypeResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /record-types [post]
func (h *RecordTypeHandler) CreateRecordType(c *gin.Context) {
	var req dto.CreateRecordTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	recordType, err := h.recordTypeService.CreateRecordType(&req)
	if err != nil {
		logger.Errorf("Failed to create record type: %v", err)
		if err == service.ErrRecordTypeExists {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Record type already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create record type",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Record type created successfully: %s", recordType.Name)
	c.JSON(http.StatusCreated, dto.Response{
		Code:    http.StatusCreated,
		Message: "Record type created successfully",
		Data:    recordType,
	})
}

// UpdateRecordType godoc
// @Summary      Update a record type
// @Description  Update an existing record type with new information
// @Tags         record-types
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "Record Type ID"
// @Param        request  body      dto.UpdateRecordTypeRequest  true  "Record type update request"
// @Success      200  {object}  dto.Response{data=dto.RecordTypeResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /record-types/{id} [put]
func (h *RecordTypeHandler) UpdateRecordType(c *gin.Context) {
	id := c.Param("id")
	var recordTypeID uint
	if _, err := fmt.Sscanf(id, "%d", &recordTypeID); err != nil {
		logger.Warnf("Invalid record type ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record type ID",
		})
		return
	}

	var req dto.UpdateRecordTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
		return
	}

	recordType, err := h.recordTypeService.UpdateRecordType(recordTypeID, &req)
	if err != nil {
		logger.Errorf("Failed to update record type: %v", err)
		if err == service.ErrRecordTypeNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record type not found",
			})
			return
		}
		if err == service.ErrRecordTypeExists {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Record type already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update record type",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Record type updated successfully: %s", recordType.Name)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Record type updated successfully",
		Data:    recordType,
	})
}

// DeleteRecordType godoc
// @Summary      Delete a record type
// @Description  Delete a record type by its ID (soft delete)
// @Tags         record-types
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Record Type ID"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /record-types/{id} [delete]
func (h *RecordTypeHandler) DeleteRecordType(c *gin.Context) {
	id := c.Param("id")
	var recordTypeID uint
	if _, err := fmt.Sscanf(id, "%d", &recordTypeID); err != nil {
		logger.Warnf("Invalid record type ID: %s", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid record type ID",
		})
		return
	}

	if err := h.recordTypeService.DeleteRecordType(recordTypeID); err != nil {
		logger.Errorf("Failed to delete record type: %v", err)
		if err == service.ErrRecordTypeNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Record type not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete record type",
			Detail:  err.Error(),
		})
		return
	}

	logger.Infof("Record type deleted successfully: ID %d", recordTypeID)
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "Record type deleted successfully",
	})
}
