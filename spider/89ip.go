package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func Spider89IP() []model.IP {
	var ips []model.IP
	var url = "https://www.89ip.cn/index_%d.html"
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	pageNum := 1
	c.OnHTML(`div.layui-form > table > tbody`, func(e *colly.HTMLElement) {
		//fmt.Printf("正在爬取第 %d 页\n", pageNum)
		// 当前页中没有ip地址
		if e.DOM.Find("tr").Index() < 0 {
			//fmt.Println("没有发现ip")

			return
		}
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			// ip, port, location, verifyTime
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				info = append(info, strings.TrimSpace(element.Text))
			})

			t, err := time.Parse("2006/01/02 15:04:05", info[4])
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
				Anonymous:     "未知",
				Location:      info[2],
				Type:          "http,https",
				VerifyTime:    t,
				POST:          true,
				ResponseSpeed: -1,
			})
		})

		pageNum++
		c.Visit(fmt.Sprintf(url, pageNum))
	})

	c.Visit(fmt.Sprintf(url, pageNum))

	return ips
}
