package StaticDir

import (
	"errors"
	. "github.com/kolonse/KolonseWeb/HttpLib"
	. "github.com/kolonse/KolonseWeb/Type"

	"os"
	"path"

	"reflect"
	"strings"
)

// The algorithm uses at most sniffLen bytes to make its decision.
//const sniffLen = 512

var (
	DefaultStaticDir    = "public"
	DefaultLikeDirIndex = false
	DefaultHopeGzipExt  = make([]string, 0)
)

// 取出静态目录参数 默认 public
func parseStaticDir(opt ...interface{}) string {
	if len(opt) >= 1 {
		if opt[0] == nil { // 如果是 nil 那么置为默认值
			return DefaultStaticDir
		}
		if reflect.TypeOf(opt[0]).Kind() != reflect.String {
			panic(errors.New("param 1 not string"))
		}
		return opt[0].(string)
	}
	// 如果没有输入参数 那么返回默认的
	return DefaultStaticDir
}

func parseLikeDirIndex(opt ...interface{}) bool {
	if len(opt) >= 2 {
		if opt[1] == nil { // 如果是 nil 那么置为默认值
			return DefaultLikeDirIndex
		}
		if reflect.TypeOf(opt[1]).Kind() != reflect.Array {
			panic(errors.New("param 2 not bool"))
		}
		return opt[1].(bool)
	}
	// 如果没有输入参数 那么返回默认的
	return DefaultLikeDirIndex
}

func parseHopeGzipExt(opt ...interface{}) []string {
	if len(opt) >= 3 {
		if opt[2] == nil { // 如果是 nil 那么置为默认值
			return DefaultHopeGzipExt
		}
		if reflect.TypeOf(opt[2]).Kind() != reflect.Array {
			panic(errors.New("param 3 not []string"))
		}
		return opt[2].([]string)
	}
	// 如果没有输入参数 那么返回默认的
	return DefaultHopeGzipExt
}

/**
*	静态文件中间件
*	opt[0] 静态目录
*	opt[1] 是否喜欢静态目录索引
*	opt[2] 希望对指定后缀文件进行压缩 eg:[".js","*.css"]
 */
func NewMiddleWare(opt ...interface{}) DoStep {
	// 解析传入参数
	dir := parseStaticDir(opt...)
	//likeDirIndex := parseLikeDirIndex(opt...)
	hopeGzipExt := parseHopeGzipExt(opt...)
	return func(req *Request, res *Response, next Next) {
		// 静态文件支持
		/**
		*	1.根据path拼接 dir 进行静态文件查找 如果能查找到 那么直接返回数据 调用next 进行下一步路路由
		 */
		// 静态文件只支持 GET 和 HEAD 方法
		if req.Method != "GET" && req.Method != "HEAD" {
			next() //  进入下一步路由
			return
		}
		requestPath := path.Clean(req.Path)
		file := path.Join(dir, requestPath)
		finfo, err := os.Stat(file)
		if err != nil { // 如果文件不存在 那么直接进行下一个路由
			next() //  进入下一步路由
			return
		}
		//if the request is dir and DirectoryIndex is false then
		if finfo.IsDir() {
			//if !likeDirIndex {
			//	exception("403", ctx)
			//	return
			//} else if ctx.Input.Request.URL.Path[len(ctx.Input.Request.URL.Path)-1] != '/' {
			//	http.Redirect(ctx.ResponseWriter, ctx.Request, ctx.Input.Request.URL.Path+"/", 302)
			//	return
			//}
			// 支持目录索引功能暂时不实现
			next()
			return
		}
		isStaticFileToCompress := false
		if hopeGzipExt != nil && len(hopeGzipExt) > 0 {
			for _, statExtension := range hopeGzipExt {
				if strings.HasSuffix(strings.ToLower(file), strings.ToLower(statExtension)) {
					isStaticFileToCompress = true
					break
				}
			}
		}
		if isStaticFileToCompress {
			// 压缩逻辑暂时不实现
			next()
			return
		} else {
			serveFile(res.ResponseWriter, req.Request, file)
		}
	}
}

// name is '/'-separated, not filepath.Separator.
