package huaweipush

import "encoding/json"

type AndroidMessagePyload struct {
	HPush string `json:"hps"`
	AndroidHPush
}

type AndroidHPush struct {
	androidMsg `json:"hps.msg"`
	androidext `json:"hps.ext"`
}

type androidMsg struct {
	msgtype       int32                   `json:"hps.msg.type"` //1为透传异步消息，3为系统通知异步消息
	androidbody   `json:"hps.msg.body"`   //透传消息可以为字符串，不必是json对象
	androidaction `json:"hps.msg.action"` //消息点击动作
}

type androidext struct {
	biTag string `json:"hps.ext.biTag"`
}

type androidbody struct {
	content string `json:"hps.msg.body.content"`
	title   string `json:"hps.msg.body.title"`
}

type androidaction struct {
	actiontype          int32 `json:"hps.msg.action.type"` //1为自定义行为，行为由intent定义，2为打开url由url定义，3为打开APP
	androidactionparama `json:"hps.msg.action.parama"`
}
type androidactionparama struct {
	intent     string `json:"hps.msg.action.parama.intent"`
	url        string `json:"hps.msg.action.parama.url"`
	appPkgName string `json:"hps.msg.action.parama.appPkgName"` //需要拉起的应用包名，必须和注册推送的包名一致
}

func NewAndroidSysMessage(actiontype int32, title, content, parama, bitag string) *AndroidMessagePyload {
	var st AndroidMessagePyload
	st.msgtype, st.actiontype, st.title, st.content, st.biTag = AndroidMessageType_System, actiontype, title, content, bitag
	switch actiontype {
	case AndroidActionType_Customize:
		st.intent = parama
	case AndroidActionType_OpenUrl:
		st.url = parama
	case AndroidActionType_OpenApp:
		st.appPkgName = parama
	default:
		return nil
	}
	return &st
}

func (m *AndroidMessagePyload) ToJson() (string, error) {
	tl, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(tl), nil
}
