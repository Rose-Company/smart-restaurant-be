package common

import (
	"app-noti/config"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOTPEmail(fromEmail, toEmail, otp string) error {
	apiKey := config.Config.Mail.ApiKey
	if apiKey == "" {
		return fmt.Errorf("mail api_key is not configured in config.yaml")
	}

	from := mail.NewEmail("Smart Restaurant", config.Config.Mail.FromEmail)
	to := mail.NewEmail("Recipient", toEmail)
	subject := "Your OTP Code"
	plainTextContent := fmt.Sprintf("Your OTP to verify email is: %s (expires in 5 minutes)", otp)
	htmlContent := fmt.Sprintf("<strong>Your OTP Code:</strong> <h2>%s</h2><p>This code expires in 5 minutes</p>", otp)

	log.Printf("[EMAIL] Sending OTP - FROM: %s (%s) | TO: %s | OTP: %s", from.Name, config.Config.Mail.FromEmail, toEmail, otp)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Printf("[EMAIL ERROR] Failed to send: %v", err)
		return err
	}

	log.Printf("[EMAIL SUCCESS] Status Code: %d | Body: %s", response.StatusCode, response.Body)
	return nil
}
