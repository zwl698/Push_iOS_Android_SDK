package meizupush

import "github.com/ddliu/go-httpclient"

const getTaskStatistics = PUSH_API_SERVER + "/garcia/api/server/push/statistics/getTaskStatistics"

//获取任务推送统计
func GetTaskStatistics(appId string, taskId string, appKey string) PushResponse {
	taskStatisticsMap := map[string]string{
		"appId":  appId,
		"taskId": taskId,
	}

	sign := GenerateSign(taskStatisticsMap, appKey)
	taskStatisticsMap["sign"] = sign

	res, err := httpclient.Post(getTaskStatistics, taskStatisticsMap)

	return ResolvePushResponse(res, err)
}
