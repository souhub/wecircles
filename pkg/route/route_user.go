package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func MyPage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	id := user.Id
	name := user.Name
	image := user.ImagePath
	posts, err := user.PostsByUser()
	type Data struct {
		Id        int
		Name      string
		ImagePath string
		Posts     []data.Post
	}
	data := Data{
		Id:        id,
		Name:      name,
		ImagePath: image,
		Posts:     posts,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "mypage")
	err = tmp.Execute(w, data)
	return
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	user, err := data.UserByUserIdStr(vals.Get(("id")))
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	posts, err := user.PostsByUser()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	type Data struct {
		ID        int
		Name      string
		UserIdStr string
		ImagePath string
		Posts     []data.Post
	}
	data := Data{
		ID:        user.Id,
		Name:      user.Name,
		UserIdStr: user.UserIdStr,
		ImagePath: user.ImagePath,
		Posts:     posts,
	}
	session, err := session(w, r)
	// ログイン前にユーザー名クリックした場合
	if err != nil {
		tmp := parseTemplateFiles("layout", "navbar.public", "user.show")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		return
	}
	// ログイン後に自分のユーザー名をクリックした場合
	if session.UserId == user.Id {
		http.Redirect(w, r, "/mypage", 302)
		return
	}
	// ログイン後に他人のユーザー名をクリックした場合
	tmp := parseTemplateFiles("layout", "navbar.private", "user.show")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	tmp := parseTemplateFiles("layout", "user.edit", "navbar.private")
	if err := tmp.Execute(w, user); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		log.Fatal(err)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Allow the "POST" method, only
	if r.Method != "POST" {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Parse the form
	err = r.ParseForm()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	user.Name = r.PostFormValue("name")
	user.UserIdStr = r.PostFormValue("user_id_str")

	fmt.Println(user.UserIdStr)
	if err = user.Update(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	fmt.Println(user.UserIdStr)
	http.Redirect(w, r, "/mypage", 302)
}

func EditUserImage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	tmp := parseTemplateFiles("layout", "user.edit.image", "navbar.private")
	if err := tmp.Execute(w, user); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		log.Fatal(err)
	}
}

func UpdateUserImage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	user, err := session.User()
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}

	// currentDir, err := os.Getwd()
	// userImageDir := fmt.Sprintf("%s/web/img/user/user%d", currentDir, user.Id)
	// _, err = os.Stat(userImageDir)
	// if err != nil {
	// 	err = os.Mkdir(userImageDir, 0777)
	// }
	// Delete the current user's image
	err = user.DeleteUserImage()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Uplaod the new user's image.
	user.ImagePath, err = user.Upload(r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Update the user's image path from old one to the new one in the DB.
	if err = user.UpdateImage(); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/user/edit", 302)
}
