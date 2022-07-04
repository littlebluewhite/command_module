package time_template

func InitRoutes(c HandlerConfig) {

	h := Handler{
		DB: c.DB,
	}

	//set api
	g := c.R.Group("/time_template")

	g.POST("/", h.AddTimeTemplate)
}
