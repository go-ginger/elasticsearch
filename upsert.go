package elasticsearch

import (
	"github.com/go-ginger/models"
)

func (handler *DbHandler) Upsert(request models.IRequest) (err error) {
	_, err = handler.Insert(request)
	return
}
