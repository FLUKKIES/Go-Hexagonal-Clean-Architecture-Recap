package smsservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
)

// twilioSMSService ส่ง SMS ผ่าน Twilio REST API
// ไม่ใช้ Third-party SDK — ใช้ net/http โดยตรงเพื่อลด Dependency
type twilioSMSService struct {
	accountSID string
	authToken  string
	fromNumber string // Twilio Phone Number เช่น "+12025551234"
}

func NewTwilioSMSService(accountSID, authToken, fromNumber string) ports.ISMSService {
	return &twilioSMSService{
		accountSID: accountSID,
		authToken:  authToken,
		fromNumber: fromNumber,
	}
}

func (t *twilioSMSService) SendOTP(phoneNumber, otp string) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.accountSID)

	data := url.Values{}
	data.Set("To", phoneNumber)
	data.Set("From", t.fromNumber)
	data.Set("Body", fmt.Sprintf("Your verification code is: %s\nThis code expires in 5 minutes.", otp))

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		var twilioErr struct {
			Message string `json:"message"`
		}
		json.Unmarshal(body, &twilioErr)
		return fmt.Errorf("twilio error: %s", twilioErr.Message)
	}

	return nil
}
