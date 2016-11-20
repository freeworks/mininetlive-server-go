package push

import (
	. "app/common"
	logger "app/logger"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	. "github.com/bitly/go-simplejson"
	"github.com/martini-contrib/render"
)

const (
	APP_KEY           string = "5774a96d67e58e4ef7001fa5"
	APP_MASTER_SECRET string = "7negllhaock3ncnkm7z2b57gbvais5ds"
)

func getSign(method, url, body string) string {
	h := md5.New()
	io.WriteString(h, method+url+body+APP_MASTER_SECRET)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func push(data []byte) (*Json, error) {
	logger.Info("push", string(data))
	sign := getSign("POST", "http://msg.umeng.com/api/send", string(data))
	logger.Info("sign", sign)
	req, err := http.NewRequest("POST", "http://msg.umeng.com/api/send?sign="+sign, bytes.NewBuffer(data))
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
		data, _ := js.Get("ret").String()
		logger.Info(data)
		return js, nil
	} else {
		result := fmt.Sprintln("response Status:", resp.Status, ",Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Info(string(body))
		return nil, errors.New(result)
	}
}

func PushNewActivity(aid, title string) {
	PushAll(
		newPayload("notification",
			newPushBody("新的活动上线", title, "新活动上线", "go_activity", "com.kouchen.mininetlive.ui.ActivityDetailActivity"),
			map[string]string{"aid": aid}))
}

func PushLiveBegin(aid, title string) {
	PushAll(
		newPayload("notification",
			newPushBody("正在直播中", title, "直播开始啦", "go_activity", "com.kouchen.mininetlive.ui.ActivityDetailActivity"),
			map[string]string{"aid": aid}))
}

func PushLiveEnd(aid, title string) {
	PushAll(
		newPayload("notification",
			newPushBody("已经结束", title, "直播结束了", "go_activity", "com.kouchen.mininetlive.ui.ActivityDetailActivity"),
			map[string]string{"aid": aid}))
}

func PushDividend(deviceId string) {
	PushOne(
		newPayload("notification",
			newPushBody("有一笔分红奖励", "有一笔分红奖励", "有一笔分红奖励", "go_activity", "com.kouchen.mininetlive.ui.DividendListActivity"), nil), deviceId)
}

type PushBody struct {
	Ticker      string            `json:"ticker"`
	Title       string            `json:"title"`
	Text        string            `json:"text"`
	Icon        string            `json:"icon"`
	LargeIcon   string            `json:"large_icon"`
	AfterOpen   string            `json:"after_open"`
	PlayVibrate string            `json:"play_vibrate"`
	PlaySound   string            `json:"play_sound"`
	PlayLights  string            `json:"play_lights"`
	Activity    string            `json:"activity"`
	Extra       map[string]string `json:"extra"`
}

type Payload struct {
	DisplayType string            `json:"display_type"`
	Body        PushBody          `json:"body"`
	Extra       map[string]string `json:"extra"`
}

func newPushBody(ticker, title, text, after_open, activity string) PushBody {
	return PushBody{
		Ticker:      ticker,
		Title:       title,
		Text:        text,
		AfterOpen:   after_open,
		Activity:    activity,
		Icon:        "ic_small",
		LargeIcon:   "ic_large",
		PlayVibrate: "true",
		PlaySound:   "true",
		PlayLights:  "true",
	}
}

func newPayload(displayType string, body PushBody, extra map[string]string) Payload {
	return Payload{
		DisplayType: displayType,
		Body:        body,
		Extra:       extra,
	}
}

func PushOne(payload Payload, deviceToken string) {
	p, err := json.Marshal(payload)
	CheckErr(err, "PushAll Marshal ")
	if err != nil {
		return
	}
	ploadString := string(p)
	logger.Info("PushAll ", ploadString)
	//android
	data := []byte(`{
					"appkey": "` + APP_KEY + `",
					"timestamp": ` + strconv.FormatInt(time.Now().Unix(), 10) + `,
					"type": "unicast",
					"payload":` + ploadString + `,
					"device_tokens":` + deviceToken + `
                    }`)
	push(data)
}

func PushAll(payload Payload) {
	p, err := json.Marshal(payload)
	CheckErr(err, "PushAll Marshal ")
	if err != nil {
		return
	}
	ploadString := string(p)
	logger.Info("PushAll ", ploadString)
	//android
	data := []byte(`{
					"appkey": "` + APP_KEY + `",
					"timestamp": ` + strconv.FormatInt(time.Now().Unix(), 10) + `,
					"type": "broadcast",
					"payload":
					` + ploadString + `
			}`)
	push(data)
	//	//IOS
	//	payload = []byte(`{
	//			"appkey": "` + APP_KEY + `",
	//			"timestamp": ` + strconv.FormatInt(time.Now().Unix(), 10) + `,
	//			"type": "broadcast",
	//			"payload":{
	//				  "aps":{ "alert":"【新活动上线】` + pushBody.Title + `"}
	//			}
	//	}`)
	//	push(payload)
}

func TestPush(req *http.Request, r render.Render) {
	PushLiveEnd("f4bee67fad5d95ff", "test")
}
