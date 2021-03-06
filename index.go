package elasticsearch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
)

func (handler *DbHandler) indexExists(index ...string) (exists bool, err error) {
	ctx := context.Background()
	var r = esapi.IndicesExistsRequest{
		Index: index,
	}
	resp, err := r.Do(ctx, handler.DB.Client.Transport)
	if err != nil {
		return
	}
	exists = resp.StatusCode == 200
	return
}

func (handler *DbHandler) ensureIndexExists(index, body string) (err error) {
	exists, err := handler.indexExists(index)
	if exists {
		return
	}
	resp, err := handler.DB.Client.Indices.Create(
		index,
		handler.DB.Client.Indices.Create.WithContext(context.Background()),
		handler.DB.Client.Indices.Create.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return
		}
		respBody := buf.String()
		err = errors.New(fmt.Sprintf("Could not create elastic index. "+
			"Response returned status code %d", resp.StatusCode))
		log.Println(fmt.Sprintf("response body: %v", respBody))
		return
	}
	return
}
