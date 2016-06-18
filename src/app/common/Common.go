package common

import (
	logger "app/logger"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	mathRand "math/rand"
	"net/http"
	"strings"
	"time"
)

func SendSMS(mobile string) (string, error) {
	//TODO 判断电话号码
	url := "https://voice.yunpian.com/v1/voice/send.json"
	apikey := "9b11127a9701975c734b8aee81ee3526"
	vcode := GeneraVCode6()
	text := fmt.Sprintln("【云片网】您的验证码是", 1234)
	//uid := "verifyPhone"  //该条短信在您业务系统内的ID
	//callback_url
	res, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader("apikey="+apikey+"&mobile＝"+mobile+"&text="+text))
	if err != nil {
		CheckErr(err, "GetSMSCode send msg")
		return "", err
	} else {
		var result interface{}
		json.NewDecoder(res.Body).Decode(result)
		res.Body.Close()
		m := result.(map[string]interface{})
		if m["code"] == 0 {
			return vcode, nil
		} else {
			logger.Info("send msg response code ", m["code"], ", msg ", m["msg"])
			err = errors.New("send msg response err")
		}
	}
	return vcode, err
}
func GeneraVCode6() string {
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return vcode
}

func CheckErr(err error, msg string) {
	if err != nil {
		logger.Error(msg, err)
		// log.Fatalln(msg, err)
	}
}

func GeneraToken8() string {
	token, err := GenerateRandomString(8)
	CheckErr(err, "GeneraToken16")
	return token
}

func GeneraToken16() string {
	token, err := GenerateRandomString(16)
	CheckErr(err, "GeneraToken16")
	return token
}

func GeneraToken32() string {
	token, err := GenerateRandomString(32)
	CheckErr(err, "GeneraToken32")
	return token
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
