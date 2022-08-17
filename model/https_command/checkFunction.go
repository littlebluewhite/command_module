package https_command

import (
	"log"
	"new_command/util"
	"time"
)

func (hc *HttpsCommand) CheckHttpsCommand() (err error) {
	err = nil
	log.Println(hc.Method)
	switch hc.Method {
	case Get:
		hc.BodyType = ""
		hc.Body = nil
	case Post, Patch, Put:
		if hc.BodyType == "" {
			err = &util.MyError{
				When: time.Now(),
				What: "http command need body type",
			}
			return
		}
		if hc.Body == nil {
			err = &util.MyError{
				When: time.Now(),
				What: "http command need body",
			}
			return
		}
	default:
	}
	return
}
