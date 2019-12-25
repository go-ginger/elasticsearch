package elasticsearch

import (
	"github.com/go-ginger/models"
)

func (handler *DbHandler) Update(request models.IRequest) (err error) {
	_, err = handler.Insert(request)
	return
}
