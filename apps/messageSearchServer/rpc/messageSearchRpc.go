package rpc

import (
	"context"
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
	)

	//获取传入参数
	params = &services.SearchRequest{
		FromId:   res.FromId,
		ToId:     res.ToId,
		Keywords: res.Keywords,
		PageSize: res.PageSize,
		Page:     res.Page,
	}

	fmt.Println(params)

	//验证参数
	if err = validate.Struct(params); err != nil {
		return err
	}

	//调用服务方法进行搜索
	if err = m.service.Search(params);err != nil {
		return err
	}


	rsp.Total = 100
	return nil
}
