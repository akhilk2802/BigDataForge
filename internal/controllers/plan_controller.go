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

// POST - Create Plan
func (controller *PlanController) CreatePlan(c *gin.Context) {

	// Validate the incoming JSON against the schema
	if !validators.ValidatePlanSchema(c) {
		return
	}

	// Delegate to service
	controller.Service.CreatePlan(c)
}

// GET - Retrieve Plan
func (controller *PlanController) GetPlan(c *gin.Context) {
	controller.Service.GetPlan(c)
}

// DELETE - Delete Plan
func (controller *PlanController) DeletePlan(c *gin.Context) {
	controller.Service.DeletePlan(c)
}

// GET - Conditional Retrieve Plan
func (controller *PlanController) ConditionalGetPlan(c *gin.Context) {
	controller.Service.ConditionalGetPlan(c)
}
