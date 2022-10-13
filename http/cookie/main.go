package main

import "net/http"

func main() {
	http.HandleFunc("/login", MockLoginHandler)
	http.HandleFunc("/user", MockGetUserInfoHandler)
	http.ListenAndServe(":18080", nil)
}

func MockLoginHandler(w http.ResponseWriter, r *http.Request) {
	// 创建新的本地cookie，正常逻辑应该是：生成 - 本地存储 - 写入cookie，此处demo
	cookie := http.Cookie{Name: "localCookie", Value: "GoLang", Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
	w.Write([]byte("设置cookie成功"))
}

func MockGetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 读取cookie
	cookie, err := r.Cookie("localCookie")
	if err == nil {
		cookieValue := cookie.Value
		// 将数据写入http连接中
		w.Write([]byte("cookie的值为:" + cookieValue))
	} else {
		w.Write([]byte("读取cookie出错:" + err.Error()))
	}
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "localCookie", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	w.Write([]byte("<删除cookie成功"))
}
