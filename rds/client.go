package rds

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client //直接继承自redis.Client

	Config *Config
}

func (this *Client) WithConfig(config *Config) *Client {
	this.Config = config
	this.Open()
	return this
}

func (this *Client) Open() {
	this.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", this.Config.Host, this.Config.Port),
		Password: this.Config.Password,
		DB:       this.Config.Db,

		//go-redis 底层维护了一个连接池，不需要手动管理。默认情况下， go-redis 连接池大小为 runtime.GOMAXPROCS * 10，
		//在大多数情况下默认值已经足够使用，且设置太大的连接池几乎没有什么用

		//钩子函数
		OnConnect: func(ctx context.Context, conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
			//fmt.Printf("conn=%v\n", conn)
			hlog.Debugf("build redis new conn %v", conn)
			return nil
		},
	})
}

var (
	clientIns *Client
)

func NewClient() (c *Client) {
	c = &Client{}
	c.Config = DefaultConfig()
	return
}

func ClientIns() (client *Client) {
	if clientIns == nil {
		clientIns = &Client{}
		clientIns.Open()
		//clientIns.Client = redis.NewClient(&redis.Options{
		//	//Network:  "tcp", //tcp or unix，default tcp
		//	Addr:     fmt.Sprintf("%s:%d", global.AppConf.Redis.Host, global.AppConf.Redis.Port),
		//	Password: global.AppConf.Redis.Password,
		//	DB:       global.AppConf.Redis.Db,
		//
		//	//go-redis 底层维护了一个连接池，不需要手动管理。默认情况下， go-redis 连接池大小为 runtime.GOMAXPROCS * 10，
		//	//在大多数情况下默认值已经足够使用，且设置太大的连接池几乎没有什么用
		//
		//	//钩子函数
		//	OnConnect: func(ctx context.Context, conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		//		//fmt.Printf("conn=%v\n", conn)
		//		hlog.Debugf("build redis new conn %v", conn)
		//		return nil
		//	},
		//})
	}
	return clientIns
}
