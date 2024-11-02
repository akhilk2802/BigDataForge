package controllers

import (
	"BigDataForge/internal/services"
	"BigDataForge/internal/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type PlanController struct {
	Service *services.PlanService
}

func NewPlanController(redisClient *redis.Client) *PlanController {
	return &PlanController{
		Service: services.NewPlanService(redisClient),
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
