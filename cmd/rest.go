package cmd

import (
	"context"
	"fmt"
	"gogin-template/baselib/middleware"
	"gogin-template/bootstrap"
	"gogin-template/internal/controller"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gogin-template/docs"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "A Rest API",
	Long:  "This is a Rest API.",
	Run: func(cmd *cobra.Command, args []string) {
		Rest()
	},
}

// @title           Gin Template for REST APIs
// @version         1.0
// @description     This is a template for developing REST APIs using the Gin framework.

// @host      localhost:8080
// @BasePath  /api/v1

func Rest() *cobra.Command {
	cfg := bootstrap.Init()
	cfg.UpdateLogger(cfg.Logger().WithField("component", "rest"))
	cfg.Logger().Info("running rest")
	config := cfg.GetConfig()

	servName := os.Getenv("SERVICE_NAME")

	// Setup Gin-Gonic
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	ginEngine := gin.Default()
	ginEngine.RedirectTrailingSlash = true
	ginEngine.RemoveExtraSlash = true
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(cors.New(corsConfig))
	ginEngine.Use(otelgin.Middleware(servName))
	ginEngine.Use(middleware.LoggingMiddleware(cfg))
	ginEngine.Use(middleware.ExceptionMiddleware(cfg))

	// Create Health
	controller.NewHealthController(ginEngine, cfg)

	// Get Configs

	// Rest Clients

	// Repositories

	// Services

	// Controllers

	// Define path for Swaggo
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the Server
	port := config.GetString("server.port")
	if port == "" {
		cfg.Logger().Fatal("server port has not been set")
	}

	server := &http.Server{
		Handler: ginEngine,
		Addr:    fmt.Sprintf(":%s", port),
	}

	go func() {
		cfg.Logger().Infof("server starting at %s", port)
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				cfg.Logger().Info("server stopped")
			} else {
				cfg.Logger().Fatal(fmt.Sprintf("failed to start server %s", err))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	cfg.Logger().Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	cfg.Logger().Info("shutting down the server...")
	cfg.Close()

	if err := server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		cfg.Logger().Fatal(fmt.Sprintf("failed to gracefully shut down the server %s", err))
	}

	return nil
}
