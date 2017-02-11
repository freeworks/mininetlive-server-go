package wxpub

import (
	. "app/common"
	logger "app/logger"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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

var tag = "[WXPub]"

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type AccessTokenErrorResponse struct {
	Errcode float64
	Errmsg  string
}

func fetchAccessToken(c *cache.Cache) (string, error) {
	wxPubAccessToken, found := c.Get("wxPubAccessToken")
	if !found {
		requestLine := strings.Join([]string{accessTokenFetchUrl,
			"?grant_type=client_credential&appid=",
			appID,
			"&secret=",
			appSecret}, "")

		resp, err := http.Get(requestLine)

		if err != nil || resp.StatusCode != http.StatusOK {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}

		logger.Info(tag, "[fetchAccessToken]", string(body))

		if bytes.Contains(body, []byte("access_token")) {
			atr := AccessTokenResponse{}
			err = json.Unmarshal(body, &atr)
			if err != nil {
				return "", err
			}
			c.Set("wxPubAccessToken", atr.AccessToken, time.Duration(int64(atr.ExpiresIn))*time.Second)
			return atr.AccessToken, nil
		} else {
			return "", errors.New("get access_token fail")
		}
	} else {
		return wxPubAccessToken.(string), nil
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
	logger.Info(tag, "[pushCustomMsg]", accessToken, toUser, msg)
	csMsg := &CustomServiceMsg{
		ToUser:  toUser,
		MsgType: "text",
		Text:    TextMsgContent{Content: msg},
	}

	body, err := json.MarshalIndent(csMsg, " ", "  ")
	if err != nil {
		CheckErr(tag, "[pushCustomMsg]", "", err)
		return err
	}
	logger.Info(tag, "[pushCustomMsg]", string(body))
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
	logger.Info(tag, "[pushCustomMsg]", string(respBody))
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
	logger.Info(tag, "[getShorturl]", string(body))
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
	logger.Info(tag, "[getShorturl]", string(respBody))
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
	logger.Info(tag, "[parseTextRequestBody]", string(body))
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
			if "bd" == msg || "提现" == msg {
				accessToken, err := fetchAccessToken(c)
				if accessToken != "" && err == nil {
					shorturl, err := getShorturl(accessToken, "http://www.weiwanglive.com/bind-phone.html?id="+openId)
					if err == nil {
						err = pushCustomMsg(accessToken, openId, "打开以下链接，绑定手机，才可以完成提现哦!"+shorturl)
						if err != nil {
							logger.Info(tag, "[RecvWXPubMsg]", "Push custom service message err:", err)
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
	openId := req.PostFormValue("openId")
	phone := req.PostFormValue("phone")
	vCode := req.PostFormValue("vcode")
	if openId == "" {
		r.JSON(200, Resp{1013, "openId不能为空", nil})
		return
	}
	if phone == "" {
		r.JSON(200, Resp{1013, "手机号不能为空", nil})
		return
	}
	if vCode == "" {
		r.JSON(200, Resp{1013, "验证码不能为空", nil})
		return
	}
	if cacheVCode, found := c.Get(phone); found {
		if cacheVCode.(string) == vCode {
			var wxPub WXPub
			err := dbmap.SelectOne(&wxPub, "SELECT * FROM t_wxpub WHERE openid=?", openId)
			if err != nil || wxPub.Phone == "" {
				w := newWXPub(openId, phone)
				err = dbmap.Insert(&w)
				if err != nil {
					r.JSON(200, Resp{1002, "手机号已绑定微信号,请检查", nil})
					return
				} else {
					r.JSON(200, Resp{0, "绑定手机成功", nil})
					return
				}
			} else {
				if wxPub.Phone != phone {
					_, err := dbmap.Exec("UPDATE t_wxpub SET phone = ? WHERE openid = ?", phone, openId)
					CheckErr(tag, "[BindWxPubPhone]", "Update nickname get failed", err)
					if err != nil {
						r.JSON(200, Resp{1002, "手机号已绑定其他微信号，不能再绑定", nil})
						return
					}
				}
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

type Config struct {
	AppID     string `json:"appId"`
	TimeStamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

type JSTokenResult struct {
	Errcode    int64
	Errmsg     string
	Ticket     string
	Expires_in int64
}

func GetConfig(w http.ResponseWriter, req *http.Request, c *cache.Cache, r render.Render, dbmap *gorp.DbMap) {
	req.ParseForm()
	url := req.PostFormValue("url")
	logger.Info(tag, "[GetConfig]", "url:", url)
	var ticketStr string
	ticketStr, err := getJsToken(c)
	logger.Info(tag, "[GetConfig]", "ticketStr:", ticketStr)
	if err != nil {
		CheckErr(tag, "[GetConfig]", "getJsToken", err)
		r.JSON(200, Resp{-1, "fail", "GetConfig  getJsToken"})
		return
	}
	nonceStr := RandomStr(16)
	timestamp := time.Now().Unix()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, url)
	sort.Strings([]string{str})
	h := sha1.New()
	io.WriteString(h, str)
	sigStr := fmt.Sprintf("%x", h.Sum(nil))
	logger.Info(tag, "[GetConfig]", "sigStr", str, "->", sigStr)
	config := Config{
		AppID:     appID,
		NonceStr:  nonceStr,
		TimeStamp: timestamp,
		Signature: sigStr,
	}
	r.JSON(200, Resp{0, "success", config})
	return
}

func getJsToken(c *cache.Cache) (string, error) {
	accessToken, err := fetchAccessToken(c)
	logger.Info(tag, "[getJsToken]", accessToken)
	if err == nil && accessToken != "" {
		if ticket, found := c.Get("jsticket"); !found {
			logger.Info(tag, "[getJsToken]", "ticket not found and get ticket from wx")
			resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi")
			if err != nil {
				CheckErr(tag, "[getJsToken]", "", err)
				return "", err
			}
			defer resp.Body.Close()
			respBody, _ := ioutil.ReadAll(resp.Body)
			logger.Info(tag, "[getJsToken]", string(respBody))
			if bytes.Contains(respBody, []byte("ticket")) {
				result := JSTokenResult{}
				err = json.Unmarshal(respBody, &result)
				CheckErr(tag, "[getJsToken]", "json.Unmarshal", err)
				if err != nil {
					return "", err
				}
				if result.Errcode == 0 {
					c.Set("jsticket", result.Ticket, time.Duration(int64(result.Expires_in))*time.Second)
					return result.Ticket, nil
				} else {
					return "", nil
				}
			} else {
				return "", errors.New("result error not Contains ticket ")
			}
		} else {
			return ticket.(string), nil
		}
	} else {
		return "", err
	}
}
