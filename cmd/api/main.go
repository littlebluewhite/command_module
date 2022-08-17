package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/http"
	"new_command/app"
	"new_command/app/database"
	"new_command/app/time_server"
	"new_command/config"
	_ "new_command/docs"
	"new_command/middleware"
	"new_command/model"
	"new_command/model/schedule"
	"new_command/pkg/logFile"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

// @title           Schedule-Command swagger API
// @version         1.0
// @description     This is a schedule-command server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Wilson
// @contact.url    https://github.com/littlebluewhite
// @contact.email  wwilson008@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5487

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
				mainLog.Error().Println("Error occurred while closing the DB :", closeErr)
			}
		}
	}()

	// go Cache
	c := cache.New(5*time.Minute, 10*time.Minute)

	// gin start
	ginFile, err := os.OpenFile("./log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		mainLog.Error().Fatal("can not open log file: " + err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(ginFile, os.Stdout)
	r := gin.Default()

	// cors middleware
	r.Use(middleware.CORSMiddleware())
	gin.SetMode(gin.ReleaseMode)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// convert model config
	modelConfig := app.ModelConfig{
		DB:     DB,
		Router: r,
		Cache:  c,
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

	// Time Server
	timeServer := time_server.NewTimeServer[schedule.Schedule](c, "schedules", 1*time.Second)
	go timeServer.Start(ctx)

	// API server Start
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mainLog.Error().Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	mainLog.Info().Println("API server shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		mainLog.Error().Fatal("Server forced to shutdown: ", err)
	}

	mainLog.Info().Println("Server exiting")
}
