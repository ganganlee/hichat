package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"hichat.zozoo.net/apps/gatewayServer/models"
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
		new(models.UserMessage0),
		new(models.UserMessage1),
		new(models.UserMessage2),
		new(models.UserMessage3),
		new(models.UserMessage4),
		new(models.UserMessage5),
		new(models.UserMessage6),
		new(models.UserMessage7),
		new(models.UserMessage8),
		new(models.UserMessage9),
		new(models.UserMessageA),
		new(models.UserMessageB),
		new(models.UserMessageC),
		new(models.UserMessageD),
		new(models.UserMessageE),
		new(models.UserMessageF),
		new(models.UserMessageG),
		new(models.UserMessageH),
		new(models.UserMessageI),
		new(models.UserMessageJ),
		new(models.UserMessageK),
		new(models.UserMessageL),
		new(models.UserMessageM),
		new(models.UserMessageN),
		new(models.UserMessageO),
		new(models.UserMessageP),
		new(models.UserMessageQ),
		new(models.UserMessageR),
		new(models.UserMessageS),
		new(models.UserMessageT),
		new(models.UserMessageU),
		new(models.UserMessageV),
		new(models.UserMessageW),
		new(models.UserMessageX),
		new(models.UserMessageY),
		new(models.UserMessageZ),
	)

	//连接池
	engine.SetMaxIdleConns(config.MaxActive)
	engine.SetMaxIdleConns(config.MaxIdle)

	AppOrm = engine
	return engine, nil

}
