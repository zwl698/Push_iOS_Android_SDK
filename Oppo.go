package PushHttps

import (
	"TKGoBase/IO/Coon/MsgMgr"
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushHttps/OppoSDK"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"k8s.io/apimachinery/pkg/util/uuid"
	"net"
	"sync/atomic"
)

var OppoCount int32

func PushMsg2Oppo(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	req := TKDic.PushOppoMsg{}
	err := binary.Read(bytes.NewBuffer(data[:36]), binary.LittleEndian, &req)
	if err != nil {
		tklog.WriteErrorlog("PushMsg2Oppo anysize fail!,err=%s", err.Error())
		return
	}
	if header.Length != (36 + req.TokenList.DWSuffixSize + req.Payload.DWSuffixSize) {
		tklog.WriteErrorlog("PushMsg2Oppo messageLength check fail!")
		return
	}
	atomic.AddInt32(&OppoCount, int32(req.DWTokenCount))
	secret, ok := M_mapCertMgr[req.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", req.DWAPPID)
		return
	}

	err = oppopush.Push(req.DWTargetType, string(data[req.TokenList.DWSuffixOffSet-24:req.TokenList.DWSuffixOffSet+req.TokenList.DWSuffixSize-24]),
		string(data[req.Payload.DWSuffixOffSet-24:req.Payload.DWSuffixOffSet+req.Payload.DWSuffixSize-24]), secret.Userinfo.OPPO_Password, secret.Userinfo.OPPO_MasterSecret, req.DWAPPID)

	if err != nil {
		tklog.WriteInfolog("Push OppoMsgFail!err=%s", err.Error())
		PushMsgMgr.PushOppoResult(&req, false)
	} else {
		PushMsgMgr.PushOppoResult(&req, true)
		//tklog.WriteInfolog("PushOppoResult success")
	}
}

func SaveMsg2Oppo(header *TKBaseDic.TKheader, data []byte, _conn *net.TCPConn) {
	defer Tool.PanicRecover()
	if header.Length < 4 {
		tklog.WriteErrorlog("SaveMsg2Oppo headerlength err.length=%d.", header.Length)
		return
	}
	var dwappid uint32
	err := binary.Read(bytes.NewBuffer(data[:4]), binary.LittleEndian, &dwappid)
	if err != nil {
		tklog.WriteErrorlog("SaveMsg2Oppo anysize fail!,err=%s", err.Error())
		return
	}
	secret, ok := M_mapCertMgr[dwappid]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", dwappid)
		return
	}
	err, msgid := oppopush.DealSaveMsg([]byte(data[4:]), secret.Userinfo.OPPO_Password, secret.Userinfo.OPPO_MasterSecret, dwappid)

	if err != nil || msgid == "" {
		tklog.WriteInfolog("Push SaveMsg2Oppo!err=%s", err.Error())
		msgid = "error"
		PushMsgMgr.Ackoppomsgid(msgid, _conn)
	} else {
		//tklog.WriteInfolog(msgid)
		PushMsgMgr.Ackoppomsgid(msgid, _conn)
	}
}

func PushSFFMsg2Oppo(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	req := TKDic.PushOppoSFFMsg{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		tklog.WriteErrorlog("PushSFFMsg2Oppo json prase err.err=%s.data=%s.", err.Error(), string(data))
		return
	}
	atomic.AddInt32(&OppoCount, int32(req.DWTokenCount))
	secret, ok := M_mapCertMgr[req.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", req.DWAPPID)
		return
	}

	err = oppopush.PushOneUser(req.TokenList,
		req.Payload, secret.Userinfo.OPPO_Password, secret.Userinfo.OPPO_MasterSecret, req.DWAPPID)

	if err != nil {
		tklog.WriteInfolog("Push OppoMsgFail!err=%s", err.Error())
		PushMsgMgr.PushDetail(req.DWAPPID, false, req.Guid, err.Error())
	} else {
		PushMsgMgr.PushDetail(req.DWAPPID, true, req.Guid, "push success")
	}
}

func OppoTESTmSG() {
	req := TKDic.PushOppoSFFMsg{}
	req.Guid = string(uuid.NewUUID())
	req.Payload = `{"app_message_id":"c73b2ef6-5b3f-11e9-acdc-0050568fb505","title":"zhaowltest","sub_title":"JJ","content":"eeee很好1234","click_action_type":0,"click_action_activity":"","click_action_url":"","action_parameters":"{\"EnterByPush\":\"502000208\",\"LGame\":\"1001,2\"}","show_time_type":0,"show_start_time":0,"show_end_time":0,"off_line":true,"off_line_ttl":42667,"push_time_type":0,"push_start_time":0,"time_zone":"","fix_speed":0,"fix_speed_rate":0,"network_type":0,"call_back_url":"https://msgdx.srv.jj.cn/msgdx/api/MSG/OPPO/PushUserOPPORecptInfo","call_back_parameter":"502000208"}`
	req.DWAPPID = 10000


	//buf := new(bytes.Buffer)
	//he := MsgMgr.GetReqTKHeader(1)
	//he.Length = uint32(len(req.Payload)+4)
	//err := binary.Write(buf, binary.LittleEndian, &req.DWAPPID)
	//err = binary.Write(buf, binary.LittleEndian, []byte(req.Payload))
	//if err != nil {
	//	tklog.WriteErrorlog("write data Fail!err=%s", err.Error())
	//	return
	//}
	//SaveMsg2Oppo(&he,buf.Bytes(),nil)

	req.TokenList = "CN_b82de998ecc16baed163d117938faca1"

	req.DWTokenCount = 1
	req.DWTargetType = 2
	d, _ := json.Marshal(req)
	he := MsgMgr.GetReqTKHeader(1)
	he.Length = uint32(len(d))
	PushSFFMsg2Oppo(&he, d)
}
