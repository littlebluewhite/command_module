package ping

func InitRoutes(c HandlerConfig) {

	h := Handler{
		DB: c.DB,
		C:  c.C,
	}

	//set api
	g := c.R.Group("/ping")

	g.GET("/", h.GetPing)
	g.GET("/list", h.GetListPing)
}
