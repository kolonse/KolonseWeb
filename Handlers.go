package KolonseWeb

import (
	. "KolonseWeb/HttpLib"
	. "KolonseWeb/Type"
	"net/http"
)

//
//
/**
* 定义 Handler HTTP 有请求时的入口函数
* Steps 有请求时进行路由的步骤
 */
type Handlers struct {
	Steps []interface{}
}

//func (h *Handlers) Walkflow(stepIndex int, req *Request, res *Response) {
//	switch Steps[i].(type) {
//		case *RoutesHandlers:
//			routesHandlers := *(value.(*RoutesHandlers))
//			routesHandlers.Do(req, res)
//			isGoNext = routesHandlers.IsGoNext()
//		case *MiddleWares:
//			middleWares := *value.(*MiddleWares)
//			middleWares.Do(req, res, middleWares.Next)
//			isGoNext = middleWares.IsGoNext()
//		}
//	}
//}

func (h *Handlers) WalkSteps(req *Request, res *Response) {
	for _, value := range h.Steps {
		isGoNext := false
		switch value.(type) {
		case *RoutesHandlers:
			routesHandlers := *(value.(*RoutesHandlers))
			routesHandlers.Do(req, res)
			isGoNext = routesHandlers.IsGoNext()
		case *MiddleWares:
			middleWares := *value.(*MiddleWares)
			middleWares.Do(req, res, middleWares.Next)
			isGoNext = middleWares.IsGoNext()
		}
		if !isGoNext {
			break
		}
	}
}

func (h *Handlers) Use(do DoStep) {
	middleWares := NewMiddleWares()
	middleWares.Do = do
	h.Steps = append(h.Steps, middleWares)
}

func (h *Handlers) GetRoutesHandlers() *RoutesHandlers {
	index := -1
	for i, value := range h.Steps {
		switch value.(type) {
		case *RoutesHandlers:
			index = i
		}
		if index != -1 {
			break
		}
	}

	if index != -1 { // 存在 那么就不需要进行创建对象直接调用
		return h.Steps[index].(*RoutesHandlers)
	} else { // 不存在就需要进行创建
		routesHandlers := NewRoutesHandlers()
		h.Steps = append(h.Steps, routesHandlers)
		return routesHandlers
	}
}

func (h *Handlers) Get(patter string, do DoStep) {
	// 先要遍历判断一下是否存在了 路由解析对象
	routesHandlers := h.GetRoutesHandlers()
	routesHandlers.Get(patter, do)
}

func (h *Handlers) Post(patter string, do DoStep) {
	// 先要遍历判断一下是否存在了 路由解析对象
	routesHandlers := h.GetRoutesHandlers()
	routesHandlers.Post(patter, do)
}

func (h *Handlers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	/**
	 *
	 */
	h.WalkSteps(toKolonseRequest(req), toKolonseResponse(&res))
}

func toKolonseRequest(req *http.Request) *Request {
	request := NewRequest()
	request.Request = req
	request.Path = request.URL.Path
	return request
}

func toKolonseResponse(res *http.ResponseWriter) *Response {
	response := NewResponse()
	response.ResponseWriter = *res
	return response
}

func NewHandler() *Handlers {
	return &Handlers{}
}
