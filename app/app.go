package app

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type ModelConfig struct {
	DB     *gorm.DB
	Router *gin.Engine
	Cache  *cache.Cache
}
