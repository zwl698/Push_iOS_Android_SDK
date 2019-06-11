package MiPush

type Sender struct {
	Appsecret               string
	Restricted_package_name string
}

func (s *Sender) SendMSG(miapi string, _msg string) ([]byte, error) {
	msg := "restricted_package_name=" + s.Restricted_package_name + "&" + _msg
	return Post(miapi, msg, s.Appsecret)

}

func GetParam(index int, data []byte) string {
	switch index {
	case 1:
		{
			return "registration_id=" + string(data)
		}
	case 5:
		{
			return ""
		}
	}
	return ""
}
