# go-meizu-push-sdk
[![GoDoc Reference][go-doc-image]][go-doc] [![Build Status][travis-image]][travis] [![Coverage Status][coveralls-image]][coveralls] [![Go Report Card][go-report-image]][go-report] [![License][license-image]][license]

## Installation
``` 
   go get github.com/MEIZUPUSH/go-meizu-push-sdk
```

## 订阅相关接口

* API 列表

```
    func ChangeAllSwitch(appId string, pushId string, subSwitch bool, appKey string) PushResponse
    func ChangeRegisterSwitch(appId string, pushId string, msgType string, subSwitch bool, appKey string) PushResponse
    func GetRegisterSwitch(appId string, pushId string, appKey string) PushResponse
    func GetSubAlias(appId string, pushId string, appKey string) PushResponse
    func GetSubTags(appId string, pushId string, appKey string) PushResponse
    func Register(appId string, appKey string, deviceId string) PushResponse
    func ResolvePushResponse(res *httpclient.Response, err error) PushResponse
    func SubscribeAlias(appId string, pushId string, alias string, appKey string) PushResponse
    func SubscribeTags(appId string, pushId string, tags string, appKey string) PushResponse
    func UnSubAllTags(appId string, pushId string, appKey string) PushResponse
    func UnSubscribeAlias(appId string, pushId string, appKey string) PushResponse
    func UnSubscribeTags(appId string, pushId string, tags string, appKey string) PushResponse

```
**NOTE:** 以test单元测试的方法说明其中一种api的使用方式

* 获取订阅开关状态

```go
func TestGetRegisterSwitch(t *testing.T) {
	message := GetRegisterSwitch(APP_ID,PUSH_ID,APP_KEY)
	fmt.Println("current message "+message.message)
	if(message.code != 200){
		t.Error("Status Code not 200")
	}
}

```

## 推送相关接口

* API 列表

```
    func PushNotificationMessageByPushId(appId string, pushIds string, messageJson string, appKey string) PushResponse
    func PushThroughMessageByPushId(appId string, pushIds string, messageJson string, appKey string) PushResponse
```
**NOTE:** 以test单元测试的方法说明其中一种api的使用方式

* 推送透传消息

```go
func TestPushThroughMessageByPushId(t *testing.T) {
	messageJson := `{"test_throught_message": "message"}`
	message := PushThroughMessageByPushId(APP_ID,PUSH_ID,buildThroughMessage(messageJson),APP_KEY)

	fmt.Println("TestPushThroughMessageByPushId ",message.message)

	if message.code != 200 {
		t.Error("Status Code not 200")
	}
}    
```

* 推送通知栏消息

```go
func TestPushNotificationMessageByPushId(t *testing.T) {
    //使用通知栏构建工具方法快速构建通知栏json
	json := BuildNotificationMessage().
		noticeBarType(2).
		noticeTitle("标题go").
		noticeContent("测试内容").toJson()
	message := PushNotificationMessageByPushId(APP_ID,PUSH_ID,json,APP_KEY)

	fmt.Println("TestPushNotificationMessageByPushId ",message.message)

	if message.code != 200 {
		t.Error("Status Code not 200")
	}

}
```

**NOTE:**  详情请参考[meizu-push-godoc](https://godoc.org/github.com/MEIZUPUSH/go-meizu-push-sdk)


[travis]: https://travis-ci.org/comsince/go-meizu-push-sdk
[travis-image]: https://travis-ci.org/comsince/go-meizu-push-sdk.svg?branch=master

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[coveralls-image]: https://coveralls.io/repos/github/comsince/go-meizu-push-sdk/badge.svg?branch=master
[coveralls]: https://coveralls.io/github/comsince/go-meizu-push-sdk?branch=master

[go-doc-image]:https://godoc.org/github.com/mattn/go-sqlite3?status.svg
[go-doc]:https://godoc.org/github.com/MEIZUPUSH/go-meizu-push-sdk


[go-report-image]:https://goreportcard.com/badge/github.com/MEIZUPUSH/go-meizu-push-sdk
[go-report]:https://goreportcard.com/report/github.com/MEIZUPUSH/go-meizu-push-sdk