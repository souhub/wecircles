package data

import "testing"

var user=User{
	Name: "Taro",
	UserIdStr: "taroId",
	Email: "taro@gmail.com",
	Password: "taroPass",
	ImagePath: "default.png",
}

func TestCreate(t *testing.T){
	if err:=user.Create(); err!=nil{
		t.Fatal(err)
	}

	gotUser,err:=UserByUserIdStr(user.UserIdStr)
	if err!=nil{
		t.Fatal(err)
	}

	got :=gotUser.Id
	want:=1
	if got!=want{
		t.Errorf("expected %d but got %d",want,got)
	}
}
