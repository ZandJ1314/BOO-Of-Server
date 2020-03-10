package servers

import (
	"boo/lib/conf"
	"boo/lib/tools"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type WXLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

//https://developers.weixin.qq.com/miniprogram/dev/api/wx.getUserInfo.html

type Watermark struct {
	AppID     string `json:"appid"`
	TimeStamp int64  `json:"timestamp"`
}

type WXUserInfo struct {
	OpenID    string    `json:"openId,omitempty"`
	NickName  string    `json:"nickName"`
	AvatarUrl string    `json:"avatarUrl"`
	Gender    int       `json:"gender"`
	Country   string    `json:"country"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	UnionID   string    `json:"unionId,omitempty"`
	Language  string    `json:"language"`
	Watermark Watermark `json:"watermark,omitempty"`
}

type ResUserInfo struct {
	UserInfo      WXUserInfo `json:"userInfo"`
	RawData       string     `json:"rawData"`
	Signature     string     `json:"signature"`
	EncryptedData string     `json:"encryptedData"`
	IV            string     `json:"iv"`
}

func Login(code string,resUserInfo *ResUserInfo) (*WXUserInfo,error) {
	//appid和secret是在微信公众平台上获取
	appId := conf.GinAdminconfig.Weixin.GetWeiXinAppid()
	secret := conf.GinAdminconfig.Weixin.GetWeiXinSecret()
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	//合成url
	url = fmt.Sprintf(url,appId,secret,code)
	//创建http的get的请求
	req,err := http.Get(url)
	if err != nil {
		return nil,err
	}
	//http相应结束后进行资源回收，否则有可能造成内存泄漏
	defer req.Body.Close()
	//解析http请求的数据，进行回包
	wxResp := new(WXLoginResponse)
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}
	s := sha1.New()
	s.Write([]byte(resUserInfo.RawData + wxResp.SessionKey))
	sha1 := s.Sum(nil)
	sha1hash := hex.EncodeToString(sha1)
	if resUserInfo.Signature != sha1hash {
		return nil,fmt.Errorf("sha1加密错误")
	}
	userInfo,userErr := DecryptUserInfoData(wxResp.SessionKey, resUserInfo.EncryptedData, resUserInfo.IV)
	return userInfo,userErr
}

func DecryptUserInfoData(sessionKey string, encryptedData string, iv string) (*WXUserInfo,error) {
	sk, _ := base64.StdEncoding.DecodeString(sessionKey)
	ed, _ := base64.StdEncoding.DecodeString(encryptedData)
	i, _ := base64.StdEncoding.DecodeString(iv)
	decryptedData, err := tools.AesCBCDecrypt(ed, sk, i)
	if err != nil {
		return nil,err
	}
	var wxUserInfo WXUserInfo
	err = json.Unmarshal(decryptedData, &wxUserInfo)
	return &wxUserInfo,err
}
