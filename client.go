package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 获取认证二维码url
func (w *WeCharClient) AuthCodeUrl(state string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		w.Appid, url.QueryEscape(w.RedirectUri),
		strings.Join(w.Scope, ","), state)
}

// 获取access token
func (w *WeCharClient) Exchange(code string) (AccessToken, error) {
	reUrl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code#wechat_redirect",
		w.Appid, w.Scope, code)
	if response, err := Requests("GET", reUrl, nil); err == nil && response.StatusCode == http.StatusOK {

		body := response.Body
		defer body.Close()
		if bodyByte, err := ioutil.ReadAll(body); err == nil {
			var result AccessToken
			if err := json.Unmarshal(bodyByte, &result); err == nil {
				return result, nil
			}
		}
	}
	return AccessToken{}, errors.New("get access token fail")
}

// 刷新access token
func (w *WeCharClient) ReGetAccessToken(refreshToken string) (AccessToken, error) {
	reUrl := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s#wechat_redirect",
		w.Appid, refreshToken)
	if response, err := Requests("GET", reUrl, nil); err == nil && response.StatusCode == http.StatusOK {

		body := response.Body
		defer body.Close()
		if bodyByte, err := ioutil.ReadAll(body); err == nil {
			var result AccessToken
			if err := json.Unmarshal(bodyByte, &result); err == nil {
				return result, nil
			}
		}
	}
	return AccessToken{}, errors.New("get access token fail")
}

// 获取用户信息
func (w *WeCharClient) GetUserInfo(accessToken string, openId string, lang ...string) (UserInfo, error){

	reUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s#wechat_redirect", accessToken, openId)
	if len(lang) > 0 {
		reUrl += fmt.Sprintf("&%s", lang[0])
	}
	if response, err := Requests("GET", reUrl, nil); err == nil && response.StatusCode == http.StatusOK {

		body := response.Body
		defer body.Close()
		if bodyByte, err := ioutil.ReadAll(body); err == nil {
			var result UserInfo
			if err := json.Unmarshal(bodyByte, &result); err == nil {
				return result, nil
			}
		}
	}
	return UserInfo{}, errors.New("get user info fail")
}