package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"new_command/app"
	"new_command/app/database"
	"new_command/config"
	"new_command/middleware"
	"new_command/model"
	"new_command/pkg/logFile"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	mainLog logFile.LogFile
)

// 初始化配置
func init() {
	// log配置
	mainLog = logFile.NewLogFile("", "main.log")
}

func main() {
	mainLog.Info().Println("command module start")

	// DB start
	DB, err := database.NewDB("mySQL", "DB.log", "db")
	if err != nil {
		mainLog.Error().Println("DB Connection failed")
		panic(err)
	} else {
		mainLog.Info().Println("DB Connection successful")
	}
	defer func() {
		closeErr := database.CloseDB(DB)
		if closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Println("Error occurred while closing the DB :", closeErr)
			}
		}
	}()

	// gin start
	ginFile, err := os.OpenFile("./log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file: " + err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(ginFile, os.Stdout)
	r := gin.Default()

	// cors middleware
	r.Use(middleware.CORSMiddleware())

	gin.SetMode(gin.ReleaseMode)

	// convert model config
	modelConfig := app.ModelConfig{
		DB:     DB,
		Router: r,
	}
	model.Inject(modelConfig)

	// server config
	ServerConfig := config.NewConfig[config.ServerConfig](".", "env", "server")

	// server
	var sb strings.Builder
	sb.WriteString(":")
	sb.WriteString(ServerConfig.Port)
	srv := &http.Server{
		Addr:           sb.String(),
		Handler:        r,
		ReadTimeout:    ServerConfig.ReadTimeout * time.Second,
		WriteTimeout:   ServerConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
