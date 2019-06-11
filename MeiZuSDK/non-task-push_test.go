package meizupush

import (
	"fmt"
	"testing"
)

func TestPushThroughMessageByPushId(t *testing.T) {
	messageJson := `{"test_throught_message": "message"}`
	message := PushThroughMessageByPushId(APP_ID, PUSH_ID, buildThroughMessage(messageJson), APP_KEY)

	fmt.Println("TestPushThroughMessageByPushId ", message.message)

	if message.code != 200 {
		t.Error("Status Code not 200")
	}
}

func TestPushNotificationMessageByPushId(t *testing.T) {
	json := BuildNotificationMessage().
		noticeBarType(2).
		noticeTitle("标题go").
		noticeContent("测试内容").toJson()
	message := PushNotificationMessageByPushId(APP_ID, PUSH_ID, json, APP_KEY)

	fmt.Println("TestPushNotificationMessageByPushId ", message.message)

	if message.code != 200 {
		t.Error("Status Code not 200")
	}

}

func TestPushNotificationMessageByAlias(t *testing.T) {
	json := BuildNotificationMessage().
		noticeBarType(2).
		noticeTitle("标题").
		noticeContent("测试内容").toJson()
	message := PushNotificationMessageByAlias(APP_ID, "Android", json, APP_KEY)

	fmt.Println("TestPushNotificationMessageByAlias ", message.message)

	if message.code != 200 {
		t.Error("Status Code not 200")
	}
}
