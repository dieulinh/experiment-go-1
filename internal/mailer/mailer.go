package mailer

import (
	"crypto/tls"

	gomail "gopkg.in/gomail.v2"
)

type Config struct {
	SenderEmail string // your gmail: noreply@gmail.com
	SenderName  string // "My App"
	AppPassword string // 16-char app password from Google
}

type Mailer struct {
	config Config
	dialer *gomail.Dialer
}

func New(cfg Config) *Mailer {
	dialer := gomail.NewDialer("smtp.gmail.com", 465, cfg.SenderEmail, cfg.AppPassword)
	dialer.TLSConfig = &tls.Config{
		ServerName:         "smtp.gmail.com",
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	return &Mailer{
		config: cfg,
		dialer: dialer,
	}
}
