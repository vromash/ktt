package http

import (
	"financing-aggregator/internal/exchange"
	"financing-aggregator/internal/mapper"
	"financing-aggregator/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type ApplicationHandler struct {
	svc services.ApplicationService
}

func NewApplicationHandler(svc services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		svc: svc,
	}
}

// SubmitApplication
//
// @Summary		Submit a new financing application
// @Description Accepts a JSON body with application details, validates input, and creates a new application.
// @Security 	BearerAuth
// @Tags		applications
// @Accept		json
// @Produce		json
// @Param		application body exchange.ApplicationRequest true "Application request"
// @Success		200 {object} exchange.ApplicationResponse
// @Failure		400 {object} exchange.ErrorResponse
// @Failure		500 {object} exchange.ErrorResponse
// @Router 		/applications [post]
func (h *ApplicationHandler) SubmitApplication(c *gin.Context) {
	var req exchange.ApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, exchange.NewErrorResponse(err.Error()))
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, exchange.NewErrorResponse(err.Error()))
		return
	}

	app, err := h.svc.SubmitApplication(c.Request.Context(), mapper.MapApplicationRequestToDTO(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, exchange.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, mapper.MapApplicationDTOToResponse(app))
}

// GetApplication
//
// @Summary		Get application by ID
// @Description Returns application details and offers for the given application ID.
// @Security 	BearerAuth
// @Tags		applications
// @Produce 	json
// @Param 		id path string true "Application ID"
// @Success 	200 {object} exchange.ApplicationResponse
// @Failure 	400 {object} exchange.ErrorResponse
// @Failure 	500 {object} exchange.ErrorResponse
// @Router 		/applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, exchange.NewErrorResponse("application id is required"))
		return
	}

	app, err := h.svc.GetApplication(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, exchange.NewErrorResponse("application not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, exchange.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, mapper.MapApplicationDTOToResponse(app))
}
