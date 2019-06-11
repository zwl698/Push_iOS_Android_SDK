package PushHttps

import (
	"TKGoBase/IO/Coon/MsgMgr"
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushHttps/VIVOSDK"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync/atomic"
)

var VivoCount int32
var RecptVivoCount int32

const tkheadlen = 24

func PushMsg2Vivo(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	request := TKDic.VIVOPushReq{}
	err := binary.Read(bytes.NewBuffer(data[:32]), binary.LittleEndian, &request)
	if nil != err {
		tklog.WriteErrorlog("[vivo] read req data to request struct error : %s", err.Error())
		return
	}
	if header.Length != (32 + request.RegId.DWSuffixSize + request.TaskId.DWSuffixSize) {
		tklog.WriteErrorlog("wrong request size ,offset:%d  suflen:%d   reqlen:%d", request.TaskId.DWSuffixOffSet, request.TaskId.DWSuffixSize, len(data))
		return
	}

	atomic.AddInt32(&VivoCount, int32(request.DWTokenCount))
	msg := "{\"regIds\":%s,\"taskId\":\"%s\",\"requestId\":\"%s\"}"
	msg = fmt.Sprintf(msg, data[request.RegId.DWSuffixOffSet-tkheadlen:request.RegId.DWSuffixOffSet+request.RegId.DWSuffixSize-tkheadlen],
		data[request.TaskId.DWSuffixOffSet-tkheadlen:request.TaskId.DWSuffixOffSet+request.TaskId.DWSuffixSize-tkheadlen],
		VIVOPush.GetGuid())

	sender := VIVOPush.Sender{}
	api, ok := VIVOPush.VIVOAPI[4]
	if !ok {
		tklog.WriteErrorlog("[vivo]wrong apiid:%d", 4)
		return
	}
	sender.Api = api

	secret, ok := M_mapCertMgr[request.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", request.DWAPPID)
		return
	}

	token, err := VIVOPush.Gettoken(secret.Userinfo.VIVO_Appid, secret.Userinfo.VIVO_AppKey, secret.Userinfo.VIVO_AppSecRet, request.DWAPPID)
	if nil != err {
		tklog.WriteErrorlog("update vivo token error :%s", err.Error())
		return
	}

	sender.Authtoken = token.AuthToken
	ret, err := sender.SendMSG(msg)

	result := VIVOPush.VIVONormalRet{}
	if nil != err {
		tklog.WriteErrorlog("Push user list to vivo error :%s", err.Error())
		return
	}
	er := json.Unmarshal(ret, &result)
	if nil != er {
		tklog.WriteErrorlog("analysis Push vivo user list return value error :%s", er.Error())
		return
	}
	if result.Result > 0 {
		tklog.WriteErrorlog("Vivo Push fail.err=%s.ID=%d.", result.Desc, result.Result)
		VIVOPush.PushSendResult(&request, false)
	} else {
		VIVOPush.PushSendResult(&request, true)
	}

}

func SaveMsg2Vivo(header *TKBaseDic.TKheader, data []byte, _conn *net.TCPConn) {
	defer Tool.PanicRecover()
	if header.Length < 4 {
		tklog.WriteErrorlog("SaveMsg2Vivo headerlength err.length=%d.", header.Length)
		return
	}
	var dwappid uint32
	err := binary.Read(bytes.NewBuffer(data[:4]), binary.LittleEndian, &dwappid)
	if err != nil {
		tklog.WriteErrorlog("SaveMsg2Vivo anysize fail!,err=%s", err.Error())
		return
	}
	secret, ok := M_mapCertMgr[dwappid]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", dwappid)
		return
	}

	msg := string(data[4:])
	msg = Tool.ConvertToString(msg, "gbk", "UTF-8")
	sender := VIVOPush.Sender{}
	api, ok := VIVOPush.VIVOAPI[3]
	if !ok {
		tklog.WriteErrorlog("[vivo]wrong apiid:%d", 3)
		VIVOPush.AckVivomsgid("error", _conn)
		return
	}
	sender.Api = api

	token, err := VIVOPush.Gettoken(secret.Userinfo.VIVO_Appid, secret.Userinfo.VIVO_AppKey, secret.Userinfo.VIVO_AppSecRet, dwappid)
	if nil != err {
		tklog.WriteErrorlog("update vivo token error :%s", err.Error())
		VIVOPush.AckVivomsgid("error", _conn)
		return
	}

	tklog.WriteInfolog("vivo save msg:%s", msg)
	sender.Authtoken = token.AuthToken
	ret, err := sender.SendMSG(msg)

	result := VIVOPush.VIVONormalRet{}
	if nil != err {
		tklog.WriteErrorlog("Save MSG to vivo error :%s", err.Error())
		VIVOPush.AckVivomsgid("error", _conn)
		return
	}
	er := json.Unmarshal(ret, &result)
	if nil != er {
		tklog.WriteErrorlog("analysis Save MSG to vivo 's return value error :%s.result=%s.", er.Error(), string(ret))
		VIVOPush.AckVivomsgid("error", _conn)
		return
	}
	if result.Result == 0 {
		VIVOPush.AckVivomsgid(result.TaskId, _conn)
	} else {
		tklog.WriteErrorlog("save vivomsg fail.err=%s.code=%d.", result.Desc, result.Result)
		VIVOPush.AckVivomsgid("error", _conn)
	}
}

