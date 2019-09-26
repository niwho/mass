package utils

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DBClient struct {
	Name     string
	Server   string
	User     string
	Password string

	RealDb *gorm.DB
}

func NewDBClient(name, server, user, passwd string) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initdb(10, 1000, 250, 300*time.Second); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

// timeout ms
func NewDBClientV2(name, server, user, passwd string, maxConn int, timeout int) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initdb(maxConn, 1000, timeout, 300*time.Second); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func NewDBClientV3(name, server, user, passwd string, maxConn int, timeout int, sessionDuration time.Duration) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initdb(maxConn, 3000, timeout, sessionDuration); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func NewDBClientV4(name, server, user, passwd string, maxConn int, connectTime int, timeout int, sessionDuration time.Duration) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initdb(maxConn, connectTime, timeout, sessionDuration); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func (db *DBClient) initdb(maxConn, connectTime int, timeout int, d time.Duration) error {
	var err error
	db.RealDb, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%dms&readTimeout=%dms&writeTimeout=%dms", db.User, db.Password, db.Server, db.Name, connectTime, timeout, timeout))
	db.RealDb.DB().SetMaxOpenConns(maxConn)
	db.RealDb.DB().SetConnMaxLifetime(d)
	return err
}
