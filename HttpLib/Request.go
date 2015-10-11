package HttpLib

import (
	"github.com/kolonse/KolonseWeb/inject"
	"net/http"
	"net/url"
)

// ResponseWriter 重新定义 增加状态码属性
type Request struct {
	*http.Request
	inject.Injector
	Path string
}

func (req *Request) GetCookie(name string) string {
	cookie, err := req.Cookie(name)
	if err != nil {
		return ""
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val
}

func NewRequest() *Request {
	ret := &Request{}
	ret.Injector = inject.New()
	return ret
}
