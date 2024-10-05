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

func generateETag(plan models.Plan) string {
	planData, _ := json.Marshal(plan)
	hash := sha1.New()
	hash.Write(planData)
	return hex.EncodeToString(hash.Sum(nil))
}

// Create a new Plan
// Create a new Plan
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
		// If Redis returns an error other than "plan not found", return an internal server error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing data"})
			return
		}

		// If the plan exists, return a conflict error
		if existingPlan != "" {
			c.JSON(http.StatusConflict, gin.H{"error": "Plan already exists"})
			return
		}
	}

	// If the plan doesn't exist, proceed to create it
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

	// Generate ETag for the plan
	eTag := generateETag(plan)

	// Check for If-None-Match header (conditional GET)
	ifMatch := c.GetHeader("If-None-Match")
	if ifMatch == eTag {
		c.Status(http.StatusNotModified)
		return
	}

	// Add ETag header to the response
	c.Header("ETag", eTag)
	c.JSON(http.StatusOK, plan)
}

// Delete a Plan
// Delete a Plan
func (service *PlanService) DeletePlan(c *gin.Context) {
	planID := c.Query("id")

	// Check if the plan exists in Redis before attempting to delete it
	exists, err := service.redisClient.Exists(ctx, "plan:"+planID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check plan existence"})
		return
	}

	// If the plan does not exist, return "Data not found"
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Proceed to delete the plan if it exists
	err = service.redisClient.Del(ctx, "plan:"+planID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
}
