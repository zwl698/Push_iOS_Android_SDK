package PushHttps

import (
	"TKGoBase/IO/Log"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"encoding/json"
)

var M_mapCertMgr map[uint32]UserAllInfo
var m_CertTabVersion uint32

func init() {
	M_mapCertMgr = make(map[uint32]UserAllInfo)
}

func CertMapOnTimer() {
	version := getCertVersion()
	if version != m_CertTabVersion {
		if true == mapTab() {
			m_CertTabVersion = version
		}
	}
}

func getCertVersion() uint32 {
	return PushMsgMgr.GetDevInterfaceVersion()
}

func mapTab() bool {
	m := getCertList()
	if len(m) == 0 {
		return false
	}
	t := make(map[uint32]UserAllInfo)
	for k, v := range m {
		kt := InitDCList()
		vt := v //写这步的原因试试就知道了
		info := UserAllInfo{&vt, kt}
		t[k] = info
	}
	M_mapCertMgr = t
	IosAddObject()
	return true
}

func getCertList() map[uint32]TKDic.PWDInfo {
	data := PushMsgMgr.GetCertInfo()
	if data == nil {
		tklog.WriteFatallog("getCertList fail.")
		return nil
	}
	var vitem []TKDic.PWDInfo
	err := json.Unmarshal(data, &vitem)
	if err != nil {
		tklog.WriteErrorlog("Unmarshal err!err=%s.", err.Error())
		return nil
	}
	m := make(map[uint32]TKDic.PWDInfo)
	for i := 0; i < len(vitem); i++ {
		m[vitem[i].JJAPPID] = vitem[i]
	}
	return m
}
