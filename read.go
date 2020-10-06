package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/go-ginger/models"
	ge "github.com/go-ginger/models/errors"
	"math"
	"reflect"
)

func (handler *DbHandler) Paginate(request models.IRequest) (result *models.PaginateResult, err error) {
	req := request.GetBaseRequest()
	model := handler.GetModelsInstance()
	indexName := handler.DB.Config.IndexNamer.GetName(model)
	parseResult := handler.QueryParser.Parse(request)
	queryBytes, err := json.Marshal(parseResult.GetQuery())
	if err != nil {
		return
	}
	fmt.Println(string(queryBytes))
	queryReader := bytes.NewReader(queryBytes)
	offset := int64((req.Page - 1) * req.PerPage)
	limit := int64(req.PerPage)
	searchReqs := []func(*esapi.SearchRequest){
		handler.DB.Client.Search.WithContext(context.Background()),
		handler.DB.Client.Search.WithIndex(indexName),
		handler.DB.Client.Search.WithBody(queryReader),
		handler.DB.Client.Search.WithFrom(int(offset)),
		handler.DB.Client.Search.WithSize(int(limit)),
		handler.DB.Client.Search.WithTrackTotalHits(true),
	}
	sortResult := parseResult.GetSort()
	if sortResult != nil {
		searchReqs = append(searchReqs, handler.DB.Client.Search.WithSort(parseResult.GetSort().([]string)...))
	}
	resp, err := handler.DB.Client.Search(searchReqs...)
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			err = e
		}
	}()
	if resp.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return
		} else {
			err = ge.GetInternalServiceError(request, fmt.Sprintf("[%s] %s: %s",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			))
			return
		}
	}
	//{
	//	"took" : 4,
	//	"timed_out" : false,
	//	"_shards" : {
	//	"total" : 1,
	//		"successful" : 1,
	//		"skipped" : 0,
	//		"failed" : 0
	//},
	//	"hits" : {
	//		"total" : {
	//			"value" : 0,
	//				"relation" : "eq"
	//		},
	//		"max_score" : null,
	//	    "hits" : [ ]
	//	}
	//}
	apiResult := struct {
		Hits struct {
			Total struct {
				Value uint64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return
	}
	items := handler.GetModelsInstancePtr()
	if apiResult.Hits.Hits != nil {
		hits := make([]interface{}, 0)
		for _, hit := range apiResult.Hits.Hits {
			hits = append(hits, hit.Source)
		}
		m, e := json.Marshal(hits)
		if e != nil {
			err = e
			return
		}
		err = json.Unmarshal(m, items)
		if err != nil {
			return
		}
		items = reflect.ValueOf(items).Elem().Interface()
	}
	pageCount := uint64(math.Ceil(float64(apiResult.Hits.Total.Value) / float64(req.PerPage)))
	result = &models.PaginateResult{
		Items: items,
		Pagination: models.PaginationInfo{
			Page:       req.Page,
			PerPage:    req.PerPage,
			PageCount:  pageCount,
			TotalCount: apiResult.Hits.Total.Value,
			HasNext:    req.Page < pageCount,
		},
	}
	return
}

func (handler *DbHandler) Get(request models.IRequest) (result models.IBaseModel, err error) {
	return
}
