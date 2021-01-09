package check

import (
	"github.com/junhaideng/IPProxy/model"
	"testing"
)

func TestCheckIP(t *testing.T) {
	ip := model.IP{
		IP:   "223.82.106.253",
		Port: "3128",
	}
	if !check(ip) {
		t.Fatal("not pass")
	}
}
