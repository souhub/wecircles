package route

import (
	"net/http"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func NewCircle(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	tmp := parseTemplateFiles("layout", "circle.new", "navbar.private")
	tmp.Execute(w, nil)
}

func CreateCircle(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	circleImage, err := user.UploadCircleImage(r)
	circle := data.Circle{
		Name:      r.PostFormValue("name"),
		ImagePath: circleImage,
		Overview:  r.PostFormValue("overview"),
		Owner:     user,
	}
	if err = circle.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/edit", 302)
		return
	}
	http.Redirect(w, r, "/mypage", 302)
}

func EditCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func EditCircleImage(w http.ResponseWriter, r *http.Request) {
	return
}

func UpdateCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func UpdateCircleImage(w http.ResponseWriter, r *http.Request) {
	return
}

func DeleteCircle(w http.ResponseWriter, r *http.Request) {
	return
}
