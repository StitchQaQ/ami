package config

import (
	"log"
	"os"
)

// Config 全局配置
type Config struct {
	Server     ServerConfig
	Redis      RedisConfig
	ClickHouse ClickHouseConfig
	RPC        RPCConfig
	Log        LogConfig
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ClickHouseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type RPCConfig struct {
	UserServiceAddr   string
	BudgetServiceAddr string
}

type LogConfig struct {
	Level      string // debug, info, warn, error
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个文件最大大小(MB)
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 保留天数
	Compress   bool   // 是否压缩旧日志
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8088"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		ClickHouse: ClickHouseConfig{
			Host:     getEnv("CLICKHOUSE_HOST", "localhost"),
			Port:     getEnv("CLICKHOUSE_PORT", "9000"),
			Database: getEnv("CLICKHOUSE_DB", "dsp_logs"),
			Username: getEnv("CLICKHOUSE_USER", "default"),
			Password: getEnv("CLICKHOUSE_PASS", ""),
		},
		RPC: RPCConfig{
			UserServiceAddr:   getEnv("USER_SERVICE_ADDR", "localhost:50051"),
			BudgetServiceAddr: getEnv("BUDGET_SERVICE_ADDR", "localhost:50052"),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			FilePath:   getEnv("LOG_FILE_PATH", "logs/dsp-system.log"),
			MaxSize:    100, // 100MB
			MaxBackups: 7,   // 保留7个备份
			MaxAge:     7,   // 保留7天
			Compress:   true,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Using default value for %s: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
