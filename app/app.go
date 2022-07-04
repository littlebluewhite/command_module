package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModelConfig struct {
	DB     *gorm.DB
	Router *gin.Engine
}
