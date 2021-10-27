package services

type (
	MessageSearchService struct {
	}
)

func NewMessageSearchService() *MessageSearchService {
	return &MessageSearchService{
	}
}

//搜索服务
func (m *MessageSearchService) Search() (err error) {
	return err
}
