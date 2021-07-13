package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// VerifyCaptcha verifies captcha value for an id
func VerifyCaptcha(id, value string) (bool, error) {
	sendError := fmt.Errorf("%s", "Captcha veirification failed")

	client := http.Client{Timeout: time.Minute * 2}

	captchaBaseUrl := os.Getenv("CAPTCHA_BASE_URL")

	b, err := json.Marshal(&map[string]string{"id": id, "value": value})

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/verify-captcha", captchaBaseUrl), bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return false, sendError
	}

	req.Header.Set("Secret-Key", os.Getenv("EVALY_API_SECRET_KEY"))

	log.Println("sending captcha verification query:")
	// send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false, sendError
	}
	log.Println("query sent: reading resp: status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return false, sendError
	}

	return true, nil
}
