package email

import (
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
	"gopkg.in/gomail.v2"
)

type EmailRepository struct {
	auth *EmailAuth
}

func NewEmailRepository(auth *EmailAuth) port.EmailRepository {
	if auth == nil {
		panic("nil EmailAuth")
	}
	return &EmailRepository{
		auth: auth,
	}
}

func (r *EmailRepository) SendAll(users []entity.User, subject, content string) error {
	// 送信対象Email
	emails := []string{}
	for _, user := range users {
		emails = append(emails, user.Email().String())
	}
	m := gomail.NewMessage()
	// 送信者
	m.SetHeader("From", r.auth.EmailAddress())
	// 送信先(BCCのため自分自身を設定)
	m.SetHeader("To", r.auth.EmailAddress())
	// 件名
	m.SetHeader("Subject", subject)
	// メールの内容
	m.SetBody("text/html", content)

	// BCC設定
	m.SetHeader("Bcc", emails...)
	// メール送信
	dialer := gomail.NewDialer(
		r.auth.Host(),
		r.auth.Port(),
		r.auth.EmailAddress(),
		r.auth.Password(),
	)
	err := dialer.DialAndSend(m)
	return err
}
