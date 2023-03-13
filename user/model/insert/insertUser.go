package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
)

func main() {
	conn := sqlx.NewMysql("red_envelope_rain:123456@tcp(127.0.0.1:3306)/red_envelope_rain")
	m := NewUserModel(conn)
	base, n := 1000001, 1000000
	for i := 0; i < n; i++ {
		u := &User{
			Id:           int64(i),
			UserAccount:  strconv.Itoa(base + i),
			UserName:     "张三的第" + strconv.Itoa(i+1) + "个孩子",
			UserPassword: "123456",
		}
		_, err := m.Insert(context.Background(), u)
		if err != nil {
			return
		}
	}
}
