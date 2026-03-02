package config

import (
	"errors"
	"log"
	"strings"

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

type AdminConfig struct {
	Username string
	Password string
}

func Init() {
	// 设置配置文件的路径和文件名
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 尝试读取配置文件，如果不存在则忽略（仅警告）
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No config file found, using environment variables and defaults")
		}
	}

	// 设置环境变量前缀，并启用自动覆盖
	viper.SetEnvPrefix("ZHIHU")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置合理地默认值（防止某些环境变量缺失）
	viper.SetDefault("mysql.max_open_conns", 10)
	viper.SetDefault("mysql.max_idle_conns", 5)
	viper.SetDefault("redis.db", 0)

	// 将配置内容映射到结构体
	Config = &AppConfig{}
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatalf("Unable to decode into struct %v", err)
	}

	// 可选：打印关键配置（调试用）
	log.Println("Configuration loaded successfully")
}
