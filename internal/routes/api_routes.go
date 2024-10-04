package routes

import (
	"BigDataForge/internal/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupRoutes(router *gin.Engine, redisClient *redis.Client) {
	planController := controllers.NewPlanController(redisClient)

	api := router.Group("/api/v1")
	{
		api.POST("/plans", planController.CreatePlan)
		api.GET("/plans", planController.GetPlan)
		api.DELETE("/plans", planController.DeletePlan)
		api.GET("/plans/conditional", planController.ConditionalGetPlan)
	}
}
