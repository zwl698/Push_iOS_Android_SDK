package MiPush

import (
	"encoding/json"
	"net/url"
	"reflect"
)

type NotificationBarMessage struct {
	description             string
	title                   string
	payload                 string
	restricted_package_name string
	registration_id         string
}

type ResultMsg struct {
	Result      string          `json:"result"`
	Description string          `json:"description"`
	Data        json.RawMessage `json:"data"`
	Code        int             `json:"code"`
	Info        string          `json:"info"`
	Reason      string          `json:"reason"`
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
