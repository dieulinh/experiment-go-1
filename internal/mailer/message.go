package mailer

import (
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

type EmailOptions struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

func (m *Mailer) Send(opts EmailOptions) error {
	msg := gomail.NewMessage()

	// This is what the receiver sees as sender
	msg.SetAddressHeader("From", m.config.SenderEmail, m.config.SenderName)
	msg.SetHeader("To", opts.To)
	msg.SetHeader("Subject", opts.Subject)

	contentType := "text/plain"
	if opts.IsHTML {
		contentType = "text/html"
	}
	msg.SetBody(contentType, opts.Body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// Contact notification email
func (m *Mailer) SendContactNotification(to, from, message string) error {
	body := fmt.Sprintf(`
        <h2>New Contact Form Submission</h2>
        <p><b>From:</b> %s</p>
        <p><b>Message:</b></p>
        <p>%s</p>
    `, from, message)

	return m.Send(EmailOptions{
		To:      to,
		Subject: "New Contact Form Submission",
		Body:    body,
		IsHTML:  true,
	})
}

// Order confirmation email
func (m *Mailer) SendOrderConfirmation(to, orderID, details string) error {
	body := fmt.Sprintf(`
        <h2>Order Confirmation</h2>
        <p>Thank you for your order!</p>
        <p><b>Order ID:</b> %s</p>
        <p><b>Details:</b></p>
        <p>%s</p>
    `, orderID, details)

	return m.Send(EmailOptions{
		To:      to,
		Subject: fmt.Sprintf("Order Confirmation #%s", orderID),
		Body:    body,
		IsHTML:  true,
	})
}
