package config

import (
	"fmt"

	"github.com/Andrew-M-C/trpc-go-utils/plugin"
)

type MYSQLConfig struct {
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	MaxConns  int    `yaml:"max_open_conns"`
	MaxIdle   int    `yaml:"max_idle_conns"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parse_time"`
	Loc 	  string `yaml:"loc"`
}

type MyRedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

var DBConfig MYSQLConfig

// var RedisConfig MyRedisConfig

/*
加载配置
*/
func init() {
	fmt.Printf("init plugin\n")
	plugin.Bind("database", "mysql", &DBConfig) // 绑定配置到结构体[7](@ref)
	// plugin.Bind("database", "redis", &RedisConfig) // 绑定配置到结构体[7](@ref)
}
