package upload

import (
	. "app/common"
	config "app/config"
	logger "app/logger"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type UcloudApiClient struct {
	publicKey  string
	privateKey string
	conn       *http.Client
}

type InitResult struct {
	BlkSize  int
	Bucket   string
	Key      string
	UploadId string
}

type SignParam struct {
	HttpVerb                   string
	ContentMd5                 string
	ContentType                string
	Date                       string
	CanonicalizedUCloudHeaders string
	CanonicalizedResource      string
}

type UcloudResponse struct {
	ContentLength int64
	ContentType   string
	ContentRange  string
	Etag          string
	StatusCode    int
	XsessionId    string
	RetCode       int
	ErrMsg        string
	Content       []byte
}

func newUcloudApiClient(publicKey, privateKey string) *UcloudApiClient {
	return &UcloudApiClient{publicKey, privateKey, &http.Client{Timeout: time.Minute}}
}

func signatureUFile(privateKey string, stringToSign string) string {
	mac := hmac.New(sha1.New, []byte(privateKey))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (self SignParam) String() string {
	return self.HttpVerb + "\n" +
		self.ContentMd5 + "\n" +
		self.ContentType + "\n" +
		self.Date + "\n" +
		self.CanonicalizedUCloudHeaders +
		self.CanonicalizedResource
}

func (self *UcloudApiClient) getAuthorization(param *SignParam) (authorization string) {
	return "UCloud" + " " + self.publicKey + ":" + signatureUFile(self.privateKey, fmt.Sprint(param))
}

func getURL(fileName, bucketName, httpVerb string) string {
	if httpVerb == "PUT" {
		return "http://" + bucketName + ".ufile.ucloud.cn" + "/" + fileName
	}
	return "http://" + bucketName + ".ufile.ucloud.com.cn" + "/" + fileName
}

func (self *UcloudApiClient) HeadFile(fileName, bucketName string) (int64, bool, error) {
	resp, err := self.doHttpRequest(fileName, bucketName, "HEAD")
	if err != nil {
		return 0, false, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return resp.ContentLength, true, nil
	case http.StatusNotFound:
		return 0, false, nil
	}
	return 0, false, fmt.Errorf("Internal Server Error, ucloud resp: %+v", *resp)
}

func (self *UcloudApiClient) GetFile(fileName, bucketName string) ([]byte, error) {
	resp, err := self.doHttpRequest(fileName, bucketName, "GET")
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, errors.New("content not found on ucloud")
	case http.StatusOK:
		return resp.Content, nil
	}
	return nil, fmt.Errorf("Internal Server Error, ucloud resp: %+v", *resp)
}

func (self *UcloudApiClient) PutFile(fileName, bucketName, contentType string, data []byte) (*UcloudResponse, error) {
	resp, err := self.doHttpRequest(fileName, bucketName, "PUT", contentType, string(data))
	if err != nil || resp.StatusCode != http.StatusOK {
		time.Sleep(time.Second * 1)
		resp, err := self.doHttpRequest(fileName, bucketName, "PUT", contentType, string(data))
		if err != nil {
			return resp, fmt.Errorf("Internal Server Error: %+v", err)
		}
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode != http.StatusOK {
				return resp, fmt.Errorf("Internal Server Error: %+v", err)
			}
		}
	}
	return resp, nil
}

func parseHttpResp(httpResp *http.Response, httpVerb string) (*UcloudResponse, error) {
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	resp := &UcloudResponse{}
	resp.ContentType = httpResp.Header.Get("Content-Type")
	resp.XsessionId = httpResp.Header.Get("X-SessionId")
	resp.Etag = httpResp.Header.Get("ETag")
	resp.StatusCode = httpResp.StatusCode
	resp.ContentLength = httpResp.ContentLength

	if resp.StatusCode == http.StatusOK {
		if httpVerb == "GET" {
			resp.Content = body
			return resp, nil
		}
		return resp, nil
	}
	if resp.StatusCode == http.StatusNotFound && httpVerb == "HEAD" {
		return resp, nil
	}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (self *UcloudApiClient) doHttpRequest(fileName, bucketName, httpVerb string, args ...string) (*UcloudResponse, error) {
	var httpReq *http.Request
	var err error

	url := getURL(fileName, bucketName, httpVerb)
	signParam := &SignParam{
		HttpVerb:              httpVerb,
		CanonicalizedResource: "/" + bucketName + "/" + fileName,
	}
	if httpVerb == "PUT" {
		if len(args) != 2 {
			return nil, fmt.Errorf("wrong number of arguments. Expected: %v, Got %v", 2, len(args))
		}
		contentType := args[0]
		data := []byte(args[1])
		signParam.ContentType = contentType
		httpReq, err = http.NewRequest(httpVerb, url, bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Add("Content-Type", contentType)
	} else {
		httpReq, err = http.NewRequest(httpVerb, url, nil)
		if err != nil {
			return nil, err
		}
	}
	httpReq.Header.Add("Authorization", self.getAuthorization(signParam))

	httpResp, err := self.conn.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return parseHttpResp(httpResp, httpVerb)
}

func (self *UcloudApiClient) initVideoUpload(bucketName, fileName string) (*InitResult, error) {
	url := "http://" + config.UVideoProxy + "/init/uvideoUploads/" + fileName + "?uploads"
	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Add("Bucket", bucketName)
	signParam := &SignParam{
		HttpVerb:              "POST",
		ContentType:           "application/x-www-form-urlencoded",
		CanonicalizedResource: "/" + bucketName + "/" + fileName,
	}
	httpReq.Header.Add("Authorization", self.getAuthorization(signParam))
	httpResp, err := self.conn.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	} else {
		var result InitResult
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println(err)
			return nil, err
		} else {
			return &result, nil
		}
	}
}

