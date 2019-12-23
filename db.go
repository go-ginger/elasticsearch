package elasticsearch

import (
	elastic "github.com/elastic/go-elasticsearch/v8"
	"log"
)

type DB struct {
	Client *elastic.Client
}

var db *DB

func connect() {
	client, err := elastic.NewClient(config.ElasticConfig)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	db = &DB{
		Client: client,
	}

	res, err := db.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
}
