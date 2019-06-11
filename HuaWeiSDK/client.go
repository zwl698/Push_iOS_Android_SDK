package huaweipush

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"

	"TKGoBase/IO/Log"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/TKDic"
	"fmt"
	"golang.org/x/net/context"
	"strings"
)

type HuaweiPushClient struct {
	clientID, clientSecret string
}

func NewClient(clientID, clientSecret string) *HuaweiPushClient {
	return &HuaweiPushClient{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *HuaweiPushClient) defaultParams(params url.Values, appid uint32) (url.Values, error) {
	accessToken, err := RequestAccess(c.clientID, c.clientSecret, appid, false)
	if err != nil {
		return params, err
	}
	if accessToken.Error != 0 { //有可能是0，得测一下
		return params, errors.New(accessToken.ErrorDescription)
	}

	accessToken.AccessToken = strings.Replace(accessToken.AccessToken, "\\", "", -1)
	params.Add("access_token", accessToken.AccessToken)
	params.Add("nsp_svc", pushMessageNew)
	params.Add("nsp_ts", strconv.FormatInt(time.Now().Unix(), 10))

	return params, nil
}

func (c *HuaweiPushClient) setPushPrama(tokenlist string, outdate time.Time, msg string, appid uint32) (url.Values, error) {
	params := url.Values{}
	params, err := c.defaultParams(params, appid)
	if err != nil {
		return params, err
	}

	date := outdate.Format("2006-01-02T15:04")
	tt := Tool.ConvertToString(msg, "gbk", "UTF-8")
	params.Add("device_token_list", tokenlist)
	params.Add("payload", tt)
	params.Add("expire_time", date)

	return params, nil
}

func (c *HuaweiPushClient) Push(ctx context.Context, tokenlist string, outdate time.Time, msg string, HWFerquency *uint32, appid uint32, huaweiappid string) (error, int32) {
	params, err := c.setPushPrama(tokenlist, outdate, msg, appid)
	if err != nil {
		return err, 0
	}

	for *HWFerquency > TKDic.MaxHuaWeiPushF {
		time.Sleep(time.Second * 1)
	}
	extrap := fmt.Sprintf("{\"ver\":\"1\",\"appId\":\"%s\"}", huaweiappid)
	extrap = fmt.Sprintf("nsp_ctx=%s", url.QueryEscape(extrap))
	pushaddr := fmt.Sprintf("%s%s", pushMessageAPINewWithVersion, extrap)
	bytes, err, code := doPost(ctx, pushaddr, params)
	if err != nil || code == 503 {
		return err, code
	}

	var result AndroidPushResult
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		tklog.WriteInfolog(string(bytes))
		return err, code
	}

	if result.ResultCode == PushAndroidMsg_Success {
		return nil, code
	} else {
		if result.ErrMessage == SessionTimeoutError || result.ErrMessage == SessionInvalidError || result.ResultCode == "80300007" {
			tklog.WriteErrorlog("huawei token error.code=%s,msg=%s.", result.ResultCode, result.ErrMessage)
			RequestAccess(c.clientID, c.clientSecret, appid, true)
			return errors.New(result.ErrMessage), code
		}
	}
	return errors.New(result.ErrMessage), code
}
