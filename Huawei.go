package PushHttps

import (
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushHttps/HuaWeiSDK"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sync/atomic"
	"time"
)

//var access_token *huaweipush.HWToken
var HuaweiCount int32

func PushMsg2HuaWei(header *TKBaseDic.TKheader, data []byte, trytimes uint32) {
	defer Tool.PanicRecover()
	req := TKDic.PushHuaWeiMsg{}
	err := binary.Read(bytes.NewBuffer(data[:36]), binary.LittleEndian, &req)
	if err != nil {
		tklog.WriteErrorlog("PushMsg2HuaWei anysize fail!,err=%s", err.Error())
		return
	}
	if header.Length != (36 + req.Payload.DWSuffixSize + req.TokenList.DWSuffixSize) {
		tklog.WriteErrorlog("PushMsg2HuaWei messageLength check fail!")
		return
	}
	atomic.AddInt32(&HuaweiCount, int32(req.DWTokenCount))

	secret, ok := M_mapCertMgr[req.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", req.DWAPPID)
		return
	}
	atomic.AddUint32(&secret.Queueinfo.FHuaWeiPush, req.DWTokenCount)

	client := huaweipush.NewClient(secret.Userinfo.HW_Appid, secret.Userinfo.HW_AppSecret)
	err, code := client.Push(context.TODO(), string(data[req.TokenList.DWSuffixOffSet-24:req.TokenList.DWSuffixOffSet+req.TokenList.DWSuffixSize-24]),
		time.Unix(int64(req.ETime), 0), string(data[req.Payload.DWSuffixOffSet-24:req.Payload.DWSuffixOffSet+req.Payload.DWSuffixSize-24]), &secret.Queueinfo.FHuaWeiPushCopy, req.DWAPPID, secret.Userinfo.HW_Appid)

	if err != nil {
		tklog.WriteInfolog("Push HuaWeiMsgFail!err=%s", err.Error())
		if code == 503 && trytimes < 10 {
			stData := new(TKDic.PushDataCacheQueueSt)
			stData.DevType = TKBaseDic.AndroidCompany_HuaWei
			stData.TryTimes = trytimes + 1
			stData.Header = *header
			stData.Data = data

			secret.Queueinfo.M_QueLock.Lock()
			secret.Queueinfo.PDataCacheQueue.PushBack(stData)
			secret.Queueinfo.M_QueLock.Unlock()
		}
		PushMsgMgr.PushHuaWeiResult(&req, false)
	} else {
		PushMsgMgr.PushHuaWeiResult(&req, true)
	}
}

func PushSFFMsg2HuaWei(header *TKBaseDic.TKheader, data []byte, trytimes uint32) {
	defer Tool.PanicRecover()
	req := TKDic.PushHuaWeiSFFMsg{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		tklog.WriteErrorlog("PushSFFMsg2HuaWei json prase err.err=%s.data=%s.", err.Error(), string(data))
		return
	}

	atomic.AddInt32(&HuaweiCount, int32(req.DWTokenCount))

	secret, ok := M_mapCertMgr[req.DWAPPID]
	if !ok {
		tklog.WriteErrorlog("getappsecret fail!appid=%d.", req.DWAPPID)
		return
	}
	atomic.AddUint32(&secret.Queueinfo.FHuaWeiPush, req.DWTokenCount)

	client := huaweipush.NewClient(secret.Userinfo.HW_Appid, secret.Userinfo.HW_AppSecret)
	err, code := client.Push(context.TODO(), req.TokenList,
		time.Unix(int64(req.Etime), 0), req.Payload, &secret.Queueinfo.FHuaWeiPushCopy, req.DWAPPID, secret.Userinfo.HW_Appid)

	if err != nil {
		tklog.WriteInfolog("PushSFFMsg2HuaWei MsgFail!err=%s", err.Error())
		if code == 503 && trytimes < 10 {
			stData := new(TKDic.PushDataCacheQueueSt)
			stData.DevType = TKBaseDic.AndroidCompany_HuaWei
			stData.TryTimes = trytimes + 1
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

func TESTHUAWEI() {
	req := TKDic.PushHuaWeiSFFMsg{}
	req.DWTokenCount = 1
	req.DWAPPID = 10000
	req.Guid = string(uuid.NewUUID())
	req.Etime = 1553993285
	//	req.Payload=`{"hps":{"msg":{"type":3,"body":{"title":"test","content":"test"},"action":{"type":3,"param":{"appPkgName":"cn.jj"}}},"ext":{"biTag":"501000090","customize":[{"EnterByPush":"503000002:dnjahdudshfshfs"}]}}}`
	req.Payload = `{"hps":{"msg":{"type":3,"body":{"title":"tt","content":"test20190319"},"action":{"type":3,"param":{"appPkgName":"cn.jj"}}},"ext":{"biTag":"502000384","customize":[{"Vdetails":"jjtks://acts.jj.cn/html/activity_topic/sign_up/1903251d3d7d.html"},{"EnterByPush":"503000002:dnjahdudshfshfs"}]}}}`
	req.TokenList = `["0864100031983011300002700500CN01"]`
//req.TokenList=`["AMsVjyP0mmhibzyjWEn1OIHpPGuK7gR1taJ5iGWyQllYBCNR0Ebed8w2FFCTg116v8eYbkHHOdx86UmCH2G0qRtBxIAl_GPF8JzmeVwy0No9e56MuuVEcZCKdpkztQ7b9g"]`
	d, _ := json.Marshal(&req)
	PushSFFMsg2HuaWei(nil, d, 1)
}
