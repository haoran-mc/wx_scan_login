package db

import (
	"fmt"

	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/models"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDB *gorm.DB

func DB() *gorm.DB {
	return gormDB
}

func syncDB() {
	_ = gormDB.AutoMigrate(
		&models.Users{},
	)
}

func init() {
	mysqlDsn := fmt.Sprintf(
		config.Conf.Mysql.Dsn,
		config.Conf.Mysql.User,
		config.Conf.Mysql.Pass,
		config.Conf.Mysql.Dbname,
	)

	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	}

	gormDB, _ = gorm.Open(postgres.Open(mysqlDsn), gormConfig)
	syncDB()
}
