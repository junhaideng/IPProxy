package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"time"
)

func Spider66IP() []model.IP {
	const selector = "#main > div.containerbox.boxindex > div.layui-row.layui-col-space15 > div:nth-child(1) > table > tbody"
	const maxPageNum = 200

	var ips = make([]model.IP, 0, 2000)
	var url = "http://www.66ip.cn/%d.html"
	c := colly.NewCollector()
	c.DetectCharset = true

	var pageNum = 1
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		//println("正在访问第", pageNum, "页")
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			if i == 0 {
				return
			}
			var info []string
			element.ForEach("td", func(j int, element *colly.HTMLElement) {
				// 这里就是我们需要的内容
				info = append(info, element.Text)
			})
			//fmt.Println(info)
			t, err := time.Parse("2006年01月02日15时 验证", info[4])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err,
					"time": info[4],
					"url":  url,
				}).Error("parse time error")
				return
			}
			ips = append(ips, model.IP{
				IP:            info[0],
				Port:          info[1],
				Location:      info[2],
				Anonymous:     info[3],
				Type:          "http,https",
				VerifyTime:    t,
				POST:          true, // 默认取true
				ResponseSpeed: -1,
			})
		})

		// 继续访问下一页
		pageNum++
		if pageNum > maxPageNum {
			return
		}
		c.Visit(fmt.Sprintf(url, pageNum))
	})

	c.Visit(fmt.Sprintf(url, pageNum))

	return ips
}
