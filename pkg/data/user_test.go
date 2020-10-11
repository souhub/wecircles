package data

import "testing"

func TestUser(t *testing.T) {
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
		if got != want {
			t.Errorf("expected %s but got %s", want, got)
		}
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

		if got != want {
			t.Errorf("expected %s but got %s", want, got)
		}
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
		if got != want {
			t.Errorf("expected %s but got %s", want, got)
		}
	})

}
