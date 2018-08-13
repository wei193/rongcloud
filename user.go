package rongcloud

import (
	"encoding/json"
	"errors"
	"log"
)

//UserGetToken 获取 Token 方法
func (r *Rongcloud) UserGetToken(userID, name, portraitURI string) (token string, err error) {
	data := make(map[string]string)
	data["userId"] = userID
	data["name"] = name
	data["portraitUri"] = portraitURI

	req, err := createRequset(UserGetToken, "POST", data, nil)
	if err != nil {
		return "", err
	}
	buf, err := r.requset(req)
	if err != nil {
		return "", err
	}
	type st struct {
		Code         int    `json:"code"`
		ErrorMessage string `json:"errorMessage"`
		Token        string `json:"token"`
		UserID       string `json:"userId"`
	}
	var in st
	err = json.Unmarshal(buf, &in)
	if err != nil {
		return "", err
	}
	if in.Code != 200 || in.UserID != userID {
		log.Println(string(buf))
		return "", errors.New(in.ErrorMessage)
	}
	return in.Token, nil
}

//UserRefresh 刷新用户信息方法
func (r *Rongcloud) UserRefresh(userID, name, portraitURI string) (err error) {
	data := make(map[string]string)
	data["userId"] = userID
	data["name"] = name
	data["portraitUri"] = portraitURI
	req, err := createRequset(UserRefresh, "POST", data, nil)
	if err != nil {
		return err
	}
	buf, err := r.requset(req)
	if err != nil {
		return err
	}
	type st struct {
		Code         int    `json:"code"`
		ErrorMessage string `json:"errorMessage"`
	}
	var in st
	err = json.Unmarshal(buf, &in)
	if err != nil {
		return err
	}
	if in.Code != 200 {
		return errors.New(in.ErrorMessage)
	}
	return nil
}
