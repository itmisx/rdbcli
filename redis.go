package rdbcli

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisCli client

type client interface {
	redis.Cmdable
	Close() error
}

type Config struct {
	Cluster  bool   `mapstructure:"cluster" `
	Host     string `mapstructure:"host" `
	Port     string `mapstructure:"port" `
	Password string `mapstructure:"password"`
	Protocol string `mapstructure:"protocol"`
	Database int    `mapstructure:"database"`
	// 最小空闲连接
	MinIdleConns int `mapstructure:"min_idle_conns"`
	// 空闲时间
	IdleTimeout int `mapstructure:"idle_timeout"`
	// 连接池大小
	PoolSize int `mapstructure:"pool_size"`
	// 连接最大可用时间
	MaxConnAge int `mapstructure:"max_conn_age"`
}

// redis初始化客户端
func RedisInit(conf Config) {
	config := conf
	ctx := context.Background()
	hostMembers := strings.Split(config.Host, ",")

	// 默认闲置连接
	if conf.MinIdleConns == 0 {
		conf.MinIdleConns = 2
	}
	// 空闲超时时间，过期关闭空闲连接
	if conf.IdleTimeout == 0 || conf.IdleTimeout > 1800 {
		conf.IdleTimeout = 1800
	}
	// 默认连接池数量为2
	if conf.PoolSize == 0 {
		conf.PoolSize = 10
	}
	// 连接的生命周期为300秒
	if conf.MaxConnAge == 0 || conf.MaxConnAge > 3600 {
		conf.MaxConnAge = 3600
	}

	// 非集群
	if len(hostMembers) <= 1 && !config.Cluster {
		rdb := redis.NewClient(&redis.Options{
			Addr:         config.Host + ":" + config.Port,
			Password:     config.Password,
			DB:           config.Database,
			MinIdleConns: config.MinIdleConns,
			IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
			PoolSize:     config.PoolSize,
			MaxConnAge:   time.Second * time.Duration(config.MaxConnAge),
		})
		res, err := rdb.Ping(ctx).Result()
		if strings.ToLower(res) != "pong" || err != nil {
			panic("redis init failed!")
		}
		redisCli = rdb
		return
	}
	// 集群
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        hostMembers,
		Password:     config.Password,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
		PoolSize:     config.PoolSize,
		MaxConnAge:   time.Second * time.Duration(config.MaxConnAge),
	})
	res, err := rdb.Ping(ctx).Result()
	if strings.ToLower(res) != "pong" || err != nil {
		panic("redis init failed!")
	}
	redisCli = rdb
}

// 获取redis cli对象
func Cli() redis.Cmdable {
	return redisCli
}

func Close() {
	redisCli.Close()
}
