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
	"reflect"
	"strings"
	"time"

	. "github.com/bitly/go-simplejson"
	"github.com/coopernurse/gorp"
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
			c.Set("easemob_token", token.Access_Token, time.Second*time.Duration(token.Expires_In-60))
		}
	}
	return access_token
}

func request(method, url string, postJson []byte, c *cache.Cache) (*Json, error) {
	access_token := GetAccessToken(c)
	if access_token == "" {
		return nil, errors.New(url + "get token fail")
	}
	logger.Info(method, url, postJson)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(postJson))
	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		js, _ := NewJson(data)
		return js, nil
	} else {
		result := fmt.Sprintln("response Status:", resp.Status, ",Headers:", resp.Header)
		data, _ := ioutil.ReadAll(resp.Body)
		logger.Info(string(data))
		return nil, errors.New(result)
	}
}

func RegisterUser(username string, c *cache.Cache) (string, error) {
	url := "https://a1.easemob.com/mininetlive/mininetlive/users"
	var jsonStr = []byte(`{"username":"` + username + `","password":"123456"}`)
	js, err := request("POST", url, jsonStr, c)
	if err == nil {
		return js.Get("entities").GetIndex(0).Get("uuid").String()
	} else {
		return "", err
	}
}

func JoinGroup(groupId, username string, c *cache.Cache) {
	url := "https://a1.easemob.com/mininetlive/mininetlive/chatgroups/" + groupId + "/users/" + username
	_, error := request("POST", url, nil, c)
	CheckErr(error, "AddUserJoinGroup")
}

func LeaveGroup(groupId, username string, c *cache.Cache) {
	url := "https://a1.easemob.com/mininetlive/mininetlive/chatgroups/" + groupId + "/users/" + username
	_, error := request("DELETE", url, nil, c)
	CheckErr(error, "AddUserJoinGroup")
}

func GetGroupOnlineUserCount(c *cache.Cache, dbmap *gorp.DbMap) {
	logger.Info("GetGroupOnlineUserCount....")
	js, err := request("GET", "https://a1.easemob.com/mininetlive/mininetlive/chatgroups", nil, c)
	if err == nil {
		groups, err := js.Get("data").Array()
		if err == nil {
			for _, g := range groups {
				groupInfo := g.(map[string]interface{})
				groupId := groupInfo["groupid"]
				affiliations := groupInfo["affiliations"]
				_, err := dbmap.Exec("UPDATE t_activity SET online_count = ? WHERE group_id = ? ", affiliations.(json.Number).String(), groupId.(string))
				CheckErr(err, "update activity online_count")
			}
		} else {
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}
}

func GetGroupMemberCount(groupId string, c *cache.Cache) (int, error) {
	js, err := request("GET", "https://a1.easemob.com/mininetlive/mininetlive/chatgroups/"+groupId, nil, c)
	if err == nil {
		groups, err := js.Get("data").Array()
		logger.Info(groups)
		if err == nil {
			group := groups[0]
			groupInfo := group.(map[string]interface{})
			count, _ := groupInfo["affiliations_count"].(json.Number).Int64()
			logger.Info("GetGroupMemberCount", count)
			return int(count), nil
		} else {
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}
	return 0, err
}

func GetGroupMemberList(groupId string, c *cache.Cache) ([]string, error) {
	js, err := request("GET", "https://a1.easemob.com/mininetlive/mininetlive/chatgroups/"+groupId, nil, c)
	if err == nil {
		groups, err := js.Get("data").Array()
		logger.Info(groups)
		if err == nil {
			group := groups[0]
			groupInfo := group.(map[string]interface{})
			members := groupInfo["affiliations"].([]interface{})
			logger.Info(members)
			uids := make([]string, len(members)-1)
			logger.Info(reflect.TypeOf(members))
			index := 0

			for _, value := range members {
				mapValue := value.(map[string]interface{})

				for key, value := range mapValue {
					logger.Info(key, value, reflect.TypeOf(key), reflect.TypeOf(value))
					if key == "member" {
						uids[index] = value.(string)
						index++
					}
				}
			}

			return uids, nil
		} else {
			logger.Error(err)
		}
	} else {
		logger.Error(err)
	}
	return nil, err
}

func CreateGroup(owner, title string, c *cache.Cache) (string, error) {
	//create group
	url := "https://a1.easemob.com/mininetlive/mininetlive/chatgroups"
	var jsonStr = []byte(`{"groupname":"` + title + `",
	"desc":"` + title + `",
	"public":true,
	"approval":false,
	"owner":"` + owner + `",
	"maxusers":10000}`)
	js, err := request("POST", url, jsonStr, c)
	if err == nil {
		groupId, _ := js.Get("data").Get("groupid").String()
		logger.Info("group id->", groupId)
		return groupId, nil
	} else {
		return "", err
	}
}
