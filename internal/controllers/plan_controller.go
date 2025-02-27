package controllers

import (
	"BigDataForge/internal/elastic"
	"BigDataForge/internal/services"
	"BigDataForge/internal/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type PlanController struct {
	Service *services.PlanService
}

func NewPlanController(redisClient *redis.Client, esFactory *elastic.Factory) *PlanController {
	return &PlanController{
		Service: services.NewPlanService(redisClient, esFactory),
	}
}

func (controller *PlanController) CreatePlan(c *gin.Context) {

	if !validators.ValidatePlanSchema(c) {
		return
	}

	controller.Service.CreatePlan(c)
}

func (controller *PlanController) GetPlan(c *gin.Context) {
	controller.Service.GetPlan(c)
}

func (controller *PlanController) DeletePlan(c *gin.Context) {
	controller.Service.DeletePlan(c)
}

func (controller *PlanController) PatchPlan(c *gin.Context) {
	if !validators.ValidatePlanSchema(c) {
		return
	}
	controller.Service.PatchPlan(c)
}

func (controller *PlanController) UpdatePlan(c *gin.Context) {
	if !validators.ValidatePlanSchema(c) {
		return
	}
	controller.Service.UpdatePlan(c)
}

func (controller *PlanController) SearchPlans(c *gin.Context) {
	controller.Service.SearchPlans(c)
}
