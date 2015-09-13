// KolonseWebUse project KolonseWebUse.go
package KolonseWeb

import (
	. "KolonseWeb/HttpLib"
	. "KolonseWeb/Type"
	"KolonseWeb/middleWare/StaticDir"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestKolonseWebGet(t *testing.T) {
	go func() { // 启动WEB
		test := make(map[string]string)
		test["status"] = "ok"
		DefaultApp.Get("/Get", func(req *Request, res *Response, next Next) {
			res.End("ok")
		})
		DefaultApp.Listen("127.0.0.1", 12346)
	}()
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://127.0.0.1:12346/Get")
	if err != nil {
		t.Error("should not err")
	}
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	fmt.Println("result: " + string(buff))
	if string(buff) != "ok" {
		t.Error("should be ok")
	}
}

func TestKolonseWebPost(t *testing.T) {
	go func() { // 启动WEB
		test := make(map[string]string)
		test["status"] = "ok"
		DefaultApp.Post("/Post", func(req *Request, res *Response, next Next) {
			res.Json(test)
		})
		DefaultApp.Listen("127.0.0.1", 12347)
	}()
	time.Sleep(1 * time.Second)
	buf := bytes.NewBuffer([]byte("{}"))
	resp, err := http.Post("http://127.0.0.1:12347/Post", "text/json", buf)
	if err != nil {
		t.Error("should not err")
	}
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	fmt.Println("result: " + string(buff))
	test := make(map[string]string)
	json.Unmarshal(buff, &test)
	if test["status"] != "ok" {
		t.Error("should be ok")
	}
}

func TestNewMiddleWareStaticDir(t *testing.T) {
	go func() { // 启动WEB
		test := make(map[string]string)
		test["status"] = "ok"
		DefaultApp.Use(StaticDir.NewMiddleWare("testStaticDir", false, nil))
		DefaultApp.Listen("127.0.0.1", 12348)
	}()
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://127.0.0.1:12348/test.txt")
	if err != nil {
		t.Error("should not err")
	}
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	fmt.Println("result: " + string(buff))
	if string(buff) != "ok" {
		t.Error("should be ok")
	}
}
