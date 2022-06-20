package config

import (
	"github.com/spf13/viper"
	"log"
)

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DB       string `mapstructure:"DB_DB"`
}

// LoadConfig 讀取.env環境變數檔
func LoadConfig[T any](myPath string, fileName string) (Config T) {
	// 若有同名環境變量則使用環境變量
	viper.AddConfigPath(myPath)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("can not load config: " + err.Error())
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal("can not load config: " + err.Error())
	}
	return
}
