package VIVOPush

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

var TaskID map[string]string
var m_mapToken map[uint32]*tokenInfo

type tokenInfo struct {
	AuthToken    string
	TokenEndtime int64
}

//var AuthToken string
//var TokenEndtime time.Time

//func init() { //程序启动时更新token
//	UpdateAuthtoken()
//}

type NotificationBarMessage struct {
	description             string
	title                   string
	payload                 string
	restricted_package_name string
	registration_id         string
}

type VIVONormalRet struct {
	Result int    `json:"result"`
	Desc   string `json:"desc"`
	TaskId string `json:"taskId"`
}

type Auth struct {
	AppId     int    `json:"appId"`
	AppKey    string `json:"appKey"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type AuthRet struct {
	Result    int    `json:"result"`
	Desc      string `json:"desc"`
	AuthToken string `json:"authToken"`
}

type PushRet struct {
	Result int    `json:"result"`
	Desc   string `json:"desc"`
	TaskId string `json:"taskId"`
}

type PushStatisticRet struct {
	Result     int          `json:"result"`
	Desc       string       `json:"desc"`
	Statistics []Statistics `json:"statistics"`
}

type Statistics struct {
	TaskID  string `json:"taskId"`
	Send    int    `json:"send"`
	Receive int    `json:"receive"`
	Display int    `json:"display"`
	Click   int    `json:"click"`
}

func Sign(authtoken *Auth, AppSecret string) {
	var data string
	data += strconv.Itoa(authtoken.AppId)
	data += authtoken.AppKey
	data += fmt.Sprintf("%d", authtoken.Timestamp)
	data += AppSecret
	md5sum := md5.Sum([]byte(data))
	authtoken.Sign = fmt.Sprintf("%x", md5sum)
}

//将结构体转换成url.values  成功：返回url.values   失败：nil
func Struct2Urlvalues(s interface{}) *url.Values {
	val := url.Values{}
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil
	}
	tp := reflect.TypeOf(s)
	for i := 0; i < v.NumField(); i++ {
		val.Add(tp.Field(i).Name, v.Field(i).String())
	}
	return &val

}

//更新授权的token
func Gettoken(appid uint32, appkeys, appsecrets string, jjid uint32) (*tokenInfo, error) {
	tokenInstance, ok := m_mapToken[jjid]
	if !ok {
		tokenInstance = &tokenInfo{
			AuthToken:    "",
			TokenEndtime: 0,
		}
	}
	nowSeconds := time.Now().Unix()
	if tokenInstance.TokenEndtime > nowSeconds && tokenInstance.AuthToken != "" {
		return tokenInstance, nil
	}

	auth := Auth{}
	var appsecret string
	auth.AppId, auth.AppKey, appsecret = int(appid), appkeys, appsecrets
	auth.Timestamp = time.Now().UnixNano() / 1e6
	Sign(&auth, appsecret)
	reqjson, err := json.Marshal(auth)
	//fmt.Println(string(reqjson))
	if nil != err {
		return nil, err
	}
	ret, err := Post(VIVOAPI[1], string(reqjson), "")
	if nil != err {
		return nil, err
	}
	authret := AuthRet{}
	err = json.Unmarshal(ret, &authret)
	if nil != err {
		return nil, err
	}
	if len(authret.AuthToken) == 0 {
		return nil, errors.New(authret.Desc)
	}
	tokenInstance.AuthToken, tokenInstance.TokenEndtime = authret.AuthToken, time.Now().Unix()+80000000

	//fmt.Println(AuthToken)
	return tokenInstance, nil

}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return fmt.Sprint(time.Now().UnixNano())

	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
