package config

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/lughong/gin-api-demo/util"

	"github.com/fsnotify/fsnotify"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ModOption func(c *Config)

type Config struct {
	Name string
}

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

// initConfig parse and read the config file content.
func (c Config) initConfig() error {
	viper.SetConfigFile(c.Name)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("MALL")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// initLog create log path if the log path is not exists.
// set the log auto back up
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

// watchConfig monitor the config file.
// hot reload if modify the config file content.
func (c Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("Config file changed: %s", e.Name)
	})
}
