package HttpLib

import (
	"KolonseWeb/inject"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	if err != nil {
		panic(err)
	}
}

func (res *Response) SetCookie(name string, value string, others ...interface{}) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = url.QueryEscape(value)

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	cookie.Path = "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	}

	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		default:
			if others[3] != nil {
				cookie.Secure = true
			}
		}
	}

	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	res.Header().Add("Set-Cookie", cookie.String())
}

func NewResponse() *Response {
	ret := &Response{}
	ret.Injector = inject.New()
	return ret
}
