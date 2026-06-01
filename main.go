package main

import (
	"Nux-xader/kaiweb/device"
	"Nux-xader/kaiweb/gvars"
	"Nux-xader/kaiweb/menu"
	"Nux-xader/kaiweb/utils"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print(utils.Banner)
	os.Mkdir("temp", os.ModePerm)
	var teleIdPath = "teleid.txt"
	ec := device.InitDeviceID()

	if *gvars.DeviceID != "" {
		fmt.Println(" [*] Device ID:", *gvars.DeviceID)
		name, ec := device.GetName(utils.LoadLicense())

		if ec == "" && name != "" {
			if !utils.IsFileExists(teleIdPath) {
				utils.Input(" [!] Please create teleid.txt file")
				return
			}
			teleID, err := utils.Readf(teleIdPath)
			if err != nil {
				utils.Input(" [!] Failed to read teleid.txt: " + err.Error())
				return
			}
			cleanTeleID := strings.Trim(*teleID, " \n\r")
			if !utils.IsDigit(cleanTeleID) || len(cleanTeleID) < 4 {
				utils.Input(" [!] Invalid teleid.txt")
				return
			}
			gvars.TeleID = &cleanTeleID
			menu.Buy()
		} else if ec != "" {
			fmt.Println(" [!] " + ec)
		} else {
			fmt.Println(" [!] Invalid license.")
		}
	} else {
		fmt.Println(" [!] " + ec)
	}

	utils.Input(" [*] Press enter to exit")
}
