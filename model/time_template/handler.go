package time_template

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"new_command/pkg/response"
	"strings"
)

type Handler struct {
	DB       *gorm.DB
	Response response.Response
	C        *cache.Cache
}

// GetTimeTemplates
// @Summary Show all time templates
// @Description Get all time templates
// @Tags time_template
// @Produce json
// @Success 200 {array} time_template.SwaggerResponse
// @Router /time_template/api [get]
func (h *Handler) GetTimeTemplates(c *gin.Context) {
	var timeTemplates []TimeTemplate

	if result := h.DB.Preload("TimeData").Find(&timeTemplates); result.Error != nil {
		h.Response.HttpResponse(c, http.StatusNotFound, result.Error).Printf(
			response.DefaultMessage, http.StatusNotFound, result.Error)
		return
	}

	h.Response.HttpResponse(c, http.StatusOK, &timeTemplates).Printf(response.DefaultMessage, http.StatusOK, &timeTemplates)
}

// GetTimeTemplateById swagger
// @Summary Show time templates
// @Description Get time templates by id
// @Tags time_template
// @Produce json
// @Param id path int true "time template id"
// @Success 200 {object} time_template.SwaggerResponse
// @Router /time_template/api/{id} [get]
func (h *Handler) GetTimeTemplateById(c *gin.Context) {
	var timeTemplate TimeTemplate
	id := c.Param("id")
	h.DB.Preload(clause.Associations).First(&timeTemplate, id)
	if timeTemplate.ID == 0 {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}
	h.Response.HttpResponse(c, http.StatusOK, &timeTemplate).Printf(
		response.DefaultMessage, http.StatusOK, &timeTemplate)
}

// AddTimeTemplate swagger
// @Summary Create time templates
// @Tags time_template
// @Accept json
// @Produce json
// @Param time_template body time_template.SwaggerCreate true "time template body"
// @Success 200 {string} string "created success"
// @Router /time_template/api [post]
func (h *Handler) AddTimeTemplate(c *gin.Context) {
	entry := &TimeTemplate{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	if err := entry.CheckTimeTemplate(); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entry).Error; err != nil {
			return err
		}
		return nil
	})
	if e != nil {
		h.Response.HttpResponse(c, 403, e).Printf(response.DefaultMessage, 403, e)
		return
	}
	h.Response.HttpResponse(c, http.StatusCreated, gin.H{
		"message": "created success",
	}).Printf(response.DefaultMessage, http.StatusCreated, gin.H{
		"message": "created success"})
}

// UpdateTimeTemplate swagger
// @Summary Update time templates
// @Tags time_template
// @Accept json
// @Produce json
// @Param id path int true "time template id"
// @Param time_template body time_template.SwaggerUpdate true "modify time template body"
// @Success 200 {object} time_template.SwaggerResponse
// @Router /time_template/api/{id} [patch]
func (h *Handler) UpdateTimeTemplate(c *gin.Context) {
	entry := &TimeTemplatePatch{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	var timeTemplate TimeTemplate
	id := c.Param("id")
	h.DB.Preload(clause.Associations).First(&timeTemplate, id)
	if timeTemplate.ID == 0 {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}
	timeTemplate.UpdateData(*entry)
	fmt.Println(timeTemplate.TimeData)
	if err := timeTemplate.CheckTimeTemplate(); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&timeTemplate).Error; err != nil {
			return err
		}
		if err := tx.Save(timeTemplate.TimeData).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		h.Response.HttpResponse(c, http.StatusBadRequest, err).Printf(
			response.DefaultMessage, http.StatusBadRequest, err)
		return
	}
	h.Response.HttpResponse(c, http.StatusOK, &timeTemplate).Printf(
		response.DefaultMessage, http.StatusOK, &timeTemplate)
}

// DeleteTimeTemplate swagger
// @Summary Delete time templates
// @Tags time_template
// @Produce json
// @Param id path int true "time template id"
// @Success 200 {string} string "delete successfully"
// @Router /time_template/api/{id} [delete]
func (h *Handler) DeleteTimeTemplate(c *gin.Context) {
	idString := c.Param("id")
	var timeTemplate TimeTemplate
	h.DB.Preload(clause.Associations).First(&timeTemplate, idString)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&timeTemplate).Error; err != nil {
			return err
		}
		if err := tx.Delete(timeTemplate.TimeData).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		h.Response.HttpResponse(c, http.StatusBadRequest, err).Printf(
			response.DefaultMessage, http.StatusBadRequest, err)
		return
	}
	var sb strings.Builder
	sb.WriteString("id: ")
	sb.WriteString(idString)
	sb.WriteString(" has been deleted successfully")
	h.Response.HttpResponse(c, http.StatusOK, gin.H{"message": sb.String()}).Printf(
		response.DefaultMessage, http.StatusOK, gin.H{"message": sb.String()})
}
