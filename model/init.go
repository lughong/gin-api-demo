package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var DB *Database

type Database struct {
	Conn *sqlx.DB
}

func Init() (*Database, error) {
	conn, err := openDB(
		viper.GetString("database.driver"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.addr"),
		viper.GetString("database.dbname"),
	)
	if err != nil {
		return nil, err
	}

	DB = &Database{
		Conn: conn,
	}

	return DB, nil
}

func openDB(driver, user, password, addr, dbname string) (*sqlx.DB, error) {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
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

	db, err := sqlx.Connect(driver, dataSourceName)
	if err != nil {
		logrus.Errorf("Database connection failed. Database name: %s. %s", dbname, err)
		return nil, err
	}

	setupDB(db)

	return db, nil
}

func setupDB(db *sqlx.DB) {
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误
	db.SetMaxOpenConns(viper.GetInt("database.maxOpenConns"))

	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.SetMaxIdleConns(viper.GetInt("database.maxIdleConns"))
}

func (db *Database) Close() {
	_ = db.Conn.Close()
}
