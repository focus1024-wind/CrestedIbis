package initizalize

import (
	"CrestedIbis/src/apps/audit_log"
	"CrestedIbis/src/apps/ipc/ipc_alarm"
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/apps/site"
	"CrestedIbis/src/apps/system/user"
	"CrestedIbis/src/config/model"
	"CrestedIbis/src/global"
	"CrestedIbis/src/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// InitDatabase 初始化数据库配置
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

// MySql数据库配置
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
		panic(fmt.Sprintf("连接数据库 %s 失败: [%s]", datasource.Type, err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("连接数据库 %s 失败: [%s]", datasource.Type, err.Error()))
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
		panic(fmt.Sprintf("连接数据库 %s 失败: [%s]", datasource.Type, err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("连接数据库 %s 失败: [%s]", datasource.Type, err.Error()))
	}

	sqlDB.SetMaxIdleConns(datasource.MaxIdle)
	sqlDB.SetMaxOpenConns(datasource.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

// InitDbTable 初始化数据库表
func InitDbTable() {
	_ = global.Db.AutoMigrate(
		&audit_log.AuditLogLogin{},
		&ipc_device.IpcDevice{},
		&ipc_device.IpcChannel{},
		&ipc_alarm.IpcAlarm{},
		&ipc_alarm.IpcRecord{},
		&user.RoleGroup{},
		&user.SysUser{},
		&utils.CasbinRule{},
		&site.Site{})
}
