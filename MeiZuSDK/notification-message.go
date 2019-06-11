package meizupush

import (
	"encoding/json"
)

//"noticeBarInfo": {
//     "noticeBarType": 通知栏样式(0, "标准"),(2, "安卓原生")【int 非必填，值为0】
//     "title": 推送标题, 【string 必填，字数限制1~32字符】
//     "content": 推送内容, 【string 必填，字数限制1~100字符】
//}
type NoticeBarInfo struct {
	NoticeBarType int    `json:"noticeBarType"`
	Title         string `json:"title"`
	Content       string `json:"content"`
}

//"noticeExpandInfo": {
//        "noticeExpandType": 展开方式 (0, "标准"),(1, "文本")【int 非必填，值为0、1】
//        "noticeExpandContent": 展开内容, 【string noticeExpandType为文本时，必填】
//    }

type NoticeExpandInfo struct {
	NoticeExpandType    int    `json:"noticeExpandType"`
	NoticeExpandContent string `json:"noticeExpandContent"`
}

//"clickTypeInfo": {
//        "clickType": 点击动作 (0,"打开应用"),(1,"打开应用页面"),(2,"打开URI页面"),(3, "应用客户端自定义")【int 非必填,默认为0】
//        "url": URI页面地址, 【string clickType为打开URI页面时，必填, 长度限制1000字节】
//        "parameters":参数 【JSON格式】【非必填】
//        "activity":应用页面地址 【string clickType为打开应用页面时，必填, 长度限制1000字节】
//        "customAttribute":应用客户端自定义【string clickType为应用客户端自定义时，必填， 输入长度为1000字节以内】
//    }
type ClickTypeInfo struct {
	ClickType       int                    `json:"clickType"`
	Url             string                 `json:"url"`
	Parameters      map[string]interface{} `json:"parameters"`
	Activity        string                 `json:"activity"`
	CustomAttribute string                 `json:"customAttribute"`
}

//"pushTimeInfo": {
//        "offLine": 是否进离线消息(0 否 1 是[validTime]) 【int 非必填，默认值为1】
//        "validTime": 有效时长 (1到72 小时内的正整数) 【int offLine值为1时，必填，默认24】
//    }
type PushTimeInfo struct {
	OffLine   int `json:"offLine"`
	ValidTime int `json:"validTime"`
}

//"advanceInfo": {
//        "suspend":是否通知栏悬浮窗显示 (1 显示  0 不显示) 【int 非必填，默认1】
//        "clearNoticeBar":是否可清除通知栏 (1 可以  0 不可以) 【int 非必填，默认1】
//        "fixDisplay":是否定时展示 (1 是  0 否) 【int 非必填，默认0】
//        "fixStartDisplayTime": 定时展示开始时间(yyyy-MM-dd HH:mm:ss) 【str 非必填】
//        "fixEndDisplayTime ": 定时展示结束时间(yyyy-MM-dd HH:mm:ss) 【str 非必填】
//        "notificationType": {
//            "vibrate":  震动 (0关闭  1 开启) ,  【int 非必填，默认1】
//            "lights":   闪光 (0关闭  1 开启), 【int 非必填，默认1】
//            "sound":   声音 (0关闭  1 开启), 【int 非必填，默认1】
//        }
//    }
type AdvanceInfo struct {
	Suspend             int              `json:"suspend"`
	ClearNoticeBar      int              `json:"clearNoticeBar"`
	FixDisplay          int              `json:"fixDisplay"`
	FixStartDisplayTime string           `json:"fixStartDisplayTime"`
	FixEndDisplayTime   string           `json:"fixEndDisplayTime"`
	NotificationType    NotificationType `json:"notificationType"`
}

type NotificationType struct {
	Vibrate int `json:"vibrate"`
	Lights  int `json:"lights"`
	Sound   int `json:"sound"`
}

type NotificationMessage struct {
	NoticeBarInfo    NoticeBarInfo    `json:"noticeBarInfo"`
	NoticeExpandInfo NoticeExpandInfo `json:"noticeExpandInfo"`
	ClickTypeInfo    ClickTypeInfo    `json:"clickTypeInfo"`
	PushTimeInfo     PushTimeInfo     `json:"pushTimeInfo"`
	AdvanceInfo      AdvanceInfo      `json:"advanceInfo"`
}

