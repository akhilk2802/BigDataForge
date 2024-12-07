package elastic

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

type Factory struct{}

func (f *Factory) NewClient() (*elasticsearch.Client, error) {
	elasticURL := os.Getenv("ELASTICSEARCH_URL")
	if elasticURL == "" {
		elasticURL = "http://localhost:9200"
	}
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticURL,
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to create Elasticsearch client: %s", err)
		return nil, err
	}
	return client, nil
}
