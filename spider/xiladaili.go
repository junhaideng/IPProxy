package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

func SpiderXiLaDaiLi() []model.IP {
	var ips []model.IP
	// 高匿
	var high = "http://www.xiladaili.com/high/%d/"
	// http代理
	var http = "http://www.xiladaili.com/http/%d/"
	// https代理
	var https = "http://www.xiladaili.com/https/%d/"
	pageNum := 1
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	c.OnHTML(`table > tbody`, func(e *colly.HTMLElement) {
		//fmt.Println("爬取第", pageNum, "页")
		if e.DOM.Find("tr").Index() < 0 {
			//fmt.Println("没有发现ip")
			return
		}
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				info = append(info, element.Text)
			})
			ip, port, _ := net.SplitHostPort(info[0])

			t, err := time.Parse("2006年1月2日 15:4", info[6])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":         err,
					"verify-time": info[6],
				}).Error("parse verify time error")
				return
			}

			ips = append(ips, model.IP{
				IP:            ip,
				Port:          port,
				Anonymous:     info[2],
				Location:      info[3],
				VerifyTime:    t,
				Type:          strings.ToLower(strings.TrimRight(info[1], "代理")),
				POST:          true,
				ResponseSpeed: -1,
			})

		})
		// 最多爬取40页，自行修改
		if pageNum > 40 {
			return
		}
		pageNum++
		c.Visit(fmt.Sprintf(high, pageNum))
		time.Sleep(time.Second)

		c.Visit(fmt.Sprintf(http, pageNum))
		time.Sleep(time.Second)

		c.Visit(fmt.Sprintf(https, pageNum))

	})

	c.Visit(fmt.Sprintf(high, pageNum))
	c.Visit(fmt.Sprintf(http, pageNum))
	c.Visit(fmt.Sprintf(https, pageNum))

	return ips
}
