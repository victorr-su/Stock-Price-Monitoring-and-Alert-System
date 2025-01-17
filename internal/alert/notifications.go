package alert

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

// EmailData contains the dynamic data to populate the email template
type EmailData struct {
	StockSymbol     string
	PercentageChange float64
	Price           float64
	Time            string
}

// SendEmail sends an email to the specified recipient with a subject and body.
func SendEmail(to, subject string, data EmailData) error {
	// Get email password from environment variable
	password := os.Getenv("EMAIL_PASSWORD")

	// Create a new email message
	mail := gomail.NewMessage()

	// Set the sender and recipient details
	mail.SetHeader("From", "su.victor21@gmail.com") // Sender's email
	mail.SetHeader("To", to)                         // Recipient's email
	mail.SetHeader("Subject", subject)               // Email subject

	// Parse and execute the HTML email template
	tmpl, err := template.New("emailTemplate").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Stock Price Alert</title>
		</head>
		<body>
			<h1>Stock Price Alert for {{.StockSymbol}}</h1>
			<p>The stock {{.StockSymbol}} has changed by {{printf "%.4f" .PercentageChange}}% within the last 24 hours.</p>
			<p>Current Price: ${{.Price}}</p>
			<p>Time of Change: {{.Time}}</p>
			<p>Best Regards,<br>Your Stock Monitoring System</p>
		</body>
		</html>
	`)

	if err != nil {
		log.Printf("Error parsing email template: %v", err)
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	// Execute the template and write the result to the buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Set the email body as HTML
	mail.SetBody("text/html", body.String())

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
