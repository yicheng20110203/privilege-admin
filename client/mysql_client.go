package client

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "gitlab.ceibsmoment.com/c/mp/config"
    "sync"
)

var (
    _db *gorm.DB
    mux sync.RWMutex
)

func GetMysqlDb() *gorm.DB {
    if _db != nil && _db.DB().Stats().OpenConnections > 0 {
        return _db
    }

    var err error
    cfg := config.Cfg

    mux.Lock()
    defer mux.Unlock()
    //连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
    _db, err = gorm.Open("mysql", cfg.Mysql.Dsn)
    if err != nil {
        panic("连接数据库失败, error=" + err.Error())
    }

    _db.DB().SetMaxIdleConns(20)
    _db.DB().SetMaxOpenConns(50)

    if cfg.Mysql.Debug == 1 {
        _db = _db.LogMode(true).Debug()
    }

    return _db
}