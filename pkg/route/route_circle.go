package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func Circle(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	UserIdStr := vals.Get("id")
	user, err := session(w, r)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	if user.UserIdStr == UserIdStr {
		user, err := data.UserByUserIdStr(UserIdStr)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/login", 302)
			return
		}
		circle, err := data.GetCirclebyUser(UserIdStr)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/circle/new", 302)
			return
		}
		type Data struct {
			User   data.User
			Circle data.Circle
		}
		data := Data{
			User:   user,
			Circle: circle,
		}
		tmp := parseTemplateFiles("layout", "navbar.private", "circle")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	} else {
		user, err := data.UserByUserIdStr(UserIdStr)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/circle/new", 302)
			return
		}
		circle, err := data.GetCirclebyUser(UserIdStr)
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/circle/new", 302)
			return
		}
		type Data struct {
			User   data.User
			Circle data.Circle
		}
		data := Data{
			User:   user,
			Circle: circle,
		}
		tmp := parseTemplateFiles("layout", "navbar.private", "circle.public")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
	}
}

// func Circle(w http.ResponseWriter, r *http.Request) {
// 	vals := r.URL.Query()
// 	UserIdStr := vals.Get("id")
// 	circle, err := data.GetCirclebyUser(UserIdStr)
// 	if err != nil {
// 		_, err := session(w, r)
// 		if err != nil {
// 			tmp := parseTemplateFiles("layout", "navbar.public", "circle.public")
// 			if err := tmp.Execute(w, circle); err != nil {
// 				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 			}
// 			return
// 		}
// 		http.Redirect(w, r, "/circle/new", 302)
// 		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 		return
// 	} else {
// 		session, err := session(w, r)
// 		if err != nil {
// 			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 			http.Redirect(w, r, "/login", 302)
// 			return
// 		}
// 		user, err := session.User()
// 		if err != nil {
// 			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 			http.Redirect(w, r, "/login", 302)
// 			return
// 		}
// 		type Data struct {
// 			User   data.User
// 			Circle data.Circle
// 		}
// 		data := Data{
// 			User:   user,
// 			Circle: circle,
// 		}
// 		if user.UserIdStr == UserIdStr {
// 			tmp := parseTemplateFiles("layout", "navbar.private", "circle")
// 			if err := tmp.Execute(w, data); err != nil {
// 				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 			}
// 			return
// 		} else {

// 			tmp := parseTemplateFiles("layout", "navbar.private", "circle.public")
// 			if err := tmp.Execute(w, circle); err != nil {
// 				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
// 			}
// 		}

// 	}
// }

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
		Name:       r.PostFormValue("name"),
		ImagePath:  circleImage,
		Overview:   r.PostFormValue("overview"),
		Category:   r.PostFormValue("category"),
		OwnerID:    user.Id,
		OwnerIDStr: user.UserIdStr,
	}
	if err = circle.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	url := fmt.Sprintf("/circle?id=%s", user.UserIdStr)
	http.Redirect(w, r, url, 302)
}

func EditCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func UpdateCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func DeleteCircle(w http.ResponseWriter, r *http.Request) {
	return
}
