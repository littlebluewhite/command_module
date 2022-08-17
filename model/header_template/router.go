package header_template

import (
	"github.com/patrickmn/go-cache"
)

func InitRoutes(c HandlerConfig) {
	var schedules []HeaderTemplate
	cacheMap := map[int]HeaderTemplate{}
	for i := 0; i < len(schedules); i++ {
		entry := schedules[i]
		cacheMap[entry.ID] = entry
	}
	c.C.Set("header_templates", cacheMap, cache.NoExpiration)

	h := Handler{
		DB:       c.DB,
		Response: c.Response,
		C:        c.C,
	}

	//set api
	g := c.R.Group("/header_template")

	g.GET("/api", h.GetHeaderTemplates)
	g.GET("/api/:id", h.GetHeaderTemplateById)
	g.POST("/api", h.AddHeaderTemplate)
	g.PATCH("/api/:id", h.UpdateHeaderTemplate)
	g.DELETE("/api/:id", h.DeleteHeaderTemplate)
}
