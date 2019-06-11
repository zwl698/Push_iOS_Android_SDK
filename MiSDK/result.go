package MiPush

import (
	"TKGoBase/IO/Coon/MsgMgr"
	tklog "TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"encoding/binary"
	"encoding/json"
)

//解析请求的返回
func ResultResolve(data []byte, err error) (*TKDic.MiPushAck, []byte) {
	ret := TKDic.MiPushAck{}
	var extdate []byte
	if nil != err {
		ret.DwState = 1
		ret.Extra.DWSuffixOffSet = 36
		ret.Extra.DWSuffixSize = 0
		return &ret, nil
	}

	result := ResultMsg{}
	er := json.Unmarshal(data, &result)
	if nil != er {
		ret.DwState = 1
		ret.Extra.DWSuffixOffSet = 36
		ret.Extra.DWSuffixSize = 0
		return &ret, nil
	}

	ret.DwState = uint32(result.Code)
	if result.Code > 0 { //失败，返回服务器端的错误回传的错误消息
		ret.Extra.DWSuffixOffSet = 36
		ret.Extra.DWSuffixSize = uint32(len(result.Reason))
		extdate = []byte(result.Reason)

	} else { //成功，返回查询的消息体，或者推送的 消息id
		ret.Extra.DWSuffixOffSet = 36
		ret.Extra.DWSuffixSize = uint32(len(result.Data))
		extdate = []byte(result.Data)
	}

	return &ret, extdate

}

func ResultS(data []byte, err error) bool {
	if nil != err {
		return false
	}
	result := ResultMsg{}
	er := json.Unmarshal(data, &result)
	if nil != er {
		tklog.WriteErrorlog("unmarsha xiaomi result err:%s", er.Error())
		return false
	}
	if result.Code > 0 { //失败，返回服务器端的错误回传的错误消息
		tklog.WriteErrorlog("push xiaomi fail!code=%d.reason=%s,desc=%s.", result.Code, result.Reason, result.Description)
		return false
	}
	return true
}

func PushSendResult(Req *TKDic.MiPushReq, ret bool) {
	var Ack TKDic.TagReqMsgInterface2AndroidPushStateMsg
	Ack.DWBID, Ack.DWMsgType, Ack.DWFirm = Req.DwBid, Req.DwMsgType, TKBaseDic.AndroidCompany_XiaoMi
	if ret {
		Ack.DWSuccess, Ack.DWFail = Req.DWTokenCount, 0
	} else {
		Ack.DWSuccess, Ack.DWFail = 0, Req.DWTokenCount
	}
	Ack.Header = MsgMgr.GetReqTKHeader(TKDic.TKID_MSGINTERFACE2MSGANDROID_PUSHRESULT)
	Ack.Header.Length = 20
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &Ack)
	if err != nil {
		tklog.WriteErrorlog("write data Fail!err=%s", err.Error())
		return
	}
	h, _, s := PushMsgMgr.AndroidSMsgMgr.SendMsgWithAck(buf)
	if !s {
		if h != nil {
			tklog.WriteErrorlog("PushXiaoMiResult Fail!send=%d,param=%d", s, h.Param)
		} else {
			tklog.WriteErrorlog("PushXiaoMiResult Success!")
		}
	}

}

/*
func PushSendBcstResult(Req *TKDic.MiPushBcstReq,ret bool)  {
	var Ack TKDic.TagReqMsgInterface2AndroidPushStateMsg
	Ack.DWBID, Ack.DWMsgType, Ack.DWFirm = Req.DwBid, Req.DwMsgType, TKDic.AndroidCompany_XiaoMi
	if ret {
		Ack.DWSuccess, Ack.DWFail = Req.DWTokenCount, 0
	} else {
		Ack.DWSuccess, Ack.DWFail = 0, Req.DWTokenCount
	}
	Ack.Header = MsgMgr.GetReqTKHeader(TKDic.TKID_MSGINTERFACE2MSGANDROID_PUSHRESULT)
	Ack.Header.Length = 20
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &Req)
	if err != nil {
		tklog.WriteErrorlog("write data Fail!err=%s", err.Error())
		return
	}
	h, _, s := PushMsgMgr.AndroidSMsgMgr.SendMsgWithAck(buf)
	if !s {
		if h != nil {
			tklog.WriteErrorlog("PushXiaoMiBcstResult Fail!send=%d,param=%d", s, h.Param)
		} else {
			tklog.WriteErrorlog("PushXiaoMiBcstResult Success!")
		}
	}

}

*/