func BuildNotificationMessage() NotificationMessage {
	var notification NotificationMessage
	notification.NoticeBarInfo = NoticeBarInfo{2, "", ""}
	notification.NoticeExpandInfo = NoticeExpandInfo{0, ""}
	clickParameterMap := map[string]interface{}{"test": "test"}
	notification.ClickTypeInfo = ClickTypeInfo{0, "", clickParameterMap, "", ""}
	notification.PushTimeInfo = PushTimeInfo{0, 24}
	notificationType := NotificationType{1, 1, 1}
	notification.AdvanceInfo = AdvanceInfo{0, 0, 0, "", "", notificationType}
	return notification
}

func (message NotificationMessage) noticeBarType(barType int) NotificationMessage {
	message.NoticeBarInfo.NoticeBarType = barType
	return message
}

func (message NotificationMessage) noticeTitle(title string) NotificationMessage {
	message.NoticeBarInfo.Title = title
	return message
}

func (message NotificationMessage) noticeContent(content string) NotificationMessage {
	message.NoticeBarInfo.Content = content
	return message
}

func (message NotificationMessage) noticeExpandType(noticeExpandType int) NotificationMessage {
	message.NoticeExpandInfo.NoticeExpandType = noticeExpandType
	return message
}

func (message NotificationMessage) noticeExpandContent(noticeExpandContent string) NotificationMessage {
	message.NoticeExpandInfo.NoticeExpandContent = noticeExpandContent
	return message
}

func (message NotificationMessage) noticeClickType(clickType int) NotificationMessage {
	message.ClickTypeInfo.ClickType = clickType
	return message
}

func (message NotificationMessage) noticeClickUrl(url string) NotificationMessage {
	message.ClickTypeInfo.Url = url
	return message
}

func (message NotificationMessage) noticeClickParams(params map[string]interface{}) NotificationMessage {
	message.ClickTypeInfo.Parameters = params
	return message
}

func (message NotificationMessage) noticeClickActivity(activity string) NotificationMessage {
	message.ClickTypeInfo.Activity = activity
	return message
}

func (message NotificationMessage) noticeClickCustomAttribute(attribute string) NotificationMessage {
	message.ClickTypeInfo.CustomAttribute = attribute
	return message
}

func (message NotificationMessage) offLine(offline int) NotificationMessage {
	message.PushTimeInfo.OffLine = offline
	return message
}

func (message NotificationMessage) validTime(validTime int) NotificationMessage {
	message.PushTimeInfo.ValidTime = validTime
	return message
}

func (message NotificationMessage) suspend(suspend int) NotificationMessage {
	message.AdvanceInfo.Suspend = suspend
	return message
}

func (message NotificationMessage) clearNoticeBar(clear bool) NotificationMessage {
	message.AdvanceInfo.ClearNoticeBar = bool2Int(clear)
	return message
}

func (message NotificationMessage) fixDisplay(fixDisplay bool) NotificationMessage {
	message.AdvanceInfo.FixDisplay = bool2Int(fixDisplay)
	return message
}

func (message NotificationMessage) fixStartDisplayTime(fixStartDisplayTime string) NotificationMessage {
	message.AdvanceInfo.FixStartDisplayTime = fixStartDisplayTime
	return message
}

func (message NotificationMessage) fixEndDisplayTime(fixEndDisplayTime string) NotificationMessage {
	message.AdvanceInfo.FixEndDisplayTime = fixEndDisplayTime
	return message
}

func (message NotificationMessage) vibrate(vibrate bool) NotificationMessage {
	message.AdvanceInfo.NotificationType.Vibrate = bool2Int(vibrate)
	return message
}

func (message NotificationMessage) lights(lights bool) NotificationMessage {
	message.AdvanceInfo.NotificationType.Lights = bool2Int(lights)
	return message
}

func (message NotificationMessage) sound(sound bool) NotificationMessage {
	message.AdvanceInfo.NotificationType.Sound = bool2Int(sound)
	return message
}

func (message NotificationMessage) toJson() string {
	j, err := json.Marshal(message)
	if err != nil {
		return ""
	}
	return string(j)
}

func bool2Int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
