package huaweipush

import (
	"encoding/json"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

type HWToken struct {
	AccessToken      string `json:"access_token"`
	ExpireIn         int64  `json:"expires_in"` // expires_in秒后token过期
	ExpireAt         int64  `json:"expires_at"` //不是华为的返回字段
	Scope            string `json:"scope"`
	Error            int32  `json:"error"`
	ErrorDescription string `json:"error_description"`
}

//var tokenInstance *HWToken
var m_HWToken map[uint32]*HWToken

func init() {
	m_HWToken = make(map[uint32]*HWToken)
}

func RequestAccess(clientID, clientSecret string, appid uint32, ReGetToken bool) (*HWToken, error) {
	var bnew bool = false
	tokenInstance, ok := m_HWToken[appid]
	if !ok {
		tokenInstance = &HWToken{
			AccessToken: "",
			ExpireAt:    0,
		}
		bnew = true
	}
	nowSeconds := time.Now().Unix()
	if !ReGetToken {
		if tokenInstance.ExpireAt > nowSeconds && tokenInstance.AccessToken != "" {
			return tokenInstance, nil
		}
	}

	form := url.Values{}
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)
	form.Add("grant_type", "client_credentials")
	bytes, err, _ := doPost(context.Background(), accessTokenAPINew, form)
	if err != nil {
		return nil, err
	}
	var newToken HWToken
	err = json.Unmarshal(bytes, &newToken)
	if err != nil {
		return nil, err
	}
	newToken.ExpireAt = nowSeconds + newToken.ExpireIn - 600
	tokenInstance = &newToken
	// invalid the token
	time.AfterFunc(time.Second*time.Duration(tokenInstance.ExpireIn), func() {
		tokenInstance.AccessToken = ""
	})
	if bnew {
		m_HWToken[appid] = tokenInstance
	}

	return tokenInstance, nil
}
