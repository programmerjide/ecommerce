// Description: This file contains the configuration structures and loading logic for the e-commerce application.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port    string `default:"8080"`
	GinMode string `default:"debug"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `default:"localhost"`
	Port     string `default:"5432"`
	User     string `default:"user"`
	Password string `default:"password"`
	DBName   string `default:"ecommerce"`
	SSLMode  string `default:"disable"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret              string `default:"your_secret_key"`
	ExpiresIn           time.Duration
	RefreshTokenExpires time.Duration
}

// AWSConfig holds AWS-related configuration
type AWSConfig struct {
	Region          string `default:"us-east-1"`
	AccessKeyID     string `default:"your_access_key_id"`
	SecretAccessKey string `default:"your_secret_access_key"`
	S3BucketName    string `default:"your_bucket_name"`
	S3Endpoint      string `default:"http://localstack:4566"`
}

// UploadConfig holds file upload-related configuration
type UploadConfig struct {
	Path        string `default:"./uploads"`
	MaxFileSize int64  `default:"1048576"`
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "ecommerce"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", "your_secret_key"),
			ExpiresIn:           time.Duration(getEnvAsInt("JWT_EXPIRES_IN", 15)) * time.Minute,
			RefreshTokenExpires: time.Duration(getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 7)) * 24 * time.Hour,
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "your_access_key_id"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "your_secret_access_key"),
			S3BucketName:    getEnv("AWS_S3_BUCKET_NAME", "your_bucket_name"),
			S3Endpoint:      getEnv("AWS_S3_ENDPOINT", "http://localstack:4566"),
		},
		Upload: UploadConfig{
			Path:        getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize: int64(getEnvAsInt("UPLOAD_MAX_FILE_SIZE", 1048576)),
		},
	}
	return cfg, nil
}

func getEnvAsInt(s string, i int) int {
	if value := os.Getenv(s); value != "" {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return i
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
