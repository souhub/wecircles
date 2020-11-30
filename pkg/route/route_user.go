package route

import (
	"fmt"
	"net/http"
	"os"

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
	myUser, err := session.User()
	posts, err := myUser.PostsByUser()
	data := Data{
		MyUser: myUser,
		Posts:  posts,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "user", "posts")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// GET /user/show
// Get the users
func User(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	user, err := data.UserByUserIdStr(id)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		url := fmt.Sprintf("/user?id=%s", user.UserIdStr)
		http.Redirect(w, r, url, 302)
		return
	}
	posts, err := user.PostsByUser()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		User:            user,
		Posts:           posts,
		ImagePathPrefix: os.Getenv("IMAGE_PATH"),
	}
	session, err := session(w, r)
	// ログイン前にユーザー名クリックした場合
	if err != nil {
		tmp := parseTemplateFiles("layout", "navbar.public", "user", "posts")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		return
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	// ログイン後に自分のユーザー名をクリックした場合
	if session.UserId == user.Id {
		// http.Redirect(w, r, "/mypage", 302)
		myUser, err := session.User()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		posts, err := myUser.PostsByUser()
		data := Data{
			MyUser:          myUser,
			Posts:           posts,
			ImagePathPrefix: os.Getenv("IMAGE_PATH"),
		}
		tmp := parseTemplateFiles("layout", "navbar.private", "user", "posts")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			return
		}
		return
	}
	data = Data{
		MyUser:          myUser,
		User:            user,
		Posts:           posts,
		ImagePathPrefix: os.Getenv("IMAGE_PATH"),
	}
	// ログイン後に他人のユーザー名をクリックした場合
	tmp := parseTemplateFiles("layout.mypage", "navbar.private", "mypage.header", "index", "posts")
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		MyUser:          myUser,
		ImagePathPrefix: os.Getenv("IMAGE_PATH"),
	}
	tmp := parseTemplateFiles("layout", "user.edit", "navbar.private")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// POST /user/update
// Update the user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	myUser, err := session.User()
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Allow the "POST" method, only
	if r.Method != "POST" {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Parse the form
	if err = r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	userImagePath, err := myUser.Upload(r)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	myUser.ImagePath = userImagePath
	myUser.Name = r.PostFormValue("name")
	myUser.UserIdStr = r.PostFormValue("user_id_str")

	if err = myUser.Update(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if err = myUser.UpdatePostUserIdStr(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if err = myUser.UpdateChatUserIdStr(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	url := fmt.Sprintf("/user?id=%s", myUser.UserIdStr)
	http.Redirect(w, r, url, 302)
}

// GET /user/edit/image
// Get the user edit image page
func EditUserImage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	tmp := parseTemplateFiles("layout", "user.edit.image", "navbar.private")
	if err := tmp.Execute(w, myUser); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

// POST /user/update/image
// Update the user image
func UpdateUserImage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	myUser, err := session.User()
	if err != nil {
		logging.Fatal(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	err = myUser.DeleteUserImage()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Uplaod the new myUser's image.
	myUser.ImagePath, err = myUser.Upload(r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}

	// Update the myUser's image path from old one to the new one in the DB.
	if err = myUser.UpdateImage(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	http.Redirect(w, r, "/myUser/edit", 302)
}

// POST /user/delete
// Delete the user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if myUser.Password != data.Encrypt(r.FormValue("password")) {
		http.Redirect(w, r, "/user/delete/confirm", 302)
		return
	}
	if err = myUser.DeleteUserImageDir(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/user/delete/confirm", 302)
		return
	}
	if err = myUser.DeletePosts(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if err = myUser.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if err := session.Delete(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	http.Redirect(w, r, "/signup", 302)
	return
}

// GET /user/delete/confirm
// Get the user-delete-confirm page
func DeleteUserConfirm(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
	}
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	data := Data{
		MyUser: myUser,
	}
	tmp := parseTemplateFiles("layout", "user.delete.confirm", "navbar.private")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
}
