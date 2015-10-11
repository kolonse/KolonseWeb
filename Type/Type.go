package Type

import (
	. "github.com/kolonse/KolonseWeb/HttpLib"
)

// next 函数定义 模拟 nodejs express 框架
// 执行 next 函数 中间件就会继续向下执行
type Next func()
type DoStep func(req *Request, res *Response, next Next)

func DefaultDoStep(req *Request, res *Response, next Next) {
}
