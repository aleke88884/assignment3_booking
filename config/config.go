package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// StorageConfig holds storage configuration (S3/MinIO/GCS)
type StorageConfig struct {
	Type            string // "s3", "minio", "local"
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Region          string
	UseSSL          bool
	PublicURL       string
}

// Load loads configuration from environment or defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "smartbooking"),
		},
		Storage: StorageConfig{
			Type:            getEnv("STORAGE_TYPE", "minio"),
			Endpoint:        getEnv("STORAGE_ENDPOINT", "http://minio:9000"),
			AccessKeyID:     getEnv("STORAGE_ACCESS_KEY", "minioadmin"),
			SecretAccessKey: getEnv("STORAGE_SECRET_KEY", "minioadmin"),
			BucketName:      getEnv("STORAGE_BUCKET", "smartbooking"),
			Region:          getEnv("STORAGE_REGION", "us-east-1"),
			UseSSL:          getEnvAsBool("STORAGE_USE_SSL", false),
			PublicURL:       getEnv("STORAGE_PUBLIC_URL", "http://localhost:9000"),
		},
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt получает целочисленное значение переменной окружения
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool получает булевое значение переменной окружения
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
