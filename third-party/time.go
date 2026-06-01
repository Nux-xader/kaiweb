package thrid_party

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"Nux-xader/kaiweb/gvars"
	"Nux-xader/kaiweb/utils"
)

type TimestampResponse struct {
	Timestamp int64   `json:"timestamp"`
	Signature *string `json:"signature"`
}

func CurrentTime() (timestamp *int64, statusCode *int) {
	nonce, err := utils.GenerateNonceHexStr()
	if err != nil {
		return
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest(
		http.MethodGet,
		"http://debora.it/sinedya/api/timestamp",
		nil,
	)
	if err != nil {
		return
	}

	req.Header.Set("X-Nonce", *nonce)
	req.Header.Set("X-Device-Id", *gvars.DeviceID)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var parsedResp TimestampResponse
	if err = json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return
	}

	timestampStr := strconv.FormatInt(parsedResp.Timestamp, 10)

	if parsedResp.Signature == nil ||
		utils.GenerateKey(timestampStr, timestampStr+"|"+*nonce) != *parsedResp.Signature {
		utils.DangerPrint("Violation detected")
		os.Exit(0)
		return
	}

	statusCode_ := resp.StatusCode

	return &parsedResp.Timestamp, &statusCode_
}
