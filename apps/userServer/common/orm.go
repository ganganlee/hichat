package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"hichat.zozoo.net/apps/userServer/model"
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
		new(model.User),             //用户表
		new(model.UserFriends),      //用户好友表
		new(model.UserGroups),       //用户群表
		new(model.UserGroupMembers), //用户群成员表
	)

	//连接池
	engine.SetMaxIdleConns(config.MaxActive)
	engine.SetMaxIdleConns(config.MaxIdle)

	AppOrm = engine
	return engine, nil

}
