package spider

import (
	"encoding/json"
	"fmt"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	FreeIP    []Info `json:"freeips"`
	PrivateIp []Info `json:"privateips"`
}

type Info struct {
	Speed      string `json:"connect_speed"`
	Country    string `json:"country"`
	Ip         string `json:"ip"`
	Level      string `json:"level"`
	Port       string `json:"port"`
	Type       string `json:"type"`
	VerifyTime string `json:"verify_time"`
}

// Deprecated: 似乎在维护，目前死了
func SpiderSuperFastIP() []model.IP {
	var ips []model.IP
	var page = 1
	for {
		url := fmt.Sprintf("https://api.superfastip.com/ip/freeip?page=%d", page)
		page++
		resp, err := http.Get(url)
		if err != nil {
			logrus.WithField("url", url).Error("get response error")
			break
		}
		var res Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			logrus.WithField("err", err).Error("decode response body error: ")
			continue
		}
		resp.Body.Close()
		if len(res.FreeIP) == 0 {
			break
		}
		for _, ip := range res.FreeIP {
			t, err := time.Parse("2006-01-02 15:04:05", ip.VerifyTime)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":         err,
					"verify-time": ip.VerifyTime,
					"url":         url,
				}).Error("parse time error")
				continue
			}

			ips = append(ips, model.IP{
				IP:            ip.Ip,
				Port:          ip.Port,
				Anonymous:     ip.Level,
				Location:      ip.Country,
				VerifyTime:    t,
				Type:          strings.ToLower(strings.Replace(ip.Type, "/", ",", -1)),
				POST:          true,
				ResponseSpeed: -1,
			})
		}
	}
	return ips
}
