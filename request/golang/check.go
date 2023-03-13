package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (u *User) check() error {
	if time.Now().Unix() >= u.refreshAfter {
		err := u.login()
		if err != nil {
			return err
		}
	}
	req, _ := http.NewRequest("POST", checkStatusUrl, strings.NewReader(""))
	req.Header.Add("Host", host+":7054")
	req.Header.Add("Authorization", u.accessToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	all, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var responseBody CheckResponseBody
	err = json.Unmarshal(all, &responseBody)
	if err != nil {
		return err
	}
	if responseBody.Code != 0 {
		return errors.New(responseBody.Msg)
	}
	u.remaining = responseBody.Data.Remaining
	u.balance = responseBody.Data.Balance
	fmt.Println(fmt.Sprintf("[检查] 用户%s检查自己的状态，余额：%d，剩余次数：%d",
		u.uid, u.balance, u.remaining))
	return nil
}
