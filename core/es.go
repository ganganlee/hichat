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
func (e *Elasticsearch) Search(index string, keywords string, fromId string, toId string, page int, pageSize int) (total int64, list []*elastic.SearchHit, err error) {

	var (
		res  *elastic.SearchResult
		from int
	)

	from = (page - 1) * pageSize
	//控制搜索条件
	fromIdQuery := elastic.NewMatchQuery("from_id", fromId)
	toIdQuery := elastic.NewMatchQuery("to_id", toId)
	keywordsQuery := elastic.NewMatchQuery("content", keywords)
	b := elastic.NewBoolQuery()
	b.Must(fromIdQuery, toIdQuery, keywordsQuery)
	q := elastic.NewQueryRescorer(b)

	//控制高亮
	h := elastic.NewHighlight()
	h.Fields(elastic.NewHighlighterField("content"))

	res, err = e.client.Search().Index(index).Query(q).Highlight(h).Source("content").Size(pageSize).From(from).Do(context.TODO())
	if err != nil {
		return 0, nil, err
	}

	//总条数
	total = res.Hits.TotalHits.Value
	list = res.Hits.Hits

	return total, list, err
}
