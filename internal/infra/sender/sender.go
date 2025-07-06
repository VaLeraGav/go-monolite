package sender

type Sender struct {
	emailClient *EmailClient
	smsClient   *SmsClient
}

func New() *Sender {
	return &Sender{
		emailClient: NewEmailClient(),
		smsClient:   NewSmsClient(),
	}
}

func (s *Sender) SendEmailCode(email, code string) error {
	return s.emailClient.Send(email, "Ваш код подтверждения: "+code)
}

func (s *Sender) SendSmsCode(phone, code string) error {
	return s.smsClient.Send(phone, "Код: "+code)
}
