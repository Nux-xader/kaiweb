package utils

import (
	"net/http"
)

func SendNotification(text, chat_id string) {
	for range 5 {
		req, err := http.NewRequest("GET", "https://api.telegram.org/bot7902642355:AAEgOYpnysSlxW3MyErXpy6y9PR3d1VHQTY/sendMessage", nil)
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
