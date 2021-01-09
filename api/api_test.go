// TODO delete
package api

import (
	"github.com/junhaideng/IPProxy/model"
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
)

type response struct{
	Ips  []model.IP `json:"data"`
	Code int        `json:"code"`
}

func TestRun(t *testing.T){
	w := httptest.NewRecorder()	
	req := httptest.NewRequest("GET", "/get_ip", nil)
	router.ServeHTTP(w, req)
	r := w.Result()
	if r.StatusCode != http.StatusOK{
		t.Fatal("响应状态码错误")
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var res response
	if err := decoder.Decode(&res); err != nil{
		t.Fatal("解析响应内容错误")
	}
	t.Log(res)
}