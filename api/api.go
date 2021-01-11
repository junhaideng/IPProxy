package api

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/junhaideng/IPProxy/dao"
	"github.com/junhaideng/IPProxy/schedule"
	"github.com/junhaideng/IPProxy/spider"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"net"
	"net/http"
	"time"
)

var start = time.Now()

type param struct {
	Num int64 `form:"num"`  // 数量
	Sort bson.M `form:"sort"` // 排序方式，与mongodb中的一致
	Filter bson.M `form:"filter"`  // 过滤条件，和mongodb的过滤一致
}



// 获取代理ip
func getIp(c *gin.Context){
	var req param
	err := c.ShouldBind(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg": "请求参数不合法",
		})
		return
	}
	fmt.Printf("%#v\n", req)
	ips, err := dao.GetLimit(req.Num, req.Filter, req.Sort)
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
	router.POST("/get_ip", getIp)
	router.GET("/healthy", healthy)

	pprof.Register(router)

	host := viper.GetString("api.host")
	port := viper.GetString("api.port")
	router.Run(net.JoinHostPort(host, port))
}
