package validators

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

const schemaPath = "./internal/schemas/plan_schema.json"

// Validate incoming JSON against schema
func ValidatePlanSchema(c *gin.Context) bool {
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)

	// Read the raw request body for validation
	body, err := ioutil.ReadAll(c.Request.Body)
	// fmt.Printf("Here is the payload Data : %v", string(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return false
	}

	// Load the request body into the validator
	documentLoader := gojsonschema.NewStringLoader(string(body))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Schema validation failed", "details": err.Error()})
		return false
	}

	// Print specific validation errors
	if !result.Valid() {
		var validationErrors []string
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.String())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": validationErrors})
		return false
	}

	// Reset the request body for further processing
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return true
}
