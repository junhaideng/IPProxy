package spider

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SpiderKxDaiLi() []model.IP {

	var ips []model.IP
	var url = "http://www.kxdaili.com/dailiip/%d/%d.html"
	var page = 1
	c := colly.NewCollector()
	c.DetectCharset = true
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"

	c.OnHTML(`div.hot-product-content > table > tbody`, func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			var info []string
			element.ForEach("td", func(i int, element *colly.HTMLElement) {
				info = append(info, element.Text)
			})
			temp := strings.Split(info[4], " ")[0]
			speed, err := strconv.ParseFloat(temp, 64)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":           err,
					"response-time": info[4],
					"url":           url,
				}).Error("parse response time error")
				return
			}

			t, err := parseTime(info[6])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err,
					"time": info[6],
					"url":  url,
				}).Error("parse time error")
				return
			}

			ips = append(ips, model.IP{
				IP:            info[0],
				Port:          info[1],
				Anonymous:     info[2],
				Location:      info[5],
				VerifyTime:    time.Now().Add(-t),
				Type:          strings.ToLower(info[3]),
				POST:          true,
				ResponseSpeed: time.Millisecond * time.Duration(speed*1000),
			})
		})
		page++
		c.Visit(fmt.Sprintf(url, 0, page))
		c.Visit(fmt.Sprintf(url, 1, page))

	})

	// http://www.kxdaili.com/dailiip/2/1.html  -> 普匿 最后的数字为页数
	// http://www.kxdaili.com/dailiip/1/1.html  -> 高匿
	c.Visit(fmt.Sprintf(url, 0, page))
	c.Visit(fmt.Sprintf(url, 1, page))

	return ips
}

// parse time, ignore some case
// (time.Duration, error)
func parseTime(s string) (time.Duration, error) {
	var ERR = errors.New("parse time error, maybe just some ignored case")
	var t time.Duration
	var reg1 = regexp.MustCompile(`(\d+)天(\d+)小时前`)
	var reg2 = regexp.MustCompile(`(\d+)小时(\d+)分前`)
	var reg3 = regexp.MustCompile(`(\d+)分(\d+)秒前`)
	dayIndex := strings.Index(s, "天")
	hourIndex := strings.Index(s, "小时")
	if dayIndex > -1 {
		match := reg1.FindStringSubmatch(s)
		if len(match) == 3 {
			day, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, err
			}
			hour, err := strconv.Atoi(match[2])
			if err != nil {
				return 0, err
			}
			t = time.Hour*24*time.Duration(day) + time.Hour*time.Duration(hour)
			return t, nil
		}
		return 0, ERR

	} else if hourIndex > -1 {
		match := reg2.FindStringSubmatch(s)
		if len(match) == 3 {
			hour, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, err
			}
			minute, err := strconv.Atoi(match[2])
			if err != nil {
				return 0, err
			}
			t = time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute)
			return t, nil
		}
		return 0, ERR

	} else {
		match := reg3.FindStringSubmatch(s)
		if len(match) == 3 {
			minute, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, err
			}
			second, err := strconv.Atoi(match[2])
			if err != nil {
				return 0, err
			}
			t = time.Minute*time.Duration(minute) + time.Second*time.Duration(second)
			return t, nil
		}
		return 0, ERR
	}
}
