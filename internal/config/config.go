package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Db         `yaml:"db"`
}

type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:8080"`
	Timeout     int    `yaml:"timeout" env-default:"20"`
	BearerToken string `yaml:"bearer_token" env-default:"2"`
	DbPort      string `yaml:"db_port"`
}

type Db struct {
	Option       string `yaml:"option"`
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	ExternalPort string `yaml:"port"`
	InternalPort string `yaml:"port"`
	NameDb       string `yaml:"name_db"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
}

func MustInit(configPath string) *Config {
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Env: MustGetEnv("ENV"),
		HTTPServer: HTTPServer{
			Address:     MustGetEnv("HTTP_SERVER_HOST") + ":" + MustGetEnv("HTTP_SERVER_PORT"),
			Timeout:     MustGetEnvAsInt("HTTP_TIMEOUT"),
			BearerToken: MustGetEnv("HTTP_ADMIN_BEARER_TOKEN"),
			DbPort:      MustGetEnv("HTTP_DB_PORT"),
		},
		Db: Db{
			Option:       MustGetEnv("DB_OPTION"),
			Driver:       MustGetEnv("DB_DRIVER"),
			Host:         MustGetEnv("DB_HOST"),
			ExternalPort: MustGetEnv("DB_EXTERNAL_PORT"),
			InternalPort: MustGetEnv("DB_INTERNAL_PORT"),
			NameDb:       MustGetEnv("DB_NAME"),
			User:         MustGetEnv("DB_USER"),
			Password:     MustGetEnv("DB_PASSWORD"),
		},
	}
}

func Path(workDir string, envFile string) string {
	if envFile == "" {
		return filepath.Join(workDir, ".env")
	}

	return filepath.Join(workDir, envFile)
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("no variable in env: %s", key)
	}
	return value
}

func MustGetEnvAsInt(name string) int {
	valueStr := MustGetEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return -1
}

func GetConfigPathFromTest(envFile string) string {
	projectRoot, err := findProjectRoot()
	if err != nil {
		log.Fatalf("failed to get config path: %s", err)
	}
	return Path(projectRoot, envFile)
}

func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// TODO: пока так
	re := regexp.MustCompile(`^(.*?/go-monolite)`)
	match := re.FindStringSubmatch(currentDir)
	if len(match) < 2 {
		return "", fmt.Errorf("project root not found in path: %s", currentDir)
	}

	return match[1], nil
}
