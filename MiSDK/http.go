package MiPush

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//Post post request to xiaomi server
func Post(url string, msg string, appsecret string) ([]byte, error) {

	client := http.Client{Timeout: time.Second * 15} //设置超时时间为15秒
	request, err := http.NewRequest("POST", url, strings.NewReader(msg))
	if nil != err {
		return nil, err
	}
	secret := "key=" + appsecret
	request.Header.Set("Authorization", secret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var redotimes = 0
redo:
	response, err := client.Do(request)
	redotimes++
	if nil != err {
		if redotimes < 3 {
			client.Timeout = time.Second * 7
			goto redo
		}
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(data))
	if nil != err {
		return nil, err
	}
	return data, nil

}
