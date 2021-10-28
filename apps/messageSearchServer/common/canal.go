package common

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"hichat.zozoo.net/core"
	"log"
	"os"
	"strconv"
	"time"
)

//canal-go服务

type Message struct {
	Id          int       `json:"id"`
	FromId      string    `json:"from_id"`
	ToId        string    `json:"to_id"`
	MsgType     string    `json:"msg_type"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content"`
	CreateTime  time.Time `json:"create_time"`
}

//连接canal
func InitCanal(address string, port int, username string, password string, destination string, soTimeOut int32, idleTimeOut int32) {

	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	connector := client.NewSimpleCanalConnector(address, port, username, password, destination, soTimeOut, idleTimeOut)
	err := connector.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//监听全库全表 ".*\\..*"
	//监听单库全表 "databaseName\\..*"
	//监听单库单表 "databaseName.user"
	//多规则使用 "databaseName\\..*,databaseName2.user1,databaseNam3.user"
	err = connector.Subscribe("hichat.message\\..*")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		printEntry(message.Entries)
	}
}

func printEntry(entrys []pbe.Entry) {

	var err error
	for _, entry := range entrys {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)

		if err = proto.Unmarshal(entry.GetStoreValue(), rowChange); err != nil {
			fmt.Printf("Fatal error: %s\n", err.Error())
			os.Exit(1)
		}

		eventType := rowChange.GetEventType()

		for _, rowData := range rowChange.GetRowDatas() {
			if eventType == pbe.EventType_DELETE {
				//删除操作
				val := printColumn(rowData.GetBeforeColumns())
				if err = core.Es.DelDoc(AppCfg.Es.Index, val["id"]); err != nil {
					fmt.Println("删除文档失败失败 err:", err.Error())
					continue
				}

			} else if eventType == pbe.EventType_INSERT {
				//新增操作
				val := printColumn(rowData.GetAfterColumns())
				message := new(Message)
				message.FromId = val["from_id"]
				message.ToId = val["to_id"]
				message.MsgType = val["msg_type"]
				message.ContentType = val["content_type"]
				message.Content = val["content"]
				message.Id, _ = strconv.Atoi(val["id"])

				message.CreateTime, err = time.Parse("2006-01-02 15:04:05", val["create_time"])
				if err != nil {
					fmt.Println("时间转换失败 err:", err.Error())
					continue
				}

				b, _ := json.Marshal(message)
				if err = core.Es.AddDoc(AppCfg.Es.Index, val["id"], string(b)); err != nil {
					fmt.Println("添加数据到es失败 err:", err.Error())
				}
			} else {
				//修改操作
				val := printColumn(rowData.GetAfterColumns())
				var m map[string]interface{}
				m = make(map[string]interface{}, 0)

				m["from_id"] = val["from_id"]
				m["to_id"] = val["to_id"]
				m["msg_type"] = val["msg_type"]
				m["content_type"] = val["content_type"]
				m["content"] = val["content"]
				m["create_time"], err = time.Parse("2006-01-02 15:04:05", val["create_time"])

				if err != nil {
					fmt.Println("时间转换失败 err:", err.Error())
					continue
				}

				if err = core.Es.Update(AppCfg.Es.Index, val["id"], m); err != nil {
					fmt.Println("修改数据到es失败 err:", err.Error())
				}
			}
		}
	}
}

//获取table修改记录
func printColumn(columns []*pbe.Column) map[string]string {
	var val map[string]string

	val = make(map[string]string, 0)
	for _, col := range columns {
		val[col.GetName()] = col.GetValue()
	}
	return val
}
