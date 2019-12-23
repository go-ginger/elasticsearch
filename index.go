package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strings"
)

func indexExists(index ...string) (exists bool, err error) {
	ctx := context.Background()
	var r = esapi.IndicesExistsRequest{
		Index: index,
	}
	resp, err := r.Do(ctx, db.Client.Transport)
	if err != nil {
		return
	}
	exists = resp.StatusCode == 200
	return
}

func ensureIndexExists(index, body string) (err error) {
	exists, err := indexExists(index)
	if exists {
		return
	}
	resp, err := db.Client.Indices.Create(
		index,
		db.Client.Indices.Create.WithContext(context.Background()),
		db.Client.Indices.Create.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Could not create elastic index. "+
			"Response returned status code %d", resp.StatusCode))
		return
	}
	return
}
