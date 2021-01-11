package main

import (
	"github.com/junhaideng/IPProxy/api"
	_ "github.com/junhaideng/IPProxy/conf"
	_ "github.com/junhaideng/IPProxy/log"
)

func main() {
	api.Run()
}
