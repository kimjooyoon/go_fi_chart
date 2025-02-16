package config

import "time"

// Config 애플리케이션 전체 설정입니다.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig HTTP 서버 설정입니다.
type ServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

// MetricsConfig 메트릭 설정입니다.
type MetricsConfig struct {
	Enabled         bool          `yaml:"enabled"`
	CollectInterval time.Duration `yaml:"collectInterval"`
	ExportInterval  time.Duration `yaml:"exportInterval"`
	RetentionPeriod time.Duration `yaml:"retentionPeriod"`
}

// DatabaseConfig 데이터베이스 설정입니다.
type DatabaseConfig struct {
	Driver          string        `yaml:"driver"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Name            string        `yaml:"name"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	ConnectTimeout  time.Duration `yaml:"connectTimeout"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

// NewDefaultConfig 기본 설정값을 가진 Config를 생성합니다.
func NewDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "localhost",
			Port:            8080,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     120 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		},
		Metrics: MetricsConfig{
			Enabled:         true,
			CollectInterval: 10 * time.Second,
			ExportInterval:  30 * time.Second,
			RetentionPeriod: 24 * time.Hour,
		},
		Database: DatabaseConfig{
			Driver:          "postgres",
			Host:            "localhost",
			Port:            5432,
			ConnectTimeout:  5 * time.Second,
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
	}
}
