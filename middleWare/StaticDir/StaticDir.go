package StaticDir

import (
	. "KolonseWeb/HttpLib"
	. "KolonseWeb/Type"
	"net/http"
	"os"
	"path"
	"strings"
)

/**
*	静态文件中间件
*	dir 静态目录
*	likeDirIndex 是否喜欢静态目录索引
*	hopeGzipExt 希望对指定后缀文件进行压缩 eg:[".js","*.css"]
 */
func NewMiddleWare(dir string, likeDirIndex bool, hopeGzipExt []string) DoStep {
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
			http.ServeFile(res.ResponseWriter, req.Request, file)
		}
	}
}
