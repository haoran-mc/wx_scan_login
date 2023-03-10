package db

import (
	"github.com/haoran-mc/wx_scan_login/web-back-end/internal/model"
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
		&model.Users{},
	)
}

func init() {
	pgDsn := config.Conf.Postgres.Dsn

	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	}

	gormDB, _ = gorm.Open(postgres.Open(pgDsn), gormConfig)

	syncDB()
}
