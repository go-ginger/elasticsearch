package elasticsearch

import (
	"github.com/go-ginger/models"
)

type Index struct {
	Name    string
	Mapping interface{}
	Setting interface{}
}

type Config struct {
	models.IConfig

	Indexes []Index
}

var config Config

func InitializeConfig(input interface{}) {
	config = Config{
	}
}
