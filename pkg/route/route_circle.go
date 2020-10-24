package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
)

func MyCircle(w http.ResponseWriter, r *http.Request) {
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
	// ここでサークルを作成済みかチェックする
	circle, err := data.GetCirclebyUser(user.UserIdStr)
	// サークルを持っていなければエラー発生し、/circle/newに飛ばされる
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	data := Data{
		User:   user,
		Circle: circle,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circle.private")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func Circle(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
	session, err := session(w, r)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	user, err := session.User()
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	// ログイン済かつかつidが自分のものの場合
	if user.UserIdStr == id {
		// サークル持っているかどうかはMyCircleハンドラで振り分ける
		http.Redirect(w, r, "/mycircle", 302)
		return
	}
	// ログイン済かつidが他人のもの場合
	// user, err := data.UserByUserIdStr(id)
	// if err != nil {
	// 	logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	// 	http.Redirect(w, r, "/circle/new", 302)
	// 	return
	// }
	circle, err := data.GetCirclebyUser(id)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
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

func Circles(w http.ResponseWriter, r *http.Request) {
	circles, err := data.Circles()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
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
		User:    user,
		Circles: circles,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circles")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func MembershipsCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func CircleManage(w http.ResponseWriter, r *http.Request) {
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
	circle, err := data.GetCirclebyUser(user.UserIdStr)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		User:   user,
		Circle: circle,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circle.manage")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
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
	data := Data{
		User: user,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circle.new")
	tmp.Execute(w, data)
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
	circle, err := data.GetCirclebyUser(user.UserIdStr)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	data := Data{
		User:   user,
		Circle: circle,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circle.edit")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

func UpdateCircle(w http.ResponseWriter, r *http.Request) {
	return
}

func DeleteCircle(w http.ResponseWriter, r *http.Request) {
	return
}
