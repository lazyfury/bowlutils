package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DriverMySQL      = "mysql"
	DriverPostgreSQL = "postgres"
)

var (
	DefaultDriver  = DriverPostgreSQL
	DefaultDrivers = []string{DriverMySQL, DriverPostgreSQL}
)

type DBConfig struct {
	Driver string
	DSN    string
}

func NewDB(driver string, dsn string, confs ...gorm.Option) *gorm.DB {
	if driver == "" || driver == "auto" {
		driver = DefaultDriver
	}

	confs = append([]gorm.Option{
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		},
	}, confs...)

	var DB *gorm.DB
	var err error
	switch driver {
	case DriverMySQL:
		DB, err = gorm.Open(mysql.Open(dsn), confs...)
	case DriverPostgreSQL:
		DB, err = gorm.Open(postgres.Open(dsn), confs...)
	default:
		panic("unsupported driver: " + driver)
	}
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return DB
}
