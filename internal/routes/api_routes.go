package routes

import (
	"BigDataForge/internal/controllers"
	"BigDataForge/internal/elastic"
	"BigDataForge/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupRoutes(router *gin.Engine, redisClient *redis.Client, esFactory *elastic.Factory) {
	planController := controllers.NewPlanController(redisClient, esFactory)

	api := router.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware()) // Apply AuthMiddleware to protect all routes in this group
	{
		api.POST("/plans", planController.CreatePlan)
		api.GET("/plans", planController.GetPlan)
		api.DELETE("/plans", planController.DeletePlan)
		api.PATCH("/plans", planController.PatchPlan)
		api.PUT("/plans", planController.UpdatePlan)
		api.POST("/search", planController.SearchPlans)
	}
}
