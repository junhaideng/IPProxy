package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func SpiderJiangXianLi() []model.IP {
	var ips []model.IP
	c := colly.NewCollector()

	var url = "https://ip.jiangxianli.com/?page=%d"
	pageNum := 1

	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	c.OnHTML(`body > div.layui-layout.layui-layout-admin > div.layui-row > div.layui-col-md9.ip-tables > div.layui-form > table > tbody`, func(e *colly.HTMLElement) {
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
			t, err := time.Parse("2006-01-02 15:04:05", info[9])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err,
					"time": info[9],
					"url":  url,
				}).Error("parse time error")
				return
			}
			var speed time.Duration
			tempIndex := strings.Index(info[7], "毫秒")
			if tempIndex > -1 {
				s, err := strconv.Atoi(info[7][:tempIndex])
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"err":           err,
						"response-time": info[7],
						"url":           url,
						"type":          "millisecond",
					}).Error("parse response time error")
					return
				}
				speed = time.Millisecond * time.Duration(s)
			} else {
				index := strings.Index(info[7], "秒")
				if index > -1 {
					s, err := strconv.ParseFloat(info[7][:index], 64)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"err":           err,
							"response-time": info[7],
							"url":           url,
							"type":          "second",
						}).Error("parse response time error")
						return
					}
					speed = time.Millisecond * time.Duration(s*1000)
				} else {
					logrus.WithFields(logrus.Fields{
						"response-time": info[7],
						"url":           url,
					}).Error("can not parse response time")
					return
				}
			}
			ips = append(ips, model.IP{
				IP:            info[0],
				Port:          info[1],
				Anonymous:     info[2],
				Location:      info[4],
				VerifyTime:    t,
				Type:          strings.ToLower(info[3]),
				POST:          true,
				ResponseSpeed: speed,
			})
		})

		pageNum++
		c.Visit(fmt.Sprintf(url, pageNum))
	})

	c.Visit(fmt.Sprintf(url, pageNum))

	return ips
}
