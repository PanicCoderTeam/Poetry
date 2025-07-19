package infrastructure

import (
	"fmt"
	"poetry/src/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	URI string
)

func InitDB() error {
	cfg := config.DBConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)
	fmt.Printf("dsn:%s\n", dsn)
	URI = dsn
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启 SQL 日志
	})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)  // 最大空闲连接
	sqlDB.SetMaxOpenConns(cfg.MaxConns) // 最大活跃连接
	return nil
}
