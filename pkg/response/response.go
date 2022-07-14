package response

import (
	"github.com/gin-gonic/gin"
	"log"
	"new_command/pkg/logFile"
)

type Response interface {
	HttpResponse(c *gin.Context, statusCode int, message any) *log.Logger
}

const DefaultMessage = "Response code: %v, Message: %v"

type response struct {
	logFile logFile.LogFile
}

func NewResponse(dirPath string, fileName string) (r Response) {
	file := logFile.NewLogFile(dirPath, fileName)
	r = &response{logFile: file}
	return
}

func (r *response) HttpResponse(c *gin.Context, statusCode int, message any) (l *log.Logger) {
	firstCode := statusCode / 100
	switch firstCode {
	case 2:
		c.JSON(statusCode, message)
		l = r.logFile.Info()
	default:
		c.AbortWithStatusJSON(statusCode, message)
		l = r.logFile.Error()
	}
	return
}
