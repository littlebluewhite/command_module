package header_template

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"new_command/pkg/response"
	"strconv"
	"strings"
)

type Handler struct {
	DB       *gorm.DB
	Response response.Response
	C        *cache.Cache
}

// GetHeaderTemplates swagger
// @Summary Show all header templates
// @Description Get all header templates
// @Tags header_template
// @Produce json
// @Success 200 {array} header_template.SwaggerResponse
// @Router /header_template/api [get]
func (h *Handler) GetHeaderTemplates(c *gin.Context) {
	var headerTemplates []HeaderTemplate
	cacheMap := map[int]HeaderTemplate{}
	if x, found := h.C.Get("header_templates"); found {
		cacheMap = x.(map[int]HeaderTemplate)
	} else {
		h.Response.HttpResponse(c, http.StatusNotFound, "cache error").Printf(
			response.DefaultMessage, http.StatusNotFound, "cache error")
		return
	}
	for _, value := range cacheMap {
		headerTemplates = append(headerTemplates, value)
	}

	h.Response.HttpResponse(c, http.StatusOK, &headerTemplates).Printf(response.DefaultMessage, http.StatusOK, &headerTemplates)
}

// GetHeaderTemplateById swagger
// @Summary Show header templates
// @Description Get header templates by id
// @Tags header_template
// @Produce json
// @Param id path int true "header template id"
// @Success 200 {object} header_template.SwaggerResponse
// @Router /header_template/api/{id} [get]
func (h *Handler) GetHeaderTemplateById(c *gin.Context) {
	id := c.Param("id")
	cacheMap := map[int]HeaderTemplate{}
	if x, found := h.C.Get("header_templates"); found {
		cacheMap = x.(map[int]HeaderTemplate)
	} else {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "cache error"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "cache error"})
		return
	}
	idInt, _ := strconv.Atoi(id)
	schedule, ok := cacheMap[idInt]
	if !ok {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}

	h.Response.HttpResponse(c, http.StatusOK, &schedule).Printf(
		response.DefaultMessage, http.StatusOK, &schedule)
}

// AddHeaderTemplate swagger
// @Summary Create header templates
// @Tags header_template
// @Accept json
// @Produce json
// @Param header_template body header_template.SwaggerCreate true "header template body"
// @Success 200 string string
// @Router /header_template/api [post]
func (h *Handler) AddHeaderTemplate(c *gin.Context) {
	entry := &HeaderTemplate{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entry).Error; err != nil {
			return err
		}
		if x, found := h.C.Get("header_templates"); found {
			cacheMap := x.(map[int]HeaderTemplate)
			cacheMap[entry.ID] = *entry
			h.C.Set("header_templates", cacheMap, cache.NoExpiration)
		} else {
			return errors.New("cache error")
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

// UpdateHeaderTemplate  swagger
// @Summary Update header templates
// @Tags header_template
// @Accept json
// @Produce json
// @Param id path int true "header template id"
// @Param header_template body header_template.SwaggerUpdate true "modify header template body"
// @Success 200 {object} header_template.SwaggerResponse
// @Router /header_template/api/{id} [patch]
func (h *Handler) UpdateHeaderTemplate(c *gin.Context) {
	entry := &HeaderTemplatePatch{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	var headerTemplates HeaderTemplate
	id := c.Param("id")
	h.DB.First(&headerTemplates, id)
	if headerTemplates.ID == 0 {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}
	headerTemplates.UpdateData(*entry)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&headerTemplates).Error; err != nil {
			return err
		}
		if x, found := h.C.Get("header_templates"); found {
			cacheMap := x.(map[int]HeaderTemplate)
			cacheMap[headerTemplates.ID] = headerTemplates
			h.C.Set("header_templates", cacheMap, cache.NoExpiration)
		} else {
			return errors.New("cache error")
		}
		return nil
	})
	if err != nil {
		h.Response.HttpResponse(c, http.StatusBadRequest, err).Printf(
			response.DefaultMessage, http.StatusBadRequest, err)
		return
	}
	h.Response.HttpResponse(c, http.StatusOK, &headerTemplates).Printf(
		response.DefaultMessage, http.StatusOK, &headerTemplates)
}

// DeleteHeaderTemplate swagger
// @Summary Delete header templates
// @Tags header_template
// @Produce json
// @Param id path int true "header template id"
// @Success 200 string string
// @Router /header_template/api/{id} [delete]
func (h *Handler) DeleteHeaderTemplate(c *gin.Context) {
	idString := c.Param("id")
	var headerTemplate HeaderTemplate
	h.DB.Preload(clause.Associations).First(&headerTemplate, idString)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&headerTemplate).Error; err != nil {
			return err
		}
		if x, found := h.C.Get("header_templates"); found {
			cacheMap := x.(map[int]HeaderTemplate)
			delete(cacheMap, headerTemplate.ID)
			h.C.Set("header_templates", cacheMap, cache.NoExpiration)
		} else {
			return errors.New("cache error")
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
