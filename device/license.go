package device

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"Nux-xader/kaiweb/gvars"
	thrid_party "Nux-xader/kaiweb/third-party"
	"Nux-xader/kaiweb/utils"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/yusufpapurcu/wmi"
)

const (
	secretKey  = "kaiweb_nux_v4.1.1"
	maxSession = 20
)

var SessionId = utils.MD5Hash(time.Now().String())

func IsLimited() bool {
	resp, err := resty.New().SetTimeout(10*time.Second).R().
		SetHeader("X-Device-Id", *gvars.DeviceID).
		SetHeader("X-Session-Id", SessionId).
		Get("http://debora.it/sinedya/api/ping")

	if err != nil {
		if strings.Contains(err.Error(), "debora.it") || strings.Contains(err.Error(), "sinedya") {
			utils.WarningPrint("Something different with server response")
			return false
		} else {
			utils.WarningPrint(err.Error())
			return true
		}
	}

	totSession, err := strconv.Atoi(resp.String())
	if err != nil {
		utils.DangerPrint(err.Error())
		return true
	}

	return totSession > maxSession
}

func Hash(input []byte) string {
	var salt = secretKey
	if len(input) > 1 {
		salt = string(input[:2]) + salt
	}
	combinedInput := append([]byte(salt), input...)
	hash := sha256.New()
	hash.Write(combinedInput)
	return hex.EncodeToString(hash.Sum(nil))
}

func key(salt, deviceID string) (result string, c string) {
	result += Hash([]byte(salt))
	result += Hash([]byte(secretKey))
	result += Hash([]byte(deviceID))

	runes := []rune(result)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	result = Hash([]byte(strings.ToUpper(Hash([]byte(string(runes))))))

	return
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("pkcs7: data is empty")
	}
	if length%blockSize != 0 {
		return nil, fmt.Errorf("pkcs7: data is not a multiple of block size")
	}
	paddingLen := int(data[length-1])
	if paddingLen == 0 || paddingLen > blockSize {
		return nil, fmt.Errorf("pkcs7: invalid padding length")
	}
	for i := range paddingLen {
		if data[length-1-i] != byte(paddingLen) {
			return nil, fmt.Errorf("pkcs7: invalid padding")
		}
	}
	return data[:length-paddingLen], nil
}

func GetName(data string) (name string, c string) {
	if len(data) > 500 {
		return
	}
	parsedData := strings.Split(data, "||")
	if len(parsedData) != 2 {
		return
	}
	if len(parsedData[1]) != 19 && len(parsedData[1]) != 10 {
		return
	}

	key, c := key(parsedData[1], *gvars.DeviceID)
	if c != "" {
		return
	}

	block, err := aes.NewCipher([]byte(key[16-5 : 32-5]))
	if err != nil {
		c = "BKCP 0x10001"
		return
	}

	cipherText, err := base64.StdEncoding.DecodeString(parsedData[0])
	if err != nil {
		c = "BDC 0x10001"
		return
	}

	cbc := cipher.NewCBCDecrypter(block, []byte(key[8+4 : 24+4])[:aes.BlockSize])

	plainText := make([]byte, len(cipherText))
	cbc.CryptBlocks(plainText, cipherText)

	plainText, err = pkcs7Unpad(plainText, aes.BlockSize)
	if err != nil {
		return
	}

	result := strings.Split(string(plainText), "||")
	if len(result) != 3 {
		return
	}
	if len(result[2]) != 19 && len(result[2]) != 10 {
		return
	}
	if result[1] != *gvars.DeviceID {
		return
	}
	expDate, err := time.Parse("2006-01-02", result[2])
	if err != nil {
		return
	}

	var currentTime *int64
	for currentTime == nil {
		var statusCode *int
		currentTime, statusCode = thrid_party.CurrentTime()
		if currentTime == nil && statusCode == nil {
			color.Red(" [!] Please check your internet connection")
		}
		time.Sleep(2 * time.Second)
	}

	gvars.LicenseExpTime = expDate.Unix()
	if *currentTime > expDate.Unix() {
		return
	}

	name = result[0]
	return
}

func InitDeviceID() (c string) {
	var (
		genUid       []byte
		serialNumber string
	)

	{
		userConfDir, err := os.UserConfigDir()
		if err != nil {
			c = "UCF 0x01011"
			return
		}

		uidPath := filepath.Join(userConfDir, ".uid.bin")
		genUid, err = os.ReadFile(uidPath)
		if err != nil {
			data := make([]byte, 16)
			_, err := rand.New(rand.NewSource(time.Now().UnixNano())).Read(data)
			if err != nil {
				c = "GUID 0x01011"
				return
			}
			genUid = data

			err = os.WriteFile(
				uidPath,
				genUid, 0644,
			)
			if err != nil {
				c = "SGUID 0x01011"
				return
			}
		}
		serialNumber += string(genUid)
	}

	{
		var biosInfo []struct {
			SerialNumber string
		}
		err := wmi.Query("SELECT SerialNumber FROM Win32_BIOS", &biosInfo)
		if err == nil {
			if len(biosInfo) > 0 {
				serialNumber += biosInfo[0].SerialNumber
			}
		}
	}

	{
		var csProductInfo []struct {
			UUID string
		}
		err := wmi.Query("SELECT UUID FROM Win32_ComputerSystemProduct", &csProductInfo)
		if err == nil {
			if len(csProductInfo) > 0 {
				serialNumber += csProductInfo[0].UUID
			}
		}
	}

	{
		var procInfo []struct {
			ProcessorId string
		}
		err := wmi.Query("SELECT ProcessorId FROM Win32_Processor", &procInfo)
		if err == nil {
			if len(procInfo) > 0 {
				serialNumber += procInfo[0].ProcessorId
			}
		}
	}

	{
		var bbInfo []struct {
			SerialNumber string
		}
		err := wmi.Query("SELECT SerialNumber FROM Win32_BaseBoard", &bbInfo)
		if err == nil {
			if len(bbInfo) > 0 {
				serialNumber += bbInfo[0].SerialNumber
			}
		}
	}

	result := Hash([]byte(serialNumber))
	result += Hash([]byte(serialNumber[:3]))
	result = strings.ToUpper(Hash([]byte(result))[:24])

	gvars.DeviceID = &result

	return
}
