package KolonseWeb

import (
	. "KolonseWeb/HttpLib"
	. "KolonseWeb/Type"
	"net/http"
)

type RoutesHandler struct {
	Do       DoStep
	nextStep bool
}

func (routesHandler *RoutesHandler) Next() {
	routesHandler.nextStep = true
}

func (routesHandler *RoutesHandler) IsGoNext() bool {
	return routesHandler.nextStep
}

func NewRoutesHandler() *RoutesHandler {
	return &RoutesHandler{
		Do:       DefaultDoStep,
		nextStep: false,
	}
}

type RoutesHandlers struct {
	Routes   map[string]map[string]*RoutesHandler
	nextStep bool
}

func (routesHandlers *RoutesHandlers) Get(patter string, do DoStep) {
	_, ok := routesHandlers.Routes["GET"]
	if !ok {
		routesHandlers.Routes["GET"] = make(map[string]*RoutesHandler)
	}

	routesHandler := NewRoutesHandler()
	routesHandler.Do = do
	routesHandlers.Routes["GET"][patter] = routesHandler
}

func (routesHandlers *RoutesHandlers) Post(patter string, do DoStep) {
	_, ok := routesHandlers.Routes["POST"]
	if !ok {
		routesHandlers.Routes["POST"] = make(map[string]*RoutesHandler)
	}

	routesHandler := NewRoutesHandler()
	routesHandler.Do = do
	routesHandlers.Routes["POST"][patter] = routesHandler
}

func (routesHandlers *RoutesHandlers) Do(req *Request, res *Response) {
	methodRoutesHandlers, ok := routesHandlers.Routes[req.Method]
	if !ok {
		routesHandlers.nextStep = true
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	routesHandler, ok := methodRoutesHandlers[req.Path]
	if !ok {
		routesHandlers.nextStep = true
		res.WriteHeader(http.StatusNotFound)
		return
	}
	routesHandler.Do(req, res, routesHandler.Next)
	routesHandlers.nextStep = routesHandler.IsGoNext()
}

func (routesHandlers *RoutesHandlers) IsGoNext() bool {
	return routesHandlers.nextStep
}

func NewRoutesHandlers() *RoutesHandlers {
	ret := &RoutesHandlers{}
	ret.Routes = make(map[string]map[string]*RoutesHandler)
	ret.nextStep = false
	return ret
}
