package core

import (
	"github.com/go-redis/redis"
	"time"
)

//redis

//var RedisPool
var CLusterClient *redis.ClusterClient

//默认过期时间
const DefaultExpire = 30 * 24 * 60 * 60 * time.Second

func RedisClusterConn(address []string, maxActive int, minIdle int, password string) *redis.ClusterClient {
	CLusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        address,
		PoolSize:     maxActive, //最大连接数
		MinIdleConns: minIdle,   //最小连接数
		Password:     password,  //连接密码
	})

	return CLusterClient
}
