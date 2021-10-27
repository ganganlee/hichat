package common

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"log"
	"os"
	"time"
)

//canal-go服务

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

	for _, entry := range entrys {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)

		eventType := rowChange.GetEventType()
		header := entry.GetHeader()

		fmt.Println("|-|-|-|-|-|-|-|-|-|-|-|-|")
		//获取操作的日志文件
		fmt.Println("log:", header.GetLogfileName())
		//获取操作的行号
		fmt.Println("offset:", header.GetLogfileOffset())
		//获取操作的数据库名称
		fmt.Println("database:", header.GetSchemaName())
		//获取操作的表名称
		fmt.Println("table:", header.GetTableName())
		//获取操作类型
		fmt.Println("operation:", header.GetEventType())

		for _, rowData := range rowChange.GetRowDatas() {
			if eventType == pbe.EventType_DELETE {
				//删除操作
				printColumn(rowData.GetBeforeColumns())
			} else if eventType == pbe.EventType_INSERT {
				//新增操作
				printColumn(rowData.GetAfterColumns())
			} else {
				//修改操作
				fmt.Println("-------> before")
				printColumn(rowData.GetBeforeColumns())
				fmt.Println("-------> after")
				printColumn(rowData.GetAfterColumns())
			}
		}
	}
}

func printColumn(columns []*pbe.Column) {
	for _, col := range columns {
		fmt.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
