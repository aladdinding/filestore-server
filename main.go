package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 静态资源处理
	//http.Handle("/static/",http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))
	//
	//// 文件存取接口
	//http.HandleFunc("/file/upload",handler.HTTPInterceptor(handler.UploadHandler))
	fmt.Println(strconv.Atoi("10"))
}
