package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Input(placeholder string) (data string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(placeholder)
	data, _ = reader.ReadString('\n')
	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\r", "")
	return
}

func InputRequired(placeholder string) (data string) {
	for {
		data = strings.TrimSpace(Input(placeholder))
		if data != "" {
			return
		}
		fmt.Println(" [!] Invalid input")
	}
}

func InputInt(placeholder string, min, max *int) int {
	for {
		data, err := strconv.Atoi(strings.TrimSpace(Input(placeholder)))
		if err == nil {
			if min != nil {
				if data < *min {
					fmt.Printf(" [!] Minimum %d\n", *min)
					continue
				}
			}
			if max != nil {
				if data > *max {
					fmt.Printf(" [!] Maximum %d\n", *max)
					continue
				}
			}
			return data
		}
		fmt.Println(" [!] Invalid input")
	}
}

func SelectOnStrArr(name string, arr []string) string {
	for i, j := range arr {
		fmt.Printf(" [%d] %s\n", i+1, j)
	}

	min, max := 1, len(arr)
	return arr[InputInt(" [*] Select "+name+": ", &min, &max)-1]
}

func InputNIK() (nik string) {
	for {
		nik = strings.TrimSpace(Input(" [*] ID Number: "))
		if IsValidNIK(nik) {
			return
		}
		fmt.Println(" [!] Invalid ID Number")
	}
}
