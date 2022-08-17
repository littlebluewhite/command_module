package command

import (
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

// GetCommands swagger
// @Summary Show all commands
// @Description Get all commands
// @Tags command
// @Produce json
// @Success 200 {array} command.SwaggerResponse
// @Router /command/api [get]
func (h *Handler) GetCommands(c *gin.Context) {
	var commands []Command

	if result := h.DB.Preload("HttpsCommand").Find(&commands); result.Error != nil {
		h.Response.HttpResponse(c, http.StatusNotFound, result.Error).Printf(
			response.DefaultMessage, http.StatusNotFound, result.Error)
		return
	}

	h.Response.HttpResponse(c, http.StatusOK, &commands).Printf(response.DefaultMessage, http.StatusOK, &commands)
}

// GetCommandById swagger
// @Summary Show commands
// @Description Get commands by id
// @Tags command
// @Produce json
// @Param id path int true "command id"
// @Success 200 {object} command.SwaggerResponse
// @Router /command/api/{id} [get]
func (h *Handler) GetCommandById(c *gin.Context) {
	var command Command
	id := c.Param("id")
	h.DB.Preload(clause.Associations).First(&command, id)
	if command.ID == 0 {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}
	h.Response.HttpResponse(c, http.StatusOK, &command).Printf(
		response.DefaultMessage, http.StatusOK, &command)
}

// AddCommand swagger
// @Summary Create commands
// @Tags command
// @Accept json
// @Produce json
// @Param command body command.SwaggerCreate true "command body"
// @Success 200 {string} string "created success"
// @Router /command/api [post]
func (h *Handler) AddCommand(c *gin.Context) {
	entry := &Command{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	if err := entry.CheckCommand(); err != nil {
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

// UpdateCommand swagger
// @Summary Update commands
// @Tags command
// @Accept json
// @Produce json
// @Param id path int true "command id"
// @Param command body command.SwaggerUpdate true "modify command body"
// @Success 200 {object} command.SwaggerResponse
// @Router /command/api/{id} [patch]
func (h *Handler) UpdateCommand(c *gin.Context) {
	entry := &CommandPatch{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	var command Command
	id := c.Param("id")
	h.DB.Preload(clause.Associations).First(&command, id)
	if command.ID == 0 {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}
	command.UpdateData(*entry)
	if err := command.CheckCommand(); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&command).Error; err != nil {
			return err
		}
		switch command.Protocol {
		case Http:
			if err := tx.Save(command.HttpsCommand).Error; err != nil {
				return err
			}
		default:
		}
		return nil
	})
	if err != nil {
		h.Response.HttpResponse(c, http.StatusBadRequest, err).Printf(
			response.DefaultMessage, http.StatusBadRequest, err)
		return
	}
	h.Response.HttpResponse(c, http.StatusOK, &command).Printf(
		response.DefaultMessage, http.StatusOK, &command)
}

// DeleteCommand swagger
// @Summary Delete commands
// @Tags command
// @Produce json
// @Param id path int true "command id"
// @Success 200 {string} string "delete successfully
// @Router /command/api/{id} [delete]
func (h *Handler) DeleteCommand(c *gin.Context) {
	idString := c.Param("id")
	var command Command
	h.DB.Preload(clause.Associations).First(&command, idString)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		switch command.Protocol {
		case Http:
			if err := tx.Delete(command.HttpsCommand).Error; err != nil {
				return err
			}
		default:
		}
		if err := tx.Delete(&command).Error; err != nil {
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
