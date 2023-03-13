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

func (u *User) open() error {
	if time.Now().Unix() >= u.refreshAfter {
		err := u.login()
		if err != nil {
			return err
		}
	}
	req, _ := http.NewRequest("POST", openEnvelopeUrl, strings.NewReader(""))
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
	var responseBody OpenResponseBody
	err = json.Unmarshal(all, &responseBody)
	if err != nil {
		return err
	}
	if responseBody.Code != 0 {
		return errors.New(responseBody.Msg)
	}
	u.remaining = responseBody.Data.Remaining
	u.balance = responseBody.Data.Balance
	fmt.Println(fmt.Sprintf("[红包] 用户%s参加抽奖获得%d，余额%d，剩余%d次机会",
		u.uid, responseBody.Data.Amount, u.balance, u.remaining))
	return nil
}
