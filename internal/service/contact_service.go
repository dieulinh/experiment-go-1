package service

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

var ErrContactMailNotConfigured = errors.New("contact mail is not configured")

type ContactService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	smtpFrom     string
	contactEmail string
}

func NewContactServiceFromEnv() *ContactService {
	return &ContactService{
		smtpHost:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		smtpPort:     strings.TrimSpace(os.Getenv("SMTP_PORT")),
		smtpUsername: strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		smtpPassword: strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		smtpFrom:     strings.TrimSpace(os.Getenv("SMTP_FROM")),
		contactEmail: strings.TrimSpace(os.Getenv("CONTACT_EMAIL")),
	}
}

func (s *ContactService) RecipientEmail() string {
	return s.contactEmail
}

func (s *ContactService) IsConfigured() bool {
	return s.smtpHost != "" &&
		s.smtpPort != "" &&
		s.smtpUsername != "" &&
		s.smtpPassword != "" &&
		s.contactEmail != ""
}

func (s *ContactService) Send(name, email, subject, message string) error {
	if !s.IsConfigured() {
		return ErrContactMailNotConfigured
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid sender email: %w", err)
	}

	from := s.smtpFrom
	if from == "" {
		from = s.smtpUsername
	}

	cleanName := sanitizeHeader(name)
	cleanEmail := sanitizeHeader(email)
	cleanSubject := sanitizeHeader(subject)
	body := strings.TrimSpace(message)

	mailBody := strings.Join([]string{
		fmt.Sprintf("New contact form message from %s", cleanName),
		"",
		fmt.Sprintf("Name: %s", cleanName),
		fmt.Sprintf("Email: %s", cleanEmail),
		fmt.Sprintf("Subject: %s", cleanSubject),
		"",
		"Message:",
		body,
	}, "\r\n")

	msg := []byte(strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", s.contactEmail),
		fmt.Sprintf("Reply-To: %s", cleanEmail),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		fmt.Sprintf("Subject: Contact form: %s", cleanSubject),
		"",
		mailBody,
	}, "\r\n"))

	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	return smtp.SendMail(net.JoinHostPort(s.smtpHost, s.smtpPort), auth, from, []string{s.contactEmail}, msg)
}

func sanitizeHeader(value string) string {
	replacer := strings.NewReplacer("\r", " ", "\n", " ")
	return replacer.Replace(strings.TrimSpace(value))
}