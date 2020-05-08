package redis

import (
	"github.com/go-redis/redis"
	"github.com/wong-winnie/library/dao/config"
)

/*
	//命令执行失败重试次数(默认不重试)
	MaxRetries: 3,
	// 空闲连接超时时间，超过超时时间的空闲连接会被关闭。
	// 如果设置成0，空闲连接将不会被关闭
	// 应该设置一个比redis服务端超时时间更短的时间
	IdleTimeout: time.Second * 30,
*/

type RedisSingleMgr struct {
	Conn *redis.Client
}

func InitRedisSingle(cfg *config.RedisCfg) *RedisSingleMgr {
	return &RedisSingleMgr{
		Conn: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
		}),
	}
}

type RedisDoubleMgr struct {
	Conn *redis.ClusterClient
}

func InitRedisDouble(cfg *config.RedisCfg) *RedisDoubleMgr {
	return &RedisDoubleMgr{
		Conn: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: cfg.Address,
		}),
	}
}

/*
	队列(左出右进)
*/
func (mgr *RedisSingleMgr) Push(key string, values ...interface{}) *redis.IntCmd {
	return mgr.Conn.RPush(key, values...)
}

func (mgr *RedisSingleMgr) Pop(key string) *redis.StringCmd {
	return mgr.Conn.LPop(key)
}

func (mgr *RedisSingleMgr) Range(key string, start, stop int64) *redis.StringSliceCmd {
	return mgr.Conn.LRange(key, start, stop)
}
