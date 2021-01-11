package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	fmt.Println("initialize configuration")
	viper.AddConfigPath("conf")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	setDefault()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func setDefault() {
	// 数据库配置
	viper.SetDefault("database.mongodb.port", "27017")
	viper.SetDefault("database.mongodb.host", "127.0.0.1")
	viper.SetDefault("database.mongodb.db", "IP")

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.mode", "console")
	viper.SetDefault("log.filename", "proxy.log")
	viper.SetDefault("log.max-size", 5)

	// 定时任务配置
	viper.SetDefault("schedule.interval", 1)

	// api 接口配置
	viper.SetDefault("api.host", "127.0.0.1")
	viper.SetDefault("api.port", "8000")
}
