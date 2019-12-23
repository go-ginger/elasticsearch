package elasticsearch

import (
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-ginger/models"
)

type Index struct {
	Name    string
	Mapping interface{}
	Setting interface{}
}

type Config struct {
	models.IConfig

	ElasticConfig elastic.Config
	Indexes       []Index
}

var config Config

func InitializeConfig(input Config) {
	config = input
	Initialize()
}
