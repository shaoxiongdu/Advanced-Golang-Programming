package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func SetRedisLog(logger *log.Logger) {
	redis.SetLogger(logger)
}

// NewRedis -
func NewRedis(conf RedisConf) (clis []*redis.Client, error error) {
	if len(conf.Hosts) == 0 {
		j, _ := json.Marshal(conf)
		return nil, fmt.Errorf("invalid redis conf: %s", j)
	}

	db, _ := strconv.Atoi(conf.DefaultDB)

	for _, host := range conf.Hosts {
		options := redis.Options{
			//连接信息
			Network:  "tcp",                                    //网络类型，tcp or unix，默认tcp
			Addr:     fmt.Sprintf("%v:%v", host.IP, host.Port), //主机名+冒号+端口，默认localhost:6379
			Password: conf.Passowrd,                            //密码
			DB:       db,                                       // redis数据库index

			//连接池容量及闲置连接数量
			PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
			MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

			//超时
			DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
			ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
			WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
			PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

			//闲置连接检查包括IdleTimeout，MaxConnAge
			IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
			IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
			MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

			//命令执行失败时的重试策略
			MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
			MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
		}

		clis = append(clis, redis.NewClient(&options))
	}

	return clis, nil
}

// NewRedisSentine -
func NewRedisSentine(conf RedisConf) (clis []*redis.Client, err error) {
	if len(conf.Hosts) == 0 {
		j, _ := json.Marshal(conf)
		return nil, fmt.Errorf("invalid redis conf: %s", j)
	}

	hosts := []string{}
	for _, host := range conf.Hosts {
		hosts = append(hosts, fmt.Sprintf("%v:%v", host.IP, host.Port))
	}

	db, _ := strconv.Atoi(conf.DefaultDB)

	failoverOptions := redis.FailoverOptions{
		MasterName:    conf.MasterName,
		SentinelAddrs: hosts,
		Password:      conf.Passowrd, //密码
		DB:            db,            // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	}

	clis = append(clis, redis.NewFailoverClient(&failoverOptions))
	return clis, nil
}
