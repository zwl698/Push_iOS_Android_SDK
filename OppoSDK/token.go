package oppopush

import (
	"TKGoBase/IO/Log"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var m_mapToken map[uint32]*token

func init() {
	m_mapToken = make(map[uint32]*token)
}

func RequestAccess(AppKey, mastersecret string, appid uint32, ReApplytoken bool) (*token, error) {
	nowSeconds := time.Now().UnixNano() / 1e6

	var bnew bool = false
	tokenInstance, ok := m_mapToken[appid]
	if !ok {
		tokenInstance = &token{
			authtoken: "",
			endtime:   0,
		}
		bnew = true
	}
	if !ReApplytoken {
		if tokenInstance.endtime > nowSeconds && tokenInstance.authtoken != "" {
			return tokenInstance, nil
		}
	}

	strtime := strconv.FormatInt(nowSeconds, 10)
	form := url.Values{}
	form.Add("app_key", AppKey)
	form.Add("sign", getsha256(AppKey, strtime, mastersecret))
	form.Add("timestamp", strtime)

	bytes, err := doPost(context.Background(), access_token, form)
	if err != nil {
		return nil, err
	}
	var newToken tokenrespon

	err = json.Unmarshal(bytes, &newToken)
	if err != nil {
		return nil, err
	}

	if newToken.Code != 0 {
		return nil, errors.New(newToken.Message)
	}

	if newToken.Data.Auth_token == "" {
		return nil, errors.New("token is nil!")
	}

	tokenInstance.authtoken, tokenInstance.endtime = newToken.Data.Auth_token, newToken.Data.Create_time+80000000
	if bnew {
		m_mapToken[appid] = tokenInstance
	}
	return tokenInstance, nil
}

func getsha256(appkey, time, mastersecret string) string {
	h := sha256.New()
	tt := appkey + time + mastersecret

	h.Write([]byte(tt))
	return hex.EncodeToString(h.Sum(nil))
}

func doPost(ctx context.Context, url string, form url.Values) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	tryTime := 0
tryAgain:
	res, err = ctxhttp.Do(ctx, client, req)
	if err != nil {
		tklog.WriteErrorlog("oppo push post err:%s,trytimes:%d", err.Error(), tryTime)
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}
		tryTime += 1
		if tryTime < 3 {
			goto tryAgain
		}
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	str := string(result)
	str, err = strconv.Unquote(str)
	if err != nil {
		str = string(result)
	}
	return []byte(str), nil
}
