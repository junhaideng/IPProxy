package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func SpiderJiangXianLi() []model.IP {
	const selector = `body > div.layui-layout.layui-layout-admin > div.layui-row > div.layui-col-md9.ip-tables > div.layui-form > table > tbody`
	const maxPageNum = 10

	var ips []model.IP
	c := colly.NewCollector()

	var url = "https://ip.jiangxianli.com/?page=%d"
	pageNum := 1

	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEachWithBreak("td", func(i int, element *colly.HTMLElement) bool {
				// 去掉最后的无关信息
				if i > 9 {
					return false
				}
				info = append(info, strings.TrimSpace(element.Text))
				return true
			})
			if len(info) < 10 {
				return
			}
			t, err := time.Parse("2006-01-02 15:04:05", info[9])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err,
					"time": info[9],
					"url":  url,
				}).Error("parse time error")
				return
			}
			ips = append(ips, model.IP{
				IP:            info[0],
				Port:          info[1],
				Anonymous:     info[2],
				Location:      info[4],
				VerifyTime:    t,
				Type:          strings.ToLower(info[3]),
				POST:          true,
				ResponseSpeed: -1,
			})
		})

		pageNum++
		if pageNum > maxPageNum {
			return
		}
		c.Visit(fmt.Sprintf(url, pageNum))
	})

	c.Visit(fmt.Sprintf(url, pageNum))

	return ips
}
