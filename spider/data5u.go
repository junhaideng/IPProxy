package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
)

// Deprecated: can not use now 2020/1/9
// 这里只有20个免费ip地址
func SpiderData5u() []model.IP {
	var ips []model.IP
	var url = "http://www.data5u.com/"
	c := colly.NewCollector()
	c.AllowURLRevisit = true
	c.DetectCharset = true

	c.OnHTML(`.wlist > ul >li:nth-child(2)`, func(e *colly.HTMLElement) {
		e.ForEach("ul", func(i int, element *colly.HTMLElement) {
			// 第一个是表头不需要记录
			if i == 0 {
				return
			}
			// ip, port, anonymous, type, location, response time, verifyTime
			var info []string
			element.ForEach("span", func(j int, element *colly.HTMLElement) {
				// 这里就是我们需要的内容
				info = append(info, element.Text)
			})

			//TODO 插入数据库
			fmt.Println(info)
		})
	})

	c.Visit(url)
	fmt.Println(ips)
	return ips
}
