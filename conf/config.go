package conf

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Sever   Server  `mapstructure:"server"`
	Redis   Redis   `mapstructure:"redis"`
	Logger  Logger  `mapstructure:"logger"`
	Adapter Adapter `mapstructure:"adapter"`
}

type Server struct {
	Name          string `mapstructure:"name"`
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	SecretKey     string `mapstructure:"secretKey"`
	RetryAfter    int64  `mapstructure:"retryAfter"`
	BanDuration   int64  `mapstructure:"banDuration"`
	ValidateLimit int64  `mapstructure:"validateLimit"`
	WriteTimeout  int64  `mapstructure:"writeTimeout"`
	ReadTimeout   int64  `mapstructure:"readTimeout"`
}

type Adapter struct {
	URL      string `mapstructure:"url"`
	Login    string `mapstructure:"login"`
	Password string `mapstructure:"password"`
	Text     string `mapstructure:"text"`
	Account  string `mapstructure:"account"`
	Timeout  int64  `mapstructure:"timeout"`
}

type Logger struct {
	WriteToFile bool   `mapstructure:"writeToFile"`
	Format      string `mapstructure:"format"`
}

type Redis struct {
	RedisAddr    string `mapstructure:"redisAddr"`
	MinIdleConns int    `mapstructure:"minIdleConns"`
	PoolSize     int    `mapstructure:"poolSize"`
	PoolTimeout  int    `mapstructure:"poolTimeout"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
}

func InitConfigs() *Config {
	path := fetchConfigPath()

	if path == "" {
		log.Fatalf("config path is empty %s", path)
	}

	var config Config

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read configs %s", err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to unmarshal configs %s", err.Error())
	}

	return &config
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
