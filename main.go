package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

// targets.yaml dosyasını okumak için
type Config struct {
	Targets []string `yaml:"targets"`
}

const (
	Tor_adresi    = "socks5://127.0.0.1:9050" // tor servisi
	cikti_klasoru = "output"                  // çıktılar buraya
	durum_raporu  = "scan_report.log"         // durum raporu
	bekleme       = 90                        // bekleme
)

func main() {
	//log dosyasını oluşturma
	f, err := os.Create(durum_raporu)
	if err != nil {
		log.Fatalf("Log dosyası oluşturulamadı: %v", err)
	}
	f.Close()

	logMsg("[INIT] Thor Scraper Başlatılıyor...")
	logMsg("[INIT] Log sistemi aktif: " + durum_raporu)

	//output klasörünü oluşturma
	if _, err := os.Stat(cikti_klasoru); os.IsNotExist(err) {
		os.Mkdir(cikti_klasoru, 0755)
	}

	//hedefleri ekrana bastır
	targets, err := readTargets("targets.yaml")
	if err != nil {
		logMsg(fmt.Sprintf("[FATAL] targets.yaml okunamadı: %v", err))
		return
	}
	logMsg(fmt.Sprintf("[INIT] %d adet hedef URL yüklendi.", len(targets)))

	// 4. KONSOL MENÜSÜ (ARAYÜZ)
	for {
		fmt.Println("\n#########################################")
		fmt.Println("#       THOR SCRAPER - SECIM EKRANI     #")
		fmt.Println("#########################################")

		for i, url := range targets {
			fmt.Printf("[%d] %s\n", i+1, url)
		}
		fmt.Println("-----------------------------------------")
		fmt.Println("[A] Tüm hedefleri tara")
		fmt.Println("[Q] Çıkış")
		fmt.Print("Seçiminiz (Numara veya Harf): ")

		var input string
		fmt.Scanln(&input)

		//çıkış
		if strings.ToLower(input) == "q" {
			logMsg("[INFO] Çıkış yapılıyor.")
			break
		}

		//toplu tarama
		if strings.ToLower(input) == "a" {
			logMsg("[INFO] Toplu tarama başlatıldı.")
			for _, url := range targets {
				// PDF Madde 1.2: URL Temizleme (Whitespace trimming)
				cleanURL := strings.TrimSpace(url)
				runScraper(cleanURL)
			}
			logMsg("[INFO] Toplu tarama tamamlandı.")
			continue
		}

		//tekil seçim
		var choice int
		fmt.Sscanf(input, "%d", &choice)

		if choice > 0 && choice <= len(targets) {
			selectedURL := strings.TrimSpace(targets[choice-1])
			logMsg(fmt.Sprintf("[USER] Tekil Hedef Seçildi: %s", selectedURL))
			runScraper(selectedURL)
		} else {
			fmt.Println("❌ Geçersiz seçim! Tekrar dene.")
		}
	}
}

// scraper ayarları
func runScraper(targetURL string) {
	//chrome a tor u entegre etmece
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(Tor_adresi),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(bekleme))
	defer cancel()

	var html_icerigi string
	var sayfa_basligi string
	var ekran_goruntusu []byte

	logMsg(fmt.Sprintf("[SCAN] Bağlanılıyor: %s ...", targetURL))

	//işlemler işte
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(time.Second*15),
		chromedp.Title(&sayfa_basligi),
		chromedp.OuterHTML("html", &html_icerigi),
		chromedp.FullScreenshot(&ekran_goruntusu, 90),
	)

	if err != nil {
		//hata durumunda durma logla ve devam et
		logMsg(fmt.Sprintf("[FAIL] Tarama Hatası (%s): %v", targetURL, err))
		return
	}

	logMsg(fmt.Sprintf("[SUCCESS] Bağlantı Başarılı! Site Başlığı: %s", sayfa_basligi))

	//verileri kaydet ilgili yerlere
	saveData(targetURL, sayfa_basligi, html_icerigi, ekran_goruntusu)
}

func saveData(urlStr, title, htmlData string, screenshotData []byte) {

	cleanFolderName := sanitizeFilename(urlStr)

	siteDir := filepath.Join(cikti_klasoru, cleanFolderName)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		os.Mkdir(siteDir, 0755)
	}

	htmlPath := filepath.Join(siteDir, "index.html")
	writeWithBufio(htmlPath, []byte(htmlData))
	logMsg(fmt.Sprintf("[SAVE] HTML kaydedildi: %s", htmlPath))

	pngPath := filepath.Join(siteDir, "screenshot.png")
	writeWithBufio(pngPath, screenshotData)
	logMsg(fmt.Sprintf("[SAVE] Screenshot kaydedildi: %s", pngPath))

	infoPath := filepath.Join(siteDir, "site_info.txt")
	infoText := fmt.Sprintf("URL: %s\nTitle: %s\nScan Date: %s", urlStr, title, time.Now().Format(time.RFC1123))
	writeWithBufio(infoPath, []byte(infoText))
}

func writeWithBufio(filename string, data []byte) {
	file, err := os.Create(filename)
	if err != nil {
		logMsg(fmt.Sprintf("[ERR] Dosya oluşturma hatası: %v", err))
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(data)
	if err != nil {
		logMsg(fmt.Sprintf("[ERR] Yazma hatası: %v", err))
		return
	}

	err = writer.Flush()
	if err != nil {
		logMsg(fmt.Sprintf("[ERR] Flush hatası: %v", err))
	}
}

// loglama işi
func logMsg(message string) {

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fullMsg := fmt.Sprintf("[%s] %s", timestamp, message)
	fmt.Println(fullMsg)

	f, err := os.OpenFile(durum_raporu, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Log dosyasına erişilemedi: %v\n", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.WriteString(fullMsg + "\n")
	writer.Flush()
}

func sanitizeFilename(name string) string {
	// http ve slashları temizle
	name = strings.ReplaceAll(name, "http://", "")
	name = strings.ReplaceAll(name, "https://", "")
	name = strings.ReplaceAll(name, "/", "")

	reg, _ := regexp.Compile("[^a-zA-Z0-9.-]+")
	safe := reg.ReplaceAllString(name, "_")

	if len(safe) > 150 {
		safe = safe[:150]
	}
	return safe
}

func readTargets(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	return config.Targets, nil
}
