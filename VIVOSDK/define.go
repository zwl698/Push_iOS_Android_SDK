package VIVOPush

import "time"

var VIVOAPI map[int]string
var InterfaceEnable map[int]time.Time //接口编号-》上次接口使用时间    用以限制高频次使用

func init() {
	VIVOAPI = make(map[int]string)
	VIVOAPI[0] = "https://api-push.vivo.com.cn/"                              //vivo 服务器的地址
	VIVOAPI[1] = "https://api-push.vivo.com.cn/message/auth"                  //鉴权接口（获取authtoken）
	VIVOAPI[2] = "https://api-push.vivo.com.cn/message/send"                  //单条消息发送接口
	VIVOAPI[3] = "https://api-push.vivo.com.cn/message/saveListPayload"       //广播群推接口
	VIVOAPI[4] = "https://api-push.vivo.com.cn/message/pushToList"            //群推时发送用户的接口
	VIVOAPI[5] = ""                                                           //"https://api-push.vivo.com.cn/message/all"                   //全量发送接口
	VIVOAPI[6] = "https://api-push.vivo.com.cn/report/getStatistics?taskIds=" //消息推送统计接口
	//	VIVOAPI[7]="https://api-push.vivo.com.cn/report/getStatistics"		//回执接收接口     //又第三方自己指定回执接收地址

	InterfaceEnable = make(map[int]time.Time) //vivo有超量限制，当接口使用超量后，为避免接口推送量降级，当日不再高频使用该接口

}
