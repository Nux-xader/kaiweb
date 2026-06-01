package utils

import (
	"os"
	"strings"
)

func Writef(path string, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(data)
	return nil
}

func Readf(path string) (*string, error) {
	byteData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(byteData)
	return &content, nil
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func LoadLicense() string {
	license, err := Readf("license.txt")
	if err != nil {
		Input(" [+] Failed load your license file: " + err.Error())
		os.Exit(0)
	}

	return strings.Trim(*license, " \n\r")
}
