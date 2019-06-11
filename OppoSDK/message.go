package oppopush

import (
	"TKGoBase/IO/Log"
	"TKGoBase/Tool"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

func gettestmessage() (url.Values, error) {
	params := url.Values{}
	params.Add("app_message_id", "2")
	params.Add("title", "test12")
	params.Add("sub_title", "zhaowl1")
	params.Add("content", "tttttttt21")
	return params, nil
}

func DealSaveMsg(msg []byte, appkey, mastersecret string, appid uint32) (error, string) {
	//可在此处解json拼url
	var oppomsg Oppomsg
	token, err := RequestAccess(appkey, mastersecret, appid, false)
	if err != nil {
		return err, ""
	}

	msg2 := Tool.ConvertToString(string(msg), "gbk", "UTF-8")
	tklog.WriteInfolog(msg2)

	err = json.Unmarshal([]byte(msg2), &oppomsg)
	if err != nil {
		return err, ""
	}

	param, err := buildupmsg(&oppomsg)
	if err != nil {
		return err, ""
	}
	param.Add("auth_token", token.authtoken)
	return saveMessageID(param, appkey, mastersecret, appid)
}

func saveMessageID(parma url.Values, appkey, mastersecret string, appid uint32) (error, string) {

	bytes, err := doPost(context.TODO(), save_message, parma)
	if err != nil {
		return err, ""
	}
	var result pushBcstresult
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return err, ""
	}
	if result.Code != 0 && result.Code != 32 {
		tklog.WriteInfolog("code=%d,msg=%s.", result.Code, result.Message)
		if result.Code == Oppo_InvaildToken {
			RequestAccess(appkey, mastersecret, appid, true)
		}
		return errors.New("saveMessageID Fail!"), ""
	}
	return nil, result.Data.Messageid
}

func buildupmsg(pmsg *Oppomsg) (url.Values, error) {
	param := url.Values{}
	if pmsg.Title == "" || pmsg.Sub_title == "" || pmsg.Content == "" {
		return param, errors.New("buildupmsg error!")
	}

	param.Add("title", pmsg.Title)
	param.Add("sub_title", pmsg.Sub_title)
	param.Add("content", pmsg.Content)

	if pmsg.App_message_id != "" {
		param.Add("app_message_id", pmsg.App_message_id)
	}

	if pmsg.Click_action_type != 0 {
		param.Add("content", strconv.Itoa(int(pmsg.Click_action_type)))
	}
	if pmsg.Click_action_activity != "" {
		param.Add("click_action_activity", pmsg.Click_action_activity)
	}
	if pmsg.Click_action_url != "" {
		param.Add("click_action_url", pmsg.Click_action_url)
	}
	if pmsg.Action_parameters != "" {
		param.Add("action_parameters", pmsg.Action_parameters)
	}
	if pmsg.Show_time_type != 0 {
		param.Add("show_time_type", strconv.Itoa(int(pmsg.Show_time_type)))
	}
	if pmsg.Show_start_time != 0 {
		param.Add("show_start_time", strconv.FormatInt(pmsg.Show_start_time, 10))
	}
	if pmsg.Show_end_time != 0 {
		param.Add("show_end_time", strconv.FormatInt(pmsg.Show_end_time, 10))
	}
	if pmsg.Off_line != true {
		param.Add("off_line", "0")
	}
	if pmsg.Off_line_ttl != 3600 && pmsg.Off_line_ttl != 0 {
		param.Add("off_line_ttl", strconv.Itoa(int(pmsg.Off_line_ttl)))
	}
	if pmsg.Push_time_type != 0 {
		param.Add("push_time_type", strconv.Itoa(int(pmsg.Push_time_type)))
	}
	if pmsg.Push_start_time != 0 {
		param.Add("push_start_time", strconv.FormatInt(pmsg.Push_start_time, 10))
	}
	//if pmsg.time_zone != "GMT+08:00" {
	//	param.Add("time_zone", pmsg.time_zone)
	//}
	if pmsg.Fix_speed != 0 {
		param.Add("fix_speed", strconv.Itoa(int(pmsg.Fix_speed)))
	}
	if pmsg.Fix_speed_rate != 0 {
		param.Add("fix_speed_rate", strconv.FormatInt(pmsg.Fix_speed_rate, 10))
	}
	if pmsg.Network_type != 0 {
		param.Add("network_type", strconv.Itoa(int(pmsg.Network_type)))
	}
	if pmsg.Call_back_url != "" {
		param.Add("call_back_url", pmsg.Call_back_url)
	}
	if pmsg.Call_back_parameter != "" {
		param.Add("call_back_parameter", pmsg.Call_back_parameter)
	}
	return param, nil
}
