package rpc

import (
	"context"
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
	return nil
}
