package data

import (
	"testing"
)

func TestSession(t *testing.T) {

	setupSession := func(t *testing.T) (session Session) {
		t.Helper()
		err := user.Create()
		if err != nil {
			t.Fatal(err)
		}
		gotUser, err := UserByUserIdStr(user.UserIdStr)
		if err != nil {
			t.Fatal(err)
		}
		session, err = gotUser.CreateSession()
		if err != nil {
			t.Fatal(err, "- Failed to create the session.")
		}
		return
	}

	// Check test
	t.Run("Check", func(t *testing.T) {
		reset(t)
		session := setupSession(t)
		_, err := session.Check()
		if err != nil {
			t.Error(err, "- Failed to check the session.")
		}
	})

	// User test
	t.Run("User", func(t *testing.T) {
		reset(t)
		session := setupSession(t)
		gotUser, err := session.User()
		if err != nil {
			t.Error(err, "- Failed to get the session from the user.")
		}
		want := user.UserIdStr
		got := gotUser.UserIdStr
		assertCorrectMessage(t, want, got)
	})

	// Delete test
	t.Run("Delete", func(t *testing.T) {
		reset(t)
		session := setupSession(t)
		err := session.Delete()
		if err != nil {
			t.Error(err, "- Failed to delete the session.")
		}
	})

	// ResetSessions test
	t.Run("Reset", func(t *testing.T) {
		reset(t)
		users := []User{
			{
				Name:      "Taro",
				UserIdStr: "taroId",
				Email:     "taro@gmail.com",
				Password:  "taroPass",
				ImagePath: "default.png",
			},
			{
				Name:      "Hana",
				UserIdStr: "hanaId",
				Email:     "hana@gmail.com",
				Password:  "hanaPass",
				ImagePath: "default.png",
			},
		}
		sessions := []Session{}
		for _, user := range users {
			err := user.Create()
			if err != nil {
				t.Fatal(err)
			}
			gotUser, err := UserByUserIdStr(user.UserIdStr)
			if err != nil {
				t.Fatal(err)
			}
			session, err := gotUser.CreateSession()
			if err != nil {
				t.Fatal(err, "- Failed to create the session.")
			}
			sessions = append(sessions, session)
		}
		if err := ResetSessions(); err != nil {
			t.Error(err, "- Failedt to reset the session table.")
		}
	})
}
