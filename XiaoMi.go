package PushHttps

import (
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushHttps/MiSDK"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync/atomic"
)

var XiaoMiCount int32

func PushMsg2XiaoMi(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	request := TKDic.MiPushReq{}
	err := binary.Read(bytes.NewBuffer(data[:32]), binary.LittleEndian, &request)
	if nil != err {
		tklog.WriteErrorlog("read req data to request struct error : %s", err.Error())
		return
	}

	if header.Length != (32 + request.Msg.DWSuffixSize + request.UserList.DWSuffixSize) {
		tklog.WriteErrorlog("PushMsg2XiaoMi messageLength check fail!")
		return
	}
	pushsender := MiPush.Sender{}
	secret, ok := M_mapCertMgr[request.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", request.DWAPPID)
		return
	}
	pushsender.Restricted_package_name, pushsender.Appsecret = secret.Userinfo.XM_PackName, secret.Userinfo.XM_AppSecret
	if request.UserList.DWSuffixOffSet+request.UserList.DWSuffixSize-tkheadlen != uint32(len(data)) {
		tklog.WriteErrorlog("[Mipush] wrong msg param offset:%d  size:%d", request.Msg.DWSuffixOffSet, request.Msg.DWSuffixSize)
		return
	}
	atomic.AddInt32(&XiaoMiCount, int32(request.DWTokenCount))
	msg := string(data[request.Msg.DWSuffixOffSet-tkheadlen : request.Msg.DWSuffixOffSet+request.Msg.DWSuffixSize-tkheadlen])
	msg = Tool.ConvertToString(msg, "gbk", "UTF-8")
	msg = msg + "&" + MiPush.GetParam(1, data[request.UserList.DWSuffixOffSet-tkheadlen:request.UserList.DWSuffixOffSet+request.UserList.DWSuffixSize-tkheadlen])

	ret, err := pushsender.SendMSG(MiPush.MiAPI[1], msg)
	if err != nil {
		tklog.WriteInfolog("SendMSG err!err=%s.", err.Error())
	}
	var resultstate bool
	resultstate = MiPush.ResultS(ret, err)
	MiPush.PushSendResult(&request, resultstate)
}

func PushSFFMsg2XiaoMi(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	request := TKDic.PushXiaoMiSFFMsg{}
	err := json.Unmarshal(data, &request)
	if err != nil {
		tklog.WriteErrorlog("PushSFFMsg2MeiZu json prase err.err=%s.data=%s.", err.Error(), string(data))
		return
	}
	pushsender := MiPush.Sender{}
	secret, ok := M_mapCertMgr[request.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", request.DWAPPID)
		return
	}
	pushsender.Restricted_package_name, pushsender.Appsecret = secret.Userinfo.XM_PackName, secret.Userinfo.XM_AppSecret

	atomic.AddInt32(&XiaoMiCount, int32(request.DWTokenCount))
	msg := request.Msg
	msg = Tool.ConvertToString(msg, "gbk", "UTF-8")
	msg = msg + "&" + MiPush.GetParam(1, []byte(request.UserList))

	ret, err := pushsender.SendMSG(MiPush.MiAPI[1], msg)
	if err == nil {
		result := MiPush.ResultMsg{}
		err = json.Unmarshal(ret, &result)
		if nil != err {
			tklog.WriteErrorlog("unmarsha xiaomi result err:%s", err.Error())
			PushMsgMgr.PushDetail(request.DWAPPID, false, request.Guid, err.Error())
		}
		if result.Code > 0 { //失败，返回服务器端的错误回传的错误消息
			tklog.WriteErrorlog("push xiaomi fail!code=%d.reason=%s,desc=%s.", result.Code, result.Reason, result.Description)
			PushMsgMgr.PushDetail(request.DWAPPID, false, request.Guid, result.Description)
		} else {
			PushMsgMgr.PushDetail(request.DWAPPID, true, request.Guid, "push success")
		}
	} else {
		PushMsgMgr.PushDetail(request.DWAPPID, false, request.Guid, err.Error())
	}
}

func Testmi() {
	////ttttttt
	//var tt []string
	//m:=make(map[string]string,4)
	//m["LGame"]="1001,2"
	//m["EnterByPush"]="111"
	//tt= append(tt, "cn.jj")
	//t:=xiaomipush.NewClient("LSE3xBpxIt3GxeSkwOHuiQ==",tt)
	//var s xiaomipush.Message
	//s.Payload="123451"
	//s.Title="zhaowltest11"
	//s.Description="test11"
	//s.PassThrough=0
	//s.NotifyType=-1
	//s.RestrictedPackageName="cn.jj"
	//s.Extra=m
	//t.Send(context.TODO(),&s,"JaZlEnbZ6y2jbHmG19Mum5mbOVWmYWsUalhJE+/3SleQ8Dc87X0lyThqobkaaRZn")

	///////tttttttt

	request := TKDic.PushXiaoMiSFFMsg{}
	request.Guid = "ddddddddd"
	request.DWAPPID = 10000
	request.DWTokenCount = 1
	request.UserList = "JaZlEnbZ6y2jbHmG19Mum5mbOVWmYWsUalhJE+/3SleQ8Dc87X0lyThqobkaaRZn"
	request.Msg = `description=tt22%E6%96%B9%E6%B3%951234&extra.RealName=&extra.callback=https%3A%2F%2Fmsgdx.srv.jj.cn%2Fmsgdx%2Fapi%2FMSG%2FXiaoMi%2FPushUserXiaoMiRecptInfo&extra.callback.param=502000443&extra.callback.type=3&extra.notify_effect=1&extra.notify_foreground=1&notify_type=-1&pass_through=0&payload=tt22%E6%96%B9%E6%B3%951234&time_to_live=44148000&title=test22%E5%8F%91%E5%8F%91%E5%8F%911233`
	//	request.Msg=Tool.ConvertToString(request.Msg,"UTF-8","gbk")
	d, _ := json.Marshal(&request)
	PushSFFMsg2XiaoMi(nil, d)
}
