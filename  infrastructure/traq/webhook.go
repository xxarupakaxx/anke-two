package traq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"
)

type Webhook struct {
}

func (w *Webhook) PostMessage(message string) error {
	url := "https://q.trap.jp/api/v3/webhooks/" + os.Getenv("TRAQ_WEBHOOK_ID")

	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	req.Header.Set("X-TRAQ-Signature", calcHMACHSA1(message))

	query := netUrl.Values{}
	query.Add("embed", "1")
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body:%w", err)
		}
	}(resp.Body)

	sb := &strings.Builder{}
	_, err = io.Copy(sb, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Message sent to %s, message: %s, response: %s\n", url, message, sb.String())

	return nil
}

func calcHMACHSA1(message string) string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
