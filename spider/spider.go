package spider

import (
	"github.com/junhaideng/IPProxy/dao"
	"github.com/junhaideng/IPProxy/model"
)

var (
	Spiders []Spider
)

// 注册所有的爬虫
func init() {
	var spiders = []func() []model.IP{
		SpiderXiLaDaiLi, Spider66IP, Spider89IP, SpiderJiangXianLi,
		SpiderXiLaDaiLi, SpiderIHuan}

	for _, s := range spiders {
		Spiders = append(Spiders, WrapSpider(s))
	}

}

type Spider interface {
	// 爬取网页，返回IP地址
	Scrap() []model.IP
	// 处理获取到的IP地址
	Start()
}

func WrapSpider(f func() []model.IP) Spider {
	return &spider{f: f}
}

type spider struct {
	f func() []model.IP
}

func (s *spider) Scrap() []model.IP {
	return s.f()
}

func (s *spider) Start() {
	for _, ip := range s.f() {
		go dao.InsertOne(ip)
	}
}
