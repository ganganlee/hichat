package services

import (
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"hichat.zozoo.net/apps/messageSearchServer/common"
	"hichat.zozoo.net/core"
	"strings"
	"time"
)

type (
	MessageSearchService struct {
	}

	//搜索请求
	SearchRequest struct {
		FromId   string `json:"from_id" validate:"required"`
		ToId     string `json:"to_id" validate:"required"`
		Keywords string `json:"keywords" validate:"required"`
		IsGroup  bool   `json:"is_group"`
		Page     uint32 `json:"page"`
		PageSize uint32 `json:"page_size"`
	}

	SearchResponse struct {
		FromId      string    `json:"from_id"`
		MsgType     string    `json:"msg_type"`
		ContentType string    `json:"content_type"`
		CreateTime  time.Time `json:"create_time"`
		Content     string    `json:"content"`
	}
)

func NewMessageSearchService() *MessageSearchService {
	return &MessageSearchService{
	}
}

//搜索服务
func (m *MessageSearchService) Search(res *SearchRequest) (total int64, list []*SearchResponse, err error) {
	var (
		index  = common.AppCfg.Es.Index
		search *elastic.SearchResult
	)

	//控制搜索条件
	q := elastic.NewBoolQuery()
	keywordsQuery := elastic.NewMatchQuery("content", res.Keywords)
	q.Must(keywordsQuery)

	//判断搜索内容，如果为群聊时搜索只需要判断to_id,如果为私聊时需要根据双方id查找
	if res.IsGroup == false {
		//私聊
		q1 := elastic.NewBoolQuery()
		q1.Must(
			elastic.NewMatchQuery("from_id", res.FromId),
			elastic.NewMatchQuery("to_id", res.ToId),
		)
		q2 := elastic.NewBoolQuery()
		q2.Must(
			elastic.NewMatchQuery("from_id", res.ToId),
			elastic.NewMatchQuery("to_id", res.FromId),
		)

		q.Should(q1, q2)
	} else {
		//群聊
		q1 := elastic.NewBoolQuery()
		q1.Must(
			elastic.NewMatchQuery("to_id", res.ToId),
		)
		q.Must(q1)
	}

	//分页控制
	from := (int(res.Page) - 1) * int(res.PageSize)

	//控制高亮
	h := elastic.NewHighlight()
	h.Fields(elastic.NewHighlighterField("content"))

	search, err = core.Es.Search(index, q, int(res.PageSize), from, h)
	if err != nil {
		return 0, nil, err
	}

	//组织返回数据
	list = make([]*SearchResponse, 0)

	for _, val := range search.Hits.Hits {
		var item = new(SearchResponse)

		if err = json.Unmarshal(val.Source, item); err != nil {
			continue
		}

		item.Content = strings.Join(val.Highlight["content"], "")
		list = append(list, item)
	}

	return search.Hits.TotalHits.Value, list, err
}
