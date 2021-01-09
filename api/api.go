package api

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/junhaideng/IPProxy/dao"
	"github.com/junhaideng/IPProxy/schedule"
	"github.com/junhaideng/IPProxy/spider"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"strconv"
	"time"
)

var start = time.Now()

// 获取代理ip
func getIp(c *gin.Context){
	var limit int64
	num, ok := c.GetQuery("num")
	if !ok {
		limit = 0
	}
	limit, err := strconv.ParseInt(num, 10, 64)
	if err != nil{
		limit = 0
	}
	ips, err := dao.GetLimit(limit)
	if err != nil{
		logrus.Error("get limit proxy ip error: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,  // 表示失败
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"num": len(ips),
		"data": ips,
	})
}

func healthy(c *gin.Context){
	ips, err := dao.GetAll()
	if err != nil{
		logrus.Error("get all proxy ip error: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1, 
			"msg": err,
		})
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"ipNum": len(ips),
		"totalTime": time.Since(start).String(),
	})
}

var router *gin.Engine

func Run() {
	scheduler := schedule.NewScheduler()
	for _, s := range spider.Spiders {
		scheduler.Register(s)
	}
	scheduler.Start()

	router = gin.Default()
	router.GET("/get_ip", getIp)
	router.GET("/healthy", healthy)

	pprof.Register(router)

	host := viper.GetString("api.host")
	port := viper.GetString("api.port")
	router.Run(net.JoinHostPort(host, port))
}
