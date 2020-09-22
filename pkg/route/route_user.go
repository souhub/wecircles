package route

import (
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
	posts, err := user.PostsByUser()
	type Data struct {
		Name  string
		Posts []data.Post
	}
	data := Data{
		Name:  name,
		Posts: posts,
	}
	tmp := parseTemplateFiles("layout", "navbar.private", "mypage")
	err = tmp.Execute(w, data)
	return
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	vals := r.URL.Query()
	user_id_str := vals.Get(("id"))
	user, err := data.UserByUserIdStr(user_id_str)
	if err != nil {
		logging.Warn("Failed to find users.")
	}
	if session.UserId != user.Id {
		name := user.Name
		posts, err := user.PostsByUser()
		if err != nil {
			logging.Warn("Failed to find posts.")
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
			logging.Warn("Failed to execute templates.")
		}
	} else {
		user, err := session.User()
		if err != nil {
			logging.Warn("Failed to find a user.")
		}
		name := user.Name
		posts, err := user.PostsByUser()
		if err != nil {
			logging.Warn("Failed to find posts.")
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
			logging.Warn("Failed to execute templates.")
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
		logging.Warn("User")
	}
	tmp := parseTemplateFiles("layout", "user.edit", "navbar.private")
	if err := tmp.Execute(w, user); err != nil {
		logging.Warn("Failed to execute templates.")
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	session, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	user, err := session.User()
	if err != nil {
		logging.Warn("Failed to find your user account from the session.")
	}
	//メソッドをPOSTのみ許可
	if r.Method != "POST" {
		log.Fatal("許可されていないメソッド")
	}

	//formから送信されたファイルを解析
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		logging.Warn("ファイルのアップロード失敗")
	}
	//アップロードされたファイル名を取得
	uploadedFileName := fileHeader.Filename
	//アップロードされたファイルを置くパスを設定
	imagePath := "view/img/user/" + uploadedFileName

	//imagePathにアップロードされたファイルを保存
	saveImage, err := os.Create(imagePath)
	if err != nil {
		logging.Warn("ファイルの確保失敗")
	}

	//保存用ファイルにアップロードされたファイルを書き込む
	_, err = io.Copy(saveImage, file)
	if err != nil {
		logging.Warn("アップロードしたファイルの書き込み失敗")
	}

	//saveImageとfileを最後に閉じる
	defer saveImage.Close()
	defer file.Close()

	err = r.ParseForm()
	if err != nil {
		logging.Warn("Failed to parse form")
	}

	// attr := map[string]interface{}{
	// 	"Name":    r.PostFormValue("name"),
	// 	"ImgPass": uploadedFileName,
	// }
	user.Name = r.PostFormValue("name")
	user.ImgPass = uploadedFileName
	if err = user.Update(); err != nil {
		logging.Warn("Failed to update your user account.")
	}
	http.Redirect(w, r, "/", 302)
}
