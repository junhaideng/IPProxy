package main

import (
	"github.com/junhaideng/IPProxy/api"
	"github.com/junhaideng/IPProxy/conf"
	"github.com/junhaideng/IPProxy/dao"
	"github.com/junhaideng/IPProxy/log"
)

func main() {
	conf.Init()
	log.Init()
	dao.Init()
	api.Run()
}
