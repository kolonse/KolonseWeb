package HttpLib

import (
	"KolonseWeb/inject"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ResponseWriter 重新定义 增加状态码属性
type Response struct {
	http.ResponseWriter
	inject.Injector
}

func (res *Response) End(resString ...interface{}) {
	io.WriteString(res, fmt.Sprint(resString...))
}

func (res *Response) Json(obj interface{}) {
	if obj == nil {
		panic(errors.New("param should not be nil"))
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(res)
	err := encoder.Encode(obj)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}
