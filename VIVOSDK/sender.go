package VIVOPush

type Sender struct {
	Api       string
	Authtoken string
}

func (s *Sender) SendMSG(_msg string) ([]byte, error) {
	return Post(s.Api, _msg, s.Authtoken)
}

func (s *Sender) GetMSG() ([]byte, error) {
	return Get(s.Api, s.Authtoken)
}
