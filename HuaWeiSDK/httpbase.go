package huaweipush

import (
	"TKGoBase/IO/Log"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func doPost(ctx context.Context, url string, form url.Values) ([]byte, error, int32) {
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
		tklog.WriteErrorlog("huawei push post err:%s,trytime=%d.", err.Error(), tryTime)
		select {
		case <-ctx.Done():
			return nil, err, int32(res.StatusCode)
		default:
		}
		tryTime += 1
		if tryTime < 3 {
			goto tryAgain
		}
		return nil, err, int32(res.StatusCode)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err, int32(res.StatusCode)
	}
	str := string(result)
	str, err = strconv.Unquote(str)
	if err != nil {
		str = string(result)
	}
	return []byte(str), nil, int32(res.StatusCode)
}
