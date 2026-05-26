package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	GinMode    string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	JWTSecret     string
	GiteaURL      string
	GiteaAdminUser  string
	GiteaAdminPassword string

	NotifyWeComWebhook  string
	NotifyDingTalkWebhook string
	NotifyFeishuWebhook string

	PrometheusURL string

	HarborURL          string
	HarborAdminUser   string
	HarborAdminPassword string
	HarborProject      string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() *Config {
	return &Config{
		ServerPort:       getEnvAsInt("BACKEND_PORT", 8080),
		GinMode:          getEnv("GIN_MODE", "debug"),

		DBHost:           getEnv("DB_HOST", "postgres"),
		DBPort:           getEnvAsInt("DB_PORT", 5432),
		DBUser:           getEnv("DB_USER", "devops"),
		DBPassword:       getEnv("DB_PASS", "devops123"),
		DBName:           getEnv("DB_NAME", "devops_star"),

		RedisHost:         getEnv("REDIS_HOST", "redis"),
		RedisPort:         getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword:     getEnv("REDIS_PASS", "devops123"),
		RedisDB:          getEnvAsInt("REDIS_DB", 0),

		JWTSecret:         getEnv("JWT_SECRET", "devops-star-jwt-secret"),
		GiteaURL:          getEnv("GITEA_URL", "http://gitea:3000"),
		GiteaAdminUser:    getEnv("GITEA_ADMIN_USER", "admin"),
		GiteaAdminPassword: getEnv("GITEA_ADMIN_PASSWORD", "admin123"),

		NotifyWeComWebhook:  getEnv("NOTIFY_WECOM_WEBHOOK", ""),
		NotifyDingTalkWebhook: getEnv("NOTIFY_DINGTALK_WEBHOOK", ""),
		NotifyFeishuWebhook: getEnv("NOTIFY_FEISHU_WEBHOOK", ""),

		PrometheusURL: getEnv("PROMETHEUS_URL", "http://prometheus:9090"),

		HarborURL:          getEnv("HARBOR_URL", "http://harbor:8080"),
		HarborAdminUser:   getEnv("HARBOR_ADMIN_USER", "admin"),
		HarborAdminPassword: getEnv("HARBOR_ADMIN_PASSWORD", "Harbor12345"),
		HarborProject:      getEnv("HARBOR_PROJECT", "devops-star"),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 465),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
