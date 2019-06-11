package MiPush

var MiAPI map[int]string

func init() {
	MiAPI = make(map[int]string)
	MiAPI[0] = ""                                                     //
	MiAPI[1] = "https://api.xmpush.xiaomi.com/v3/message/regid"       //向regid发消息
	MiAPI[2] = "https://api.xmpush.xiaomi.com/v3/message/alias"       //向别名发消息
	MiAPI[3] = "https://api.xmpush.xiaomi.com/v3/message/topic"       //向topic发消息
	MiAPI[4] = "https://api.xmpush.xiaomi.com/v3/message/multi_topic" //向多个topic发消息
	MiAPI[5] = "https://api.xmpush.xiaomi.com/v3/message/all"         //向所有设备发消息        //上述API皆为发送单条消息的API

}
