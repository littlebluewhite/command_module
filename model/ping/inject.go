package ping

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"new_command/app"
)

type HandlerConfig struct {
	R  *gin.Engine
	DB *gorm.DB
}

func Inject(modelConfig app.ModelConfig) {
	InitRoutes(HandlerConfig{
		R:  modelConfig.Router,
		DB: modelConfig.DB,
	})
}
