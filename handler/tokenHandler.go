package handler

import (
	"encoding/json"
	"fmt"
	"github.com/zhangmingkai4315/weichat-golang-backend/config"
	"log"
	"net/http"
	"time"
)

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	ErrorCode string `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
	ExpiresAt time.Time
}

//RefreshToken 用于直接查询token，将token保存在AccessToken对象中
func (tokenObj *accessToken) RefreshToken() (err error) {
	appId := config.ConfigObj.Security.AppId
	appSecret := config.ConfigObj.Security.AppSecret
	tokenUrl := config.ConfigObj.Url.TokenRefreshUrl
	requestUrl := fmt.Sprintf("%s&appid=%s&secret=%s", tokenUrl, appId, appSecret)
	log.Println(requestUrl)
	var myClient = &http.Client{Timeout: 10 * time.Second}
	res, err := myClient.Get(requestUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(tokenObj)
	if err != nil {
		return err
	}
	tokenObj.ExpiresAt = time.Now().Add(time.Second * time.Duration(tokenObj.ExpiresIn))
	log.Printf("%+v", tokenObj)
	return nil
}

func (tokenObj *accessToken) GetToken() (string, error) {
	//如果token存在且到期时间大于现在的时间，则直接返回
	if tokenObj.Token == "" || tokenObj.ExpiresAt.Before(time.Now()) {
		err := tokenObj.RefreshToken()
		if err != nil {
			return "", err
		}
	}
	log.Printf("Token is %s\n", tokenObj.Token)
	return tokenObj.Token, nil
}

var Token accessToken

func init() {
	Token.GetToken()
}
