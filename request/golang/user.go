package main

import (
	"fmt"
	"time"
)

type User struct {
	uid          string
	password     string
	name         string
	remaining    int
	balance      int
	accessToken  string
	accessExpire int64
	refreshAfter int64
}

func NewUser(uid string) *User {
	return &User{
		uid:      uid,
		password: "123456",
	}
}

func (u *User) work() {
	_ = u.login()
	_ = u.check()
	for u.remaining > 0 {
		time.Sleep(4 * time.Second)
		_ = u.open()
	}
	fmt.Println(fmt.Sprintf("[结束] 用户%s参加抽奖共获得%d", u.uid, u.balance))
}
