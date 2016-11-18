package wxpub

import (
	. "app/common"
	logger "app/logger"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

const (
	token               = "zhujiang123"
	appID               = "wx52016ab80d994351"
	appSecret           = "fbfbef30831019e7262fa7581dc14dca"
	accessTokenFetchUrl = "https://api.weixin.qq.com/cgi-bin/token"
)

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type AccessTokenErrorResponse struct {
	Errcode float64
	Errmsg  string
}

func fetchAccessToken() (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchUrl,
		"?grant_type=client_credential&appid=",
		appID,
		"&secret=",
		appSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", 0.0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0.0, err
	}

	logger.Info(string(body))
	//Json Decoding
	if bytes.Contains(body, []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return "", 0.0, err
		}
		return atr.AccessToken, atr.ExpiresIn, nil
	} else {
		fmt.Println("return err")
		ater := AccessTokenErrorResponse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return "", 0.0, err
		}
		return "", 0.0, fmt.Errorf("%s", ater.Errmsg)
	}
}

type CustomServiceMsg struct {
	ToUser  string         `json:"touser"`
	MsgType string         `json:"msgtype"`
	Text    TextMsgContent `json:"text"`
}

type ShourtUrl struct {
	Action  string `json:"action"`
	LongUrl string `json:"long_url"`
}

type ShourtUrlResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	ShortUrl string `json:"short_url"`
}

type TextMsgContent struct {
	Content string `json:"content"`
}

func pushCustomMsg(accessToken, toUser, msg string) error {
	csMsg := &CustomServiceMsg{
		ToUser:  toUser,
		MsgType: "text",
		Text:    TextMsgContent{Content: msg},
	}

	body, err := json.MarshalIndent(csMsg, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	postReq, err := http.NewRequest("POST",
		"https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+accessToken,
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	respBody, _ := ioutil.ReadAll(resp.Body)
	logger.Info(string(respBody))
	if err != nil {
		logger.Error(err)
		return err
	}
	resp.Body.Close()

	return nil
}

func getShorturl(accessToken, longurl string) (string, error) {
	csMsg := &ShourtUrl{
		Action:  "long2short",
		LongUrl: longurl,
	}
	body, err := json.MarshalIndent(csMsg, " ", "  ")
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))
	postReq, err := http.NewRequest("POST",
		"https://api.weixin.qq.com/cgi-bin/shorturl?access_token="+accessToken,
		bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	logger.Info(string(respBody))
	if bytes.Contains(respBody, []byte("short_url")) {
		atr := ShourtUrlResponse{}
		err = json.Unmarshal(respBody, &atr)
		if err != nil {
			return "", err
		}
		return atr.ShortUrl, nil
	}
	return "", errors.New("long2shourt error:" + string(respBody))
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

//recvtextmsg_unencrypt.go
func parseTextRequestBody(r *http.Request) *TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println(string(body))
	requestBody := &TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

func RecvWXPubMsg(render render.Render, c *cache.Cache, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//	if !validateUrl(w, r) {
	//		log.Println("Wechat Service: this http request is not from Wechat platform!")
	//		return
	//	}
	if r.Method == "POST" {
		textRequestBody := parseTextRequestBody(r)
		if textRequestBody != nil {
			msg := strings.ToLower(textRequestBody.Content)
			openId := textRequestBody.FromUserName
			if "bd" == msg {
				var token string
				wxPubAccessToken, found := c.Get("wxPubAccessToken")
				if !found {
					// Fetch access_token
					accessToken, expiresIn, err := fetchAccessToken()
					token = accessToken
					if err != nil {
						log.Println("Get access_token error:", err)
						return
					}
					c.Set("wxPubAccessToken", accessToken, time.Duration(int64(expiresIn))*time.Second)
					fmt.Println(accessToken, expiresIn)
				} else {
					token = wxPubAccessToken.(string)
				}
				if wxPubAccessToken != "" {
					shorturl, err := getShorturl(token, "http://106.75.19.205/bind-phone.html?id"+openId)
					if err == nil {
						err = pushCustomMsg(token, openId, shorturl)
						if err != nil {
							log.Println("Push custom service message err:", err)
							return
						}
					}
				}
			}
		}
	} else if r.Method == "Get" {
		if r.Form["echostr"] != nil && len(r.Form["echostr"]) > 0 {
			io.WriteString(w, r.Form["echostr"][0])
		}
	}
}

func GetVCodeForWxPub(req *http.Request, c *cache.Cache, r render.Render) {
	req.ParseForm()
	phone := req.PostFormValue("phone")
	//TODO 校验
	vCode, err := SendSMS(phone)
	if err != nil {
		r.JSON(200, Resp{1009, "获取验证码失败", nil})
	} else {
		c.Set(phone, vCode, 60*time.Second)
		r.JSON(200, Resp{0, "获取验证码成功", nil})
	}
}

type WXPub struct {
	Id      int64     `db:"id"`
	OpenId  string    `db:"openid"`
	Phone   string    `db:"phone"`
	Created time.Time `db:"create_time"`
}

func newWXPub(openId, phone string) WXPub {
	return WXPub{
		Created: time.Now(),
		OpenId:  openId,
		Phone:   phone,
	}
}

func BindWxPubPhone(req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	openId := req.PostFormValue("vcode")
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	if openId == "" {
		r.JSON(200, Resp{1013, "openId不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			var wxPubs []WXPub
			_, err := dbmap.Select(&wxPubs, "SELECT * FROM t_wxpub WHERE openid=?", openId)
			if wxPubs == nil {
				w := newWXPub(openId, phone)
				err = dbmap.Insert(&w)

			} else {
				wxPubs[0].Phone = phone
				_, err = dbmap.Update(&wxPubs[0])

			}
			CheckErr(err, "BindWxPubPhone")
			if err != nil {
				r.JSON(200, Resp{1002, "绑定手机失败，服务器异常", nil})
				return
			} else {
				r.JSON(200, Resp{0, "绑定手机成功", nil})
				return
			}
		} else {
			r.JSON(200, Resp{1010, "输入验证码有误,请重新输入", nil})
			return
		}
	} else {
		r.JSON(200, Resp{1011, "验证码过期,请重新获取验证码", nil})
		return
	}
}