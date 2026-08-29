```text
# Thor Scraper - Tor Network CTI Tool

Bu proje, Sibervatan "Yıldız CTI" eğitimi kapsamında geliştirilmiş, **Go (Golang)** tabanlı bir Siber Tehdit İstihbaratı (CTI) toplama aracıdır.

## Proje Amacı
Tor ağındaki .onion uzantılı sitelerden hacker temalı olanları hedef alarak o sitelerden html içeriği ekran görüntüsü verilerini çekmektedir.

## Teknik Özellikler
- **Tor Proxy Entegrasyonu:** `127.0.0.1:9050` üzerinden anonim bağlantı (SOCKS5 Proxy) sağlar ve IP sızıntısını önler.
- **Headless Chrome (Chromedp):** Modern web teknolojilerini render ederek tam sayfa ekran görüntüsü alır.
- **Dinamik Klasörleme:** Her hedefi kendi URL ismine (`example.onion/`) göre ayrı klasörlerde saklar.
- **Raporlama:** Tarama sonuçlarını anlık olarak `scan_report.log` dosyasına ve terminale işler.

## Kurulum ve Kullanım

### Gereksinimler
- Go
- Tor Service (Arka planda çalışıyor olmalı)
- Linux / macOS / Windows


```
### Çalıştırma Adımları

1. **Repoyu Klonlayın:**

```bash
git clone (https://github.com/Muzoovy4606/CTI-THOR-SCRAPER
cd CTI-THOR-SCRAPER

```


2. **Bağımlılıkları Yükleyin:**

```bash
go mod tidy

```

3. **targets.yaml Dosyasını Oluşturun:**
Güvenlik ve gizlilik nedeniyle (Ransomware/Hacker grubu linkleri içerdiği için) **targets.yaml** dosyası repoya yüklenmemiştir. Projeyi çalıştırmadan önce ana dizinde `targets.yaml` dosyası oluşturup hedeflerinizi içine eklemelisiniz.
4. **Tor Servisini Başlatın (Linux/Arch):**

```bash
sudo systemctl start tor

```

*(Not: Tor servisinin 9050 portunda çalıştığından emin olun.)*

5. **Aracı Çalıştırın:**

```bash
go run main.go

```

### Linux Binary Kullanımı (Derlenmiş Dosya)

Kodları tekrar derlemekle uğraşmadan, repo içerisinde gelen hazır Linux çalıştırılabilir dosyasını (binary) kullanmak için şu komutları uygulayın:

```bash
chmod +x thor-scraper
./thor-scraper

```

## 📂 Çıktı Yapısı (Output)

Program çalıştığında `output/` klasörü altında şu yapıyı oluşturur:

```text
output/
├── example.onion/
│   ├── index.html        # Sitenin HTML kaynak kodu
│   ├── screenshot.png    # Tam sayfa ekran görüntüsü
│   └── site_info.txt     # Meta veriler (Başlık, Tarama Tarihi)
└── scan_report.log       # Detaylı durum raporu (SUCCESS/FAIL kayıtları)

```

<p align="center">
SiberVatan CTI çalışması kapsamında <a href="https://www.google.com/search?q=https://github.com/Muzoovy4606">Muzoovy</a> tarafından geliştirilmiştir.
</p>

```

```
