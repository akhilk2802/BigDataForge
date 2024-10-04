package services

import (
	"context"
	"encoding/json"
	"net/http"

	"BigDataForge/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type PlanService struct {
	redisClient *redis.Client
}

func NewPlanService(redisClient *redis.Client) *PlanService {
	return &PlanService{redisClient: redisClient}
}

// Create a new Plan
func (service *PlanService) CreatePlan(c *gin.Context) {
	var plan models.Plan
	if err := c.BindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	planID := plan.ObjectID
	planJSON, err := json.Marshal(plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process plan"})
		return
	}

	err = service.redisClient.Set(ctx, "plan:"+planID, planJSON, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store plan in Redis"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Plan created", "planId": planID})
}

// Retrieve a Plan
func (service *PlanService) GetPlan(c *gin.Context) {
	planID := c.Query("id")
	planJSON, err := service.redisClient.Get(ctx, "plan:"+planID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}

	var plan models.Plan
	err = json.Unmarshal([]byte(planJSON), &plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse plan data"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// Delete a Plan
func (service *PlanService) DeletePlan(c *gin.Context) {
	planID := c.Query("id")
	err := service.redisClient.Del(ctx, "plan:"+planID).Err()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
}

// Conditional Retrieve Plan (example condition: "planType" == "inNetwork")
func (service *PlanService) ConditionalGetPlan(c *gin.Context) {
	planID := c.Query("id")
	planJSON, err := service.redisClient.Get(ctx, "plan:"+planID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}

	var plan models.Plan
	err = json.Unmarshal([]byte(planJSON), &plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse plan data"})
		return
	}

	// Example condition: return plan only if planType is "inNetwork"
	if plan.PlanType == "inNetwork" {
		c.JSON(http.StatusOK, plan)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Condition not met"})
	}
}
