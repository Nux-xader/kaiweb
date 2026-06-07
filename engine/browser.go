package engine

import (
	"Nux-xader/kaiweb/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const (
	cdpPort            = 9222
	chromeBin          = `C:\Program Files\Google\Chrome\Application\chrome.exe`
	userData           = `C:\temp\chrome_master_isolated`
	chromiumSuffixPath = "Application\\chrome.exe"
)

var (
	appData           = strings.TrimRight(os.Getenv("APPDATA"), "Roaming")
	chromiumOriDir    = appData + "Local\\Chromium\\"
	chromiumModDir    = appData + "Local\\com.windm.co\\"
	chromiumShortcut  = appData + "Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Chromium.lnk"
	chromiumShortcut2 = strings.TrimRight(appData, "AppData\\") + "\\Desktop\\Chromium.lnk"
)

var cdpUrl = fmt.Sprintf("http://127.0.0.1:%d/json/version", cdpPort)

type Browser struct {
	Ctx         *rod.Browser
	Page        *rod.Page
	chromeCmd   *exec.Cmd
	userDataDir string
}

func isCDPReady() bool {
	client := &http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get(cdpUrl)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func getWebSocketURL() string {
	resp, err := http.Get(cdpUrl)
	if err != nil {
		log.Fatal("ws not ready:", err)
	}
	defer resp.Body.Close()

	var info struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Fatal("failed read ws data:", err)
	}
	return info.WebSocketDebuggerURL
}

