package main

import "net/http"

func main() {
	http.HandleFunc("/readcookie", ReadCookie)
	http.HandleFunc("/writecookie", WriteCookie)
	http.HandleFunc("/deletecookie", DeleteCookie)
	http.ListenAndServe(":9090", nil)
}

func WriteCookie(w http.ResponseWriter, r *http.Request) {
	// 创建新的本地cookie
	cookie := http.Cookie{Name: "localCookie", Value: "GoLang", Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
	w.Write([]byte("设置cookie成功"))
}

func ReadCookie(w http.ResponseWriter, r *http.Request) {
	// 读取cookie
	cookie, err := r.Cookie("localCookie")
	if err == nil {
		cookieValue := cookie.Value
		// 将数据写入http连接中
		w.Write([]byte("cookie的值为：" + cookieValue))
	} else {
		w.Write([]byte("读取cookie出错：" + err.Error()))
	}
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "localCookie", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	w.Write([]byte("<删除cookie成功"))
}
