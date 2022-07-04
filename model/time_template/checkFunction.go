package time_template

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func checkAddTimeTemplate(c *gin.Context) (entry *TimeTemplate, err error) {
	entry = &TimeTemplate{}
	if err = c.ShouldBindJSON(entry); err != nil {
		fmt.Println(err)
		return nil, err
	}
	ch := make(chan error)
	go func(entry *TimeTemplate, ch chan error) {
		ch <- entry.CheckRepeatType()
	}(entry, ch)
	if err = <-ch; err != nil {
		return nil, err
	}
	return
}
