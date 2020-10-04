package route

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func MyPage(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	name := user.Name
	image := user.ImagePath
	posts, err := user.PostsByUser()
	type Data struct {
		Name      string
		ImagePath string
		Posts     []data.Post
	}
	data := Data{
		Name:      name,
		ImagePath: image,
		Posts:     posts,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "mypage")
	err = tmp.Execute(w, data)
	return
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	vals := r.URL.Query()
	// user_id_str := vals.Get(("id"))
	user, err := data.UserByUserIdStr(vals.Get(("id")))
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	if session.UserIdStr != user.UserIdStr {
		name := user.Name
		posts, err := user.PostsByUser()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		type Data struct {
			Name  string
			Posts []data.Post
		}
		data := Data{
			Name:  name,
			Posts: posts,
		}
		tmp := parseTemplateFiles("layout", "navbar.private", "user.show")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	} else {
		user, err := session.User()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		name := user.Name
		posts, err := user.PostsByUser()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		type Data struct {
			Name  string
			Posts []data.Post
		}
		data := Data{
			Name:  name,
			Posts: posts,
		}
		tmp := parseTemplateFiles("layout", "navbar.private", "mypage")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
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
	http.Redirect(w, r, "/", 302)
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

	// Get the file sent form the form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	// Get the uploaded file's name from the file.
	uploadedFileName := fileHeader.Filename
	// Set the uploaded file's path
	imagePath := "web/img/user/" + uploadedFileName

	// Save the uploaded file to "imagePath"
	saveImage, err := os.Create(imagePath)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Write the uploaded file to the file for saving.
	_, err = io.Copy(saveImage, file)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}

	// Close the "saveImage" and "file"
	defer saveImage.Close()
	defer file.Close()

	user.ImagePath = uploadedFileName

	if err = user.UpdateImage(); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", 302)
}
