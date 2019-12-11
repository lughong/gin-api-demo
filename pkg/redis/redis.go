package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var Pool *redis.Pool

// Init 初始化redis连接池
func Init() {
	Pool = &redis.Pool{
		MaxIdle:     viper.GetInt("redis.maxIdle"),
		MaxActive:   viper.GetInt("redis.maxActive"),
		IdleTimeout: time.Duration(viper.GetInt64("redis.idleTimeout")),
		Dial: func() (redis.Conn, error) {
			// 链接redis
			c, err := redis.Dial(
				viper.GetString("redis.protocol"),
				viper.GetString("redis.host")+viper.GetString("redis.port"),
			)
			if err != nil {
				return nil, err
			}

			// 进行校验，如果设置了密码
			password := viper.GetString("redis.password")
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

			// 选择操作库
			if _, err := c.Do("SELECT", viper.GetInt64("redis.db")); err != nil {
				_ = c.Close()
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Set 保存一个值并设定过期时间
func Set(key, value string, seconds int) error {
	conn := Pool.Get()
	_ = conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return err
	}

	return nil
}

// Get 获取
func Get(key string) (string, error) {
	conn := Pool.Get()
	_ = conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	if err == redis.ErrNil {
		return "", nil
	}

	return reply, nil
}

// Del 删除
func Del(key string) (bool, error) {
	conn := Pool.Get()
	_ = conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// Exists 判断是否存在
func Exists(key string) bool {
	conn := Pool.Get()
	_ = conn.Close()

	reply, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return reply
}
