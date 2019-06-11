package meizupush

import "testing"

func TestGenerateSign(t *testing.T) {
	changeRegisterMap := map[string]string{
		"appId":     APP_ID,
		"pushId":    PUSH_ID,
		"msgType":   "0",
		"subSwitch": "0",
	}
	sign := GenerateSign(changeRegisterMap, APP_KEY)
	if sign == "5bd43df385bec6f236c6417d437741e7" {
		t.Error("sign error")
	}

}

func TestPushParamMD5(t *testing.T) {
	md5Str := PushParamMD5("appId=100999msgType=0pushId=862891030007404100999subSwitch=1531732bc45324098978bf41c6954c09e")
	if md5Str != "5bd43df385bec6f236c6417d437741e7" {
		t.Error("md5 error")
	}

	md5Str2 := PushParamMD5("appId=100999pushId=862891030007404100999msgType=0subSwitch=1531732bc45324098978bf41c6954c09e")
	if md5Str2 == "5bd43df385bec6f236c6417d437741e7" {
		t.Error("md5 error")
	}
}
