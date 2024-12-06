package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type Factory struct{}

func (f *Factory) NewClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to create Elasticsearch client: %s", err)
		return nil, err
	}
	return client, nil
}
