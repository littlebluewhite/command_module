package command

import "fmt"

func (c *Command) Execute() (err error) {
	switch c.Protocol {
	case Http:
		if c.HttpsCommand != nil {
			result := c.HttpsCommand.Execute()
			fmt.Println(result)
		}
	default:
	}
	return
}
