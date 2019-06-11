package VIVOPush

import (
	"TKGoBase/IO/Coon/MsgMgr"
	tklog "TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
)

func ResultSlove(data []byte, err error, api int) (*TKDic.VIVOPushAck, []byte) {
	switch api {
	case 0:
		{

			break
		}
	case 1:
		{

			break
		}
	case 2:
		{
			return NorResultResolve(data, err)
		}
	case 3:
		{
			return NorResultResolve(data, err)
		}
	case 4:
		{
			return NorResultResolve(data, err)
		}
	case 5:
		{
			return NorResultResolve(data, err)
		}
	case 6:
		{
			//todo 二期完成 解析统计数据
		}
	default:

	}

	return nil, nil
}

//解析请求的返回
func NorResultResolve(data []byte, err error) (*TKDic.VIVOPushAck, []byte) {
	ret := TKDic.VIVOPushAck{}
	ret.Msg.DWSuffixOffSet = 44
	var extdate []byte
	if nil != err {
		ret.DwState = 1
		ret.Msg.DWSuffixSize = 0
		return &ret, nil
	}

	result := VIVONormalRet{}
	er := json.Unmarshal(data, &result)
	if nil != er {
		ret.DwState = 1
		ret.Msg.DWSuffixSize = 0
		return &ret, nil
	}

	ret.DwState = uint32(result.Result)
	if result.Result > 0 { //失败，返回服务器端的错误回传的错误消息
		ret.Msg.DWSuffixSize = uint32(len(result.Desc))
		extdate = []byte(result.Desc)

	} else { //成功，返回查询的消息体，或者推送的 消息id
		length := len(result.TaskId)
		if length > 0 { //返回taskid
			ret.Msg.DWSuffixSize = uint32(len(result.TaskId))
			extdate = []byte(result.TaskId)
		} else { //推送名单接口，返回成功的描述
			ret.Msg.DWSuffixSize = uint32(len(result.Desc))
			extdate = []byte(result.Desc)
		}

	}
	return &ret, extdate

}

func PushSendResult(Req *TKDic.VIVOPushReq, ret bool) {
	var Ack TKDic.TagReqMsgInterface2AndroidPushStateMsg
	Ack.DWBID, Ack.DWMsgType, Ack.DWFirm = Req.DWBID, Req.DWMsgType, TKBaseDic.AndroidCompany_Vivo
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
			tklog.WriteErrorlog("PushVivoResult Fail!send=%d,param=%d", s, h.Param)
		} else {
			tklog.WriteErrorlog("PushVivoResult Success!")
		}
	}

}

func AckVivomsgid(msgID string, _conn *net.TCPConn) {
	Header := MsgMgr.GetAckTKHeader(TKDic.TKID_MSGANDROID2MSGINTERFACE_SAVEVIVOMSG)
	Header.Length = uint32(len(msgID))
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &Header)
	if err != nil {
		tklog.WriteErrorlog("write data Fail!err=%s", err.Error())
		return
	}

	err = binary.Write(buf, binary.LittleEndian, []byte(msgID))
	if err != nil {
		tklog.WriteErrorlog("write data Fail!err=%s", err.Error())
		return
	}

	_, err = _conn.Write(buf.Bytes())
	if err != nil {
		tklog.WriteErrorlog("send back vivo taskid err=%s,taskid = %d", err.Error(), msgID)
	}
}

func AckVivoStatistics(tp uint32, data []byte, _conn *net.TCPConn) {
	ret := PushStatisticRet{}
	ack := TKDic.ACKStatistic{}

	ack.Header.Origine = MsgMgr.GetOriange()
	ack.Header.Type = TKBaseDic.TK_ACK | tp
	var sli_statistic []TKDic.Statistic
	var taskid []string
	var leng uint32
	err := json.Unmarshal(data, &ret)
	if nil != err {
		tklog.WriteErrorlog("Unmarshal statistic data error :%s", err.Error())
		return
	}
	if ret.Result > 0 {
		ack.Count = 0
		ack.Header.Length = 4
	} else {
		size := len(ret.Statistics)
		ack.Count = uint32(size)
		for i := 0; i < size; i++ {
			sta := TKDic.Statistic{}
			sta.Receive = uint32(ret.Statistics[i].Receive)
			sta.Send = uint32(ret.Statistics[i].Send)
			sta.Click = uint32(ret.Statistics[i].Click)
			if ret.Statistics[i].Display < 0 {
				ret.Statistics[i].Display *= -1
			}
			sta.Display = uint32(ret.Statistics[i].Display)
			sli_statistic = append(sli_statistic, sta)
			taskid = append(taskid, ret.Statistics[i].TaskID)
		}

		if size > 0 {
			sli_statistic[0].Taskid.DWSuffixOffSet = 24
			sli_statistic[0].Taskid.DWSuffixSize = uint32(len(taskid[0]))
			leng += sli_statistic[0].Taskid.DWSuffixSize
		}

		for i := 1; i < size; i++ {
			sli_statistic[i].Taskid.DWSuffixOffSet = 24
			sli_statistic[i].Taskid.DWSuffixSize = uint32(len(taskid[i]))
			leng += sli_statistic[i].Taskid.DWSuffixSize
		}
		if size > 0 {
			ack.Header.Length = leng + uint32(4+size*24)
		} else {
			ack.Header.Length = 4
		}

	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &ack)
	size := len(ret.Statistics)
	for i := 0; i < size; i++ {
		binary.Write(buf, binary.LittleEndian, &(sli_statistic[i]))
		binary.Write(buf, binary.LittleEndian, []byte(taskid[i]))
	}
	_, err = _conn.Write(buf.Bytes())
	if nil != err {
		tklog.WriteErrorlog("write statistic to android error :", err.Error())
	}
}

func StatisticWriteFailAck(tp uint32, _conn *net.TCPConn) {
	ack := TKDic.ACKStatistic{}
	ack.Count = 0
	ack.Header.Origine = MsgMgr.GetOriange()
	ack.Header.Type = TKBaseDic.TK_ACK | tp
	ack.Header.Param = TKBaseDic.TK_ACKRESULT_FAILED
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &ack)
	_, err := _conn.Write(buf.Bytes())
	if nil != err {
		tklog.WriteErrorlog("write statistic to android error :", err.Error())
	}
}
