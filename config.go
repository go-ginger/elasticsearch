package elasticsearch

import (
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-ginger/helpers/namer"
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
	IndexNamer    namer.INamer
}
