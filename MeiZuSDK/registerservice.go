package meizupush

import (
	"fmt"
	"github.com/ddliu/go-httpclient"
)

const (
	getRegisterSwitch    = PUSH_API_SERVER + "/garcia/api/server/message/getRegisterSwitch"
	changeRegisterSwitch = PUSH_API_SERVER + "/garcia/api/server/message/changeRegisterSwitch"
	changeAllSwitch      = PUSH_API_SERVER + "/garcia/api/server/message/changeAllSwitch"
	subscribeAlias       = PUSH_API_SERVER + "/garcia/api/server/message/subscribeAlias"
	unSubscribeAlias     = PUSH_API_SERVER + "/garcia/api/server/message/unSubscribeAlias"
	getSubAlias          = PUSH_API_SERVER + "/garcia/api/server/message/getSubAlias"
	subscribeTags        = PUSH_API_SERVER + "/garcia/api/server/message/subscribeTags"
	unSubscribeTags      = PUSH_API_SERVER + "/garcia/api/server/message/unSubscribeTags"
	getSubTags           = PUSH_API_SERVER + "/garcia/api/server/message/getSubTags"
	unSubAllTags         = PUSH_API_SERVER + "/garcia/api/server/message/unSubAllTags"
)

const (
	APP_ID  = "100999"
	APP_KEY = "531732bc45324098978bf41c6954c09e"
	PUSH_ID = "862891030007404100999"

	NOTIFICATION_SWITCH    = "0"
	THROUGH_MESSAGE_SWITCH = "1"
)

//获取订阅开关状态
func GetRegisterSwitch(appId string, pushId string, appKey string) PushResponse {
	registerMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
	}

	sign := GenerateSign(registerMap, appKey)
	registerMap["sign"] = sign

	for k, v := range registerMap {
		fmt.Println(k, ":", v)
	}

	res, err := httpclient.Get(getRegisterSwitch, registerMap)

	return ResolvePushResponse(res, err)

}

//修改订阅开关状态
func ChangeRegisterSwitch(appId string, pushId string, msgType string, subSwitch bool, appKey string) PushResponse {
	var subSwitchinner string
	if subSwitch {
		subSwitchinner = "1"
	} else {
		subSwitchinner = "0"
	}
	changeRegisterMap := map[string]string{
		"appId":     appId,
		"pushId":    pushId,
		"msgType":   msgType,
		"subSwitch": subSwitchinner,
	}

	sign := GenerateSign(changeRegisterMap, appKey)
	changeRegisterMap["sign"] = sign

	res, err := httpclient.Post(changeRegisterSwitch, changeRegisterMap)

	return ResolvePushResponse(res, err)
}

//修改所有开关状态
func ChangeAllSwitch(appId string, pushId string, subSwitch bool, appKey string) PushResponse {
	var subSwitchinner string
	if subSwitch {
		subSwitchinner = "1"
	} else {
		subSwitchinner = "0"
	}
	changeRegisterMap := map[string]string{
		"appId":     appId,
		"pushId":    pushId,
		"subSwitch": subSwitchinner,
	}

	sign := GenerateSign(changeRegisterMap, appKey)
	changeRegisterMap["sign"] = sign

	res, err := httpclient.Post(changeAllSwitch, changeRegisterMap)

	return ResolvePushResponse(res, err)
}

//别名订阅
func SubscribeAlias(appId string, pushId string, alias string, appKey string) PushResponse {
	subscribeAliasMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
		"alias":  alias,
	}

	sign := GenerateSign(subscribeAliasMap, appKey)
	subscribeAliasMap["sign"] = sign

	res, err := httpclient.Post(subscribeAlias, subscribeAliasMap)

	return ResolvePushResponse(res, err)
}

//取消别名订阅
func UnSubscribeAlias(appId string, pushId string, appKey string) PushResponse {
	unSubscribeAliasMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
	}

	sign := GenerateSign(unSubscribeAliasMap, appKey)
	unSubscribeAliasMap["sign"] = sign

	res, err := httpclient.Post(unSubscribeAlias, unSubscribeAliasMap)

	return ResolvePushResponse(res, err)
}

//获取订阅别名
func GetSubAlias(appId string, pushId string, appKey string) PushResponse {
	getSubAliasMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
	}

	sign := GenerateSign(getSubAliasMap, appKey)
	getSubAliasMap["sign"] = sign

	res, err := httpclient.Post(getSubAlias, getSubAliasMap)

	return ResolvePushResponse(res, err)
}

//标签订阅
func SubscribeTags(appId string, pushId string, tags string, appKey string) PushResponse {
	subtagsMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
		"tags":   tags,
	}

	sign := GenerateSign(subtagsMap, appKey)
	subtagsMap["sign"] = sign

	res, err := httpclient.Post(subscribeTags, subtagsMap)

	return ResolvePushResponse(res, err)
}

//取消标签订阅
func UnSubscribeTags(appId string, pushId string, tags string, appKey string) PushResponse {
	unSubtagsMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
		"tags":   tags,
	}

	sign := GenerateSign(unSubtagsMap, appKey)
	unSubtagsMap["sign"] = sign

	res, err := httpclient.Post(unSubscribeTags, unSubtagsMap)

	return ResolvePushResponse(res, err)
}

//获取订阅标签
func GetSubTags(appId string, pushId string, appKey string) PushResponse {
	getSubtagsMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
	}

	sign := GenerateSign(getSubtagsMap, appKey)
	getSubtagsMap["sign"] = sign

	res, err := httpclient.Post(getSubTags, getSubtagsMap)

	return ResolvePushResponse(res, err)
}

//取消订阅所有标签
func UnSubAllTags(appId string, pushId string, appKey string) PushResponse {
	unSubAlltagsMap := map[string]string{
		"appId":  appId,
		"pushId": pushId,
	}

	sign := GenerateSign(unSubAlltagsMap, appKey)
	unSubAlltagsMap["sign"] = sign

	res, err := httpclient.Post(unSubAllTags, unSubAlltagsMap)

	return ResolvePushResponse(res, err)
}
