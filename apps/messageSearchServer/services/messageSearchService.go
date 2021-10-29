package services

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"hichat.zozoo.net/apps/messageSearchServer/common"
	"hichat.zozoo.net/core"
)

type (
	MessageSearchService struct {
	}

	//搜索请求
	SearchRequest struct {
		FromId   string `json:"from_id" validate:"required"`
		ToId     string `json:"to_id" validate:"required"`
		Keywords string `json:"keywords" validate:"required"`
		Page     uint32 `json:"page"`
		PageSize uint32 `json:"page_size"`
	}
)

func NewMessageSearchService() *MessageSearchService {
	return &MessageSearchService{
	}
}

//搜索服务
func (m *MessageSearchService) Search(res *SearchRequest) (err error) {
	var (
		index = common.AppCfg.Es.Index
		total int64
		list  []*elastic.SearchHit
	)

	fmt.Println(res)

	total, list, err = core.Es.Search(index, res.Keywords, res.FromId, res.ToId, int(res.Page), int(res.PageSize))
	if err != nil {
		return err
	}

	fmt.Println(list)
	fmt.Println(total)
	return err
}
