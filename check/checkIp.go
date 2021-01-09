// 该程序可以简单的检查ip是否可用
package check

import (
	"fmt"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)


func check(ip model.IP) bool {
	var start time.Time
	start = time.Now()

	proxy := ip.ProxyURL()
	if proxy == nil {
		return false
	}
	var client = http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
			DialContext: (&net.Dialer{
				Timeout: 30 * time.Second,
			}).DialContext,
		},
	}
	// req, err := http.NewRequest("POST", "http://httpbin.org/post", nil)
	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("build request error")
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.WithField("err", err).Error("send http request error")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("共用时: ", time.Now().Sub(start))
		fmt.Println("PASS")
		return true
	} else {
		fmt.Printf("%#v", resp)
		fmt.Println("ERROR")
		return false
	}
}
