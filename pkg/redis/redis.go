package redis

import (
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	pool *redis.Pool
}

func NewPool(p *redis.Pool) *Redis {
	return &Redis{
		pool: p,
	}
}

// Set 保存一个值并设定过期时间
func (r *Redis) Set(key, value string, seconds int) error {
	conn := r.pool.Get()

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
func (r *Redis) Get(key string) (string, error) {
	conn := r.pool.Get()

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
func (r *Redis) Del(key string) (bool, error) {
	conn := r.pool.Get()

	return redis.Bool(conn.Do("DEL", key))
}

// Exists 判断是否存在
func (r *Redis) Exists(key string) bool {
	conn := r.pool.Get()

	reply, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return reply
}
