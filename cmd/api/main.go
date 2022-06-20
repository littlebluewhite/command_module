package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"new_command/app/database"
	"new_command/config"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Trace        *log.Logger
	Info         *log.Logger
	Error        *log.Logger
	ServerConfig *config.ServerConfig
)

// 初始化配置
func init() {
	// log配置
	newPath := filepath.Join(".", "log")
	_ = os.MkdirAll(newPath, os.ModePerm)
	file, err := os.OpenFile("./log/main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file: " + err.Error())
	}

	Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func init() {
	//server config
	ServerConfig = config.LoadConfig[*(config.ServerConfig)]("./evn", "server")
}

func main() {
	Info.Println("command module start")

	// DB start
	DB, err := database.NewDB()
	if err != nil {
		Error.Println("DB Connection failed")
		panic(err)
	} else {
		Info.Println("DB Connection successful")
	}
	defer database.CloseDB(DB)

	// gin start
	ginFile, err := os.OpenFile("./log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file: " + err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(ginFile, os.Stdout)
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	var sb strings.Builder
	sb.WriteString(":")
	sb.WriteString(ServerConfig.Port)
	s := &http.Server{
		Addr:           sb.String(),
		Handler:        r,
		ReadTimeout:    ServerConfig.ReadTimeout * time.Second,
		WriteTimeout:   ServerConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal("Server can not run: " + err.Error())
	}
}
