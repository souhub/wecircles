package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/route"
)

func main() {
	currentRootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir := fmt.Sprintf("%s/web/", currentRootDir)
	files := http.FileServer(http.Dir(dir))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", route.Index)
	http.HandleFunc("/login", route.Login)
	http.HandleFunc("/signup", route.Signup)
	http.HandleFunc("/signup_account", route.SignupAccount)
	http.HandleFunc("/authenticate", route.Authenticate)
	http.HandleFunc("/logout", route.Logout)

	http.HandleFunc("/post/new", route.NewPost)
	http.HandleFunc("/post/create", route.CreatePost)
	http.HandleFunc("/posts/manage", route.PostsManage)
	http.HandleFunc("/post", route.ShowPost)
	http.HandleFunc("/post/edit", route.EditPost)
	http.HandleFunc("/post/update", route.UpdatePost)
	http.HandleFunc("/post/delete", route.DeletePost)

	// http.HandleFunc("/mypage", route.MyPage)
	// http.HandleFunc("/mycircle", route.MyCircle)

	http.HandleFunc("/user/edit", route.EditUser)
	http.HandleFunc("/user/update", route.UpdateUser)
	http.HandleFunc("/user", route.User)
	http.HandleFunc("/user/delete/confirm", route.DeleteUserConfirm)
	http.HandleFunc("/user/delete", route.DeleteUser)

	http.HandleFunc("/circle", route.Circle)
	http.HandleFunc("/circles", route.Circles)
	http.HandleFunc("/circle/manage", route.CircleManage)
	http.HandleFunc("/circle/manage/members", route.CircleManageMembers)
	http.HandleFunc("/circle/new", route.NewCircle)
	http.HandleFunc("/circle/create", route.CreateCircle)
	http.HandleFunc("/circle/edit", route.EditCircle)
	http.HandleFunc("/circle/update", route.UpdateCircle)
	http.HandleFunc("/circle/delete", route.DeleteCircle)

	http.HandleFunc("/circle/memberships", route.MembershipsCircles)
	http.HandleFunc("/circle/membership/create", route.MembershipsCircleCreate)
	http.HandleFunc("/circle/membership/delete", route.DeleteMembership)
	http.HandleFunc("/circle/membership/delete/byowner", route.DeleteMembershipByOwner)
	http.HandleFunc("/circle/tweets", route.TweetsCircle)
	http.HandleFunc("/circle/settings", route.SettingsCircle)

	http.HandleFunc("/chat/create", route.CreateChat)
	http.HandleFunc("/chat/delete", route.DeleteChat)

	http.ListenAndServe(":80", nil)
}
