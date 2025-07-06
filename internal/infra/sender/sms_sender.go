package sender

type SmsClient struct{}

func NewSmsClient() *SmsClient {
	return &SmsClient{}
}

func (s *SmsClient) Send(phone string, mess string) error {
	return nil
}
