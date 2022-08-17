package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"new_command/pkg/logFile"
	"time"
)

var (
	pingLog logFile.LogFile
)

func init() {
	pingLog = logFile.NewLogFile("model", "ping.log")
}

type Handler struct {
	DB *gorm.DB
	C  *cache.Cache
}

func (h *Handler) GetPing(c *gin.Context) {
	example := c.MustGet("example").(string)
	c.JSON(200, gin.H{
		"message": example,
	})
	pingLog.Info().Println("example: ", example)
}

func (h *Handler) GetListPing(c *gin.Context) {
	data := []map[string]interface{}{
		{
			"name": "wilson",
			"age":  5,
			"time": time.Now(),
		},
		{
			"name": "phoebe",
			"age":  4,
		},
	}
	c.JSON(200, data)
}
