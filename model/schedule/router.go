package schedule

import "C"
import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

func InitRoutes(c HandlerConfig) {
	var schedules []Schedule
	c.DB.Preload("Command", func(db *gorm.DB) *gorm.DB {
		return db.Preload("HttpsCommand")
	}).Preload("TimeData").Find(&schedules)
	cacheMap := map[int]Schedule{}
	for i := 0; i < len(schedules); i++ {
		entry := schedules[i]
		cacheMap[entry.ID] = entry
	}
	c.C.Set("schedules", cacheMap, cache.NoExpiration)

	h := Handler{
		DB:       c.DB,
		Response: c.Response,
		C:        c.C,
	}

	//set api
	g := c.R.Group("/schedule")

	g.GET("/api", h.GetSchedules)
	g.GET("/api/:id", h.GetScheduleById)
	g.POST("/api", h.AddSchedule)
	g.PATCH("/api/:id", h.UpdateSchedule)
	g.DELETE("/api/:id", h.DeleteSchedule)
}
