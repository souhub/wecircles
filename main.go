package main

import (
	"log"
	"net/http"

	"github.com/souhub/wecircles/pkg/route"
)

func main() {
	files := http.FileServer(http.Dir("web/img"))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", route.Index)
	http.HandleFunc("/login", route.Login)
	http.HandleFunc("/signup", route.Signup)
	http.HandleFunc("/signup_account", route.SignupAccount)
	http.HandleFunc("/authenticate", route.Authenticate)
	http.HandleFunc("/logout", route.Logout)

	http.HandleFunc("/post/new", route.NewPost)
	http.HandleFunc("/post/create", route.CreatePost)
	http.HandleFunc("/post/show", route.ShowPost)
	http.HandleFunc("/post/edit", route.EditPost)
	http.HandleFunc("/post/update", route.UpdatePost)
	http.HandleFunc("/post/delete", route.DeletePost)

	http.HandleFunc("/mypage", route.MyPage)
	http.HandleFunc("/user/edit", route.EditUser)
	http.HandleFunc("/user/update", route.UpdateUser)
	http.HandleFunc("/user/show", route.ShowUser)
	log.Fatal(http.ListenAndServe(":80", nil))
}
