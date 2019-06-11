package oppopush

import (
	"TKGoBase/IO/Log"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"net/url"
)

func Push(targettype uint32, tokenlist, messageid, appkey, mastersecret string, appid uint32) error {
	token, err := RequestAccess(appkey, mastersecret, appid, false)

	if err != nil {
		tklog.WriteErrorlog("getOppoAccecc_token Fail!err=%s", err.Error())
		return err
	}

	params := url.Values{}
	params.Add("message_id", messageid)
	params.Add("target_type", "2")
	params.Add("target_value", tokenlist)
	params.Add("auth_token", token.authtoken)

	bytes, err := doPost(context.TODO(), push_Broadcast, params)
	if err != nil {
		return err
	}
	var result pushBcstresult
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return err
	}

	if result.Code != 0 || result.Data.Err0 != "" || result.Data.Err1 != "" || result.Data.Err2 != "" ||
		result.Data.Err3 != "" || result.Data.Err4 != "" || result.Data.Err5 != "" {
		if result.Code == Oppo_InvaildToken {
			RequestAccess(appkey, mastersecret, appid, true)
		}
		return errors.New(result.Message)
	}
	return nil
}

func PushOneUser(tokenl, msg, appkey, mastersecret string, appid uint32) error {
	token, err := RequestAccess(appkey, mastersecret, appid, false)

	if err != nil {
		tklog.WriteErrorlog("getOppoAccess_token Fail!err=%s", err.Error())
		return err
	}

	stemp := fmt.Sprintf("{\"target_type\":2,\"target_value\":\"%s\",\"notification\":%s}", tokenl, msg)
	params := url.Values{}
	params.Add("message", stemp)
	params.Add("auth_token", token.authtoken)

	bytes, err := doPost(context.TODO(), push_unicast, params)
	if err != nil {
		return err
	}
	var result pushBcstresult
	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return err
	}

	if result.Code != 0 || result.Data.Err0 != "" || result.Data.Err1 != "" || result.Data.Err2 != "" ||
		result.Data.Err3 != "" || result.Data.Err4 != "" || result.Data.Err5 != "" {
		if result.Code == Oppo_InvaildToken {
			RequestAccess(appkey, mastersecret, appid, true)
		}
		return errors.New(result.Message)
	}
	return nil

}
