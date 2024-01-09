package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GrpcConfig    `yaml:"grpc"`
	Rest        Rest          `yaml:"rest"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Rest struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(confPath string) *Config {
	var cnf Config

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file does not exist: %s", err))
	}

	err := cleanenv.ReadConfig(confPath, &cnf)
	if err != nil {
		panic(fmt.Sprintf("err read config: %s", err))
	}

	return &cnf

}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
