// Package check 检查ip是否可用
package check

import (
	"github.com/junhaideng/IPProxy/dao"
	"github.com/junhaideng/IPProxy/model"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

var pool *ants.PoolWithFunc

func init() {
	var err error
	pool, err = ants.NewPoolWithFunc(100, checkIP)
	if err != nil {
		panic(err)
	}
}

type Arg struct {
	Proxy func(r *http.Request) (*url.URL, error)
	IP    model.IP
}

// check ip with specified proxy
func checkIP(arg interface{}) {
	argument := arg.(Arg)
	proxy, ip := argument.Proxy, argument.IP
	start := time.Now()
	var client = http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
			DialContext: (&net.Dialer{
				Timeout: time.Minute,
			}).DialContext,
		},
	}
	// req, err := http.NewRequest("POST", "http://httpbin.org/post", nil)
	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("build request error")
		deleteIP(ip)
		return
	}
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithField("err", err).Error("send http request error")
		deleteIP(ip)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		ip.ResponseSpeed = time.Now().Sub(start)
		ip.VerifyTime = time.Now()
		err := dao.Update(bson.M{"_id": ip.ID}, bson.M{"$set": bson.M{"response_speed": ip.ResponseSpeed, "verifyTime": ip.VerifyTime}})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":  err,
				"ip":   ip.IP,
				"port": ip.Port,
				"id":   ip.ID.Hex(),
			}).Error("update ip error")
			deleteIP(ip)
			return
		}
		logrus.WithFields(logrus.Fields{
			"ip":   ip.IP,
			"port": ip.Port,
			"id":   ip.ID.Hex(),
		}).Info("update ip success")
	} else {
		// 这个ip不可用
		deleteIP(ip)
	}
}

func deleteIP(ip model.IP) {
	err := dao.Delete(bson.M{
		"_id": ip.ID.Hex(),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"ip":   ip.IP,
			"port": ip.Port,
			"id":   ip.ID.Hex(),
		}).Error("delete ip error")
		return
	}
	logrus.WithFields(logrus.Fields{
		"ip":   ip.IP,
		"port": ip.Port,
		"id":   ip.ID.Hex(),
	}).Info("delete ip success")
}

func CheckSingleIP(ip model.IP) {
	proxy := ip.ProxyURL()

	if proxy == nil {
		return
	}
	for _, p := range proxy {
		err := pool.Invoke(Arg{p, ip})
		if err != nil {
			logrus.Error(err)
			continue
		}
	}
}

func CheckIP() {
	ips, err := dao.GetAll()
	if err != nil {
		logrus.WithField("err", err).Error("get all ips error")
		return
	}
	for _, ip := range ips {
		CheckSingleIP(ip)
	}
}
