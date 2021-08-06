package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// 匹配路径参数
var rePathParams, _ = regexp.Compile("/(:[a-zA-Z0-9_]+)")

// RouteError 路由错误, 实现 Error 接口
type RouteError string

func (e RouteError) Error() string {
	return string(e)
}

// HttpRouter 路由表
type HttpRouter struct {
	routeTable map[string]map[string]http.Handler
	pathRe     map[string]*regexp.Regexp
	paths      []string
}

func (r *HttpRouter) String() string {
	return fmt.Sprintf("-------------------------\n%s\n%s", strings.Join(r.paths, "\n"), "-------------------------")
}

func NewHttpRoute() *HttpRouter {
	return &HttpRouter{}
}

func (r *HttpRouter) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	handler, err := r.selectHandler(request)
	if err != nil {
		write.WriteHeader(http.StatusNotFound)
		return
	}

	handler.ServeHTTP(write, request)
}

// 根据请求路径查找处理器
// 带有路径参数路由: /players/:uid
func (r *HttpRouter) selectHandler(request *http.Request) (handler http.Handler, err error) {
	method := request.Method
	path := request.URL.Path

	// 根据请求路径和请求方式得到处理器
	for key := range r.routeTable {
		re := r.pathRe[key]
		if re.MatchString(path) {
			handler = r.routeTable[key][method]
			break
		}
	}

	if handler == nil {
		return nil, RouteError("the request path handler not found: " + path)
	}

	return handler, nil
}

// 注册处理器
func (r *HttpRouter) registerHandler(path, method string, handler http.Handler) {
	if r.routeTable == nil {
		r.routeTable = make(map[string]map[string]http.Handler)
		r.pathRe = make(map[string]*regexp.Regexp)
	}

	val := r.routeTable[path]
	if val == nil {
		val = make(map[string]http.Handler)
	}

	_, exists := val[method]
	if exists {
		log.Fatalf("duplicate route mapping: [%s] [%s]", method, path)
	}

	r.routeTable[path] = val
	r.pathRe[path] = pathToRegexp(path)

	val[method] = handler
	r.paths = append(r.paths, fmt.Sprintf("%s %s", method, path))
}

// 将路径参数变为正则表达式
func pathToRegexp(path string) *regexp.Regexp {
	re, _ := regexp.Compile(rePathParams.ReplaceAllString(path, strings.ReplaceAll(rePathParams.String(), ":", "")))
	return re
}

func (r *HttpRouter) POST(path string, handlerFunc http.HandlerFunc) {
	r.registerHandler(path, http.MethodPost, handlerFunc)
}

func (r *HttpRouter) DELETE(path string, handlerFunc http.HandlerFunc) {
	r.registerHandler(path, http.MethodDelete, handlerFunc)
}

func (r *HttpRouter) PUT(path string, handlerFunc http.HandlerFunc) {
	r.registerHandler(path, http.MethodPut, handlerFunc)
}

func (r *HttpRouter) GET(path string, handlerFunc http.HandlerFunc) {
	r.registerHandler(path, http.MethodGet, handlerFunc)
}
