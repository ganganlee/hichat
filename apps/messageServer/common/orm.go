package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"hichat.zozoo.net/apps/messageServer/model"
)

//连接数据库

var AppOrm *xorm.Engine

func OrmConn(cfg *Config) (engine *xorm.Engine, err error) {

	var (
		config = cfg.Mysql
		conn   string
	)

	conn = config.Username + ":" + config.Password + "@(" + config.Host + ")/" + config.Database + "?charset=" + config.Charset
	if engine, err = xorm.NewEngine("mysql", conn); err != nil {
		return nil, err
	}

	//控制展示sql语句开关
	engine.ShowSQL(config.ShowSql)

	//同步数据库结构
	engine.Sync2(
		new(model.UserMessage0),
		new(model.UserMessage1),
		new(model.UserMessage2),
		new(model.UserMessage3),
		new(model.UserMessage4),
		new(model.UserMessage5),
		new(model.UserMessage6),
		new(model.UserMessage7),
		new(model.UserMessage8),
		new(model.UserMessage9),
		new(model.UserMessageA),
		new(model.UserMessageB),
		new(model.UserMessageC),
		new(model.UserMessageD),
		new(model.UserMessageE),
		new(model.UserMessageF),
		new(model.UserMessageG),
		new(model.UserMessageH),
		new(model.UserMessageI),
		new(model.UserMessageJ),
		new(model.UserMessageK),
		new(model.UserMessageL),
		new(model.UserMessageM),
		new(model.UserMessageN),
		new(model.UserMessageO),
		new(model.UserMessageP),
		new(model.UserMessageQ),
		new(model.UserMessageR),
		new(model.UserMessageS),
		new(model.UserMessageT),
		new(model.UserMessageU),
		new(model.UserMessageV),
		new(model.UserMessageW),
		new(model.UserMessageX),
		new(model.UserMessageY),
		new(model.UserMessageZ),
	)

	//连接池
	engine.SetMaxIdleConns(config.MaxActive)
	engine.SetMaxIdleConns(config.MaxIdle)

	AppOrm = engine
	return engine, nil

}
