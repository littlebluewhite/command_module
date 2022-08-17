package model

import (
	"new_command/app"
	"new_command/model/command"
	"new_command/model/header_template"
	"new_command/model/ping"
	"new_command/model/schedule"
	"new_command/model/time_template"
)

func Inject(modelConfig app.ModelConfig) {
	modelConfig.Router.Use(Logger())

	ping.Inject(modelConfig)
	time_template.Inject(modelConfig)
	command.Inject(modelConfig)
	schedule.Inject(modelConfig)
	header_template.Inject(modelConfig)
}
