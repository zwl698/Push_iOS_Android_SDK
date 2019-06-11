package PushHttps

import (
	"TKGoBase/IO/Coon/MsgMgr"
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushHttps/iOSSDK"
	"TKMSGDevInterfaceService/PushHttps/iOSSDK/certificate"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"k8s.io/apimachinery/pkg/util/uuid"
	"strings"
	"sync/atomic"
	"time"
)

var M_pIosMgr *apns2.ClientManager
var M_mapiosinfo map[uint32]iosinfo
var AppleCount int32

type iosinfo struct {
	topic   string
	cert    tls.Certificate
	certpwd string
}

func init() {
	M_pIosMgr = apns2.NewClientManager()
	M_mapiosinfo = make(map[uint32]iosinfo)
}

func IosAddObject() {
	iosmgr := apns2.NewClientManager()
	iosmap := make(map[uint32]iosinfo)
	for _, app := range M_mapCertMgr {
		var err error
		var cert tls.Certificate
		if strings.HasSuffix(app.Userinfo.IOS_CertPath, "p12") {
			cert, err = certificate.FromP12File(app.Userinfo.IOS_CertPath, app.Userinfo.IOS_CertPwd)
			if err != nil {
				tklog.WriteErrorlog("Cert Error:%s.", err.Error())
				continue
			}
		} else {
			cert, err = certificate.FromPemFile(app.Userinfo.IOS_CertPath, app.Userinfo.IOS_CertPwd)
			if err != nil {
				tklog.WriteErrorlog("Cert Error:%s.", err.Error())
				continue
			}
		}

		client := iosmgr.Factory(cert).Production()
		iosmgr.Add(client)
		info := iosinfo{app.Userinfo.IOS_Topic, cert, app.Userinfo.IOS_CertPwd}
		iosmap[app.Userinfo.JJAPPID] = info
	}
	M_pIosMgr = iosmgr
	M_mapiosinfo = iosmap
}

func PushMsg2Apple(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	req := TKDic.ApplePushReq{}
	err := binary.Read(bytes.NewBuffer(data[:32]), binary.LittleEndian, &req)
	if err != nil {
		tklog.WriteErrorlog("PushMsg2Apple anysize fail!,err=%s", err.Error())
		return
	}
	if header.Length != (32 + req.Payload.DWSuffixSize + req.Token.DWSuffixSize) {
		tklog.WriteErrorlog("PushMsg2Apple messageLength check fail!")
		return
	}
	atomic.AddInt32(&AppleCount, 1)

	v, ok := M_mapiosinfo[req.DWAPPID]
	if !ok {
		tklog.WriteFatallog("M_mapiosinfo find appid fail.appid=%d.", req.DWAPPID)
		return
	}
	msg := string(data[req.Payload.DWSuffixOffSet-24 : req.Payload.DWSuffixOffSet+req.Payload.DWSuffixSize-24])
	notification := &apns2.Notification{}
	notification.Topic = v.topic
	notification.Payload = Tool.ConvertToString(msg, "gbk", "utf-8")
	notification.DeviceToken = string(data[req.Token.DWSuffixOffSet-24 : req.Token.DWSuffixOffSet+req.Token.DWSuffixSize-24])
	notification.Expiration = time.Unix(int64(req.Etime), 0)
	notification.ApnsID = string(uuid.NewUUID())
	notification.Priority = 10

	client := M_pIosMgr.Get(v.cert)
	if client == nil {
		tklog.WriteFatallog("GetClient fail.APPID=%d.", req.DWAPPID)
		return
	}

	res, err := client.Push(notification)

	if err != nil {
		if res != nil {
			if res.Sent() {
				tklog.WriteErrorlog("PushApple err but success.err=%s.", err.Error())
				PushMsgMgr.PushAppleResult(&req, true)
			} else {
				tklog.WriteInfolog("Push PushMsg2Apple!err=%s.code=%d.responerr:%s.", err.Error(), res.StatusCode, res.Reason)
				PushMsgMgr.PushAppleResult(&req, false)
			}
		} else {
			tklog.WriteInfolog("Push PushMsg2Apple!err=%s.", err.Error())
			PushMsgMgr.PushAppleResult(&req, false)
		}
	} else {
		if res != nil {
			if res.Sent() {
				PushMsgMgr.PushAppleResult(&req, true)
			} else {
				tklog.WriteErrorlog("Push Apple Ack Fail.code=%d,reason=%s.", res.StatusCode, res.Reason)
				PushMsgMgr.PushAppleResult(&req, false)
			}
		} else {
			tklog.WriteFatallog("PushApple err=nil,ack=nil.May be fail.")
			PushMsgMgr.PushAppleResult(&req, false)
		}
	}
}

