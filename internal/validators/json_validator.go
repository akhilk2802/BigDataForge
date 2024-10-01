package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

const planSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "planCostShares": {
      "type": "object",
      "properties": {
        "deductible": { "type": "integer" },
        "_org": { "type": "string" },
        "copay": { "type": "integer" },
        "objectId": { "type": "string" },
        "objectType": { "type": "string" }
      },
      "required": ["deductible", "_org", "copay", "objectId", "objectType"]
    },
    "linkedPlanServices": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "linkedService": {
            "type": "object",
            "properties": {
              "_org": { "type": "string" },
              "objectId": { "type": "string" },
              "objectType": { "type": "string" },
              "name": { "type": "string" }
            },
            "required": ["_org", "objectId", "objectType", "name"]
          },
          "planserviceCostShares": {
            "type": "object",
            "properties": {
              "deductible": { "type": "integer" },
              "_org": { "type": "string" },
              "copay": { "type": "integer" },
              "objectId": { "type": "string" },
              "objectType": { "type": "string" }
            },
            "required": ["deductible", "_org", "copay", "objectId", "objectType"]
          },
          "_org": { "type": "string" },
          "objectId": { "type": "string" },
          "objectType": { "type": "string" }
        },
        "required": ["linkedService", "planserviceCostShares", "_org", "objectId", "objectType"]
      }
    },
    "_org": { "type": "string" },
    "objectId": { "type": "string" },
    "objectType": { "type": "string" },
    "planType": { "type": "string" },
    "creationDate": { "type": "string" }
  },
  "required": ["planCostShares", "linkedPlanServices", "_org", "objectId", "objectType", "planType", "creationDate"]
}
`

func ValidatePlanSchema(c *gin.Context) bool {
	var jsonData interface{}

	// This binds the request body into jsonData and preserves the body for further use
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return false
	}
	// Convert jsonData to string for schema validation
	documentLoader := gojsonschema.NewGoLoader(jsonData)
	schemaLoader := gojsonschema.NewStringLoader(planSchema)

	// Validate the JSON against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil || !result.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format", "details": result.Errors()})
		return false
	}

	return true
}
