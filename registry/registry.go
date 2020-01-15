package registry

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/sarulabs/di"

	_logic "github.com/lughong/gin-api-demo/logic"
	_repository "github.com/lughong/gin-api-demo/repository"
)

var (
	// database config error code
	ErrDatabaseConfigNotFound   = errors.New("The database config was not found. ")
	ErrDatabaseDriverNotFound   = errors.New("The driver was not found for database config. ")
	ErrDatabaseUserNotFound     = errors.New("The user was not found for database config. ")
	ErrDatabasePasswordNotFound = errors.New("The password was not found for database config. ")
	ErrDatabaseAddrNotFound     = errors.New("The addr was not found for database config. ")
	ErrDatabaseDBNameNotFound   = errors.New("The dbname was not found for database config. ")

	// redis config error code
	ErrRedisConfigNotFound   = errors.New("The redis config was not found. ")
	ErrRedisProtocolNotFound = errors.New("The protocol was not found for redis config. ")
	ErrRedisHostNotFound     = errors.New("The host was not found for redis config. ")
	ErrRedisPortNotFound     = errors.New("The port was not found for redis config. ")
)

var config Config

type Config struct {
	DBConfig      map[string]string
	RedisConfig   map[string]string
	ContextConfig map[string]string
}

type Container struct {
	ctn di.Container
}

// CreateApp 依赖注入容器
func NewContainer(cfg Config) (*Container, error) {
	config = cfg

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
	if len(config.DBConfig) == 0 {
		return nil, ErrDatabaseConfigNotFound
	}

	driver, ok := config.DBConfig["driver"]
	if !ok {
		return nil, ErrDatabaseDriverNotFound
	}

	user, ok := config.DBConfig["user"]
	if !ok {
		return nil, ErrDatabaseUserNotFound
	}

	password, ok := config.DBConfig["password"]
	if !ok {
		return nil, ErrDatabasePasswordNotFound
	}

	addr, ok := config.DBConfig["addr"]
	if !ok {
		return nil, ErrDatabaseAddrNotFound
	}

	dbname, ok := config.DBConfig["dbname"]
	if !ok {
		return nil, ErrDatabaseDBNameNotFound
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		user,
		password,
		addr,
		dbname,
		"utf8",
		true,
		"Local",
	)

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	if moc, err := strconv.Atoi(config.DBConfig["maxOpenConns"]); err == nil {
		// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误
		db.SetMaxOpenConns(moc)
	}

	if mic, err := strconv.Atoi(config.DBConfig["maxIdleConns"]); err == nil {
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
	if len(config.RedisConfig) == 0 {
		return nil, ErrRedisConfigNotFound
	}

	protocol, ok := config.RedisConfig["protocol"]
	if !ok {
		return nil, ErrRedisProtocolNotFound
	}

	host, ok := config.RedisConfig["host"]
	if !ok {
		return nil, ErrRedisHostNotFound
	}

	port, ok := config.RedisConfig["port"]
	if !ok {
		return nil, ErrRedisPortNotFound
	}

	dbName, err := strconv.Atoi(config.RedisConfig["db"])
	if err != nil {
		dbName = 1
	}

	maxIdle, err := strconv.Atoi(config.RedisConfig["maxIdle"])
	if err != nil {
		maxIdle = 3
	}

	maxActive, err := strconv.Atoi(config.RedisConfig["maxActive"])
	if err != nil {
		maxActive = 3
	}

	idleTimeout, err := time.ParseDuration(config.RedisConfig["idleTimeout"])
	if err != nil {
		idleTimeout = 30 * time.Second
	}

	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			// 链接redis
			c, err := redis.Dial(
				protocol,
				host+port,
			)
			if err != nil {
				return nil, err
			}

			// 进行校验，如果设置了密码
			if config.RedisConfig["password"] != "" {
				if _, err := c.Do("AUTH", config.RedisConfig["password"]); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

			// 选择操作库
			if _, err := c.Do("SELECT", dbName); err != nil {
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
	timeout, err := time.ParseDuration(config.ContextConfig["timeout"])
	if err != nil {
		timeout = 5 * time.Second
	}

	userRepo := _repository.NewMysqlUserRepository(ctn.Get("mysql").(*sql.DB))
	userLogic := _logic.NewUserLogic(userRepo, timeout)
	return userLogic, nil
}
