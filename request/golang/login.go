package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (u *User) login() error {

	marshal, err := json.Marshal(&LoginRequest{
		Account:  u.uid,
		Password: u.password,
	})
	if err != nil {
		return err
	}
	loginBody := strings.NewReader(string(marshal))
	req, _ := http.NewRequest("POST", loginUrl, loginBody)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Host", host+":7050")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	all, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var responseBody LoginResponseBody
	err = json.Unmarshal(all, &responseBody)
	if err != nil {
		return err
	}
	if responseBody.Code != 0 {
		return errors.New(responseBody.Msg)
	}
	u.name = responseBody.Data.UserName
	u.accessToken = responseBody.Data.AccessToken
	u.accessExpire = responseBody.Data.AccessExpire
	u.refreshAfter = responseBody.Data.RefreshAfter
	fmt.Println(fmt.Sprintf("[登录] 用户%s进行登录", u.uid))
	return nil
}
