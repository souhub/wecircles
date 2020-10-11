package data

import (
	"database/sql"
	"testing"
)

func TestUser(t *testing.T) {

	var user = User{
		Name:      "Taro",
		UserIdStr: "taroId",
		Email:     "taro@gmail.com",
		Password:  "taroPass",
		ImagePath: "default.png",
	}

	// Test helpers
	// Reset the all of the tables
	reset := func(t *testing.T) {
		t.Helper()
		err := ResetUsers()
		if err != nil {
			t.Fatal(err)
		}
		err = ResetSessions()
		if err != nil {
			t.Fatal(err)
		}
		err = ResetPosts()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test helper
	// Output the assertions
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("expected %s but got %s", want, got)
		}
	}

	// "Create" test
	t.Run("Create", func(t *testing.T) {
		reset(t)
		if err := user.Create(); err != nil {
			t.Error(err)
		}
	})

	// CreateSession test
	t.Run("Create Session", func(t *testing.T) {
		reset(t)
		gotSession, err := user.CreateSession()
		if err != nil {
			t.Fatal(err)
		}
		got := gotSession.UserIdStr
		want := user.UserIdStr
		assertCorrectMessage(t, want, got)
	})

	// UserbyUserIdStr test
	t.Run("Get the user by UuidStr", func(t *testing.T) {
		// Reset all the tables.
		reset(t)

		// Define "want" before saving the user
		want := user.UserIdStr

		// Save the user to get it
		if err := user.Create(); err != nil {
			t.Error(err)
		}

		// Get the user from DB
		gotUser, err := UserByUserIdStr(user.UserIdStr)
		if err != nil {
			t.Fatal(err)
		}

		// Define "got" after saving the user
		got := gotUser.UserIdStr

		assertCorrectMessage(t, want, got)
	})

	// UserbyEmail test
	// Just like the test of "UserbyUserIdStr"
	t.Run("Get the user by Email", func(t *testing.T) {
		reset(t)
		want := user.UserIdStr
		if err := user.Create(); err != nil {
			t.Fatal(err)
		}
		gotUser, err := UserByEmail(user.Email)
		if err != nil {
			t.Fatal(err)
		}
		got := gotUser.UserIdStr
		assertCorrectMessage(t, want, got)
	})

	// Update test
	t.Run("Update", func(t *testing.T) {
		reset(t)
		if err := user.Create(); err != nil {
			t.Fatal(err)
		}
		gotUser, err := UserByUserIdStr(user.UserIdStr)
		if err != nil {
			t.Fatal(err)
		}
		want := "UpdatedtaroId"
		gotUser.UserIdStr = want
		if err := gotUser.Update(); err != nil {
			t.Error(err)
		}
		got := gotUser.UserIdStr
		assertCorrectMessage(t, got, want)
	})

	// Delete test
	t.Run("Delete", func(t *testing.T) {
		reset(t)
		if err := user.Create(); err != nil {
			t.Fatal(err)
		}
		gotUser, err := UserByUserIdStr(user.UserIdStr)
		if err != nil {
			t.Fatal(err)
		}
		if err := gotUser.Delete(); err != nil {
			t.Error(err)
		}
		_, err = UserByUserIdStr(gotUser.UserIdStr)
		if err != sql.ErrNoRows {
			t.Error(err, "- Failed to delete the user.")
		}
	})

}
