package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"example.c/docs"
)

// @title Config UI Example
// @version 1.0
// @BasePath /api/v1
func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		hr := v1.Group("/health")
		hc := NewHealthController()
		{
			hr.GET("", hc.Get)
			hr.POST("", hc.Set)
		}

		rc := NewReadyController()
		rr := v1.Group("/ready")
		{
			rr.GET("", rc.Get)
			rr.POST("", rc.Set)
		}
	}

	// OpenAPI reflection
	log.Printf("Serving Swagger UI at /swagger for API %s", docs.SwaggerInfo.Title)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Exiting...")
}

type Health struct {
	Health string `json:"health" binding:"required" example:"very fit"`
}
type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// Get godoc
// @Router /health [get]
// @Summary Get health
// @Success 200 {object} Health "Current health"
func (ctrl *HealthController) Get(c *gin.Context) {
	data := Health{
		Health: "foo",
	}

	c.JSON(http.StatusOK, data)
}

// Set godoc
// @Router /health [post]
// @Summary Set health
// @Param health body Health true "New health"
func (ctrl *HealthController) Set(c *gin.Context) {
	var data Health
	if err := c.BindJSON(&data); err != nil {
		return
	}

	log.Println("Got:", data)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

type Ready struct {
	Ready string `json:"ready" binding:"required"`
}
type ReadyController struct {
}

func NewReadyController() *ReadyController {
	return &ReadyController{}
}

// Get godoc
// @Summary Get readiness
// @Router /ready [get]
func (ctrl *ReadyController) Get(c *gin.Context) {
	data := Ready{
		Ready: "bar",
	}

	c.JSON(http.StatusOK, data)
}

// Set godoc
// @Summary Set readiness
// @Param ready body string true "new readiness"
// @Router /ready [post]
func (ctrl *ReadyController) Set(c *gin.Context) {
	var data Ready
	if err := c.BindJSON(&data); err != nil {
		return
	}

	log.Println("Got:", data)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
