package services

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"BigDataForge/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type PlanService struct {
	redisClient redis.Cmdable
}

func NewPlanService(redisClient *redis.Client) *PlanService {
	return &PlanService{redisClient: redisClient}
}

// Helper function to generate ETag based on Plan data
func generateETag(plan models.Plan) string {
	planData, _ := json.Marshal(plan)
	hash := sha1.New()
	hash.Write(planData)
	return hex.EncodeToString(hash.Sum(nil))
}

func (service *PlanService) CreatePlan(c *gin.Context) {
	var plan models.Plan
	if err := c.BindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	planID := plan.ObjectID

	// Check if the plan already exists in Redis
	existingPlan, err := service.redisClient.Get(ctx, "plan:"+planID).Result()
	if err != redis.Nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing data"})
			return
		}
		if existingPlan != "" {
			c.JSON(http.StatusConflict, gin.H{"error": "Plan already exists"})
			return
		}
	}

	// Store the new plan in Redis
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

	// Generate and set the ETag in the response header
	c.Header("ETag", generateETag(plan))
	c.JSON(http.StatusCreated, gin.H{"message": "Plan created", "planId": planID})
}

func (service *PlanService) GetPlan(c *gin.Context) {
	planID := c.Query("id")

	// Retrieve the plan from Redis
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

	eTag := generateETag(plan)
	// Check for the If-None-Match header to support conditional reads
	ifMatch := c.GetHeader("If-None-Match")
	if ifMatch == eTag {
		c.Status(http.StatusNotModified)
		return
	}

	// Set ETag header and return the plan
	c.Header("ETag", eTag)
	c.JSON(http.StatusOK, plan)
}

func (service *PlanService) DeletePlan(c *gin.Context) {
	planID := c.Query("id")

	// Check if the plan exists in Redis
	exists, err := service.redisClient.Exists(ctx, "plan:"+planID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check plan existence"})
		return
	}

	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Delete the plan from Redis
	err = service.redisClient.Del(ctx, "plan:"+planID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (service *PlanService) PatchPlan(c *gin.Context) {
	planID := c.Query("id")

	// Retrieve existing plan data from Redis
	planJSON, err := service.redisClient.Get(ctx, "plan:"+planID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}

	var existingPlan models.Plan
	err = json.Unmarshal([]byte(planJSON), &existingPlan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse plan data"})
		return
	}

	currentETag := generateETag(existingPlan)
	ifMatch := c.GetHeader("If-Match")
	// Use If-Match header for conditional writes
	if ifMatch != "" && ifMatch != currentETag {
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Resource has been modified"})
		return
	}

	// Parse the incoming JSON for the update
	var updatedData models.Plan
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Merge updates with the existing plan data
	if updatedData.PlanType != "" {
		existingPlan.PlanType = updatedData.PlanType
	}
	if updatedData.CreationDate != "" {
		existingPlan.CreationDate = updatedData.CreationDate
	}
	if updatedData.Org != "" {
		existingPlan.Org = updatedData.Org
	}
	if len(updatedData.LinkedPlanServices) > 0 {
		existingPlan.LinkedPlanServices = updatedData.LinkedPlanServices
	}
	if updatedData.PlanCostShares != (models.PlanCostShares{}) {
		existingPlan.PlanCostShares = updatedData.PlanCostShares
	}

	// Generate a new ETag for the updated plan
	newETag := generateETag(existingPlan)

	// Serialize the merged plan and store it back in Redis
	mergedData, err := json.Marshal(existingPlan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated plan"})
		return
	}
	err = service.redisClient.Set(ctx, "plan:"+planID, mergedData, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store updated plan"})
		return
	}

	// Set the new ETag in the response header
	c.Header("ETag", newETag)
	c.JSON(http.StatusOK, gin.H{"message": "Plan updated", "planId": planID})
}
