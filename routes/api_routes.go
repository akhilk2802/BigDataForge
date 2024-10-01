package routes

import (
	"BigDataForge/internal/controllers"
	"BigDataForge/internal/services"
	"BigDataForge/internal/storage"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	redisStore := storage.NewRedisStore()
	planService := services.NewPlanService(redisStore)

	api := router.Group("/api/v1")
	{
		api.POST("/plans", controllers.CreatePlan(planService))
		api.GET("/plans/:id", controllers.GetPlan(planService))
		api.DELETE("/plans/:id", controllers.DeletePlan(planService))
	}
}
