package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/souhub/wecircles/pkg/data"
	"github.com/souhub/wecircles/pkg/logging"
	"gopkg.in/go-playground/validator.v9"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tmp := parseTemplateFiles("login.layout", "navbar.public", "login")
	err := tmp.Execute(w, nil)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	tmp := parseTemplateFiles("login.layout", "navbar.public", "signup")
	err := tmp.Execute(w, nil)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
	}
}

func SignupAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/signup", 302)
		return
	}
	user := data.User{
		Name:      r.PostFormValue("name"),
		UserIdStr: r.PostFormValue("user_id_str"),
		Email:     r.PostFormValue("email"),
		Password:  data.Encrypt(r.PostFormValue("password")),
		ImagePath: "default.png",
	}
	// ユーザーIDとメールアドレスの一意性を保証
	existedUsers, err := data.Users()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/signup", 302)
		return
	}
	for _, existedUser := range existedUsers {
		if user.UserIdStr == existedUser.UserIdStr || user.Email == existedUser.Email {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/signup", 302)

			return
		}
	}
	validate := validator.New() //validatorインスタンス生成
	//validator実行
	if err := validate.Struct(&user); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/signup", 302)
		return
	}
	//DBに登録
	if err := user.Create(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/signup", 302)
		return
	}
	// http.Redirect(w, r, "/login", 302)
	//そのまま認証終わらせてマイページに飛ばす
	signupedUserID := user.UserIdStr
	signupedUser, err := data.UserByUserIdStr(signupedUserID)
	// userごとの画像保存フォルダ作成
	currentRootDir, err := os.Getwd()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/signup", 302)
		return
	}
	userImageDir := fmt.Sprintf("%s/web/img/user%d", currentRootDir, signupedUser.Id)
	_, err = os.Stat(userImageDir)
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		err = os.Mkdir(userImageDir, 0777)
	}
	session, err := signupedUser.CreateSession()
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/signup", 302)
		return
	}
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.Uuid,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := fmt.Sprintf("/user?id=%s", session.UserIdStr)
	http.Redirect(w, r, url, 302)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	//フォームの入力内容解析
	if err := r.ParseForm(); err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	//フォームに入力されたEmailを持つユーザーをデータベースから特定

	user, err := data.UserByUserIdStr(r.FormValue("useridstr"))
	if err != nil {
		logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
		http.Redirect(w, r, "/login", 302)
		return
	}
	//その特定したユーザーとフォームに入力されたパスワードのチェック
	if user.Password == data.Encrypt(r.FormValue("password")) {
		//認証されたらユーザーのセッションを作成しクッキーを持たせ、マイページにリダイレクトさせる
		//SHA-1でuuid作ってクッキーを作ってブラウザに渡す。
		session, err := user.CreateSession()
		if err != nil {
			logging.Warn(err, logging.GetCurrentFile(), logging.GetCurrentFileLine())
			http.Redirect(w, r, "/login", 302)
			return
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
		return
	} else {
		//認証されなかったらログインファームにリダイレクトさせる
		http.Redirect(w, r, "/login", 302)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := data.Session{Uuid: cookie.Value}
		session.Delete()
	}
	http.Redirect(w, r, "/login", 302)
}