func setupUngoogledChromium() {
	const chromeSetupFile = "temp\\chrome.exe"

	go os.Remove(chromeSetupFile)
	go os.Remove(chromiumShortcut)
	go os.Remove(chromiumShortcut2)
	if utils.IsFileExists(chromiumModDir + chromiumSuffixPath) {
		return
	}
	os.RemoveAll(chromiumOriDir)

	utils.Download(
		"https://github.com/ungoogled-software/ungoogled-chromium-windows/releases/download/148.0.7778.178-1.1/ungoogled-chromium_148.0.7778.178-1.1_installer_x64.exe",
		chromeSetupFile,
	)

	utils.OriPrint("Running Monkey Patcher for chrome")
	exec.Command(chromeSetupFile).Run()
	for {
		if utils.IsFileExists(chromiumOriDir) {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	go os.Remove(chromeSetupFile)
	go os.Remove(chromiumShortcut)
	go os.Remove(chromiumShortcut2)
	os.Rename(chromiumOriDir, chromiumModDir)
	utils.OriPrint("Successfully patch chrome")
}

func InitBrowser() (*Browser, error) {
	setupUngoogledChromium()
	userDataDir, err := os.MkdirTemp("", "chrome_userdata_*")
	if err != nil {
		return nil, fmt.Errorf("gagal buat temp dir: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		os.RemoveAll(userDataDir)
		return nil, fmt.Errorf("gagal cari port random: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	chromePath := chromiumModDir + chromiumSuffixPath
	if _, err := os.Stat(chromePath); os.IsNotExist(err) {
		os.RemoveAll(userDataDir)
		return nil, fmt.Errorf("binary Chrome tidak ditemukan")
	}

	cmd := exec.Command(
		chromePath,
		fmt.Sprintf("--remote-debugging-port=%d", port),
		fmt.Sprintf("--user-data-dir=%s", userDataDir),
		"--no-first-run", "--no-default-browser-check",
		"--disable-gpu", "--disable-dev-shm-usage",
		"--disable-backgrounding-occluded-windows",
		"--disable-renderer-backgrounding",
		"--disable-background-timer-throttling",
		"--disable-features=CalculateNativeWinOcclusion",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-background-networking", // request tetap jalan meski background
		"--enable-features=NetworkService,NetworkServiceInProcess",
	)
	cmd.Dir = filepath.Dir(chromePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		os.RemoveAll(userDataDir)
		return nil, fmt.Errorf("gagal start Chrome: %w", err)
	}

	var wsURL string
	cdpURL := fmt.Sprintf("http://127.0.0.1:%d/json/version", port)
	timeout := time.After(1 * time.Minute)
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

pollLoop:
	for {
		select {
		case <-timeout:
			cmd.Process.Kill()
			os.RemoveAll(userDataDir)
			return nil, fmt.Errorf("Endpoint tidak ready dalam 1 menit")
		case <-ticker.C:
			resp, err := http.Get(cdpURL)
			if err != nil || resp.StatusCode != http.StatusOK {
				if resp != nil {
					resp.Body.Close()
				}
				continue
			}

			var info struct {
				WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&info); err == nil && info.WebSocketDebuggerURL != "" {
				resp.Body.Close()
				wsURL = info.WebSocketDebuggerURL
				break pollLoop
			}
			resp.Body.Close()
		}
	}

	// 5️⃣ Connect rod ke Chrome
	browser := rod.New().ControlURL(wsURL).MustConnect()
	page := browser.MustPages().First()

	// Interaksi manual sesuai flow-mu
	utils.Input(" [*] Buka tab baru lalu tutup tab sebelumnya, kemudian ENTER")
	page = browser.MustPages().Last()
	_ = proto.EmulationClearDeviceMetricsOverride{}.Call(page)

	return &Browser{
		Ctx:         browser,
		Page:        page,
		chromeCmd:   cmd,
		userDataDir: userDataDir,
	}, nil
}

// func InitBrowser() *Browser {
// 	fmt.Println(1)
// 	setupUngoogledChromium()
// 	fmt.Println(chromiumModDir + chromiumSuffixPath)
// 	// if isCDPReady() {
// 	// 	wsUrl = getWebSocketURL()
// 	// } else {

// 	if !utils.IsFileExists(chromiumModDir + chromiumSuffixPath) {
// 		utils.DangerPrint("Failed load chrome.")
// 		utils.Input("")
// 		os.Exit(1)
// 		return nil
// 	}

// 	wsUrl, err := launcher.New().
// 		// Bin(chromeBin).
// 		Bin(chromiumModDir+chromiumSuffixPath).
// 		RemoteDebuggingPort(cdpPort).
// 		// Set("user-data-dir", userData).
// 		Headless(false).
// 		Set("start-maximized").

// 		// ✅ Tambahkan flags yang membantu:
// 		Set("disable-blink-features", "AutomationControlled").
// 		Set("disable-gpu").
// 		Set("disable-dev-shm-usage").
// 		Set("no-sandbox").
// 		Set("no-first-run").
// 		Set("no-default-browser-check").
// 		Set("disable-background-networking").
// 		Set("disable-features", "IsolateOrigins,site-per-process,NetworkService").
// 		Set("log-level", "3").
// 		Launch()
// 	// }

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	browser := rod.New().ControlURL(wsUrl).MustConnect() //.MustIncognito()
// 	page := browser.MustPage()
// 	utils.Input(" [*] Buka tab baru lalu tutup tab sebelumnya, kemudian ENTER")
// 	page = browser.MustPages().Last()
// 	_ = proto.EmulationClearDeviceMetricsOverride{}.Call(page)
// 	return &Browser{Ctx: browser, Page: page}
// }

func (b *Browser) PatchSoldOut() {
	go func() {
		defer func() {
			// Tangkap panic agar tidak membunuh seluruh program saat page navigasi/detach
			if r := recover(); r != nil {
				log.Printf("[PatchSoldOut] recovered from panic: %v", r)
			}
		}()

		// EachEvent akan memblokir goroutine ini hingga page ditutup/error
		b.Page.EachEvent(func(_ *proto.PageLoadEventFired) {
			defer func() {
				// Tangkap panic per-event juga, agar loop EachEvent tidak ikut mati
				if r := recover(); r != nil {
					log.Printf("[PatchSoldOut] recovered from event panic: %v", r)
				}
			}()

			// Hanya patch jika sedang di halaman /search
			info, err := b.Page.Info()
			if err != nil || !strings.Contains(info.URL, "/search") {
				return
			}

			_, err = b.Page.Eval(`() => {
			    var MARKER = 'data-kaifixed';

			    var habisLinks = document.querySelectorAll('.habis');
			    if (!habisLinks.length) return { status: 'skip', reason: 'already-modified-or-none' };

			    var count = 0;
			    for (var i = 0; i < habisLinks.length; i++) {
			        var link = habisLinks[i];
			        if (link.hasAttribute(MARKER)) continue;

			        var form = link.closest('form');
			        var formId = form ? form.id : '';
			        if (!formId) continue;

			        link.classList.remove('habis');
			        link.classList.add('card-schedule');
			        link.style.cursor = 'pointer';
			        link.style.opacity = '1';
			        link.style.pointerEvents = 'auto';

			        link.setAttribute('onclick', 'document.getElementById("' + formId + '").submit(); return false;');
			        link.setAttribute(MARKER, 'true');
			        count++;
			    }
			    return { status: 'success', modified: count };
			}`)
			if err != nil {
				log.Printf("[PatchSoldOut] eval err: %v", err)
			}
		})()
	}()
}

func (b *Browser) WaitAndGrabImage(selector string) (*string, error) {
	for {
		// 1. Tunggu elemen muncul
		el, err := b.Page.Timeout(7 * time.Second).Element(selector)
		if err != nil {
			return nil, fmt.Errorf("elemen tidak ditemukan: %w", err)
		}

		// 2. Cek status load gambar
		res, err := el.Eval(`() => this.complete && this.naturalWidth > 0`)
		if err != nil {
			return nil, fmt.Errorf("gagal cek status gambar: %w", err)
		}
		loaded := res.Value.Bool()

		if loaded {
			// 3. Ekstrak via canvas
			resCanvas, err := el.Eval(`() => {
				const c = document.createElement('canvas');
				c.width = this.naturalWidth;
				c.height = this.naturalHeight;
				c.getContext('2d').drawImage(this, 0, 0);
				return c.toDataURL('image/png');
			}`)
			if err != nil {
				return nil, fmt.Errorf("gagal render canvas: %w", err)
			}
			dataURL := resCanvas.Value.String()

			// Decode base64
			parts := strings.SplitN(dataURL, ",", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("format dataURL tidak valid")
			}
			b64Img := parts[1]
			return &b64Img, nil
		}

		log.Println("⏳ Gambar belum load → klik & jeda...")
		if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
			return nil, fmt.Errorf("gagal klik elemen: %w", err)
		}
		time.Sleep(125 * time.Millisecond)
	}
}
