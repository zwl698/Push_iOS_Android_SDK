package oppopush

const (
	push_base      = "https://api.push.oppomobile.com/"
	access_token   = "https://api.push.oppomobile.com/server/v1/auth"
	push_Broadcast = "https://api.push.oppomobile.com/server/v1/message/notification/broadcast"
	push_unicast   = "https://api.push.oppomobile.com/server/v1/message/notification/unicast"
	push_Batch     = "https://api.push.oppomobile.com/server/v1/message/notification/unicast_batch"
	save_message   = "https://api.push.oppomobile.com/server/v1/message/notification/save_message_content"
)

var Oppo_InvaildToken int32 = 11

type getOPPOToken struct {
	app_key      string
	sign         string
	mastersecret string
	timestamp    int64
}

type tokenrespon struct {
	Code    int32     `json:"code"`
	Message string    `json:"message"`
	Data    tokendata `json:"data"`
}

type tokendata struct {
	Auth_token  string `json:"auth_token"`
	Create_time int64  `json:"create_time"`
}

type token struct {
	authtoken string
	endtime   int64
}

type pushBcstresult struct {
	Code    int32        `json:"code"`
	Message string       `json:"message"`
	Data    pushBcstdate `json:"data"`
}

type pushBcstdate struct {
	Messageid string `json:"message_id"`
	Taskid    string `json:"task_id"`
	Err0      string `json:"10000"`
	Err1      string `json:"10001"`
	Err2      string `json:"10002"`
	Err3      string `json:"10003"`
	Err4      string `json:"10004"`
	Err5      string `json:"10005"`
}

//type pushMcstresult struct {
//	Code    int32     `json:"code"`
//	Message string    `json:"message"`
//	Data    pushMcstdate `json:"data"`
//}
//
//type pushMcstdate struct {
//	Messageid int32   `json:"message_id"`
//	Taskid     int32  `json:"task_id"`
//	Err1     string    `json:"10000"`
//	Err2    string   `json:"10002"`
//	Err3 string  `json:"10004"`
//	Err4 string   `json:"10005"`
//}
type Oppomsg struct {
	App_message_id        string `json:"app_message_id"`
	Title                 string `json:"title"`
	Sub_title             string `json:"sub_title"`
	Content               string `json:"content"`
	Click_action_type     int32  `json:"click_action_type"`
	Click_action_activity string `json:"click_action_activity"`
	Click_action_url      string `json:"click_action_url"`
	Action_parameters     string `json:"action_parameters"`
	Show_time_type        int32  `json:"show_time_type"`
	Show_start_time       int64  `json:"show_start_time"`
	Show_end_time         int64  `json:"show_end_time"`
	Off_line              bool   `json:"off_line"`
	Off_line_ttl          int32  `json:"off_line_ttl"`
	Push_time_type        int32  `json:"push_time_type"`
	Push_start_time       int64  `json:"push_start_time"`
	Time_zone             string `json:"time_zone"`
	Fix_speed             int32  `json:"fix_speed"`
	Fix_speed_rate        int64  `json:"fix_speed_rate"`
	Network_type          int32  `json:"network_type"`
	Call_back_url         string `json:"call_back_url"`
	Call_back_parameter   string `json:"call_back_parameter"`
}
