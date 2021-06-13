package util

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var err error

func init() {
	logs.Info("[gorm] init gorm")
	initDB()
	ShowDB()
	//initTaskShowDB()
}

func initDB() {
	Url, _ := web.AppConfig.String("sql.Url")
	DriverName, _ := web.AppConfig.String("sql.dsn")
	_, _ = web.AppConfig.String("log.ORMLog")
	RunMode, _ := web.AppConfig.String("runmode")
	logs.Trace("connect to ", Url, DriverName)
	var logInfo logger.LogLevel
	// Silent、Error、Warn、Info
	if RunMode == "dev" {
		logInfo = logger.Info
	} else {
		logInfo = logger.Error
	}
	ormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logInfo,     // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// set default database
	//dsn := "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=3306 sslmode=disable TimeZone=Asia/Shanghai"
	DB, err = gorm.Open(postgres.Open(DriverName), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 全局单数表名
		},
	},
	)
	if err != nil {
		logs.Emergency("CANNOT CONNECTION TO POSTGRES :", err.Error())
	}
	logs.Info("connection to postgresSQL success")
}

func ShowDB() {
	type Result struct {
		Version string
	}
	var result Result
	DB.Raw("SELECT VERSION()").Scan(&result)
	logs.Info("[gorm] RAW `SELECT VERSION()` :", result.Version) // slene
	sqlDB, errDB := DB.DB()
	if errDB != nil {
		logs.Error("[gorm] ShowDB() :", errDB)
	}
	stats := sqlDB.Stats()

	logs.Trace("[gorm] ShowDB()", fmt.Sprintf("\n\t %+v", stats))

}
