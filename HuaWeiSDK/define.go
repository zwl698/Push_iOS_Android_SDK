package huaweipush

const (
	accessTokenAPINew            = "https://login.cloud.huawei.com/oauth2/v2/token"
	pushMessageAPINewWithVersion = "https://api.push.hicloud.com/pushsend.do?"
	pushMessageAPINew            = "https://api.push.hicloud.com/pushsend.do"
	pushMessageNew               = "openpush.message.api.send"
)

const (
	SessionTimeoutError = "session timeout"
	SessionInvalidError = "invalid session"
)

//AccessToken错误定义
const GetAccessToken_Success int32 = 200        //成功
const GetAccessToken_SrvSystemError int32 = 500 //服务器系统错误
//业务级错误嘛查询详见：https://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush_agent.html?page=hmssdk_huaweipush_api_reference_agent_s1

//PushAndroid错误定义
const PushAndroidMsg_Success string = "80000000" //成功
//其他错误：https://developer.huawei.com/consumer/cn/service/hms/catalog/huaweipush_agent.html?page=hmssdk_huaweipush_api_reference_agent_s2#%E8%A1%A85-3%20pushmsg%E5%BA%94%E7%94%A8%E7%BA%A7%E9%94%99%E8%AF%AF%E7%A0%81

//安卓消息定义
const AndroidMessageType_App int32 = 1    //透传消息
const AndroidMessageType_System int32 = 3 //系统消息
//安卓消息动作定义
const AndroidActionType_Customize int32 = 1 //自定义行为
const AndroidActionType_OpenUrl int32 = 2   //打开url
const AndroidActionType_OpenApp int32 = 3   //打开APP

type AndroidPushResult struct {
	ResultCode string `json:"code"` //80000000 成功
	ErrMessage string `json:"msg"`
	RequestID  string `json:"requestId"` //请求标识，由PUSH服务器唯一分配
	Prama      string `json:"ext"`       //扩展信息(暂未使用)
}
