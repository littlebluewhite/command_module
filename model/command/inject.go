package command

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"new_command/app"
	"new_command/pkg/response"
)

type HandlerConfig struct {
	R        *gin.Engine
	DB       *gorm.DB
	Response response.Response
	C        *cache.Cache
}

func Inject(modelConfig app.ModelConfig) {
	resp := response.NewResponse("model", "command.log")
	InitRoutes(HandlerConfig{
		R:        modelConfig.Router,
		DB:       modelConfig.DB,
		Response: resp,
		C:        modelConfig.Cache,
	})
}
