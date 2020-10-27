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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	// ここでサークルを作成済みかチェックする
	circle, err := data.GetCirclebyUser(myUser.UserIdStr)
	// サークルを持っていなければエラー発生し、/circle/newに飛ばされる
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	data := Data{
		MyUser: myUser,
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
	myUser, err := session.User()
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	// ログイン済かつかつidが自分のものの場合
	if myUser.UserIdStr == id {
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
	owner, err := circle.GetOwner()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	membershipValid, err := checkMembership(myUser, circle)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		MyUser:          myUser,
		User:            owner,
		Circle:          circle,
		MembershipValid: membershipValid,
	}
	tmp := parseTemplateFiles("layout.mypage", "navbar.mypage", "mypage.header", "mypage.circle")
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	data := Data{
		MyUser:  myUser,
		Circles: circles,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circles")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func CircleManage(w http.ResponseWriter, r *http.Request) {
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
	circle, err := data.GetCirclebyUser(myUser.UserIdStr)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		MyUser: myUser,
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	data := Data{
		MyUser: myUser,
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	currentRootDir, err := os.Getwd()
	circleImageDir := fmt.Sprintf("%s/web/img/user%d/circles/mycircle", currentRootDir, myUser.Id)
	_, err = os.Stat(circleImageDir)
	if err != nil {
		err = os.MkdirAll(circleImageDir, 0777)
	}
	circleImage, err := myUser.UploadCircleImage(r)
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
		OwnerID:    myUser.Id,
		OwnerIDStr: myUser.UserIdStr,
	}
	if err = circle.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/new", 302)
		return
	}
	createdCircle, err := data.GetCirclebyUser(myUser.UserIdStr)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	membership := data.Membership{
		UserID:   createdCircle.OwnerID,
		CircleID: createdCircle.ID,
	}
	if err := membership.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	url := fmt.Sprintf("/circle?id=%s", myUser.UserIdStr)
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
	circle, err := data.GetCirclebyUser(myUser.UserIdStr)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	circleImage, err := circle.Upload(r)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	circle.ImagePath = circleImage
	circle.Name = r.PostFormValue("name")
	circle.Overview = r.PostFormValue("overview")
	circle.Category = r.PostFormValue("category")
	if err := circle.Update(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	http.Redirect(w, r, "/circle/manage", 302)
}

func DeleteCircle(w http.ResponseWriter, r *http.Request) {
	return
}
