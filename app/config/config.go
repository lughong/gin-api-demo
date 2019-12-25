package config

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/lughong/gin-api-demo/app/util"

	"github.com/fsnotify/fsnotify"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ModOption 在新建结构体的时候调用该函数修改默认值
type ModOption func(c *Config)

// Config 配置文件结构体，存储配置文件路径
type Config struct {
	Name string
}

// New 返回初始化后的Config结构体指针
func New(modOptions ...ModOption) *Config {
	c := Config{
		Name: "",
	}

	for _, f := range modOptions {
		f(&c)
	}

	return &c
}

func (c Config) Load() error {
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := c.initLog(); err != nil {
		return err
	}

	c.watchConfig()

	return nil
}

// initConfig 解析和读取配置文件内容，这里使用的是yaml格式的配置文件
func (c Config) initConfig() error {
	viper.SetConfigFile(c.Name)
	viper.SetConfigType("yaml")

	// 设置环境变量前缀，可以通过环境变量来覆盖配置文件的值。
	// 优先级 explicit call to Set>flag>env>config>key/value store>default
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MALL")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// initLog 初始化日志
func (c Config) initLog() error {
	basePath := viper.GetString("log.path")
	if basePath != "" {
		_, _ = util.CreateDir(basePath)
	}

	fileName := viper.GetString("log.fileName")
	logFile := path.Join(basePath, fileName)

	src, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	// 设置日志输出
	logrus.SetOutput(src)

	level := viper.Get("log.logger_level")
	switch level {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// 设置日志的备份操作
	writer, err := rotates.New(
		logFile+".%Y%m%d%H",
		rotates.WithMaxAge(60*24*time.Hour),
		rotates.WithRotationTime(time.Duration(3600)*time.Second),
	)
	if err != nil {
		return err
	}

	writerMap := lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}

	lfHook := lfshook.NewHook(
		writerMap,
		&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	)

	logrus.AddHook(lfHook)

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return nil
}

// watchConfig 监控配置文件，如果配置文件发生了改变则进行热加载。
func (c Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("Config file changed: %s", e.Name)
	})
}