func PushSFFMsg2Vivo(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	request := TKDic.PushVivoSFFMsg{}
	err := json.Unmarshal(data, &request)
	if err != nil {
		tklog.WriteErrorlog("PushSFFMsg2Vivo json prase err.err=%s.data=%s.", err.Error(), string(data))
		return
	}

	atomic.AddInt32(&VivoCount, int32(request.DWTokenCount))
	msg := fmt.Sprintf("{\"regId\":\"%s\",%s,\"requestId\":\"%s\"}", request.RegId, request.Payload, request.Guid)

	sender := VIVOPush.Sender{}
	api, ok := VIVOPush.VIVOAPI[2]
	if !ok {
		tklog.WriteErrorlog("[vivo]wrong apiid:%d", 4)
		return
	}
	sender.Api = api

	secret, ok := M_mapCertMgr[request.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", request.DWAPPID)
		return
	}

	token, err := VIVOPush.Gettoken(secret.Userinfo.VIVO_Appid, secret.Userinfo.VIVO_AppKey, secret.Userinfo.VIVO_AppSecRet, request.DWAPPID)
	if nil != err {
		tklog.WriteErrorlog("update vivo token error :%s", err.Error())
		return
	}

	sender.Authtoken = token.AuthToken
	ret, err := sender.SendMSG(msg)

	if nil != err {
		tklog.WriteErrorlog("Push user list to vivo error :%s", err.Error())
		PushMsgMgr.PushDetail(request.DWAPPID, false, request.Guid, err.Error())
		return
	}
	result := VIVOPush.VIVONormalRet{}
	err = json.Unmarshal(ret, &result)
	if nil != err {
		tklog.WriteErrorlog("analysis Push vivo user list return value error :%s", err.Error())
		PushMsgMgr.PushDetail(request.DWAPPID, false, request.Guid, err.Error())
		return
	}
	if result.Result > 0 {
		tklog.WriteErrorlog("Vivo Push fail.err=%s.ID=%d.", result.Desc, result.Result)
		PushMsgMgr.PushDetail(request.DWAPPID, false, request.Guid, result.Desc)
	} else {
		PushMsgMgr.PushDetail(request.DWAPPID, true, request.Guid, "push success")
	}
}

func GetVivoStatistic(header *TKBaseDic.TKheader, data []byte, _conn *net.TCPConn) {
	defer Tool.PanicRecover()
	atomic.AddInt32(&RecptVivoCount, 1)

	sender := VIVOPush.Sender{}
	api, ok := VIVOPush.VIVOAPI[6]
	if !ok {
		tklog.WriteErrorlog("[vivo]wrong apiid:%d", 6)
		VIVOPush.StatisticWriteFailAck(header.Type, _conn)
		return
	}
	api += string(data)
	//api += "532976725354147840,532979536653459456"
	sender.Api = api

	secret, ok := M_mapCertMgr[header.Param]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", header.Param)
		return
	}

	token, err := VIVOPush.Gettoken(secret.Userinfo.VIVO_Appid, secret.Userinfo.VIVO_AppKey, secret.Userinfo.VIVO_AppSecRet, header.Param)
	if nil != err {
		tklog.WriteErrorlog("update vivo token error :%s", err.Error())
		VIVOPush.StatisticWriteFailAck(header.Type, _conn)
		return
	}

	sender.Authtoken = token.AuthToken
	ret, err := sender.GetMSG()
	if nil != err {
		tklog.WriteErrorlog("get vivo statistics msg error :%s", err.Error())
		VIVOPush.StatisticWriteFailAck(header.Type, _conn)
		return
	}
	VIVOPush.AckVivoStatistics(header.Type, ret, _conn)
}

func Vivotestpush() {
	h:=MsgMgr.GetReqTKHeader(1)
	h.Param=10000
	request := TKDic.PushVivoSFFMsg{}
	request.Payload =`{"title":"test22发发发1233","content":"tt22方法1234","isBusinessMsg":0,"notifyType":4,"timeToLive":44148,"skipType":1,"clientCustomMap":{"EnterByPush":"502000443","LGame":"1001,2"},"requestId":"18d3f1f5-5b3e-11e9-98f3-005056912135","extra":{"callback":"https://msgdx.srv.jj.cn/msgdx/api/MSG/VIVO/PushUserVIVORecptInfo","callback.param":"502000243"}}`

	buf := new(bytes.Buffer)
	he := MsgMgr.GetReqTKHeader(1)
	he.Length = uint32(len(request.Payload)+4)
	err := binary.Write(buf, binary.LittleEndian, &h.Param)
	err = binary.Write(buf, binary.LittleEndian, []byte(request.Payload))
	if err != nil {
		tklog.WriteErrorlog("write data Fail!err=%s", err.Error())
		return
	}
	SaveMsg2Vivo(&he,buf.Bytes(),nil)

	GetVivoStatistic(&h,[]byte("561151469521809408"),nil)

	request.Guid = VIVOPush.GetGuid()

	request.RegId = "15517820401031121941926"
	request.DWTokenCount = 1
	request.DWAPPID = 10004
	d, _ := json.Marshal(request)
	he = MsgMgr.GetReqTKHeader(1)
	he.Length = uint32(len(d))
	PushSFFMsg2Vivo(&he, d)
}
