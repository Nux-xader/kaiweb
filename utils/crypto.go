package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
)

var salt = "415dded27a9e61b20367be37e875dea3"

func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

func GenerateNonceHexStr() (*string, error) {
	nonceByte, err := GenerateNonce()
	if err != nil {
		return nil, err
	}

	result := hex.EncodeToString(nonceByte)
	return &result, nil
}

func GenerateKey(secondarySalt string, data string) string {
	// Step 1: Character interleaving
	var part1, part2 strings.Builder
	for i, c := range data {
		if i%2 == 0 {
			part1.WriteRune(c)
		} else {
			part2.WriteRune(c)
		}
	}

	// Step 2: Reverse secondary salt
	reversedSalt := reverse(secondarySalt)

	// Step 3: XOR secondary salt with primary salt bytes
	saltInt := int64(0)
	for i, c := range secondarySalt {
		saltInt ^= int64(c) << (i % 8 * 8)
	}
	primarySaltBytes := []byte(salt)
	for i := 0; i < len(primarySaltBytes) && i < 8; i++ {
		saltInt ^= int64(primarySaltBytes[i]) << (i * 8)
	}
	obfuscatedSalt := strconv.FormatInt(saltInt, 16)

	// Step 4: Bit rotation on data bytes
	dataBytes := []byte(data)
	rotationAmount := len(secondarySalt) % 8
	if rotationAmount == 0 {
		rotationAmount = 1
	}
	rotatedData := rotate(dataBytes, rotationAmount)

	// Step 5: Build complex data by interleaving all obfuscated components
	halfSalt := len(salt) / 2
	var complexData strings.Builder
	complexData.WriteString(part1.String())
	complexData.WriteString(obfuscatedSalt)
	complexData.WriteString(salt[:halfSalt])
	complexData.WriteString(part2.String())
	complexData.WriteString(reversedSalt)
	complexData.WriteString(salt[halfSalt:])
	complexData.WriteString(hex.EncodeToString(rotatedData))
	complexData.WriteString(secondarySalt)

	// Step 6: Hash ONCE at the end - efficient but still unpredictable
	return MD5Hash(complexData.String())
}

func MD5Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func rotate(data []byte, positions int) []byte {
	n := len(data)
	positions = positions % n
	result := make([]byte, n)
	for i := range n {
		result[i] = data[(i+positions)%n]
	}
	return result
}
