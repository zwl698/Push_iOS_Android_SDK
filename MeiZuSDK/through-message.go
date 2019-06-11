package meizupush

import "encoding/json"

//{
//    "title": 推送标题,
//    "content": 推送内容, json 格式
//    "pushTimeInfo": {
//        "offLine": 是否进离线消息 0 否 1 是[validTime]
//        "validTime": 有效时长 (0- 72 小时内的正整数)
//    }
//    "advanceInfo": {
//        "fixSpeed": 是否定速推送 0 否  1 是 (fixSpeedRate 定速速率)
//        "fixSpeedRate": 定速速率
//    }
//}

//这里所有的字段都是头字母大写
type ThroughMessage struct {
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	PushTimeInfo T_PushTimeInfo `json:"pushTimeInfo"`
	AdvanceInfo  T_AdvanceInfo  `json:"advanceInfo"`
}

type T_PushTimeInfo struct {
	OffLine   int `json:"offLine"`
	ValidTime int `json:"validTime"`
}

type T_AdvanceInfo struct {
	FixSpeed     int `json:"fixSpeed"`
	FixSpeedRate int `json:"fixSpeedRate"`
}

// 构建透传消息,目前只开放content字段设置，其他的默认,throughMessage必须是json格式
func buildThroughMessage(throughMessage string) string {
	var message ThroughMessage
	message.Title = "默认透传消息标题"
	message.Content = throughMessage
	message.PushTimeInfo = T_PushTimeInfo{OffLine: 1, ValidTime: 2}
	message.AdvanceInfo = T_AdvanceInfo{FixSpeed: 0, FixSpeedRate: 0}
	throughMessageJson, _ := json.Marshal(message)
	return string(throughMessageJson)
}
