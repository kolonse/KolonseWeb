package HttpLib

import (
	"net/http"
)

// ResponseWriter 重新定义 增加状态码属性
type Request struct {
	*http.Request
	Path string
}
