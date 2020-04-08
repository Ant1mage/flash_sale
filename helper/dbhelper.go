package helper

import (
	"fmt"
	"log"
	"sync"
	
	"flash-sale/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db   *gorm.DB
	lock sync.Mutex
)

func InstanceDB() *gorm.DB {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if db == nil {
			c := conf.MasterDBConf
			driverSource := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Pwd, c.DBName)
			println(driverSource)
			db, err := gorm.Open(conf.DriverName, driverSource)
			if err != nil {
				log.Fatal("db instance error: ", err)
				return nil
			}
			return db
		}
	}
	return db
}
