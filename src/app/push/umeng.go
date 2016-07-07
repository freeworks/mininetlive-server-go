package push

import (
	"time"
		"crypto/md5"
		"net/http"
		"io/ioutil"
		"bytes"
		"io"
		"fmt"
		"errors"
		"strconv"
		logger "app/logger"
	. "github.com/bitly/go-simplejson"
)

const (
	APP_KEY string = "5774a96d67e58e4ef7001fa5"
	APP_MASTER_SECRET string = "7negllhaock3ncnkm7z2b57gbvais5ds"
)

func  getSign(method,url,body string) string{
	h := md5.New()
	io.WriteString(h, method+url+body+APP_MASTER_SECRET)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func push(payload [] byte) (*Json,string) {
	logger.Info("push","PushAppointment",string(payload))
	sign := getSign("POST","http://msg.umeng.com/api/send",string(payload))
	logger.Info("sign",sign)
	req, err := http.NewRequest("POST", "http://msg.umeng.com/api/send?sign="+sign, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		js, _ := NewJson(body)
		data,_ := js.Get("ret").String()
		logger.Info(data)
		return js, nil
	} else {
		result := fmt.Sprintln("response Status:", resp.Status, ",Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Info(string(body))
		return nil, errors.New(result)
	}
}


func PushAppointment(title,aid string){
 	payload := []byte(`{
			"appkey": "`+APP_KEY+`",
			"timestamp": `+strconv.FormatInt(time.Now().Unix(), 10) +`,
			"type": "groupcast",
			"payload":{
				  "aps":{ "alert":"【新活动上线】`+title+`"}
			},
			"filter":{
				"where": {
					    "and": 
					    [
					      {"tag":"`+aid+`"}
					    ]
				}
			}
	}`)
	push(payload)
	//IOS
	payload := []byte(`{
			"appkey": "`+APP_KEY+`",
			"timestamp": `+strconv.FormatInt(time.Now().Unix(), 10) +`,
			"type": "broadcast",
			"payload":{
				  "aps":{ "alert":"【新活动上线】`+title+`"}
			}
	}`)
	push(payload)
}

func PushNewActivity(title string) {
	//android
 	payload := []byte(`{
			"appkey": "`+APP_KEY+`",
			"timestamp": `+strconv.FormatInt(time.Now().Unix(), 10) +`,
			"type": "broadcast",
			"payload":{
				"display_type":"notification",
				"body":{
					"ticker":"新活动上线",
					"title":"有新的活动即将上线!",
					"text":"`+title+`",
					"after_open":"go_app",
				}
			}
	}`)
	push(payload)
	//IOS
	payload := []byte(`{
			"appkey": "`+APP_KEY+`",
			"timestamp": `+strconv.FormatInt(time.Now().Unix(), 10) +`,
			"type": "broadcast",
			"payload":{
				  "aps":{ "alert":"【新活动上线】`+title+`"}
			}
	}`)
	push(payload)
}
