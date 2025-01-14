package alert

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail sends an email to the specified recipient with a subject and body.
func SendEmail(to, subject, body string) error {
	password := os.Getenv("EMAIL_PASSWORD");
	// Create a new email message
	mail := gomail.NewMessage()

	// Set the sender and recipient details
	mail.SetHeader("From", "su.victor21@gmail.com") // Sender's email
	mail.SetHeader("To", to) // Recipient's email
	mail.SetHeader("Subject", subject) // Email subject
	mail.SetBody("text/plain", body) // Email body content

	// Set up the SMTP client
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "su.victor21@gmail.com", password)

	// Send the email
	if err := dialer.DialAndSend(mail); err != nil {
		// If the email fails to send, log the error and return it
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	// Return nil if the email was sent successfully
	return nil
}
