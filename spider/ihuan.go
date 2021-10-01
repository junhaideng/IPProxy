package spider

import (
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"time"
)

func SpiderIHuan() []model.IP {
	// 选择下一页的内容
	const pageSelector = "div.col-md-10 > nav > ul.pagination > li:last-child > a"
	// ip 所在的元素选择器
	const selector = "div.table-responsive > table > tbody"
	const maxPageNum = 40

	var ips []model.IP
	var url = "https://ip.ihuan.me"
	page := 0
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
	
	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	// 获取到下一页的内容
	c.OnHTML(pageSelector, func(e *colly.HTMLElement) {
		page++
		// 限制页数
		if page > maxPageNum {
			return
		}
		c.Visit(url + e.Attr("href"))
	})

	c.OnHTML(selector, func(e *colly.HTMLElement) {

		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				info = append(info, element.Text)
			})

			// 支持https, http
			var typ string
			if info[4] == "支持" {
				typ = "https,http"
			} else {
				typ = "http"
			}

			post := false
			if info[5] == "支持" {
				post = true
			}

			ips = append(ips, model.IP{
				IP:            info[0],
				Port:          info[1],
				Anonymous:     info[6],
				Location:      info[2],
				VerifyTime:    time.Time{},
				Type:          typ,
				POST:          post,
				ResponseSpeed: -1,
			})
		})

	})
	c.Visit(url)

	return ips
}
