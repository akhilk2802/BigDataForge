package main

import (
	"BigDataForge/internal/elastic"
	"BigDataForge/internal/models"
	"BigDataForge/internal/rabbitmq"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	rabbitFactory := rabbitmq.Factory{}
	elasticFactory := elastic.Factory{}

	// Connect to RabbitMQ
	conn, err := rabbitFactory.NewConnection()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a new RabbitMQ channel
	ch, err := rabbitFactory.NewChannel(conn)
	failOnError(err, "Failed to create RabbitMQ channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"plan_queue",
		false,
		true,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"myConsumer",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	// Connect to Elasticsearch
	esClient, err := elasticFactory.NewClient()
	failOnError(err, "Failed to create Elasticsearch client")

	// Ensure the Elasticsearch index and mapping exist
	ensureIndex(esClient)

	// Start message processing
	go processMessages(msgs, esClient)

	log.Println("Listening for messages. Press CTRL+C to exit.")
	select {}

}

func processMessages(msgs <-chan amqp.Delivery, esClient *elasticsearch.Client) {
	for d := range msgs {
		log.Printf("Received message: %s", d.Body)

		var plan models.Plan
		if err := json.Unmarshal(d.Body, &plan); err != nil {
			log.Printf("Failed to deserialize Plan: %s", err)
			continue
		}

		// Index the plan and its related documents
		if err := indexPlan(esClient, plan); err != nil {
			log.Printf("Failed to index Plan: %s", err)
		}
	}
}

func indexPlan(esClient *elasticsearch.Client, plan models.Plan) error {
	// Index the main plan
	plan.PlanJoin = map[string]interface{}{"name": "plan"}
	if err := indexDocument(esClient, "plans", plan.ObjectID, plan, ""); err != nil {
		return err
	}

	// Index PlanCostShares
	plan.PlanCostShares.PlanJoin = map[string]interface{}{"name": "planCostShares", "parent": plan.ObjectID}
	if err := indexDocument(esClient, "plans", plan.PlanCostShares.ObjectID, plan.PlanCostShares, plan.ObjectID); err != nil {
		return err
	}

	// Index LinkedPlanServices and related documents
	for _, linkedPlanService := range plan.LinkedPlanServices {
		linkedPlanService.PlanJoin = map[string]interface{}{"name": "linkedPlanServices", "parent": plan.ObjectID}
		if err := indexDocument(esClient, "plans", linkedPlanService.ObjectID, linkedPlanService, plan.ObjectID); err != nil {
			return err
		}

		// Index LinkedService
		linkedPlanService.LinkedService.PlanJoin = map[string]interface{}{
			"name":   "linkedService",
			"parent": linkedPlanService.ObjectID,
		}
		if err := indexDocument(esClient, "plans", linkedPlanService.LinkedService.ObjectID, linkedPlanService.LinkedService, linkedPlanService.ObjectID); err != nil {
			return err
		}

		// Index PlanServiceCostShares
		linkedPlanService.PlanserviceCostShares.PlanJoin = map[string]interface{}{
			"name":   "planServiceCostShares",
			"parent": linkedPlanService.ObjectID,
		}
		if err := indexDocument(esClient, "plans", linkedPlanService.PlanserviceCostShares.ObjectID, linkedPlanService.PlanserviceCostShares, linkedPlanService.ObjectID); err != nil {
			return err
		}
	}
	return nil
}

func indexDocument(esClient *elasticsearch.Client, indexName, documentID string, document interface{}, routing string) error {
	docJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
		Routing:    routing,
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing document ID=%s: %s", documentID, res.String())
		return fmt.Errorf("failed to index document ID=%s", documentID)
	}

	log.Printf("Document ID=%s indexed successfully", documentID)
	return nil
}

func ensureIndex(esClient *elasticsearch.Client) {
	// Create the index if it doesn't exist

	fmt.Println("Creating index")
	req := esapi.IndicesCreateRequest{Index: "plans"}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Fatalf("Failed to create index: %s", err)
	}
	defer res.Body.Close()

	fmt.Println("Index created")

	if res.IsError() {
		log.Printf("Index creation response: %s", res.String())
	}

	fmt.Println("Applying mapping")

	// Apply the mapping
	mapping := getMapping()

	fmt.Println("Mapping created")

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		log.Fatalf("Failed to serialize mapping: %s", err)
	}

	fmt.Println("Mapping serialized")

	req2 := esapi.IndicesPutMappingRequest{
		Index: []string{"plans"},
		Body:  bytes.NewReader(mappingJSON),
	}
	res2, err := req2.Do(context.Background(), esClient)
	if err != nil {
		log.Fatalf("Failed to apply mapping: %s", err)
	}
	defer res2.Body.Close()

	if res2.IsError() {
		log.Printf("Mapping application response: %s", res2.String())
	} else {
		log.Println("Mapping applied successfully")
	}
}
func getMapping() map[string]interface{} {
	return map[string]interface{}{
		"properties": map[string]interface{}{
			"plan": map[string]interface{}{
				"properties": map[string]interface{}{
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
					"planType": map[string]interface{}{
						"type": "text",
					},
					"creationDate": map[string]interface{}{
						"type":   "date",
						"format": "MM-dd-yyyy",
					},
				},
			},
			"planCostShares": map[string]interface{}{
				"properties": map[string]interface{}{
					"copay": map[string]interface{}{
						"type": "long",
					},
					"deductible": map[string]interface{}{
						"type": "long",
					},
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"linkedPlanServices": map[string]interface{}{
				"properties": map[string]interface{}{
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"linkedService": map[string]interface{}{
				"properties": map[string]interface{}{
					"_org": map[string]interface{}{
						"type": "text",
					},
					"name": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"planserviceCostShares": map[string]interface{}{
				"properties": map[string]interface{}{
					"copay": map[string]interface{}{
						"type": "long",
					},
					"deductible": map[string]interface{}{
						"type": "long",
					},
					"_org": map[string]interface{}{
						"type": "text",
					},
					"objectId": map[string]interface{}{
						"type": "keyword",
					},
					"objectType": map[string]interface{}{
						"type": "text",
					},
				},
			},
			"plan_join": map[string]interface{}{
				"type":                  "join",
				"eager_global_ordinals": "true",
				"relations": map[string]interface{}{
					"plan":               []string{"planCostShares", "linkedPlanServices"},
					"linkedPlanServices": []string{"linkedService", "planserviceCostShares"},
				},
			},
		},
	}
}
