package main

import (
	"flag"
	"github.com/junhaideng/IPProxy/api"
	_ "github.com/junhaideng/IPProxy/conf"
	_ "github.com/junhaideng/IPProxy/log"
	"github.com/spf13/viper"
)

var conf string

func init(){
	flag.StringVar(&conf, "conf", "", "configuration file")
}

func main() {
	flag.Parse()
	if conf != ""{
		viper.AddConfigPath(conf)
	}
	api.Run()
}
