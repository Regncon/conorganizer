package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Secrets struct {
	AWS_ENDPOINT_URL_S3   string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	BUCKET_NAME           string
	DB_PREFIX             string
}

type S3Config struct {
	S3_endpoint   string
	S3_access_id  string
	S3_access_key string
	S3_region     string
	S3_bucket     string
	S3_prefix     string
}

type Config struct {
	S3      S3Config
	EnvPath string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadEnv(envPath string) error {
	// Check if file exists
	fileInfo, err := os.Stat(envPath)
	if err != nil {
		return fmt.Errorf("LoadEnv failed: %f", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("LoadEnv failed: %s is a dir", fileInfo.Name())
	}

	// Read .env file and set values
	secrets, err := godotenv.Read(envPath)
	if err != nil {
		return fmt.Errorf("godotenv failed reading .env file: %f", err)
	}
	for k, v := range secrets {
		_ = os.Setenv(k, v)
	}

	// Update S3 config
	c.S3.S3_access_id = os.Getenv("AWS_ACCESS_KEY_ID")
	c.S3.S3_access_key = os.Getenv("AWS_SECRET_ACCESS_KEY")
	c.S3.S3_endpoint = os.Getenv("AWS_ENDPOINT_URL_S3")
	c.S3.S3_bucket = os.Getenv("BUCKET_NAME ")
	c.S3.S3_region = os.Getenv("AWS_REGION")
	c.S3.S3_prefix = os.Getenv("DB_PREFIX")

	// Update env path in config
	test, _ := filepath.Abs(envPath)
	c.EnvPath = envPath
	fmt.Printf("Loaded secrets from: %s", test)

	return nil
}
