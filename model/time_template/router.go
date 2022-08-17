package time_template

func InitRoutes(c HandlerConfig) {
	h := Handler{
		DB:       c.DB,
		Response: c.Response,
		C:        c.C,
	}

	//set api
	g := c.R.Group("/time_template")

	g.GET("/api", h.GetTimeTemplates)
	g.GET("/api/:id", h.GetTimeTemplateById)
	g.POST("/api", h.AddTimeTemplate)
	g.PATCH("/api/:id", h.UpdateTimeTemplate)
	g.DELETE("/api/:id", h.DeleteTimeTemplate)
}
