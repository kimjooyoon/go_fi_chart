package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(_ *testing.T) {
	// TODO: Add tests
}

func TestNewDefaultConfig(t *testing.T) {
	t.Run("기본 설정이 올바르게 생성되어야 함", func(t *testing.T) {
		// When
		config := NewDefaultConfig()

		// Then
		// 서버 설정 검증
		assert.Equal(t, "localhost", config.Server.Host)
		assert.Equal(t, 8080, config.Server.Port)
		assert.Equal(t, 5*time.Second, config.Server.ReadTimeout)
		assert.Equal(t, 10*time.Second, config.Server.WriteTimeout)
		assert.Equal(t, 120*time.Second, config.Server.IdleTimeout)
		assert.Equal(t, 5*time.Second, config.Server.ShutdownTimeout)

		// 메트릭 설정 검증
		assert.True(t, config.Metrics.Enabled)
		assert.Equal(t, 10*time.Second, config.Metrics.CollectInterval)
		assert.Equal(t, 30*time.Second, config.Metrics.ExportInterval)
		assert.Equal(t, 24*time.Hour, config.Metrics.RetentionPeriod)

		// 데이터베이스 설정 검증
		assert.Equal(t, "postgres", config.Database.Driver)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, 5432, config.Database.Port)
		assert.Equal(t, 5*time.Second, config.Database.ConnectTimeout)
		assert.Equal(t, 25, config.Database.MaxOpenConns)
		assert.Equal(t, 5, config.Database.MaxIdleConns)
		assert.Equal(t, 5*time.Minute, config.Database.ConnMaxLifetime)
	})
}

func TestConfig_Validation(t *testing.T) {
	t.Run("YAML 태그가 올바르게 설정되어야 함", func(t *testing.T) {
		config := &Config{}
		serverType := getStructTags(config, "Server")
		metricsType := getStructTags(config, "Metrics")
		databaseType := getStructTags(config, "Database")

		assert.Equal(t, "server", serverType)
		assert.Equal(t, "metrics", metricsType)
		assert.Equal(t, "database", databaseType)
	})

	t.Run("ServerConfig YAML 태그가 올바르게 설정되어야 함", func(t *testing.T) {
		config := &ServerConfig{}
		assert.Equal(t, "host", getStructTags(config, "Host"))
		assert.Equal(t, "port", getStructTags(config, "Port"))
		assert.Equal(t, "readTimeout", getStructTags(config, "ReadTimeout"))
		assert.Equal(t, "writeTimeout", getStructTags(config, "WriteTimeout"))
		assert.Equal(t, "idleTimeout", getStructTags(config, "IdleTimeout"))
		assert.Equal(t, "shutdownTimeout", getStructTags(config, "ShutdownTimeout"))
	})

	t.Run("MetricsConfig YAML 태그가 올바르게 설정되어야 함", func(t *testing.T) {
		config := &MetricsConfig{}
		assert.Equal(t, "enabled", getStructTags(config, "Enabled"))
		assert.Equal(t, "collectInterval", getStructTags(config, "CollectInterval"))
		assert.Equal(t, "exportInterval", getStructTags(config, "ExportInterval"))
		assert.Equal(t, "retentionPeriod", getStructTags(config, "RetentionPeriod"))
	})

	t.Run("DatabaseConfig YAML 태그가 올바르게 설정되어야 함", func(t *testing.T) {
		config := &DatabaseConfig{}
		assert.Equal(t, "driver", getStructTags(config, "Driver"))
		assert.Equal(t, "host", getStructTags(config, "Host"))
		assert.Equal(t, "port", getStructTags(config, "Port"))
		assert.Equal(t, "name", getStructTags(config, "Name"))
		assert.Equal(t, "user", getStructTags(config, "User"))
		assert.Equal(t, "password", getStructTags(config, "Password"))
		assert.Equal(t, "connectTimeout", getStructTags(config, "ConnectTimeout"))
		assert.Equal(t, "maxOpenConns", getStructTags(config, "MaxOpenConns"))
		assert.Equal(t, "maxIdleConns", getStructTags(config, "MaxIdleConns"))
		assert.Equal(t, "connMaxLifetime", getStructTags(config, "ConnMaxLifetime"))
	})
}

// getStructTags는 구조체 필드의 YAML 태그를 반환합니다.
func getStructTags(v interface{}, fieldName string) string {
	field, ok := reflect.TypeOf(v).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get("yaml")
}
