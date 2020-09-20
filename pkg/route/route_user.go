package route

import (
	"net/http"
	"wecircles/data"
)

func myPage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	name := user.Name
	posts, err := user.PostByUser()
	type Data struct {
		Name  string
		Posts []data.Post
	}
	data := Data{
		Name:  name,
		Posts: posts,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "mypage")
	err = tmp.Execute(w, data)
	return
}
