package emailservice

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
)

type smtpEmailService struct {
	host     string // smtp.gmail.com
	port     string // 587
	username string // your@gmail.com
	password string // App Password (ไม่ใช่ Password Gmail จริง)
	from     string // "App Name <your@gmail.com>"
}

func NewSMTPEmailService(host, port, username, password, from string) ports.IEmailService {
	return &smtpEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *smtpEmailService) SendPasswordResetEmail(toEmail, firstName, resetLink string) error {
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Password Reset Request</h2>
  <p>Hi %s,</p>
  <p>We received a request to reset your password. Click the button below to reset it.</p>
  <p><strong>This link expires in 15 minutes.</strong></p>
  <a href="%s" style="
    display: inline-block;
    padding: 12px 24px;
    background-color: #4F46E5;
    color: white;
    text-decoration: none;
    border-radius: 6px;
    margin: 16px 0;
  ">Reset Password</a>
  <p>If you didn't request this, please ignore this email. Your password won't change.</p>
  <p>For security, this link can only be used once.</p>
</body>
</html>`, template.HTMLEscapeString(firstName), resetLink)

	return s.sendEmail(toEmail, subject, body)
}

func (s *smtpEmailService) SendWelcomeEmail(toEmail, firstName string) error {
	subject := "Welcome! 🎉"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Welcome, %s!</h2>
  <p>Your account has been created successfully.</p>
  <p>Next step: Verify your phone number to unlock all features.</p>
</body>
</html>`, template.HTMLEscapeString(firstName))

	return s.sendEmail(toEmail, subject, body)
}

func (s *smtpEmailService) sendEmail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := buildMIMEMessage(s.from, to, subject, htmlBody)

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", s.host, s.port),
		auth,
		s.username,
		[]string{to},
		msg,
	)
}

func buildMIMEMessage(from, to, subject, htmlBody string) []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)
	return buf.Bytes()
}
