package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func Circle(w http.ResponseWriter, r *http.Request) {
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
	circle, err := user.GetCircle()
	type Data struct {
		User   data.User
		Circle data.Circle
	}
	data := Data{
		Circle: circle,
		User:   user,
	}
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	tmp := parseTemplateFiles("layout", "circle", "navbar.private")
	tmp.Execute(w, data)
}

func Circles(w http.ResponseWriter, r *http.Request) {
	circles, err := data.Circles()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	type Data struct {
		Circles []data.Circle
		User    data.User
	}
	session, err := session(w, r)
	if err != nil {
		data := Data{
			Circles: circles,
		}
		tmp := parseTemplateFiles("layout", "circles", "navbar.public")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		return
	}
	user, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	data := Data{
		Circles: circles,
		User:    user,
	}
	tmp := parseTemplateFiles("layout", "circles", "navbar.private")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func MembershipsCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func NewCircle(w http.ResponseWriter, r *http.Request) {
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
	tmp := parseTemplateFiles("layout", "circle.new", "navbar.private")
	tmp.Execute(w, user)
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
	currentRootDir, err := os.Getwd()
	circleImageDir := fmt.Sprintf("%s/web/img/user%d/circles/mycircle", currentRootDir, user.Id)
	_, err = os.Stat(circleImageDir)
	if err != nil {
		err = os.MkdirAll(circleImageDir, 0777)
	}
	circleImage, err := user.UploadCircleImage(r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	circle := data.Circle{
		Name:      r.PostFormValue("name"),
		ImagePath: circleImage,
		Overview:  r.PostFormValue("overview"),
		Category:  r.PostFormValue("category"),
		Owner:     user,
	}
	if err = circle.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	http.Redirect(w, r, "/mypage", 302)
}

// func ShowCircle(w http.ResponseWriter, r *http.Request) {
// 	session, err := session(w, r)
// 	if err != nil {
// 		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 		http.Redirect(w, r, "/login", 302)
// 		return
// 	}
// 	user, err := session.User()
// 	if err != nil {
// 		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 		http.Redirect(w, r, "/login", 302)
// 		return
// 	}
// 	circle, err := user.GetCircle()
// 	if err != nil {
// 		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 		return
// 	}
// 	if circle.OwnerID == user.Id {
// 		http.Redirect(w, r, "/circle", 302)
// 		return
// 	}
// 	tmp := parseTemplateFiles("layout", "circle.show", "navbar.private")
// 	tmp.Execute(w, circle)
// }

func EditCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func UpdateCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func DeleteCircle(w http.ResponseWriter, r *http.Request) {
	return
}
