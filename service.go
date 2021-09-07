package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/rpc/hello"
	"log"
)

type HelloServer struct {
}

func (h *HelloServer) Say(ctx context.Context, res *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "你好世界！"
	return nil
}

func main() {
	var (
		err error
	)

	service := micro.NewService(
		micro.Name("hello"),
		//使用etcd作为注册中心
		micro.Registry(etcd.NewRegistry(registry.Addrs("192.168.3.11:2379"))),
	)

	service.Init()

	if err = hello.RegisterHelloServiceHandler(service.Server(), new(HelloServer)); err != nil {
		panic(err)
	}

	if err = service.Run(); err != nil {
		log.Fatalf("服务启动失败，err:%v", err)
	}
}
