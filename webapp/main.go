package main

import (
	"net/http"
	"text/template"
)

//IndexHandler 去首页
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//解析模板
	t := template.Must(template.ParseFiles("index.html"))
	//执行
	t.Execute(w, "")
}

func main() {
	//设置处理静态资源,/static/会匹配以/static/开头的路径，当浏览器请求index.html页面中的style.css文件时,static前缀会被
	//替换为views/static，然后去views/static/css目录中查找style.css文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	//直接去html页面
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("views/pages"))))
	http.HandleFunc("/main", IndexHandler)
	http.ListenAndServe("8080", nil)
}
