package data

import (
	"testing"

	"github.com/google/uuid"
)

var user = User{
	Name:      "Taro",
	UserIdStr: "taroId",
	Email:     "taro@gmail.com",
	Password:  "taroPass",
	ImagePath: "default.png",
}

var post = Post{
	Uuid:          uuid.New().String(),
	Title:         "Hello",
	Body:          "Hello World",
	UserId:        user.Id,
	UserIdStr:     user.UserIdStr,
	UserName:      user.Name,
	ThumbnailPath: "default_thumbnail.jpg",
	CreatedAt:     "2020-10-10",
}

// Test helpers
// Reset the all of the tables
func reset(t *testing.T) {

	t.Helper()
	// if err := ResetUserImageDir(); err != nil {
	// 	t.Fatal(err)
	// }
	if err := ResetUsers(); err != nil {
		t.Fatal(err)
	}
	if err := ResetSessions(); err != nil {
		t.Fatal(err)
	}
	if err := ResetPosts(); err != nil {
		t.Fatal(err)
	}
	if err := ResetCircles(); err != nil {
		t.Fatal(err)
	}

}

// Test helper
// Output the assertions
func assertCorrectMessage(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("expected %s but got %s", want, got)
	}
}
