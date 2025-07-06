package sender

type EmailClient struct {
}

func NewEmailClient() *EmailClient {
	return &EmailClient{}
}

func (s *EmailClient) Send(email, mess string) error {
	return nil
}
