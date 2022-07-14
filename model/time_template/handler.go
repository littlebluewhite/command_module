package time_template

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"new_command/pkg/response"
)

type Handler struct {
	DB       *gorm.DB
	Response response.Response
}

func (h *Handler) AddTimeTemplate(c *gin.Context) {
	entry := &TimeTemplate{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	err := checkTimeTemplate(entry)
	if err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(entry).Error; err != nil {
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

func (h *Handler) GetTimeTemplates(c *gin.Context) {
	var timeTemplates []TimeTemplate

	if result := h.DB.Preload("WeeklyRepeat").Preload("MonthlyRepeat").Find(&timeTemplates); result.Error != nil {
		h.Response.HttpResponse(c, http.StatusNotFound, result.Error).Printf(
			response.DefaultMessage, http.StatusNotFound, result.Error)
		return
	}

	h.Response.HttpResponse(c, http.StatusOK, &timeTemplates).Printf(response.DefaultMessage, http.StatusOK, &timeTemplates)
}

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

func (h *Handler) UpdateTimeTemplate(c *gin.Context) {
	entry := &TimeTemplatePatch{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	timeTemplateBody := entry.ToModel()
	err := checkTimeTemplate(&timeTemplateBody)
	if err != nil {
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
	h.Response.HttpResponse(c, http.StatusOK, &timeTemplate).Printf(
		response.DefaultMessage, http.StatusOK, &timeTemplate)
}
