package model

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type IP struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	// 代理ip地址
	IP string `bson:"ip,omitempty" json:"ip,omitempty"`
	// 代理ip端口
	Port string `bson:"port,omitempty" json:"port,omitempty"`
	// 高匿，普匿，透明，未知等
	Anonymous string `bson:"anonymous,omitempty" json:"anonymous,omitempty"`
	// ip地址所在地理位置
	Location string `bson:"location,omitempty" json:"location,omitempty"`
	// 上次验证时间
	VerifyTime time.Time `bson:"verifyTime,omitempty" json:"verify_time,omitempty"`
	// 代理ip类型，HTTP, HTTPS
	Type string `bson:"type,omitempty" json:"type,omitempty"`
	// 是否支持POST请求
	POST bool `bson:"post,omitempty" json:"post,omitempty"`
	// 响应速度, -1 表示爬取的网页中没有对应的数据，并且本地还没有进行检测
	ResponseSpeed time.Duration `bson:"response_speed,omitempty" json:"response_speed,omitempty"`
}

func(ip IP) URL() ([]*url.URL, error){
	if ip.Type == "" {
		ip.Type = "http"
	}
	typs := strings.Split(ip.Type, ",")
	var uris []*url.URL
	var e error
	for _, typ := range typs{
		uri, err := url.Parse(typ + "://" + ip.IP + ":" + ip.Port)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
				"ip":  ip,
				"uri": uri,
			}).Error("ip解析错误")
			e = err
			continue
		}
		uris = append(uris, uri)
	}

	return uris, e
}

func (ip IP) ProxyURL() []func(r *http.Request) (*url.URL, error) {
	uris, err := ip.URL()
	if err != nil{
		logrus.WithFields(logrus.Fields{
			"err": err,
			"uri": uris,
		}).Error("创建代理ip错误")
		return nil
	}
	var fs []func(r *http.Request)(*url.URL, error)
	for _, uri := range uris{
		fs = append(fs, http.ProxyURL(uri))
	}
	return fs
}
