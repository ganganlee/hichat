package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"hichat.zozoo.net/apps/messageSearchServer/services"
	"hichat.zozoo.net/rpc/messageSearch"
)

type (
	MessageSearchRpc struct {
		service *services.MessageSearchService
	}
)

func NewMessageSearchRpc(s *services.MessageSearchService) *MessageSearchRpc {
	return &MessageSearchRpc{
		service: s,
	}
}

//搜索服务
func (m *MessageSearchRpc) Search(ctx context.Context, res *messageSearch.SearchRequest, rsp *messageSearch.SearchResponse) error {
	var (
		params   *services.SearchRequest
		validate = validator.New()
		err      error
		total    int64
		data     []*services.SearchResponse
		b        []byte
	)

	//获取传入参数
	params = &services.SearchRequest{
		FromId:   res.FromId,
		ToId:     res.ToId,
		Keywords: res.Keywords,
		IsGroup: res.IsGroup,
		PageSize: res.PageSize,
		Page:     res.Page,
	}

	//验证参数
	if err = validate.Struct(params); err != nil {
		return err
	}

	//调用服务方法进行搜索
	if total, data, err = m.service.Search(params); err != nil {
		return err
	}

	rsp.Total = total
	if b, err = json.Marshal(data); err != nil {
		fmt.Println(err)
		return err
	}

	rsp.Result = string(b)
	rsp.Msg = "ok"
	return nil
}
