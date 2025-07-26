package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Debug        bool          `mapstructure:"debug"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConnections  int32         `mapstructure:"max_connections"`
	MinConnections  int32         `mapstructure:"min_connections"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

type JWTConfig struct {
	SecretKey string        `mapstructure:"secret_key"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// Load загружает конфигурацию из файла и environment variables
func Load(configPath string) (*Config, error) {
	// Инициализация Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Добавляем пути поиска конфига
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Чтение конфигурационного файла
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Автоматическая привязка environment variables
	viper.SetEnvPrefix("CHAT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Парсинг конфигурации
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Валидация конфигурации
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверка сервера
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// Проверка базы данных
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}
	if c.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	// Проверка JWT
	if c.JWT.SecretKey == "" {
		return fmt.Errorf("jwt secret key is required")
	}

	// Проверка логгера
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true, "panic": true,
	}
	if !validLevels[c.Logger.Level] {
		return fmt.Errorf("invalid logger level: %s", c.Logger.Level)
	}

	validFormats := map[string]bool{"json": true, "text": true}
	if !validFormats[c.Logger.Format] {
		return fmt.Errorf("invalid logger format: %s", c.Logger.Format)
	}

	// Проверка приложения
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("invalid environment: %s", c.App.Environment)
	}

	return nil
}

// GetServerAddress возвращает адрес сервера в формате host:port
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDatabaseDSN возвращает строку подключения к базе данных
func (c *Config) GetDatabaseDSN() string {
	// Формат: postgres://username:password@host:port/database?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// IsDevelopment проверяет, является ли окружение development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction проверяет, является ли окружение production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// Print выводит текущую конфигурацию (без секретов)
func (c *Config) Print() {
	fmt.Printf("=== Application Configuration ===\n")
	fmt.Printf("App: %s v%s (%s)\n", c.App.Name, c.App.Version, c.App.Environment)
	fmt.Printf("Server: %s\n", c.GetServerAddress())
	fmt.Printf("Database: %s@%s:%d/%s\n", c.Database.Username, c.Database.Host, c.Database.Port, c.Database.Name)
	fmt.Printf("JWT Expires: %v\n", c.JWT.ExpiresIn)
	fmt.Printf("Logger: %s level, %s format\n", c.Logger.Level, c.Logger.Format)
	fmt.Printf("================================\n")
}
