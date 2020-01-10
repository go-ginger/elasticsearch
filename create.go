package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/go-ginger/models"
	"log"
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
	var indexReq *esapi.IndexRequest
	iRequestParam := request.GetTemp("esapi_request")
	if iRequestParam != nil {
		indexReq, _ = iRequestParam.(*esapi.IndexRequest)
	}
	if indexReq == nil {
		indexReq = &esapi.IndexRequest{}
	}
	indexReq.Index = indexName
	indexReq.DocumentID = fmt.Sprintf("%v", req.GetIDString())
	indexReq.Body = strings.NewReader(body)
	indexReq.Refresh = "true"

	resp, err := indexReq.Do(context.Background(), handler.DB.Client)
	if err != nil {
		return
	}
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			err = e
		}
	}()
	if resp.IsError() {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return
		}
		respBody := buf.String()
		err = errors.New(fmt.Sprintf("Insert of elasticsearch returned status code %d", resp.StatusCode))
		log.Println(fmt.Sprintf("response body: %v", respBody))
		return
	}
	_, err = handler.BaseDbHandler.Insert(req)
	return
}
