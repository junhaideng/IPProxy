package spider

import (
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
)

// Deprecated
func SpiderGouBanJia() []model.IP {
	var ips []model.IP
	var url = "http://www.goubanjia.com/"
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	c.OnHTML(`#services > div > div.row > div > div > div > table > tbody`, func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				if i == 3 {
					// 去掉其中的空白
					info = append(info, strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(element.Text, "\n", ""), "\t", ""), " ", ""))
				} else {
					info = append(info, strings.TrimSpace(element.Text))
				}
			})
			ip, port, err := net.SplitHostPort(info[0])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":     err,
					"ip:port": info[0],
					"url":     url,
				}).Error("split host port error")
				return
			}

			minute, err := strconv.Atoi(strings.Split(info[6], "分")[0])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":         err,
					"url":         url,
					"verify-time": info[6],
				}).Error("want to parse verify time error")
				return
			}
			t := time.Now()
			t = t.Add(-time.Minute * time.Duration(minute))
			ips = append(ips, model.IP{
				IP:            ip,
				Port:          port,
				Anonymous:     info[1],
				Location:      info[3],
				VerifyTime:    t,
				Type:          info[2],
				POST:          true,
				ResponseSpeed: -1,
			})
		})
	})
	c.Visit(url)

	return ips
}
