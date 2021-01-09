package schedule

import (
	"github.com/junhaideng/IPProxy/spider"
)

type Register interface{
	// 注册爬虫
	Register(spider.Spider) error
	// 注销
	Unregister(spider.Spider) error 
	// 清除所有
	Clear() error
}
