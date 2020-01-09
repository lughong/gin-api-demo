package registry

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_logic "github.com/lughong/gin-api-demo/api/user/logic"
	_repository "github.com/lughong/gin-api-demo/api/user/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/sarulabs/di"
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

	var (
		driver   string
		user     string
		password string
		addr     string
		dbname   string

		ok bool
	)

	if driver, ok = config.DBConfig["driver"]; !ok {
		return nil, ErrDatabaseDriverNotFound
	}

	if user, ok = config.DBConfig["user"]; !ok {
		return nil, ErrDatabaseUserNotFound
	}

	if password, ok = config.DBConfig["password"]; !ok {
		return nil, ErrDatabasePasswordNotFound
	}

	if addr, ok = config.DBConfig["addr"]; !ok {
		return nil, ErrDatabaseAddrNotFound
	}

	if dbname, ok = config.DBConfig["dbname"]; !ok {
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

	if maxOpenConns, ok := config.DBConfig["maxOpenConns"]; ok {
		if moc, err := strconv.Atoi(maxOpenConns); err == nil {
			// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误
			db.SetMaxOpenConns(moc)
		}
	}

	if maxIdleConns, ok := config.DBConfig["maxIdleConns"]; ok {
		if mic, err := strconv.Atoi(maxIdleConns); err == nil {
			// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
			db.SetMaxIdleConns(mic)
		}
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

	var (
		maxIdle     int
		maxActive   int
		idleTimeout int

		protocol string
		host     string
		port     string
		dbName   int

		err error

		ok bool
	)

	if protocol, ok = config.RedisConfig["protocol"]; !ok {
		return nil, ErrRedisProtocolNotFound
	}

	if host, ok = config.RedisConfig["host"]; !ok {
		return nil, ErrRedisHostNotFound
	}

	if port, ok = config.RedisConfig["port"]; !ok {
		return nil, ErrRedisPortNotFound
	}

	if db, ok := config.RedisConfig["db"]; ok {
		if dbName, err = strconv.Atoi(db); err != nil {
			dbName = 0
		}
	}

	if mi, ok := config.RedisConfig["maxIdle"]; ok {
		if maxIdle, err = strconv.Atoi(mi); err != nil {
			maxIdle = 3
		}
	}

	if ma, ok := config.RedisConfig["maxActive"]; ok {
		if maxActive, err = strconv.Atoi(ma); err != nil {
			maxActive = 3
		}
	}

	if it, ok := config.RedisConfig["idleTimeout"]; ok {
		if idleTimeout, err = strconv.Atoi(it); err != nil {
			idleTimeout = 30
		}
	}

	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout),
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
	var (
		timeout int

		err error
	)
	if t, ok := config.ContextConfig["timeout"]; ok {
		if timeout, err = strconv.Atoi(t); err != nil {
			timeout = 30
		}
	}

	userRepo := _repository.NewMysqlUserRepository(ctn.Get("mysql").(*sql.DB))
	userLogic := _logic.NewUserLogic(userRepo, time.Duration(timeout)*time.Second)
	return userLogic, nil
}
