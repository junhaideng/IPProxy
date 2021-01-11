package schedule

import (
	"github.com/go-co-op/gocron"
	"github.com/junhaideng/IPProxy/check"
	"github.com/junhaideng/IPProxy/spider"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Status string 

const (
	Idle Status = "idle"
	Running Status = "running"
	Stop Status = "stop"
)


type Scheduler interface{
	// 开始进行运行
	Start() error 
	// 停止所有爬虫
	Stop() error 
	// 返回当前状态
	Status() Status
}


type ProxySchedule struct {
	// 所有注册的爬虫
	spiders []spider.Spider
	// 当前状态
	status Status
	// scheduler
	s *gocron.Scheduler
	m sync.RWMutex
}

func (p *ProxySchedule) Start() error {
	p.m.Lock()
	defer p.m.Unlock()
	p.status = Running
	interval := viper.GetUint64("schedule.interval")
	for index, s := range p.spiders{
		logrus.Infof("start number %d spider", index)
		_, err := p.s.Every(interval).Minute().Do(s.Start)
		if err != nil{
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("start number %d spider error", index)
		}
	}

	// 每interval*len(spiders)分钟进行另一次ip的检测
	_, err := p.s.Every(interval*uint64(len(p.spiders))).Minute().Do(check.CheckIP)
	if err != nil{
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("check ip err")
	}

	// 开启任务
	p.s.StartAsync()
	return nil
}

func (p *ProxySchedule)Status() Status {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.status
}

func (p *ProxySchedule) Stop() error {
	p.m.Lock()
	defer p.m.Unlock()
	p.status = Stop
	p.s.Stop()
	return  nil
}

func NewScheduler() *ProxySchedule {
	return &ProxySchedule{
		s: gocron.NewScheduler(time.UTC),
		status:  Idle,
	}
}

func (p *ProxySchedule) Register(s spider.Spider)error {
	p.m.Lock()
	defer p.m.Unlock()
	p.spiders = append(p.spiders, s)
	return nil
}

func (p *ProxySchedule) Unregister(s spider.Spider)error {
	p.m.Lock()
	defer p.m.Unlock()
	p.spiders = append(p.spiders, s)
	return nil
}

func (p *ProxySchedule) Clear()error {
	p.m.Lock()
	defer p.m.Unlock()
	p.s.Clear()
	p.spiders = make([]spider.Spider, 0)
	p.status = Idle
	return nil
}