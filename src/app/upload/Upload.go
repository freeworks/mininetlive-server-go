package upload

import(
		"net/http"
		"github.com/martini-contrib/render"
		"log"
		"os"
		"io"
		"github.com/pborman/uuid"
		"io/ioutil"
		."app/models"
		."app/common"
		"bytes"
		"crypto/hmac"
		"crypto/sha1"
		"encoding/base64"
		"encoding/json"
		"errors"
		"fmt"
		"time"
	)

func Upload(r *http.Request, render render.Render) {
	log.Println("parsing form")
	err := r.ParseMultipartForm(100000)
	// CheckErr("upload ParseMultipartForm",err)
	if err != nil {
		render.JSON(500, "server err")
	}
	file, head, err := r.FormFile("file")
	CheckErr(err, "upload Fromfile")
	log.Println(head.Filename)
	defer file.Close()
	tempDir := "/Users/cainli/dev/mininetlvie/temp/"
	filepath := tempDir + head.Filename
	fW, err := os.Create(tempDir + head.Filename)
	CheckErr(err, "create file error")
	defer fW.Close()
	_, err = io.Copy(fW, file)
	CheckErr(err, "create file error")
	fileUUID := uuid.New()
	err = UploadToUCloudCND(filepath, fileUUID, render)
	if err == nil {
		render.JSON(200, Resp{0, "上传成功", "token :" + fileUUID})
	} else {
		render.JSON(200, Resp{3001, "上传失败", nil})
	}
}

func UploadToUCloudCND(path string, fileName string, render render.Render) error {
	publicKey := "enqyjAgoDAQm0mx6A/xk8eyxEuEJWK+LQ6n258NtsT6lARMyF+YFgA=="
	privateKey := "2e3da80f079d3362f504a5db3776a9cd41feeea2"
	u := NewUcloudApiClient(
		publicKey,
		privateKey,
	)
	contentType := "image/jpeg"
	bucketName := "mininetlive"
	data, err := ioutil.ReadFile(path)
	CheckErr(err, "ReadFile")
	resp, err := u.PutFile(fileName, bucketName, contentType, data)
	CheckErr(err, "upload ucloud")
	if err == nil {
		log.Println(resp.StatusCode)
		log.Println(string(resp.Content))
	}
	return err
}


type UcloudApiClient struct {
	publicKey  string
	privateKey string
	conn       *http.Client
}

func NewUcloudApiClient(publicKey, privateKey string) *UcloudApiClient {
	return &UcloudApiClient{publicKey, privateKey, &http.Client{Timeout: time.Minute}}
}

func signatureUFile(privateKey string, stringToSign string) string {
	mac := hmac.New(sha1.New, []byte(privateKey))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type SignParam struct {
	HttpVerb                   string
	ContentMd5                 string
	ContentType                string
	Date                       string
	CanonicalizedUCloudHeaders string
	CanonicalizedResource      string
}

func (self SignParam) String() string {
	return self.HttpVerb + "\n" +
		self.ContentMd5 + "\n" +
		self.ContentType + "\n" +
		self.Date + "\n" +
		self.CanonicalizedUCloudHeaders +
		self.CanonicalizedResource
}

func (self *UcloudApiClient) genUFileAuth(param *SignParam) (authorization string) {
	return "UCloud" + " " + self.publicKey + ":" + signatureUFile(self.privateKey, fmt.Sprint(param))
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

// func (self *UcloudApiClient) PutFile(fileName, bucketName, contentType string, data []byte) error {
// 	resp, err := self.doHttpRequest(fileName, bucketName, "PUT", contentType, string(data))
// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		time.Sleep(time.Second * 1)
// 		resp, err := self.doHttpRequest(fileName, bucketName, "PUT", contentType, string(data))
// 		if err != nil {
// 			return fmt.Errorf("Internal Server Error: %+v", err)
// 		}
// 		if resp.StatusCode != http.StatusOK {
// 			if resp.StatusCode != http.StatusOK {
// 				return fmt.Errorf("Internal Server Error: %+v", err)
// 			}
// 		}
// 	}
// 	return nil
// }

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
	httpReq.Header.Add("Authorization", self.genUFileAuth(signParam))

	httpResp, err := self.conn.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return parseHttpResp(httpResp, httpVerb)
}
