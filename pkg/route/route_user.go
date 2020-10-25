package route

import (
	"fmt"
	"log"
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

// GET /mypage
// Get the mypage
func MyPage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	posts, err := user.PostsByUser()
	data := Data{
		User:  user,
		Posts: posts,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "mypage", "posts")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// GET /user/show
// Get the users
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
	data := Data{
		User:  user,
		Posts: posts,
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	data = Data{
		MyUser: myUser,
		User:   user,
		Posts:  posts,
	}
	// ログイン後に他人のユーザー名をクリックした場合
	tmp := parseTemplateFiles("layout.mypage", "navbar.mypage", "mypage.header", "mypage.posts")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

// GET /user/edit
// Edit the user
func EditUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	type Data struct {
		User data.User
	}
	data := Data{
		User: user,
	}
	tmp := parseTemplateFiles("layout", "user.edit", "navbar.private")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		log.Fatal(err)
	}
}

// POST /user/update
// Update the user
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
	// if r.Method != "POST" {
	// 	logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	// }

	// Parse the form
	if err = r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	userImagePath, err := user.Upload(r)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	user.ImagePath = userImagePath
	user.Name = r.PostFormValue("name")
	user.UserIdStr = r.PostFormValue("user_id_str")

	if err = user.Update(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	fmt.Println(user.UserIdStr)
	http.Redirect(w, r, "/mypage", 302)
}

// GET /user/edit/image
// Get the user edit image page
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

// POST /user/update/image
// Update the user image
func UpdateUserImage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	user, err := session.User()
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	err = user.DeleteUserImage()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Uplaod the new user's image.
	user.ImagePath, err = user.Upload(r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}

	// Update the user's image path from old one to the new one in the DB.
	if err = user.UpdateImage(); err != nil {
		log.Fatal(err)
		return
	}
	http.Redirect(w, r, "/user/edit", 302)
}

// POST /user/delete
// Delete the user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if err := session.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if user.Password == data.Encrypt(r.FormValue("password")) {
		if err = user.DeleteUserImageDir(); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/user/delete/confirm", 302)
			return
		}
		if err = user.DeletePosts(); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		if err = user.Delete(); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		http.Redirect(w, r, "/signup", 302)
		return
	}
	http.Redirect(w, r, "/user/delete/confirm", 302)
}

func DeleteUserConfirm(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	tmp := parseTemplateFiles("layout", "user.delete.confirm", "navbar.private")
	if err := tmp.Execute(w, user); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
}
