package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config *AppConfig

type AppConfig struct {
	Mysql MysqlConfig
	Redis RedisConfig
	Admin AdminConfig
}

type MysqlConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

//新增管理员配置结构体

type AdminConfig struct {
	Username string
	Password string
}

func Init() {
	// 设置配置文件的路径和文件名
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// 将配置内容映射到结构体
	Config = &AppConfig{}
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatalf("Unable to decode into struct %v", err)
	}
}
