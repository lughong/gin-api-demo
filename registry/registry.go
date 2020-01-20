package registry

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/sarulabs/di"
	"github.com/spf13/viper"

	_logic "github.com/lughong/gin-api-demo/logic"
	_repository "github.com/lughong/gin-api-demo/repository"
)

var (
	// database config error code
	ErrDatabaseDriverNotFound   = errors.New("The driver was not found for database config. ")
	ErrDatabaseUserNotFound     = errors.New("The user was not found for database config. ")
	ErrDatabasePasswordNotFound = errors.New("The password was not found for database config. ")
	ErrDatabaseAddrNotFound     = errors.New("The addr was not found for database config. ")
	ErrDatabaseDBNameNotFound   = errors.New("The dbname was not found for database config. ")

	// redis config error code
	ErrRedisProtocolNotFound = errors.New("The protocol was not found for redis config. ")
	ErrRedisAddrNotFound     = errors.New("The addr was not found for redis config. ")
)

type Container struct {
	ctn di.Container
}

// CreateApp 依赖注入容器
func NewContainer() (*Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	if err := builder.Add([]di.Def{
		{
			Name:  "mysql-pool",
			Scope: di.App,
			Build: buildMysqlPool,
			Close: func(obj interface{}) error {
				return obj.(*sql.DB).Close()
			},
		},
		{
			Name:  "mysql",
			Scope: di.App,
			Build: buildMysql,
			Close: func(obj interface{}) error {
				return obj.(*sql.DB).Close()
			},
		},
		{
			Name:  "redis-pool",
			Scope: di.App,
			Build: buildRedisPool,
			Close: func(obj interface{}) error {
				return obj.(*redis.Pool).Close()
			},
		},
		{
			Name:  "redis",
			Scope: di.Request,
			Build: buildRedis,
			Close: func(obj interface{}) error {
				return obj.(*redis.Pool).Close()
			},
		},
		{
			Name:  "user-logic",
			Scope: di.App,
			Build: buildUserLogic,
		},
	}...); err != nil {
		return nil, err
	}

	return &Container{
		builder.Build(),
	}, nil
}

// Resolve 获取一个名为name的服务
func (c *Container) Resolve(name string) interface{} {
	return c.ctn.Get(name)
}

func (c *Container) Clean(name string) error {
	return c.ctn.Clean()
}

func (c *Container) Delete() error {
	return c.ctn.Delete()
}

// buildMysqlPool 链接mysql
func buildMysqlPool(ctn di.Container) (interface{}, error) {
	driver := viper.GetString("database.driver")
	if driver == "" {
		return nil, ErrDatabaseDriverNotFound
	}

	user := viper.GetString("database.user")
	if user == "" {
		return nil, ErrDatabaseUserNotFound
	}

	password := viper.GetString("database.password")
	if password == "" {
		return nil, ErrDatabasePasswordNotFound
	}

	addr := viper.GetString("database.addr")
	if addr == "" {
		return nil, ErrDatabaseAddrNotFound
	}

	dbName := viper.GetString("database.dbName")
	if dbName == "" {
		return nil, ErrDatabaseDBNameNotFound
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		user,
		password,
		addr,
		dbName,
		"utf8",
		true,
		"Local",
	)

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误
	db.SetMaxOpenConns(viper.GetInt("database.maxOpenConns"))

	if mic := viper.GetInt("database.maxIdleConns"); mic > 0 {
		// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
		db.SetMaxIdleConns(mic)
	}

	return db, err
}

// buildMysql 建立一个mysql链接
func buildMysql(ctn di.Container) (interface{}, error) {
	pool := ctn.Get("mysql-pool").(*sql.DB)
	return pool, nil
}

// buildRedisPool 链接redis
func buildRedisPool(ctn di.Container) (interface{}, error) {
	protocol := viper.GetString("redis.protocol")
	if protocol == "" {
		return nil, ErrRedisProtocolNotFound
	}

	addr := viper.GetString("redis.addr")
	if addr == "" {
		return nil, ErrRedisAddrNotFound
	}

	idleTimeout, err := time.ParseDuration(viper.GetString("redis.idleTimeout"))
	if err != nil {
		idleTimeout = 30 * time.Second
	}

	return &redis.Pool{
		MaxIdle:     viper.GetInt("redis.maxIdle"),
		MaxActive:   viper.GetInt("redis.maxActive"),
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			// 链接redis
			c, err := redis.Dial(protocol, addr)
			if err != nil {
				return nil, err
			}

			// 进行校验，如果设置了密码
			if password := viper.GetString("redis.password"); password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

			// 选择操作库
			if _, err := c.Do("SELECT", viper.GetInt("redis.dbName")); err != nil {
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
	}, nil
}

// buildRedis 建立一个redis链接
func buildRedis(ctn di.Container) (interface{}, error) {
	pool := ctn.Get("redis-pool").(*redis.Pool)
	return pool, nil
}

// buildUserLogic 建立一个user的逻辑处理实例
func buildUserLogic(ctn di.Container) (interface{}, error) {
	timeout, err := time.ParseDuration(viper.GetString("context.timeout"))
	if err != nil {
		timeout = 5 * time.Second
	}

	userRepo := _repository.NewMysqlUserRepository(ctn.Get("mysql").(*sql.DB))
	userLogic := _logic.NewUserLogic(userRepo, timeout)
	return userLogic, nil
}
