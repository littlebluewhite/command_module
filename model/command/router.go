package command

func InitRoutes(c HandlerConfig) {
	h := Handler{
		DB:       c.DB,
		Response: c.Response,
		C:        c.C,
	}

	//set api
	g := c.R.Group("/command")

	g.GET("/api", h.GetCommands)
	g.GET("/api/:id", h.GetCommandById)
	g.POST("/api", h.AddCommand)
	g.PATCH("/api/:id", h.UpdateCommand)
	g.DELETE("/api/:id", h.DeleteCommand)
}
