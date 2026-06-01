package menu

import (
	"Nux-xader/kaiweb/models"
	"Nux-xader/kaiweb/utils"
	"fmt"
	"strings"
)

func setupTotPsg() (totPsgAdult, totPsgInfant int) {
	minPsgAdult, maxPsgAdult := 1, 4
	minPsgInfant, maxPsgInfant := 0, 4
	totPsgAdult = utils.InputInt(" [*] Total Adult: ", &minPsgAdult, &maxPsgAdult)
	totPsgInfant = utils.InputInt(" [*] Total Infant: ", &minPsgInfant, &maxPsgInfant)

	return
}

func inputPsgKaiWeb(psgType string) (psg models.KAIWebPassengerData) {
	psg.PassengerIdType = "ktp"
	psg.PassengerType = psgType

	for {
		gender := strings.ToLower(strings.TrimSpace(utils.Input(" [*] Gender (f/m): ")))
		switch gender {
		case "f":
			psg.PassengerTitle = "MR"
		case "m":
			psg.PassengerTitle = "MRS."
		default:
			fmt.Println(" [!] Invalid gender")
			continue
		}
		break
	}

	for {
		name := strings.TrimSpace(utils.Input(" [*] Name: "))
		if utils.IsValidName(name) {
			psg.PassengerName = strings.ToUpper(name)
			break
		}
		fmt.Println(" [!] Invalid name")
	}

	psg.PassengerId = utils.InputNIK()
	return

}

func setupPsgKaiWeb() (psgs []models.KAIWebPassengerData) {
	totPsgAdult, totPsgInfant := setupTotPsg()

	fmt.Println(" [+] Input Adult")
	fmt.Print(" _______________________________\n")
	for range totPsgAdult {
		psgs = append(psgs, inputPsgKaiWeb("A"))
		fmt.Print(" _______________________________\n")
	}

	if totPsgInfant > 0 {
		fmt.Println(" [+] Input Infant")
		fmt.Print(" _______________________________\n")
		for range totPsgInfant {
			psgs = append(psgs, inputPsgKaiWeb("I"))
			fmt.Print(" _______________________________\n")
		}
	}

	return
}
