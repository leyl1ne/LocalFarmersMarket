package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string         `yaml:"env" env-default:"local"`
	Database       DatabaseConfig `yaml:"database"`
	GRPC           GRPCConfig     `yaml:"grpc"`
	MigrationsPath string         `yaml:"migrations_path"`
	TokenTTL       time.Duration  `yaml:"token_ttl" env-default:"1h"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver" env-default:"postgres"` //sqlite/postgres
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"dbname" env:"DB_NAME"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type GRPCConfig struct {
	Port    int           `yaml:"proto"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	if cfg.Database.Driver == "postgres" {
		if cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.DBName == "" {
			panic("database credentials are required for postgres")
		}
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
