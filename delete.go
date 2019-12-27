package elasticsearch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/go-ginger/models"
	"log"
)

func (handler *DbHandler) Delete(request models.IRequest) (err error) {
	req := request.GetBaseRequest()
	indexName := handler.DB.Config.IndexNamer.GetNameByType(handler.ModelType)
	resp, err := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%v", req.GetIDString()),
		Refresh:    "true",
	}.Do(context.Background(), handler.DB.Client)
	if err != nil {
		return
	}
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			err = e
		}
	}()
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return
		}
		respBody := buf.String()
		err = errors.New(fmt.Sprintf("Delete of elasticsearch returned status code %d", resp.StatusCode))
		log.Println(fmt.Sprintf("response body: %v", respBody))
		return
	}
	err = handler.BaseDbHandler.Delete(req)
	return
}
