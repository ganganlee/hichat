package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/apps/messageSearchServer/common"
	"hichat.zozoo.net/apps/messageSearchServer/rpc"
	"hichat.zozoo.net/apps/messageSearchServer/services"
	"hichat.zozoo.net/core"
	"hichat.zozoo.net/rpc/messageSearch"
	"log"
	"os"
	"strings"
)

func main() {
	//获取配置文件
	var (
		path string
		err  error
		cfg  *common.Config
		es   *core.Elasticsearch
	)

	//获取当前目录
	if path, err = os.Getwd(); err != nil {
		log.Fatalf("获取配置文件目录失败 err:%v\n", err)
	}
	path = path + "/config/app.json"

	//获取配置文件
	if cfg, err = common.GetConfig(path); err != nil {
		log.Fatalf("读取配置文件失败 err:%v\n", err)
	}

	//注册es服务
	es, err = core.NewEs(cfg.Es.Host, cfg.Es.Username, cfg.Es.Password)
	if err != nil {
		log.Fatalf("elasticsearch服务连接失败 err:%v\n", err)
	}

	//创建索引
	if err = es.CreateIndex(cfg.Es.Index, common.GetEsIndex()); err != nil {
		//判断索引存在不报错
		if strings.Index(err.Error(), "already exists") == -1 {
			log.Fatalf("elasticsearch创建索引失败 err:%v\n", err)
		}
	}

	//连接canal
	go common.InitCanal(cfg.CanalConfig.Address, cfg.CanalConfig.Port, cfg.CanalConfig.Username, cfg.CanalConfig.Password, cfg.CanalConfig.Destination, cfg.CanalConfig.SoTimeOut, cfg.CanalConfig.IdleTimeOut)

	//初始化微服务
	service := micro.NewService(
		micro.Name(cfg.ServerName),
		micro.Registry(etcd.NewRegistry(registry.Addrs(cfg.Etcd.Host))),
	)
	service.Init()

	//注册搜索服务
	messageSearchService := services.NewMessageSearchService()
	if err = messageSearch.RegisterSearchMessageServiceHandler(service.Server(), rpc.NewMessageSearchRpc(messageSearchService)); err != nil {
		log.Fatalf("注册微服务失败 err:%v\n", err)
	}

	//运行微服务
	if err = service.Run(); err != nil {
		log.Fatalf("微服务启动失败 err:%v\n", err)
	}
}
