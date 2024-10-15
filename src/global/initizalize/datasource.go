package initizalize

import (
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/config/model"
	"CrestedIbis/src/global"
	"CrestedIbis/src/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// InitDatabase connect database
func InitDatabase(datasource *model.Datasource) *gorm.DB {
	switch datasource.Type {
	case "mysql":
		return initMysql(datasource)
	case "postgres":
		return initPostgres(datasource)
	default:
		panic(fmt.Sprintf("Database type %s not support, Can Only Support postgres", datasource.Type))
	}
}

// postgres database config
func initMysql(datasource *model.Datasource) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		datasource.Username,
		datasource.Password,
		datasource.Host,
		datasource.Port,
		datasource.DbName,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to %s database: [%s]", datasource.Type, err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to %s database: [%s]", datasource.Type, err.Error()))
	}
	sqlDB.SetMaxIdleConns(datasource.MaxIdle)
	sqlDB.SetMaxOpenConns(datasource.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

// postgres database config
func initPostgres(datasource *model.Datasource) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		datasource.Host,
		datasource.Port,
		datasource.Username,
		datasource.Password,
		datasource.DbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to %s database: [%s]", datasource.Type, err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to %s database: [%s]", datasource.Type, err.Error()))
	}
	sqlDB.SetMaxIdleConns(datasource.MaxIdle)
	sqlDB.SetMaxOpenConns(datasource.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func InitDbTable() {
	err := global.Db.AutoMigrate(
		&ipc_device.IpcDevice{},
		&ipc_device.IpcChannel{},
		&utils.CasbinRule{})

	if err != nil {
		return
	}
}
