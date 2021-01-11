package check

import (
	"github.com/junhaideng/IPProxy/model"
	"sync"
	"testing"
)

func TestCheckIP(t *testing.T) {
	ip := model.IP{
		IP:   "223.82.106.253",
		Port: "3128",
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go CheckSingleIP(ip, &wg)
	wg.Wait()
}
