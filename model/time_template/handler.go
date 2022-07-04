package time_template

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"new_command/app/logFile"
)

var (
	timeTemplateLog logFile.LogFile
)

func init() {
	timeTemplateLog = logFile.NewLogFile("model", "time_template.log")
}

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) AddTimeTemplate(c *gin.Context) {
	entry, err := checkAddTimeTemplate(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, err)
		return
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(entry).Error; err != nil {
			timeTemplateLog.Error().Printf("err: %v, err type: %T", err, err)
			return err
		}
		timeTemplateLog.Info().Printf("add New Time Template: %v", entry.Name)
		return nil
	})
	if e != nil {
		timeTemplateLog.Error().Printf("err: %v, err type: %T", e, e)
		c.AbortWithStatusJSON(403, e)
		return
	}
	c.JSON(202, gin.H{
		"message": "success",
	})
}
