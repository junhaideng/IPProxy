// 该程序可以简单的检查ip是否可用
package check

import (
	"fmt"
	"github.com/junhaideng/IPProxy/dao"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// check ip with specified proxy
func checkIP(proxy func(r *http.Request)(*url.URL, error), ip model.IP, wg *sync.WaitGroup, limit <- chan struct{}){
	defer func() {
		wg.Done()
		<- limit
	}()
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
		return
	}
	time.Sleep(time.Duration(rand.Intn(10))*time.Second)
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithField("err", err).Error("send http request error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		ip.ResponseSpeed = time.Now().Sub(start)
		ip.VerifyTime = time.Now()
		err := dao.Update(bson.M{"_id": ip.ID.Hex()}, bson.M{"$set": bson.M{"response_speed":ip.ResponseSpeed, "verifyTime":ip.VerifyTime}})
		if err != nil{
			logrus.WithFields(logrus.Fields{
				"err": err,
				"ip": ip.IP,
				"port": ip.Port,
				"id": ip.ID.Hex(),
			}).Error("update ip error")
			return
		}
		logrus.WithFields(logrus.Fields{
			"ip": ip.IP,
			"port": ip.Port,
			"id": ip.ID.Hex(),
		}).Info("update ip success")
	} else {
		// 这个ip不可用
		err := dao.Delete(bson.M{
			"_id": ip.ID.Hex(),
		})
		if err != nil{
			logrus.WithFields(logrus.Fields{
				"err": err,
				"ip": ip.IP,
				"port": ip.Port,
				"id": ip.ID.Hex(),
			}).Error("delete ip error")
			return
		}
		logrus.WithFields(logrus.Fields{
			"ip": ip.IP,
			"port": ip.Port,
			"id": ip.ID.Hex(),
		}).Info("delete ip success")
	}
}

func CheckSingleIP(ip model.IP, wg *sync.WaitGroup) {
	defer wg.Done()
	proxy := ip.ProxyURL()

	if proxy == nil {
		return
	}
	var limit = make(chan struct{}, 100)
	var lg sync.WaitGroup
	for _, p := range proxy{
		lg.Add(1)
		limit <- struct{}{}
		go checkIP(p, ip, &lg, limit)
	}
	lg.Wait()
	
}


func CheckIP(){
	ips, err := dao.GetAll()
	var wg sync.WaitGroup
	if err != nil{
		logrus.WithField("err", err).Error("get all ips error")
		return
	}
	num := 1
	for _, ip := range ips{
		wg.Add(1)
		fmt.Printf("checking number %d ip: %s:%s\n", num, ip.IP, ip.Port)
		num ++
		go CheckSingleIP(ip, &wg)
	}
	wg.Wait()
}