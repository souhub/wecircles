package data

import (
	"testing"
)

func TestCircle(t *testing.T) {
	user := User{
		Id:        132,
		Name:      "TARO",
		UserIdStr: "taroID",
		Email:     "taro@gmail.com",
		Password:  "taroPass",
		ImagePath: "taroImage",
		CreatedAt: "2020-10-10-14",
	}

	circle := Circle{
		Name:      "testCircle",
		ImagePath: "default.jpg",
		Overview:  "たくさんのイベントを行うテニスサークルです。",
		Owner:     user,
	}

	t.Run("Create", func(t *testing.T) {
		reset(t)
		if err := circle.Create(); err != nil {
			t.Error(err, "- Failed to create the circle.")
		}
	})

	t.Run("Update", func(t *testing.T) {
		reset(t)
		if err := circle.Update(); err != nil {
			t.Error(err, "- Failed to update the circle.")
		}
	})

	// ResetUsers test
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
		circles := []Circle{
			{
				Name:      "TaroCircle",
				ImagePath: "default.png",
				Overview:  "Hello world",
				Category:  "tennis",
				Owner:     users[0],
			},
			{
				Name:      "HanaCircle",
				ImagePath: "default.png",
				Overview:  "Hello world",
				Category:  "tennis",
				Owner:     users[1],
			},
		}
		for _, circle := range circles {
			if err := circle.Create(); err != nil {
				t.Fatal(err)
			}
		}
		if err := ResetCircles(); err != nil {
			t.Error(err, "- Failed to reset the users table.")
		}
	})

}
