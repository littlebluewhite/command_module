package time_template

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"new_command/app"
	"new_command/pkg/response"
)

type HandlerConfig struct {
	R        *gin.Engine
	DB       *gorm.DB
	Response response.Response
}

func Inject(modelConfig app.ModelConfig) {
	resp := response.NewResponse("model", "time_template.log")
	InitRoutes(HandlerConfig{
		R:        modelConfig.Router,
		DB:       modelConfig.DB,
		Response: resp,
	})
}
