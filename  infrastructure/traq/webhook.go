package traq

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/domain/repository/traq"
	"gopkg.in/guregu/null.v4"
	"io"
	"log"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"
)

type webhook struct {
}

func NewWebhook() traq.IWebhook {
	return &webhook{}
}

func (w *webhook) PostMessage(message string) error {
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

func (w *webhook) CreateQuestionnaireMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit null.Time, targets []string) string {
	var resTimeLimitText string
	if resTimeLimit.Valid {
		resTimeLimitText = resTimeLimit.Time.Local().Format("2021/11/11 15:00")
	} else {
		resTimeLimitText = "なし"
	}

	var targetsMentionText string
	if len(targets) == 0 {
		targetsMentionText = "なし"
	} else {
		targetsMentionText = "@" + strings.Join(targets, " @")
	}

	return fmt.Sprintf(`### アンケート『[%s](https://anke-to.trap.jp/questionnaires/%d)』が作成されました
#### 管理者
%s
#### 説明
%s
#### 回答期限
%s
#### 対象者
%s
#### 回答リンク
https://anke-to.trap.jp/responses/new/%d`,
		title,
		questionnaireID,
		strings.Join(administrators, ","),
		description,
		resTimeLimitText,
		targetsMentionText,
		questionnaireID)

}

func calcHMACHSA1(message string) string {
	mac := hmac.New(sha1.New, []byte(os.Getenv("TRAQ_WEBHOOK_SECRET")))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
