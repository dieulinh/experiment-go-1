package handler

import (
	"net/http"
	"net/mail"
	"os"
	"strings"

	"restapi/internal/helper"
	"restapi/internal/mailer"
)

type ContactPageData struct {
	Error          string
	Success        string
	Name           string
	Email          string
	Subject        string
	Message        string
	RecipientEmail string
	MailEnabled    bool
}

func (h *Handler) ContactPage(w http.ResponseWriter, r *http.Request) {
	h.renderContactPage(w, ContactPageData{
		RecipientEmail: h.contactService.RecipientEmail(),
		MailEnabled:    h.contactService.IsConfigured(),
	})
}

func (h *Handler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	data := ContactPageData{
		Name:           strings.TrimSpace(r.FormValue("name")),
		Email:          strings.TrimSpace(r.FormValue("email")),
		Subject:        strings.TrimSpace(r.FormValue("subject")),
		Message:        strings.TrimSpace(r.FormValue("message")),
		RecipientEmail: h.contactService.RecipientEmail(),
		MailEnabled:    h.contactService.IsConfigured(),
	}

	if data.Name == "" || data.Email == "" || data.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = "Name, email, and message are required."
		h.renderContactPage(w, data)
		return
	}

	if _, err := mail.ParseAddress(data.Email); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		data.Error = "Please enter a valid email address."
		h.renderContactPage(w, data)
		return
	}

	if data.Subject == "" {
		data.Subject = "Website contact request"
	}

	appPassword := strings.TrimSpace(os.Getenv("GOOGLE_APP_PASSWORD"))
	if appPassword == "" {
		w.WriteHeader(http.StatusInternalServerError)
		data.Error = "Email service is not configured. Please set GOOGLE_APP_PASSWORD."
		h.renderContactPage(w, data)
		return
	}

	recipient := strings.TrimSpace(os.Getenv("CONTACT_EMAIL"))
	if recipient == "" {
		recipient = "linhvn09@gmail.com"
	}

	m := mailer.New(mailer.Config{
		SenderEmail: "linhunog@gmail.com",
		SenderName:  "My Shop",
		AppPassword: appPassword,
	})

	body := "Name: " + data.Name + "\n" +
		"Email: " + data.Email + "\n" +
		"Subject: " + data.Subject + "\n\n" +
		data.Message

	err := m.SendContactNotification(recipient, data.Email, body)
	if err != nil {
		h.log.Errorf("[Contact] failed to send email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.Error = "Unable to send your message right now. Please try again later."
		h.renderContactPage(w, data)
		return
	}

	h.renderContactPage(w, ContactPageData{
		Success:        "Your message has been sent. I will get back to you by email.",
		RecipientEmail: recipient,
		MailEnabled:    true,
	})
}

func (h *Handler) renderContactPage(w http.ResponseWriter, data ContactPageData) {
	helper.Render(w, data, "templates/contact.html")
}
