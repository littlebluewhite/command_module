package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"new_command/app"
)

type HandlerConfig struct {
	R  *gin.Engine
	DB *gorm.DB
	C  *cache.Cache
}

func Inject(modelConfig app.ModelConfig) {
	InitRoutes(HandlerConfig{
		R:  modelConfig.Router,
		DB: modelConfig.DB,
		C:  modelConfig.Cache,
	})
}
