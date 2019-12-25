package elasticsearch

import (
	elastic "github.com/elastic/go-elasticsearch/v8"
	"log"
)

type DB struct {
	Client *elastic.Client
	Config *Config
}

func (db *DB) Connect() {
	client, err := elastic.NewClient(db.Config.ElasticConfig)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	db.Client = client
	res, err := db.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	return
}
