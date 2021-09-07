package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"hichat.zozoo.net/rpc/hello"
	"log"
)

func main() {
	var (
		err error
		res *hello.Response
	)

	service := micro.NewService(
		micro.Registry(etcd.NewRegistry(registry.Addrs("192.168.3.11:2379"))),
		)
	h := hello.NewHelloService("hello",service.Client())

	res,err = h.Say(context.TODO(),&hello.Request{
		Name: "gangan",
	})

	if err != nil {
		log.Fatalf("调用微服务方法出错! err:%v",err)
	}

	fmt.Println(res.Msg)
}
