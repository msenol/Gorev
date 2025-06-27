package gorev

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TemplateOlustur yeni bir görev template'i oluşturur
func (vy *VeriYonetici) TemplateOlustur(template *GorevTemplate) error {
	template.ID = uuid.New().String()

	// Alanları JSON'a çevir
	alanlarJSON, err := json.Marshal(template.Alanlar)
	if err != nil {
		return fmt.Errorf("alanlar JSON'a çevrilemedi: %w", err)
	}

	// Örnek değerleri JSON'a çevir
	ornekDegerlerJSON, err := json.Marshal(template.OrnekDegerler)
	if err != nil {
		return fmt.Errorf("örnek değerler JSON'a çevrilemedi: %w", err)
	}

	sorgu := `INSERT INTO gorev_templateleri 
		(id, isim, tanim, varsayilan_baslik, aciklama_template, alanlar, ornek_degerler, kategori, aktif)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = vy.db.Exec(sorgu, template.ID, template.Isim, template.Tanim,
		template.VarsayilanBaslik, template.AciklamaTemplate,
		string(alanlarJSON), string(ornekDegerlerJSON), template.Kategori, template.Aktif)

	if err != nil {
		return fmt.Errorf("template oluşturulamadı: %w", err)
	}

	return nil
}

// TemplateListele tüm aktif template'leri listeler
func (vy *VeriYonetici) TemplateListele(kategori string) ([]*GorevTemplate, error) {
	var sorgu string
	var args []interface{}

	if kategori != "" {
		sorgu = `SELECT id, isim, tanim, varsayilan_baslik, aciklama_template, 
				alanlar, ornek_degerler, kategori, aktif 
				FROM gorev_templateleri WHERE aktif = 1 AND kategori = ? ORDER BY isim`
		args = append(args, kategori)
	} else {
		sorgu = `SELECT id, isim, tanim, varsayilan_baslik, aciklama_template, 
				alanlar, ornek_degerler, kategori, aktif 
				FROM gorev_templateleri WHERE aktif = 1 ORDER BY kategori, isim`
	}

	rows, err := vy.db.Query(sorgu, args...)
	if err != nil {
		return nil, fmt.Errorf("template'ler getirilemedi: %w", err)
	}
	defer rows.Close()

	var templates []*GorevTemplate
	for rows.Next() {
		template := &GorevTemplate{}
		var alanlarJSON, ornekDegerlerJSON string

		err := rows.Scan(&template.ID, &template.Isim, &template.Tanim,
			&template.VarsayilanBaslik, &template.AciklamaTemplate,
			&alanlarJSON, &ornekDegerlerJSON, &template.Kategori, &template.Aktif)
		if err != nil {
			return nil, fmt.Errorf("template okunamadı: %w", err)
		}

		// Alanları parse et
		if err := json.Unmarshal([]byte(alanlarJSON), &template.Alanlar); err != nil {
			return nil, fmt.Errorf("alanlar parse edilemedi: %w", err)
		}

		// Örnek değerleri parse et
		if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.OrnekDegerler); err != nil {
			return nil, fmt.Errorf("örnek değerler parse edilemedi: %w", err)
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// TemplateGetir belirli bir template'i getirir
func (vy *VeriYonetici) TemplateGetir(templateID string) (*GorevTemplate, error) {
	template := &GorevTemplate{}
	var alanlarJSON, ornekDegerlerJSON string

	sorgu := `SELECT id, isim, tanim, varsayilan_baslik, aciklama_template, 
			alanlar, ornek_degerler, kategori, aktif 
			FROM gorev_templateleri WHERE id = ?`

	err := vy.db.QueryRow(sorgu, templateID).Scan(
		&template.ID, &template.Isim, &template.Tanim,
		&template.VarsayilanBaslik, &template.AciklamaTemplate,
		&alanlarJSON, &ornekDegerlerJSON, &template.Kategori, &template.Aktif)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template bulunamadı: %s", templateID)
		}
		return nil, fmt.Errorf("template getirilemedi: %w", err)
	}

	// Alanları parse et
	if err := json.Unmarshal([]byte(alanlarJSON), &template.Alanlar); err != nil {
		return nil, fmt.Errorf("alanlar parse edilemedi: %w", err)
	}

	// Örnek değerleri parse et
	if err := json.Unmarshal([]byte(ornekDegerlerJSON), &template.OrnekDegerler); err != nil {
		return nil, fmt.Errorf("örnek değerler parse edilemedi: %w", err)
	}

	return template, nil
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (vy *VeriYonetici) TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error) {
	// Template'i getir
	template, err := vy.TemplateGetir(templateID)
	if err != nil {
		return nil, err
	}

	// Zorunlu alanları kontrol et
	for _, alan := range template.Alanlar {
		if alan.Zorunlu {
			if _, ok := degerler[alan.Isim]; !ok {
				return nil, fmt.Errorf("zorunlu alan eksik: %s", alan.Isim)
			}
		}
	}

	// Başlık oluştur
	baslik := template.VarsayilanBaslik
	for key, value := range degerler {
		baslik = strings.ReplaceAll(baslik, "{{"+key+"}}", value)
	}

	// Açıklama oluştur
	aciklama := template.AciklamaTemplate
	for key, value := range degerler {
		aciklama = strings.ReplaceAll(aciklama, "{{"+key+"}}", value)
	}

	// Varsayılan değerleri uygula
	oncelik := "orta"
	if val, ok := degerler["oncelik"]; ok {
		oncelik = val
	}

	var sonTarih *time.Time
	if val, ok := degerler["son_tarih"]; ok {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			sonTarih = &t
		}
	}

	// Etiketleri ayır
	var etiketler []string
	if val, ok := degerler["etiketler"]; ok {
		etiketler = strings.Split(val, ",")
		for i := range etiketler {
			etiketler[i] = strings.TrimSpace(etiketler[i])
		}
	}

	// Görev oluştur
	gorev := &Gorev{
		Baslik:   baslik,
		Aciklama: aciklama,
		Oncelik:  oncelik,
		Durum:    "beklemede",
	}

	// ProjeID'yi ayarla
	if val, ok := degerler["proje_id"]; ok {
		gorev.ProjeID = val
	} else {
		// Aktif projeyi kullan
		aktifProjeID, _ := vy.AktifProjeGetir()
		gorev.ProjeID = aktifProjeID
	}

	// ID ve tarihler ayarla
	gorev.ID = uuid.New().String()
	gorev.OlusturmaTarih = time.Now()
	gorev.GuncellemeTarih = time.Now()
	gorev.SonTarih = sonTarih

	// Görevi kaydet
	err = vy.GorevKaydet(gorev)
	if err != nil {
		return nil, fmt.Errorf("görev kaydedilemedi: %w", err)
	}

	// Etiketleri ayarla
	if len(etiketler) > 0 {
		etiketNesneleri, err := vy.EtiketleriGetirVeyaOlustur(etiketler)
		if err != nil {
			return nil, fmt.Errorf("etiketler oluşturulamadı: %w", err)
		}

		err = vy.GorevEtiketleriniAyarla(gorev.ID, etiketNesneleri)
		if err != nil {
			return nil, fmt.Errorf("görev etiketleri ayarlanamadı: %w", err)
		}

		gorev.Etiketler = etiketNesneleri
	}

	return gorev, nil
}

// VarsayilanTemplateleriOlustur varsayılan template'leri oluşturur
func (vy *VeriYonetici) VarsayilanTemplateleriOlustur() error {
	templates := []*GorevTemplate{
		{
			Isim:             "Bug Raporu",
			Tanim:            "Yazılım hatası bildirimi için detaylı template",
			VarsayilanBaslik: "🐛 [{{modul}}] {{baslik}}",
			AciklamaTemplate: `## 🐛 Hata Açıklaması
{{aciklama}}

## 📍 Nerede Oluşuyor?
**Modül/Bileşen:** {{modul}}
**Ortam:** {{ortam}}

## 🔄 Tekrar Üretme Adımları
{{adimlar}}

## ✅ Beklenen Davranış
{{beklenen}}

## ❌ Mevcut Davranış
{{mevcut}}

## 📸 Ekran Görüntüleri/Loglar
{{ekler}}

## 🔧 Olası Çözüm
{{cozum}}

## 📊 Öncelik: {{oncelik}}
## 🏷️ Etiketler: {{etiketler}}`,
			Alanlar: []TemplateAlan{
				{Isim: "baslik", Tip: "text", Zorunlu: true},
				{Isim: "aciklama", Tip: "text", Zorunlu: true},
				{Isim: "modul", Tip: "text", Zorunlu: true},
				{Isim: "ortam", Tip: "select", Zorunlu: true, Secenekler: []string{"development", "staging", "production"}},
				{Isim: "adimlar", Tip: "text", Zorunlu: true},
				{Isim: "beklenen", Tip: "text", Zorunlu: true},
				{Isim: "mevcut", Tip: "text", Zorunlu: true},
				{Isim: "ekler", Tip: "text", Zorunlu: false},
				{Isim: "cozum", Tip: "text", Zorunlu: false},
				{Isim: "oncelik", Tip: "select", Zorunlu: true, Varsayilan: "orta", Secenekler: []string{"dusuk", "orta", "yuksek"}},
				{Isim: "etiketler", Tip: "text", Zorunlu: false, Varsayilan: "bug"},
			},
			Kategori: "Teknik",
			Aktif:    true,
		},
		{
			Isim:             "Özellik İsteği",
			Tanim:            "Yeni özellik veya geliştirme isteği için template",
			VarsayilanBaslik: "✨ {{baslik}}",
			AciklamaTemplate: `## ✨ Özellik Açıklaması
{{aciklama}}

## 🎯 Amaç ve Faydalar
{{amac}}

## 👥 Hedef Kullanıcılar
{{kullanicilar}}

## 📋 Kabul Kriterleri
{{kriterler}}

## 🎨 UI/UX Düşünceleri
{{ui_ux}}

## 🔗 İlgili Özellikler/Modüller
{{ilgili}}

## 📊 Tahmini Efor
{{efor}}

## 🏷️ Etiketler: {{etiketler}}`,
			Alanlar: []TemplateAlan{
				{Isim: "baslik", Tip: "text", Zorunlu: true},
				{Isim: "aciklama", Tip: "text", Zorunlu: true},
				{Isim: "amac", Tip: "text", Zorunlu: true},
				{Isim: "kullanicilar", Tip: "text", Zorunlu: true},
				{Isim: "kriterler", Tip: "text", Zorunlu: true},
				{Isim: "ui_ux", Tip: "text", Zorunlu: false},
				{Isim: "ilgili", Tip: "text", Zorunlu: false},
				{Isim: "efor", Tip: "select", Zorunlu: false, Secenekler: []string{"küçük", "orta", "büyük"}},
				{Isim: "oncelik", Tip: "select", Zorunlu: true, Varsayilan: "orta", Secenekler: []string{"dusuk", "orta", "yuksek"}},
				{Isim: "etiketler", Tip: "text", Zorunlu: false, Varsayilan: "özellik"},
			},
			Kategori: "Özellik",
			Aktif:    true,
		},
		{
			Isim:             "Teknik Borç",
			Tanim:            "Refaktöring veya teknik iyileştirme için template",
			VarsayilanBaslik: "🔧 [{{alan}}] {{baslik}}",
			AciklamaTemplate: `## 🔧 Teknik Borç Açıklaması
{{aciklama}}

## 📍 Etkilenen Alan
**Alan/Modül:** {{alan}}
**Dosyalar:** {{dosyalar}}

## ❓ Neden Gerekli?
{{neden}}

## 📊 Mevcut Durum Analizi
{{analiz}}

## 🎯 Önerilen Çözüm
{{cozum}}

## ⚠️ Riskler
{{riskler}}

## 📈 Beklenen İyileştirmeler
{{iyilestirmeler}}

## ⏱️ Tahmini Süre: {{sure}}
## 🏷️ Etiketler: {{etiketler}}`,
			Alanlar: []TemplateAlan{
				{Isim: "baslik", Tip: "text", Zorunlu: true},
				{Isim: "aciklama", Tip: "text", Zorunlu: true},
				{Isim: "alan", Tip: "text", Zorunlu: true},
				{Isim: "dosyalar", Tip: "text", Zorunlu: false},
				{Isim: "neden", Tip: "text", Zorunlu: true},
				{Isim: "analiz", Tip: "text", Zorunlu: true},
				{Isim: "cozum", Tip: "text", Zorunlu: true},
				{Isim: "riskler", Tip: "text", Zorunlu: false},
				{Isim: "iyilestirmeler", Tip: "text", Zorunlu: true},
				{Isim: "sure", Tip: "select", Zorunlu: false, Secenekler: []string{"1 gün", "2-3 gün", "1 hafta", "2+ hafta"}},
				{Isim: "oncelik", Tip: "select", Zorunlu: true, Varsayilan: "orta", Secenekler: []string{"dusuk", "orta", "yuksek"}},
				{Isim: "etiketler", Tip: "text", Zorunlu: false, Varsayilan: "teknik-borç,refaktöring"},
			},
			Kategori: "Teknik",
			Aktif:    true,
		},
		{
			Isim:             "Araştırma Görevi",
			Tanim:            "Teknoloji veya çözüm araştırması için template",
			VarsayilanBaslik: "🔍 {{konu}} Araştırması",
			AciklamaTemplate: `## 🔍 Araştırma Konusu
{{konu}}

## 🎯 Araştırma Amacı
{{amac}}

## ❓ Cevaplanması Gereken Sorular
{{sorular}}

## 📚 Araştırılacak Kaynaklar
{{kaynaklar}}

## 🔄 Alternatifler
{{alternatifler}}

## ⚖️ Değerlendirme Kriterleri
{{kriterler}}

## 📅 Bitiş Tarihi: {{son_tarih}}
## 🏷️ Etiketler: {{etiketler}}`,
			Alanlar: []TemplateAlan{
				{Isim: "konu", Tip: "text", Zorunlu: true},
				{Isim: "amac", Tip: "text", Zorunlu: true},
				{Isim: "sorular", Tip: "text", Zorunlu: true},
				{Isim: "kaynaklar", Tip: "text", Zorunlu: false},
				{Isim: "alternatifler", Tip: "text", Zorunlu: false},
				{Isim: "kriterler", Tip: "text", Zorunlu: true},
				{Isim: "son_tarih", Tip: "date", Zorunlu: false},
				{Isim: "oncelik", Tip: "select", Zorunlu: true, Varsayilan: "orta", Secenekler: []string{"dusuk", "orta", "yuksek"}},
				{Isim: "etiketler", Tip: "text", Zorunlu: false, Varsayilan: "araştırma"},
			},
			Kategori: "Araştırma",
			Aktif:    true,
		},
	}

	for _, template := range templates {
		if err := vy.TemplateOlustur(template); err != nil {
			// Template zaten varsa hata verme
			if !strings.Contains(err.Error(), "UNIQUE constraint") {
				return fmt.Errorf("varsayılan template oluşturulamadı (%s): %w", template.Isim, err)
			}
		}
	}

	return nil
}
