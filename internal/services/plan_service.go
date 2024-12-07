package services

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"BigDataForge/internal/elastic"
	"BigDataForge/internal/models"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type PlanService struct {
	redisClient *redis.Client
	esClient    *elastic.Factory
}

func NewPlanService(redisClient *redis.Client, esFactory *elastic.Factory) *PlanService {
	return &PlanService{
		redisClient: redisClient,
		esClient:    esFactory,
	}
}

// Helper function to generate ETag based on Plan data
func generateETag(plan models.Plan) string {
	planData, _ := json.Marshal(plan)
	hash := sha1.New()
	hash.Write(planData)
	return hex.EncodeToString(hash.Sum(nil))
}

// Helper function to check if a plan exists in Redis
func (service *PlanService) getPlanFromRedis(planID string) (*models.Plan, error) {
	planJSON, err := service.redisClient.Get(ctx, "plan:"+planID).Result()
	if err == redis.Nil {
		return nil, nil // Plan does not exist
	}
	if err != nil {
		return nil, err
	}

	var plan models.Plan
	err = json.Unmarshal([]byte(planJSON), &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// CreatePlan handles creating a new plan
func (service *PlanService) CreatePlan(c *gin.Context) {
	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	planID := plan.ObjectID

	// Check if the plan already exists
	existingPlan, err := service.getPlanFromRedis(planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing data"})
		return
	}
	if existingPlan != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Plan already exists"})
		return
	}

	// Save the new plan in Redis
	if err := service.savePlanToRedis(planID, plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store plan"})
		return
	}

	// Generate and set the ETag
	c.Header("ETag", generateETag(plan))
	c.JSON(http.StatusCreated, gin.H{"message": "Plan created", "planId": planID})
}

// GetPlan retrieves a plan by ID
func (service *PlanService) GetPlan(c *gin.Context) {
	planID := c.Query("id")

	plan, err := service.getPlanFromRedis(planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}
	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Handle conditional read with ETag
	eTag := generateETag(*plan)
	if c.GetHeader("If-None-Match") == eTag {
		c.Status(http.StatusNotModified)
		return
	}

	// Return the plan with ETag
	c.Header("ETag", eTag)
	c.JSON(http.StatusOK, plan)
}

// DeletePlan removes a plan by ID
func (service *PlanService) DeletePlan(c *gin.Context) {
	planID := c.Query("id")

	// Check if the plan exists
	exists, err := service.redisClient.Exists(ctx, "plan:"+planID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check plan existence"})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Delete the plan
	if err := service.redisClient.Del(ctx, "plan:"+planID).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}

	c.Status(http.StatusNoContent)
}

// PatchPlan updates specific fields of a plan
func (service *PlanService) PatchPlan(c *gin.Context) {
	planID := c.Query("id")

	// Get the existing plan
	existingPlan, err := service.getPlanFromRedis(planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}
	if existingPlan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Check If-Match header for conditional writes
	ifMatch := c.GetHeader("If-Match")
	if ifMatch != "" && ifMatch != generateETag(*existingPlan) {
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "Resource has been modified"})
		return
	}

	// Update plan fields
	var updatedData models.Plan
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	service.mergePlan(existingPlan, &updatedData)

	// Save the updated plan
	if err := service.savePlanToRedis(planID, *existingPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan updated", "planId": planID})
}

// Helper to merge plans
func (service *PlanService) mergePlan(existingPlan, updatedData *models.Plan) {
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
	if !isPlanCostSharesEmpty(updatedData.PlanCostShares) {
		existingPlan.PlanCostShares = updatedData.PlanCostShares
	}
}

// Helper to check if PlanCostShares is empty
func isPlanCostSharesEmpty(planCostShares models.PlanCostShares) bool {
	return planCostShares.Copay == 0 && planCostShares.Deductible == 0 &&
		len(planCostShares.PlanJoin) == 0 && planCostShares.Org == "" &&
		planCostShares.ObjectID == "" && planCostShares.ObjectType == ""
}

// Save plan to Redis
func (service *PlanService) savePlanToRedis(planID string, plan models.Plan) error {
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	return service.redisClient.Set(ctx, "plan:"+planID, planJSON, 0).Err()
}

// SearchPlans performs a search query in Elasticsearch
func (service *PlanService) SearchPlans(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search query"})
		return
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				req.Key: req.Value,
			},
		},
	}
	fmt.Println(query)
	queryBytes, err := json.Marshal(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build query"})
		return
	}

	client, err := service.esClient.NewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Elasticsearch client"})
		return
	}

	searchReq := esapi.SearchRequest{
		Index: []string{"plans"},
		Body:  bytes.NewReader(queryBytes),
	}

	res, err := searchReq.Do(ctx, client)
	if err != nil || res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute search"})
		return
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse search results"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdatePlan updates a plan by deleting the existing one and recreating it
func (service *PlanService) UpdatePlan(c *gin.Context) {
	var updatedPlan models.Plan
	if err := c.ShouldBindJSON(&updatedPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the plan exists
	existingPlan, err := service.getPlanFromRedis(updatedPlan.ObjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}
	if existingPlan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	// Validate request schema
	if !validatePlanSchema(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	// Delete the existing plan
	if err := service.DeletePlanByID(updatedPlan.ObjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing plan"})
		return
	}

	// Create the new plan
	if err := service.savePlanToRedis(updatedPlan.ObjectID, updatedPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create updated plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan updated successfully"})
}

// DeletePlanByID deletes a plan by ID without using gin.Context
func (service *PlanService) DeletePlanByID(planID string) error {
	err := service.redisClient.Del(ctx, "plan:"+planID).Err()
	if err != nil {
		log.Printf("Failed to delete plan %s: %v", planID, err)
	}
	return err
}

// ValidatePlanSchema validates the request schema
func validatePlanSchema(c *gin.Context) bool {
	var plan models.Plan
	if err := c.ShouldBindJSON(&plan); err != nil {
		return false
	}
	return true
}
