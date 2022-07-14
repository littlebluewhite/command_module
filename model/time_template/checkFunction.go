package time_template

import (
	"fmt"
)

func checkTimeTemplate(entry *TimeTemplate) (err error) {
	ch := make(chan error, 3)
	defer close(ch)
	go func(entry *TimeTemplate, ch chan error) {
		ch <- entry.CheckRepeatType()
	}(entry, ch)
	go func(entry *TimeTemplate, ch chan error) {
		ch <- entry.CheckTime()
	}(entry, ch)
	go func(entry *TimeTemplate, ch chan error) {
		ch <- entry.CheckDate()
	}(entry, ch)
	for i := 0; i < 3; i++ {
		select {
		case e := <-ch:
			if e != nil {
				entry = nil
				err = e
				fmt.Println(err)
			}
		}
	}
	return
}
