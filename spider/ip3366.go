package spider

import (
	"fmt"
	"github.com/gocolly/colly"
)

// Deprecated: 此网站似乎全部代理ip不能使用
func SpiderIP3366() {
	var url = "http://www.ip3366.net/?page=%d"
	pageNum := 1
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	c.OnHTML(`#list > table > tbody`, func(e *colly.HTMLElement) {
		fmt.Println("正在爬取第", pageNum, "页")
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				info = append(info, element.Text)
			})
			fmt.Println(info)
		})

		// 只爬取10页内容
		if pageNum == 10 {
			return
		}
		pageNum++
		c.Visit(fmt.Sprintf(url, pageNum))

	})

	// 页链接不太好构造，需要在爬取的时候进行获取
	err := c.Visit(fmt.Sprintf(url, pageNum))
	if err != nil {
		fmt.Println(err)
	}
}
