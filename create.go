package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/go-ginger/models"
	"strings"
)

func (handler *DbHandler) Insert(request models.IRequest) (result interface{}, err error) {
	req := request.GetBaseRequest()
	marshalledBody, err := json.Marshal(req.Body)
	if err != nil {
		return
	}
	body := string(marshalledBody)
	indexName := handler.DB.Config.IndexNamer.GetName(req.Body)
	indexReq := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%v", req.Body.GetIDString()),
		Body:       strings.NewReader(body),
		Refresh:    "true",
	}
	res, err := indexReq.Do(context.Background(), handler.DB.Client)
	if err != nil {
		return
	}
	defer func() {
		e := res.Body.Close()
		if e != nil {
			err = e
		}
	}()
	_, err = handler.BaseDbHandler.Insert(req)
	return
}
