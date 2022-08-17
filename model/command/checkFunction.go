package command

import (
	"new_command/util"
	"time"
)

func (c *Command) CheckCommand() (err error) {
	err = nil
	switch c.Protocol {
	case Http:
		if c.HttpsCommand == nil {
			err = &util.MyError{
				When: time.Now(),
				What: "lose https command!",
			}
			return
		}
		err = c.HttpsCommand.CheckHttpsCommand()
	default:
	}
	return
}
