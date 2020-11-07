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
	circle, err := data.GetCirclebyUser(id)
	if err != nil {
		data := Data{
			MyUser: myUser,
		}
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		if myUser.UserIdStr == id {
			tmp := parseTemplateFiles("layout", "navbar.private", "circle.new")
			if err := tmp.Execute(w, data); err != nil {
				logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			}
		}
		url := fmt.Sprintf("/circle?id=%s", myUser.UserIdStr)
		http.Redirect(w, r, url, 302)
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
	chats, err := data.GetChats(circle.ID)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
	data := Data{
		MyUser:          myUser,
		User:            owner,
		Circle:          circle,
		MembershipValid: membershipValid,
		Chats:           chats,
	}
	// ログイン済かつかつidが自分のものの場合
	if myUser.UserIdStr == id {
		// ここでサークルを作成済みかチェックする
		// circle, err := data.GetCirclebyUser(myUser.UserIdStr)
		// // サークルを持っていなければエラー発生し、/circle/newに飛ばされる
		// if err != nil {
		// 	logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		// 	http.Redirect(w, r, "/circle/new", 302)
		// 	return
		// }
		// chats, err := data.GetChats(circle.ID)
		// if err != nil {
		// 	logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		// }
		// data := Data{
		// 	MyUser: myUser,
		// 	Circle: circle,
		// 	Chats:  chats,
		// }
		tmp := parseTemplateFiles("layout.mypage", "navbar.private", "mypage.header.private", "mypage.circle", "mypage.navbar", "mypage.chats.private")
		if err := tmp.Execute(w, data); err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		}
		return
	}
	// ログイン済かつidが他人のもの場合
	// user, err := data.UserByUserIdStr(id)
	// if err != nil {
	// 	logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	// 	http.Redirect(w, r, "/circle/new", 302)
	// 	return
	// }
	tmp := parseTemplateFiles("layout.mypage", "navbar.private", "mypage.header", "mypage.circle", "mypage.navbar", "mypage.chats")
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

func CircleManageMembers(w http.ResponseWriter, r *http.Request) {
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
	memberships, err := circle.MembershipsByCircleID()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/manage", 302)
		return
	}
	users, err := data.GetUsersByUserID(memberships)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/manage", 302)
		return
	}
	numberOfMemberships, err := circle.CountMemberships()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/circle/manage", 302)
		return
	}
	type Data struct {
		MyUser              data.User
		Circle              data.Circle
		Users               []data.User
		NumberOfMemberships int
	}
	data := Data{
		MyUser:              myUser,
		Circle:              circle,
		Users:               users,
		NumberOfMemberships: numberOfMemberships,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "circle.manage.memberships")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		return
	}
}

func TweetsCircle(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
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
	circle, err := data.GetCirclebyUser(id)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		url := fmt.Sprintf("/circle?id=%s", myUser.UserIdStr)
		http.Redirect(w, r, url, 302)
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
	tmp := parseTemplateFiles("layout.mypage", "navbar.private", "mypage.header", "mypage.navbar", "mypage.tweets")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func SettingsCircle(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	id := vals.Get("id")
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
	circle, err := data.GetCirclebyUser(id)
	if err != nil {
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		url := fmt.Sprintf("/circle?id=%s", myUser.UserIdStr)
		http.Redirect(w, r, url, 302)
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
	tmp := parseTemplateFiles("layout.mypage", "navbar.private", "mypage.header", "mypage.navbar", "mypage.settings")
	if err := tmp.Execute(w, data); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
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
		logging.Info(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
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
		TwitterID:  r.PostFormValue("twitter"),
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
	myUser, err := session.User()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	circle, err := data.GetCirclebyUser(myUser.UserIdStr)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	data := Data{
		MyUser: myUser,
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
	circle.TwitterID = r.PostFormValue("twitter")
	if err := circle.Update(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
	http.Redirect(w, r, "/circle/manage", 302)
}

func DeleteCircle(w http.ResponseWriter, r *http.Request) {
	return
}