func PushSFFMsg2Apple(header *TKBaseDic.TKheader, data []byte) {
	defer Tool.PanicRecover()
	//测试如果不通看看是否把消息转为utf-8
	req := TKDic.PushAppleSFFMsg{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		tklog.WriteErrorlog("PushSFFMsg2Apple json prase err.err=%s.data=%s.", err.Error(), string(data))
		return
	}
	atomic.AddInt32(&AppleCount, 1)

	v, ok := M_mapiosinfo[req.DWAPPID]
	if !ok {
		tklog.WriteFatallog("M_mapiosinfo find appid fail.appid=%d.", req.DWAPPID)
		return
	}
	notification := &apns2.Notification{}
	notification.Topic = v.topic
	notification.Payload = req.Payload
	notification.DeviceToken = req.Token
	notification.Expiration = time.Unix(int64(req.Etime), 0)
	notification.ApnsID = string(uuid.NewUUID())
	notification.Priority = 10

	client := M_pIosMgr.Get(v.cert)
	if client == nil {
		tklog.WriteFatallog("GetClient fail.APPID=%d.", req.DWAPPID)
		return
	}

	res, err := client.Push(notification)

	if err != nil {
		if res != nil {
			if res.Sent() {
				tklog.WriteErrorlog("PushApple err but success.err=%s.", err.Error())
				PushMsgMgr.PushDetail(req.DWAPPID, true, req.Guid, "push success")
			} else {
				tklog.WriteInfolog("Push PushMsg2Apple!err=%s.code=%d.responerr:%s.", err.Error(), res.StatusCode, res.Reason)
				PushMsgMgr.PushDetail(req.DWAPPID, false, req.Guid, res.Reason)
			}
		} else {
			tklog.WriteInfolog("Push PushMsg2Apple!err=%s.", err.Error())
			PushMsgMgr.PushDetail(req.DWAPPID, false, req.Guid, err.Error())
		}
	} else {
		if res != nil {
			if res.Sent() {
				PushMsgMgr.PushDetail(req.DWAPPID, true, req.Guid, "push success")
			} else {
				tklog.WriteErrorlog("Push Apple Ack Fail.code=%d,reason=%s.", res.StatusCode, res.Reason)
				PushMsgMgr.PushDetail(req.DWAPPID, false, req.Guid, res.Reason)
			}
		} else {
			tklog.WriteFatallog("PushApple err=nil,ack=nil.May be fail.")
			PushMsgMgr.PushDetail(req.DWAPPID, false, req.Guid, "ack nil")
		}
	}
}

func Test1() {
	req := TKDic.PushAppleSFFMsg{}
	req.DWAPPID = 10000
	req.Payload = `{"aps":{"alert":{"title":"zhaowl","body":"ttttttttt","launch-image":" Default.png"},"badge":1,"sound":"default","content-available":1,"category":"set","thread-id":"get"}}`
	//	req.Token = "6452af4f77b7cc3a318793b2c4e74ca167230b71fa2a114a27d98b9c21d26b7c"
	req.Token = "d0b9912f9f7af71e10a7ca69665c28e2a50386e9eed07d860e68674cbe2f378d"
	req.Guid = string(uuid.NewUUID())
	req.Etime = 1563846974
	//	gg := Tool.ConvertToString(req.Payload, "gbk", "utf-8")
	//	tklog.WriteInfolog(gg)
	//notification := &apns2.Notification{}
	////	notification.DeviceToken = "f063e47b8fa3ed8c7863f95a93572124c7aa80f09cf2dd9f2f409acb09d16cc2" //福军
	////	notification.DeviceToken = "de6c86753ba8edf422d11674545eedd33e94b8811b0a5d5bd8ba4ff6fb0b3ad0" //涂俊
	//notification.DeviceToken = "d0b9912f9f7af71e10a7ca69665c28e2a50386e9eed07d860e68674cbe2f378d" //二虎
	//notification.Topic = "cn.jj.TKLobby"
	//notification.Payload = []byte(`{"aps":{"alert":{"title":"二虎5","body":"二虎！你好~hello."},"badge":100,"sound":""}}`) // See Payload section below
	//
	//client := apns2.NewClient(cert).Production()
	//res, err := client.Push(notification)
	//
	//if err != nil {
	//	tklog.WriteErrorlog("Error:", err)
	//}

	//	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	h := MsgMgr.GetReqTKHeader(1)
	m, _ := json.Marshal(req)
	PushSFFMsg2Apple(&h, m)
}
