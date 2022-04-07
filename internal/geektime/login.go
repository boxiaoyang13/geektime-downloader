package geektime

import (
	"errors"
	"net/http"

	"github.com/nicoxiang/geektime-downloader/internal/client"
	pgt "github.com/nicoxiang/geektime-downloader/internal/pkg/geektime"
)

var ErrWrongPassword = errors.New("密码错误, 请尝试重新登录")
var ErrTooManyLoginAttemptTimes = errors.New("密码输入错误次数过多，已触发验证码校验，请稍后再试")

// Login call geektime login api and return auth cookies
func Login(phone, password string) ([]*http.Cookie, error) {
	var result struct {
		Code int `json:"code"`
		Data struct {
			UID  int    `json:"uid"`
			Name string `json:"nickname"`
		} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"error"`
	}

	loginResponse, err := client.NewTimeGeekAccountRestyClient().R().
		SetHeader("Referer", pgt.GeekBangAccount+"/signin?redirect=https%3A%2F%2Ftime.geekbang.org%2F").
		SetBody(
			map[string]interface{}{
				"country":   86,
				"appid":     1,
				"platform":  3,
				"cellphone": phone,
				"password":  password,
			}).
		SetResult(&result).
		Post("/account/ticket/login")

	if err != nil {
		panic(err)
	}

	if result.Code == 0 {
		var cookies []*http.Cookie
		for _, c := range loginResponse.Cookies() {
			if c.Name == "GCID" || c.Name == "GCESS" {
				cookies = append(cookies, c)
			}
		}
		return cookies, nil
	} else if result.Error.Code == -3031 {
		return nil, ErrWrongPassword
	} else if result.Error.Code == -3005 {
		return nil, ErrTooManyLoginAttemptTimes
	}
	panic(result.Error.Msg)
}
