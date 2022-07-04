package model

import (
	"new_command/app"
	"new_command/model/ping"
	"new_command/model/time_template"
)

func Inject(modelConfig app.ModelConfig) {
	modelConfig.Router.Use(Logger())

	ping.Inject(modelConfig)
	time_template.Inject(modelConfig)
}
