package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"time"
)

// TODO
func SpiderIHuan() []model.IP {
	var ips []model.IP
	var url = "https://ip.ihuan.me"
	var page string
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})
	// c.OnRequest(func(r *colly.Request){
	// 	r.Headers.Set("", "")
	// })

	// 获取到下一页的内容
	c.OnHTML(`div.col-md-10 > nav > ul > li:nth-child(8) > a`, func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		page = e.Attr("href")
	})

	c.OnHTML(`div.table-responsive > table > tbody`, func(e *colly.HTMLElement) {
		fmt.Println("page: ", page)

		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				info = append(info, element.Text)
			})
			fmt.Println(info)

			ips = append(ips, model.IP{
				IP:            info[0],
				Port:          info[1],
				Anonymous:     info[5],
				Location:      info[2],
				VerifyTime:    time.Time{},
				Type:          "",
				POST:          false,
				ResponseSpeed: -1,
			})
		})

		fmt.Println(page)
		c.Visit(url + "?page=" + page)

	})

	c.Visit(url)

	fmt.Println(ips)
	return ips

}
