package menu

import (
	"Nux-xader/kaiweb/engine"
	"Nux-xader/kaiweb/gvars"
	"Nux-xader/kaiweb/models"
	"Nux-xader/kaiweb/utils"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type EventHandler struct {
	Browser       *engine.Browser
	Psgs          []models.KAIWebPassengerData
	Delay         time.Duration
	CaptchaApi    string
	CaptchaApiKey string
	Tried         int
}

func Buy() {
	browser, err := engine.InitBrowser()
	if err != nil {
		utils.DangerPrint("load browser err: " + err.Error())
		os.Exit(1)
	}
	go browser.PatchSoldOut()

	fmt.Println("     ---")
	captchaApi := utils.SelectOnStrArr("Captcha Service", utils.CAPTCHA_APIS)
	captchaApiKey := utils.InputRequired(" [*] Enter captcha api key: ")
	fmt.Print(utils.LineSep)

	psgs := setupPsgKaiWeb()
	psgsJsonStr, _ := json.MarshalIndent(psgs, "", "    ")
	fmt.Println(" [+] The following passengers will be used:")
	fmt.Println(string(psgsJsonStr))
	utils.Input("[Press Enter to Continue] ")

	minDelay := 4
	delay := time.Millisecond * time.Duration(utils.InputInt(" [*] Delay (in miliseconds): ", &minDelay, nil))

	e := &EventHandler{
		Browser:       browser,
		Psgs:          psgs,
		Delay:         delay,
		CaptchaApi:    captchaApi,
		CaptchaApiKey: captchaApiKey,
	}

	gotIt := false
	for !gotIt {
		err = browser.Page.WaitLoad()
		if err != nil {
			utils.DangerPrint("Wait loading page err: " + err.Error())
			browser.Page.Reload()
			continue
		}

		pageInfo, err := browser.Page.Info()
		if err != nil {
			utils.DangerPrint("page info err :" + err.Error())
			continue
		}
		currentUrl := pageInfo.URL

		switch {
		case strings.Contains(currentUrl, "/passengerdata"):
			e.submitBooking()
		case strings.Contains(currentUrl, "/passengercontrol"):
			e.confirmBooking()
		case strings.Contains(currentUrl, "/paytype"):
			e.selectPayment()
		case strings.Contains(currentUrl, "/payment"):
			gotIt = true
		case currentUrl == "https://booking.kai.id/" || currentUrl == "https://booking.kai.id":
			browser.Page.Navigate("https://booking.kai.id")
		}
	}
}

func (e *EventHandler) submitBooking() {
	_, err := e.Browser.Page.Element("#pesanan")
	if err != nil {
		utils.DangerPrint("Form pesanan tidak ditemukan: " + err.Error())
		e.Browser.Page.Reload()
		return
	}

	ordererIdentity := utils.FakeIdentity()
	_, err = e.Browser.Page.Eval(`() => {
		document.querySelector('#pemesan_nama').value = '` + ordererIdentity.Name + `';
		document.querySelector('#pemesan_notandapengenal').value = '` + ordererIdentity.Id + `';
		document.querySelector('#pemesan_email').value = '` + ordererIdentity.Email + `';
		document.querySelector('#pemesan_nohp').value = '` + ordererIdentity.Phone + `';
		document.querySelector('#pemesan_alamat').value = '` + ordererIdentity.Address + `';
	}`)
	if err != nil {
		utils.DangerPrint("Pengisian data pemesan gagal: " + err.Error())
		e.Browser.Page.Reload()
		return
	}

	adultI, _ := 0, 0
	for _, psg := range e.Psgs {
		switch psg.PassengerType {
		case "A":
			i := "[" + strconv.Itoa(adultI) + "]"
			_, err = e.Browser.Page.Eval(`() => {
				const select = document.querySelectorAll('select[name="penumpang_title[]"]')` + i + `;
				select.value = '` + strings.ToUpper(psg.PassengerTitle) + `';
				select.dispatchEvent(new Event('change'));

				document.querySelectorAll('input[name="penumpang_nama[]"]')` + i + `.value = '` + psg.PassengerName + `';
				document.querySelectorAll('input[name="penumpang_notandapengenal[]"]')` + i + `.value = '` + psg.PassengerId + `';
			}`)
			if err != nil {
				utils.DangerPrint(fmt.Sprintf("Pengisian data penumpang ke %d err: %s", adultI, err.Error()))
				e.Browser.Page.Reload()
				return
			}
			adultI++

		case "I":
		}
	}

	b64Img, err := e.Browser.WaitAndGrabImage("#pesanan > div:nth-child(3) > div > div:nth-child(4) > img")
	if err != nil {
		utils.DangerPrint("gagal mendaptkan captcha : " + err.Error())
		e.Browser.Page.Reload()
		return
	}

	captchaVal, err := utils.ImageOCR3RdParty("", *b64Img, e.CaptchaApi, e.CaptchaApiKey)
	if err != nil {
		utils.DangerPrint("gagal membaca captcha : " + err.Error())
		e.Browser.Page.Reload()
		return
	}
	utils.OriPrint("Captcha terbaca: " + captchaVal)

	_, err = e.Browser.Page.Eval(`() => document.querySelector('#captcha').value = "` + captchaVal + `";`)
	if err != nil {
		utils.DangerPrint("gagal input captcha : " + err.Error())
		e.Browser.Page.Reload()
		return
	}

	_, err = e.Browser.Page.Eval(`() => {
		document.querySelector('#setuju').click();
		document.querySelector('#pesanan').submit();
	}`)
	if err != nil {
		utils.DangerPrint("Submit data gagal err: " + err.Error())
		e.Browser.Page.Reload()
		return
	}

	e.Tried += 1
	if e.Tried >= 6 {
		for n := range e.Psgs {
			e.Psgs[n].PassengerId = utils.NikRotator(e.Psgs[n].PassengerId)
		}
		e.Tried = 0
	}

	time.Sleep(e.Delay)
}

func (e *EventHandler) confirmBooking() {
	err := e.Browser.Page.Navigate("https://booking.kai.id/paytype")
	if err != nil {
		utils.DangerPrint("gagal konfirmasi booking: " + err.Error())
		e.Browser.Page.Reload()
		return
	}
}

func (e *EventHandler) selectPayment() {
	_, err := e.Browser.Page.Eval(`() => document.querySelector("#payForm5").submit()`)
	if err != nil {
		utils.DangerPrint("gagal pilih pembayaran: " + err.Error())
		e.Browser.Page.Reload()
		return
	}

	time.Sleep(3 * time.Second)
	e.Browser.Page.WaitLoad()
	time.Sleep(3 * time.Second)
	html, err := e.Browser.Page.HTML()
	if err != nil {
		utils.DangerPrint("Gagal mendapatkan html pembayaran: " + err.Error())
	}

	go utils.SendNotificationFile(html, "Dapet tiket nih, dengan penumpang: "+e.Psgs[0].PassengerId, *gvars.TeleID)
}
