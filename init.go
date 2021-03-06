package elasticsearch

import (
	"github.com/go-ginger/helpers/namer"
	"log"
)

func (handler *DbHandler) InitializeConfig(config *Config) {
	if config.IndexNamer == nil {
		config.IndexNamer = &namer.Default{}
	}
	config.IndexNamer.Initialize()
	handler.DB = &DB{
		Config: config,
	}
	handler.DB.Connect()
	for _, index := range handler.DB.Config.Indexes {
		data := `{
   "settings":
`
		if setting, ok := index.Setting.(string); ok {
			data += setting
		}
		data += `,
	"mappings":
`
		if mapping, ok := index.Mapping.(string); ok {
			data += mapping
		}
		data += `}`
		err := handler.ensureIndexExists(config.IndexNamer.GetName(index.Model), data)
		if err != nil {
			log.Fatalf("Error on Initialize, err: %v", err)
		}
	}
	handler.InsertInBackground = true
	handler.UpdateInBackground = true
	handler.DeleteInBackground = true
	handler.IsFullObjectOnUpdateRequired = true
	return
}
