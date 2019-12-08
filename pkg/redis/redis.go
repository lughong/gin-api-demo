package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var Pool *redis.Pool

func Init() {
	Pool = &redis.Pool{
		MaxIdle:     viper.GetInt("redis.maxIdle"),
		MaxActive:   viper.GetInt("redis.maxActive"),
		IdleTimeout: time.Duration(viper.GetInt64("redis.idleTimeout")),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				viper.GetString("redis.protocol"),
				viper.GetString("redis.host")+viper.GetString("redis.port"),
			)
			if err != nil {
				return nil, err
			}

			password := viper.GetString("redis.password")
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

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

func Del(key string) (bool, error) {
	conn := Pool.Get()
	_ = conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func Exists(key string) bool {
	conn := Pool.Get()
	_ = conn.Close()

	reply, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return reply
}
