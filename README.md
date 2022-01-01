# rdbcli

封装了go-redis，支持单机和集群

配置

```go
type Config struct {
    // 是否为集群模式
    Cluster  bool   `mapstructure:"cluster" `
    // 主机
    Host     string `mapstructure:"host" `
    // 端口
    Port     string `mapstructure:"port" `
    // 密码
    Password string `mapstructure:"password"`
    // 连接协议
    Protocol string `mapstructure:"protocol"`
    // 初始连接的数据库
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
```

#### 使用

```go
cli:=rdbcli.New(config)
rdbcli.Cli.Cmd()
```

> 具体命令可参考 [go-redis](https://github.com/go-redis/redishttps://)
