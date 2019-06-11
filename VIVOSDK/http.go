package VIVOPush

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//Post post request to xiaomi server
func Post(url string, msg string, token string) ([]byte, error) {

	client := http.Client{Timeout: time.Second * 15} //设置超时时间为15秒
	request, err := http.NewRequest("POST", url, strings.NewReader(msg))
	if nil != err {
		return nil, err
	}
	if token != "" { //授权时
		request.Header.Set("authToken", token)
	}
	request.Header.Set("Content-Type", "application/json")
	var redotimes = 0
redo:
	response, err := client.Do(request)
	redotimes++
	if nil != err {
		client.Timeout = time.Second * 7
		if redotimes < 3 {
			goto redo
		}
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)

	if nil != err {
		return nil, err
	}
	return data, nil

}

func Get(url string, token string) ([]byte, error) {

	client := http.Client{Timeout: time.Second * 15} //设置超时时间为15秒
	request, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return nil, err
	}
	if token != "" { //授权时
		request.Header.Set("authToken", token)
	}
	request.Header.Set("Content-Type", "application/json")
	var redotimes = 0
redo:
	response, err := client.Do(request)
	if nil != err {
		client.Timeout = time.Second * 7
		if redotimes < 3 {
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
