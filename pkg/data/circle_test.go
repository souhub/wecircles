package data

import (
	"testing"
)

func TestCircle(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		reset(t)
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
		if err := circle.Create(); err != nil {
			t.Error(err, "- Failed to create the circle.")
		}
	})
}