func (self *UcloudApiClient) multVideoUpload(bucketName, fileName, filePath string, uploadId string, blkSize int) ([]string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(blkSize)))
	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
	eTags := make([]string, totalPartsNum)
	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(blkSize), float64(fileSize-int64(i*uint64(blkSize)))))
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		url := "http://" + config.UVideoProxy + "/uvideoUploads/" + fileName + "?uploadId=" + uploadId + "&partNumber=" + strconv.FormatUint(i, 10)
		httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(partBuffer))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		httpReq.Header.Add("Content-Type", "video/mp4")
		httpReq.Header.Add("Expect", "")
		httpReq.Header.Add("Bucket", bucketName)
		signParam := &SignParam{
			HttpVerb:              "PUT",
			ContentType:           "video/mp4",
			CanonicalizedResource: "/" + bucketName + "/" + fileName,
		}
		httpReq.Header.Add("Authorization", self.getAuthorization(signParam))
		httpResp, err := self.conn.Do(httpReq)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer httpResp.Body.Close()
		body, err := ioutil.ReadAll(httpResp.Body)
		fmt.Println(string(body))
		var obj interface{} // var obj map[string]interface{}
		json.Unmarshal(body, &obj)
		m := obj.(map[string]interface{})
		fmt.Println(m)
		rPartNumber := uint64(reflect.ValueOf(m["PartNumber"]).Float())
		fmt.Println(rPartNumber)
		if rPartNumber != i {
			return nil, errors.New("unmatch partnumber")
		} else {
			eTag := httpResp.Header.Get("ETag")
			eTags[rPartNumber] = eTag
		}
	}
	return eTags, nil
}

func (self *UcloudApiClient) finishVideoUpload(bucketName, fileName, uploadId string, eTags []string) (string, error) {
	uri := "http://" + config.UVideoProxy + "/finish/uvideoUploads/" + fileName + "?uploadId=" + uploadId
	httpReq, err := http.NewRequest("POST", uri, strings.NewReader(strings.Join(eTags[:], ",")))
	if err != nil {
		return "", err
	}
	httpReq.Header.Add("Content-Type", "text/plain")
	httpReq.Header.Add("Bucket", bucketName)
	signParam := &SignParam{
		HttpVerb:              "POST",
		ContentType:           "text/plain",
		CanonicalizedResource: "/" + bucketName + "/" + fileName,
	}
	httpReq.Header.Add("Authorization", self.getAuthorization(signParam))
	httpResp, err := self.conn.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	} else {
		fmt.Println(string(body))
		// {"Bucket":"mininetliverecord","FileSize":2171328,"Key":"xiongsan_test"}
		downloadUrl := "http://" + bucketName + config.DownloadProxySuffix + "/" + url.QueryEscape(fileName) ////php rawurlencode
		return downloadUrl, nil
	}
}

func UploadImageFile(path string, fileName string) (string, error) {
	logger.Info("[Upload]", "[UploadImageFile]", "path:"+path, "fileName:"+fileName)
	u := newUcloudApiClient(
		config.UCloudPublicKey,
		config.UCloudPrivateKey,
	)
	contentType := "image/jpeg"
	bucketName := "mininetlivepub"
	data, err := ioutil.ReadFile(path)
	CheckErr("[Upload]", "[UploadToUCloudCND]", "ReadFile", err)
	resp, err := u.PutFile(fileName, bucketName, contentType, data)
	CheckErr("[Upload]", "[UploadToUCloudCND]", "upload ucloud", err)
	if err == nil {
		logger.Info("[Upload]", "[UploadToUCloudCND]", resp.StatusCode)
		logger.Info("[Upload]", "[UploadToUCloudCND]", string(resp.Content))
		return getURL(fileName, bucketName, "PUT"), nil
	} else {
		return "", err
	}
}

func UploadVideoFile(path string, fileName string) (string, error) {
	u := newUcloudApiClient(
		config.UCloudPublicKey,
		config.UCloudPrivateKey,
	)
	bucketName := config.SpaceName
	key := fileName
	resp, err := u.initVideoUpload(bucketName, key)
	if err != nil {
		logger.Info("InitVideoUpload err")
		return "", err
	} else {
		uploadId := resp.UploadId
		blkSize := resp.BlkSize
		eTags, err := u.multVideoUpload(bucketName, key, path, uploadId, blkSize)
		if err != nil {
			logger.Info("MultVideoUpload err")
			return "", err
		} else {
			downloadUrl, _ := u.finishVideoUpload(bucketName, key, uploadId, eTags)
			return downloadUrl, nil
		}
	}
}
