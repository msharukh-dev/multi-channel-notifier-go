package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Send routes notification to the appropriate service
func Send(typeSend, to, message, webhookURL string) error {
	switch typeSend {
	case "email":
		return sendEmailMailtrap(to, message)
	case "webhook":
		return sendWebhook(to, message, webhookURL)
	case "sms":
		return sendSMS(to, message)
	default:
		return fmt.Errorf("unknown type: %s", typeSend)
	}
}

/*
MAILTRAP EMAIL SENDER
Uses Mailtrap Email Sending HTTP API
Docs: https://api-docs.mailtrap.io
*/
func sendEmailMailtrap(to, message string) error {
	apiToken := os.Getenv("MAILTRAP_API_TOKEN")
	fromEmail := os.Getenv("MAILTRAP_FROM_EMAIL")

	if apiToken == "" || fromEmail == "" {
		return fmt.Errorf("mailtrap credentials not configured")
	}

	payload := map[string]interface{}{
		"from": map[string]string{
			"email": fromEmail,
			"name":  "Webhook API",
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject": "Notification from Webhook API",
		"text":    message,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		"https://send.api.mailtrap.io/api/send",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("mailtrap send failed: %s", resp.Status)
	}

	return nil
}

// Send SMS via Twilio
func sendSMS(to, message string) error {
	twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhone := os.Getenv("TWILIO_PHONE_NUMBER")

	if twilioSID == "" || twilioToken == "" || twilioPhone == "" {
		return fmt.Errorf("twilio credentials not configured")
	}

	data := fmt.Sprintf("To=%s&From=%s&Body=%s", to, twilioPhone, message)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", twilioSID),
		bytes.NewBufferString(data),
	)
	if err != nil {
		return err
	}

	req.SetBasicAuth(twilioSID, twilioToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("twilio send failed: %s", resp.Status)
	}

	return nil
}

// Send webhook POST request
func sendWebhook(webhookURL, message, clientWebhookURL string) error {
	if webhookURL == "" && clientWebhookURL == "" {
		return fmt.Errorf("webhook URL not provided")
	}

	// Use provided webhook URL or client's default
	url := webhookURL
	if url == "" {
		url = clientWebhookURL
	}

	payload := map[string]interface{}{
		"message":   message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook delivery failed: %s", resp.Status)
	}

	return nil
}
