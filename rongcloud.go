package rongcloud

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//SERVERADDR 请求地址
const (
	SERVERADDR                  = "http://api.cn.ronghub.com"
	UserGetToken                = SERVERADDR + "/user/getToken.json"
	UserRefresh                 = SERVERADDR + "/user/refresh.json"
	UserBlock                   = SERVERADDR + "/user/block.json"
	UserUnblock                 = SERVERADDR + "/user/unblock.json"
	UserBlockQuery              = SERVERADDR + "/user/block/query.json"
	UserBlacklistAdd            = SERVERADDR + "/user/blacklist/add.json"
	UserBlacklistRemove         = SERVERADDR + "/user/blacklist/remove.json"
	UserBlacklistQuery          = SERVERADDR + "/user/blacklist/query.json"
	ConversationNotificationSet = SERVERADDR + "/conversation/notification/set.json"
	ConversationNotificationGet = SERVERADDR + "/conversation/notification/get.json"
)

//Rongcloud Rongcloud
type Rongcloud struct {
	AppKey    string
	AppSecret string
}

//NewRongcloud NewRongcloud
func NewRongcloud(AppKey, AppSecret string) (rongcloud *Rongcloud) {
	return &Rongcloud{
		AppKey:    AppKey,
		AppSecret: AppSecret,
	}
}

func (r *Rongcloud) requset(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	nonce := createNonce(10, 0)
	timestamp := fmt.Sprint(time.Now().Unix())
	req.Header.Set("App-Key", r.AppKey)
	req.Header.Set("Nonce", nonce)
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Signature", signSha1(r.AppSecret, nonce, timestamp))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//createRequset 生成请求requset参数
func createRequset(surl, method string, p map[string]string, d interface{}) (req *http.Request, err error) {
	buf, err := json.Marshal(d)
	if d != nil && err != nil {
		return nil, err
	}

	val := make(url.Values)
	for k, v := range p {
		val.Add(k, v)
	}
	ContentType := ""
	switch method {
	case "GET":
		if strings.Index(surl, "?") != -1 {
			surl += "&" + val.Encode()
		} else {
			surl += "?" + val.Encode()
		}
	case "POST":
		if d == nil {
			buf = []byte(val.Encode())
			ContentType = "application/x-www-form-urlencoded"
		} else {
			if strings.Index(surl, "?") != -1 {
				surl += "&" + val.Encode()
			} else {
				surl += "?" + val.Encode()
			}
		}
	default:
		if strings.Index(surl, "?") != -1 {
			surl += "&" + val.Encode()
		} else {
			surl += "?" + val.Encode()
		}
	}

	req, err = http.NewRequest(method, surl, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", ContentType)
	return req, err
}

func createNonce(size int, random int) string {
	iRandom, Randoms, result := random, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	iAll := random > 2 || random < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if iAll { // random ikind
			iRandom = rand.Intn(3)
		}
		scope, base := Randoms[iRandom][0], Randoms[iRandom][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func signSha1(AppSecret, Nonce, Timestamp string) (Signature string) {
	t := sha1.New()
	io.WriteString(t, AppSecret+Nonce+Timestamp)
	return fmt.Sprintf("%x", t.Sum(nil))
}
