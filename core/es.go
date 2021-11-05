package core

//es操作公共方法
import (
	"context"
	"github.com/olivere/elastic/v7"
)

type Elasticsearch struct {
	client *elastic.Client
}

var Es *Elasticsearch

func NewEs(host, username, password string) (es *Elasticsearch, err error) {
	var client *elastic.Client

	//链接es服务器
	client, err = elastic.NewClient(elastic.SetURL(host), elastic.SetBasicAuth(username, password), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	Es = &Elasticsearch{
		client: client,
	}

	return Es, err
}

//创建索引
func (e *Elasticsearch) CreateIndex(index string, mapping string) (err error) {
	_, err = e.client.CreateIndex(index).BodyString(mapping).Do(context.TODO())
	return err
}

//添加文档
func (e *Elasticsearch) AddDoc(index, id, content string) (err error) {
	_, err = e.client.Index().Index(index).Id(id).BodyString(content).Do(context.TODO())
	return err
}

//删除文档
func (e *Elasticsearch) DelDoc(index, id string) (err error) {
	_, err = e.client.Delete().Index(index).Id(id).Do(context.TODO())
	return err
}

//修改文档
func (e *Elasticsearch) Update(index, id string, content map[string]interface{}) (err error) {
	_, err = e.client.Update().Index(index).Id(id).Doc(content).Do(context.TODO())
	return err
}

//搜索文档
func (e *Elasticsearch) Search(index string, q *elastic.BoolQuery, pageSize int, from int, h *elastic.Highlight) (res *elastic.SearchResult, err error) {

	res, err = e.client.Search().Index(index).Query(q).Highlight(h).Size(pageSize).From(from).Do(context.TODO())
	if err != nil {
		return nil, err
	}

	return res, nil
}
