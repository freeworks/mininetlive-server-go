package easemob

import (
	. "app/common"
	logger "app/logger"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	. "github.com/bitly/go-simplejson"
	cache "github.com/patrickmn/go-cache"
)

type AccessToken struct {
	Access_Token string `json:access_token`
	Expires_In   int    `json:expires_in`
	Application  string `json:application`
}

const (
	CLIENT_ID     = "YXA6Zp1gwDYSEeasYzMDSE8dCA"
	CLIENT_SECRET = "YXA6WWzjQOVY3DASCGa22MCALGyKzoA"
	OAUTH_API_URL = "https://a1.easemob.com/authorize"
	API_TOKEN     = "https://a1.easemob.com/mininetlive/mininetlive/token"
)

func GetTokenFromServer() (AccessToken, error) {
	var payload = []byte(`{"grant_type":"client_credentials", 
	"client_id":"` + CLIENT_ID + `", "client_secret":"` + CLIENT_SECRET + `"}`)
	req, err := http.NewRequest("POST", API_TOKEN, strings.NewReader(string(payload)))
	CheckErr(err, "get token")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return AccessToken{}, err
	}
	defer resp.Body.Close()
	jsonMap := AccessToken{}
	json.NewDecoder(resp.Body).Decode(&jsonMap)
	return jsonMap, nil
}

func GetAccessToken(c *cache.Cache) string {
	var access_token string
	if token, found := c.Get("easemob"); found {
		access_token = token.(string)
	} else {
		token, err := GetTokenFromServer()
		if err == nil {
			access_token = token.Access_Token
			c.Set("easemob_token", token.Access_Token, time.Second*time.Duration(token.Expires_In))
		}
	}
	return access_token
}

func RegisterUser(username string, c *cache.Cache) (string, error) {
	access_token := GetAccessToken(c)
	if access_token == "" {
		return "", errors.New("register user get token fail")
	}
	url := "https://a1.easemob.com/mininetlive/mininetlive/users"
	var jsonStr = []byte(`{"username":"` + username + `","password":"123456"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		js, _ := NewJson(data)
		uuid, _ := js.Get("entities").GetIndex(0).Get("uuid").String()
		logger.Info(string(data))
		return uuid, nil
	} else {
		result := fmt.Sprintln("response Status:", resp.Status, ",Headers:", resp.Header)
		logger.Info(result)
		return "", errors.New(result)
	}
}

func CreateGroup(owner, title string, c *cache.Cache) (string,error) {
	access_token := GetAccessToken(c)
	if access_token == "" {
		return "",errors.New("create group get token fail")
	}
	//create group
	url := "https://a1.easemob.com/mininetlive/mininetlive/chatgroups"
	var jsonStr = []byte(`{"groupname":"` + title + `",
	"desc":"` + title + `",
	"public":true,
	"approval":false,
	"owner":"` + owner + `",
	"maxusers":10000}`)
	logger.Info(string(jsonStr))
	logger.Info("access_token ", access_token)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		js, _ := NewJson(data)
		groupId, _ := js.Get("data").Get("groupid").String()
		logger.Info("group id->",groupId)
		return groupId,nil
	} else {
		result := fmt.Sprintln("response Status:", resp.Status, ",Headers:", resp.Header)
		logger.Info(result)
		return "",errors.New(result)
	}
}
