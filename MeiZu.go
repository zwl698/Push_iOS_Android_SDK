package PushHttps

import (
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushHttps/MeiZuSDK"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync/atomic"
)

var MeiZuCount int32

func PushMsg2MeiZu(header *TKBaseDic.TKheader, data []byte, trytimes uint32) {
	defer Tool.PanicRecover()
	req := TKDic.PushMeiZuMsg{}
	err := binary.Read(bytes.NewBuffer(data[:32]), binary.LittleEndian, &req)
	if err != nil {
		tklog.WriteErrorlog("PushMsg2MeiZu anysize fail!,err=%s", err.Error())
		return
	}
	if header.Length != (32 + req.Payload.DWSuffixSize + req.TokenList.DWSuffixSize) {
		tklog.WriteErrorlog("PushMsg2MeiZu messageLength check fail!")
		return
	}
	atomic.AddInt32(&MeiZuCount, int32(req.DWTokenCount))

	secret, ok := M_mapCertMgr[req.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", req.DWAPPID)
		return
	}
	atomic.AddUint32(&secret.Queueinfo.FMeiZuPush, req.DWTokenCount)

	appid, sercet := secret.Userinfo.MZ_Appid, secret.Userinfo.MZ_AppSecret
	code, err := meizupush.PushMeiZuMsg(appid, string(data[req.TokenList.DWSuffixOffSet-24:req.TokenList.DWSuffixOffSet+req.TokenList.DWSuffixSize-24]),
		string(data[req.Payload.DWSuffixOffSet-24:req.Payload.DWSuffixOffSet+req.Payload.DWSuffixSize-24]), sercet, &secret.Queueinfo.FMeiZuPushCopy)

	if err != nil {
		tklog.WriteInfolog("Push MeiZuMsgFail!code=%d,err=%s", code, err.Error())
		if (code == 110010 || code == 1003) && trytimes < 10 {
			stData := new(TKDic.PushDataCacheQueueSt)
			stData.DevType = TKBaseDic.AndroidCompany_MeiZu
			stData.DevType = trytimes + 1
			stData.Header = *header
			stData.Data = data

			secret.Queueinfo.M_QueLock.Lock()
			secret.Queueinfo.PDataCacheQueue.PushBack(stData)
			secret.Queueinfo.M_QueLock.Unlock()
		}
		PushMsgMgr.PushMeiZuResult(&req, false)
	} else {
		PushMsgMgr.PushMeiZuResult(&req, true)
	}
}

func PushSFFMsg2MeiZu(header *TKBaseDic.TKheader, data []byte, trytimes uint32) {
	defer Tool.PanicRecover()
	req := TKDic.PushMeiZusffMsg{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		tklog.WriteErrorlog("PushSFFMsg2MeiZu json prase err.err=%s.data=%s.", err.Error(), string(data))
		return
	}
	atomic.AddInt32(&MeiZuCount, int32(req.DWTokenCount))

	secret, ok := M_mapCertMgr[req.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", req.DWAPPID)
		return
	}
	atomic.AddUint32(&secret.Queueinfo.FMeiZuPush, req.DWTokenCount)

	appid, sercet := secret.Userinfo.MZ_Appid, secret.Userinfo.MZ_AppSecret
	code, err := meizupush.PushMeiZuMsg(appid, req.TokenList,
		req.Payload, sercet, &secret.Queueinfo.FMeiZuPushCopy)

	if err != nil {
		tklog.WriteInfolog("Push MeiZuSFFMsgFail!code=%d,err=%s", code, err.Error())
		if (code == 110010 || code == 1003) && trytimes < 10 {
			stData := new(TKDic.PushDataCacheQueueSt)
			stData.DevType = TKBaseDic.AndroidCompany_MeiZu
			stData.DevType = trytimes + 1
			stData.Header = *header
			stData.Data = data

			secret.Queueinfo.M_QueLock.Lock()
			secret.Queueinfo.PDataCacheQueue.PushBack(stData)
			secret.Queueinfo.M_QueLock.Unlock()
		}
		PushMsgMgr.PushDetail(req.DWAPPID, false, req.Guid, err.Error())
	} else {
		PushMsgMgr.PushDetail(req.DWAPPID, true, req.Guid, "push success")
	}
}
