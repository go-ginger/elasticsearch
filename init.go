package elasticsearch

import "log"

func Initialize() {
	connect()
	for _, index := range config.Indexes {
		data := `{
   "settings":
`
		if setting, ok := index.Setting.(string); ok {
			data += setting
		}
		data += `},
	"mappings":
`
		if mapping, ok := index.Mapping.(string); ok {
			data += mapping
		}
		data += `}`
		err := ensureIndexExists(index.Name, data)
		if err != nil {
			log.Fatalf("Error on Initialize, err: %v", err)
		}
	}
}
