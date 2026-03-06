package config

import (
	"errors"
	"log"
	"strings"
	"zhihu-go/pkg/logger"

	"github.com/spf13/viper"
)

var Config *AppConfig

type AppConfig struct {
	Mysql MysqlConfig
	Redis RedisConfig
	Admin AdminConfig
	Log   LogConfig
}

type MysqlConfig struct {
	DSN          string `mapstructure:"dsn"`            // 添加标签
	MaxOpenConns int    `mapstructure:"max_open_conns"` // 添加标签
	MaxIdleConns int    `mapstructure:"max_idle_conns"` // 添加标签
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`     // 添加标签
	Password string `mapstructure:"password"` // 添加标签
	DB       int    `mapstructure:"db"`       // 添加标签
}

type AdminConfig struct {
	Username string `mapstructure:"username"` // 添加标签
	Password string `mapstructure:"password"` // 添加标签
}

type LogConfig struct {
	Level  string   `mapstructure:"level"`
	Output []string `mapstructure:"output"` // 例如 ["stdout", "./logs/app.log"]
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

	// **手动绑定环境变量到配置键**（关键步骤）
	viper.BindEnv("mysql.dsn", "ZHIHU_MYSQL_DSN")
	viper.BindEnv("mysql.max_open_conns", "ZHIHU_MYSQL_MAX_OPEN_CONNS")
	viper.BindEnv("mysql.max_idle_conns", "ZHIHU_MYSQL_MAX_IDLE_CONNS")
	viper.BindEnv("redis.addr", "ZHIHU_REDIS_ADDR")
	viper.BindEnv("redis.password", "ZHIHU_REDIS_PASSWORD")
	viper.BindEnv("redis.db", "ZHIHU_REDIS_DB")
	viper.BindEnv("admin.username", "ZHIHU_ADMIN_USERNAME")
	viper.BindEnv("admin.password", "ZHIHU_ADMIN_PASSWORD")
	viper.BindEnv("log.level", "ZHIHU_LOG_LEVEL")
	viper.BindEnv("log.output", "ZHIHU_LOG_OUTPUT")

	// 设置合理地默认值（防止某些环境变量缺失）
	viper.SetDefault("mysql.max_open_conns", 10)
	viper.SetDefault("mysql.max_idle_conns", 5)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", []string{"stdout"})

	// 将配置内容映射到结构体
	Config = &AppConfig{}
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatalf("Unable to decode into struct %v", err)
	}

	logger.S().Infow("Configuration loaded successfully",
		"mysql_dsn", Config.Mysql.DSN,
		"redis_addr", Config.Redis.Addr,
	)
}
