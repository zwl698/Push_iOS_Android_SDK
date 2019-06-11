package meizupush

import (
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/TKDic"
	"errors"
	"github.com/ddliu/go-httpclient"
	"strconv"
	"time"
)

const (
	pushThroughMessageByPushId         = PUSH_API_SERVER + "/garcia/api/server/push/unvarnished/pushByPushId"
	pushThroughMessageByPushIdNew      = PUSH_API_SERVER_NEW + "/ups/api/server/push/unvarnished/pushByPushId"
	pushNotificationMessageByPushId    = PUSH_API_SERVER + "/garcia/api/server/push/varnished/pushByPushId"
	pushNotificationMessageByPushIdNew = PUSH_API_SERVER_NEW + "/ups/api/server/push/varnished/pushByPushId"
	pushThroughMessageByAlias          = PUSH_API_SERVER + "/garcia/api/server/push/unvarnished/pushByAlias"
	pushThroughMessageByAliasNew       = PUSH_API_SERVER_NEW + "/ups/api/server/push/unvarnished/pushByAlias"
	pushNotificationMessageByAlias     = PUSH_API_SERVER + "/garcia/api/server/push/varnished/pushByAlias"
	pushNotificationMessageByAliasNew  = PUSH_API_SERVER_NEW + "/ups/api/server/push/varnished/pushByAlias"
	pushAllUserNew                     = PUSH_API_SERVER_NEW + "/ups/api/server/push/pushTask/pushToApp"
	cancelPushTask                     = PUSH_API_SERVER_NEW + "/ups/api/server/push/pushTask/cancel"
	statisticsNew                      = PUSH_API_SERVER_NEW + "/ups/api/server/push/statistics/dailyPushStatics"
)

/**
 * 通过PushId推送透传消息
 */
func PushThroughMessageByPushId(appId string, pushIds string, messageJson string, appKey string) PushResponse {
	pushThroughMessageMap := map[string]string{
		"appId":       appId,
		"pushIds":     pushIds,
		"messageJson": messageJson,
	}

	sign := GenerateSign(pushThroughMessageMap, appKey)
	pushThroughMessageMap["sign"] = sign

	res, err := httpclient.Post(pushThroughMessageByPushId, pushThroughMessageMap)

	return ResolvePushResponse(res, err)
}

func PushMeziZuBroadcast(appId, appKey, messageJson string) error {
	mt := Tool.ConvertToString(messageJson, "gbk", "utf-8")
	message := PushNotificationBcst(appId, appKey, mt)
	if message.Code != "200" || message.Message != "" {
		return errors.New(message.Message)
	}
	return nil
}

func PushNotificationBcst(appId, appKey, messageJson string) PushResponse {
	pushNotificationMessageMap := map[string]string{
		"appId":       appId,
		"pushType":    "0",
		"messageJson": messageJson,
	}

	sign := GenerateSign(pushNotificationMessageMap, appKey)
	pushNotificationMessageMap["sign"] = sign

	res, err := httpclient.Post(pushAllUserNew, pushNotificationMessageMap)
	return ResolvePushResponse(res, err)
}

func PushMeiZuMsg(appId string, pushIds string, messageJson string, appKey string, MZFrequency *uint32) (int32, error) {
	mt := Tool.ConvertToString(messageJson, "gbk", "utf-8")
	message := PushNotificationMessageByPushId(appId, pushIds, mt, appKey, MZFrequency)
	if message.Code != "200" || message.Message != "" {
		code, _ := strconv.Atoi(message.Code)
		return int32(code), errors.New(message.Code + message.Message)
	}
	return 200, nil
}

//pushId推送接口（通知栏消息）
func PushNotificationMessageByPushId(appId string, pushIds string, messageJson string, appKey string, MZFrequency *uint32) PushResponse {
	pushNotificationMessageMap := map[string]string{
		"appId":       appId,
		"pushIds":     pushIds,
		"messageJson": messageJson,
	}

	sign := GenerateSign(pushNotificationMessageMap, appKey)
	pushNotificationMessageMap["sign"] = sign

	for *MZFrequency > TKDic.MaxMeiZuPushF {
		time.Sleep(time.Second * 1)
	}
	res, err := httpclient.Post(pushNotificationMessageByPushId, pushNotificationMessageMap)

	return ResolvePushResponse(res, err)
}

//别名推送接口（透传消息
func PushThroughMessageByAlias(appId string, alias string, messageJson string, appKey string) PushResponse {
	pushThroughMessageMap := map[string]string{
		"appId":       appId,
		"alias":       alias,
		"messageJson": messageJson,
	}

	sign := GenerateSign(pushThroughMessageMap, appKey)
	pushThroughMessageMap["sign"] = sign

	res, err := httpclient.Post(pushThroughMessageByAlias, pushThroughMessageMap)

	return ResolvePushResponse(res, err)
}

//别名推送接口（通知栏消息）
func PushNotificationMessageByAlias(appId string, alias string, messageJson string, appKey string) PushResponse {
	pushNotificationMessageMap := map[string]string{
		"appId":       appId,
		"alias":       alias,
		"messageJson": messageJson,
	}

	sign := GenerateSign(pushNotificationMessageMap, appKey)
	pushNotificationMessageMap["sign"] = sign

	res, err := httpclient.Post(pushNotificationMessageByAlias, pushNotificationMessageMap)

	return ResolvePushResponse(res, err)
}
