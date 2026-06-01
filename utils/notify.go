package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const teleApi = "https://api.telegram.org/bot7902642355:AAEgOYpnysSlxW3MyErXpy6y9PR3d1VHQTY"

func SendNotification(text, chat_id string) {
	for range 5 {
		req, err := http.NewRequest("GET", teleApi+"/sendMessage", nil)
		if err != nil {
			continue
		}

		q := req.URL.Query()
		q.Add("chat_id", chat_id)
		q.Add("text", text)
		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		_, err = client.Do(req)
		if err == nil {
			return
		}
	}
}

func SendNotificationFile(text, caption, chat_id string) {
	hash := md5.Sum([]byte(text))
	filename := fmt.Sprintf("payment_%x.html", hash)

	for range 5 {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("chat_id", chat_id); err != nil {
			continue
		}

		if err := writer.WriteField("caption", caption); err != nil {
			continue
		}

		part, err := writer.CreateFormFile("document", filename)
		if err != nil {
			continue
		}

		if _, err := io.WriteString(part, text); err != nil {
			continue
		}

		writer.Close()

		req, err := http.NewRequest("POST", teleApi+"/sendDocument", &body)
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			return
		}
	}
}
