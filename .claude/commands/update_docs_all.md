# Proje Dokümantasyon Denetimi ve Yenileme

## Görev
Gorev projesinin tüm dokümantasyonunu kapsamlı şekilde gözden geçir, tutarlılığı sağla ve gerekli iyileştirmeleri yap.

## İşlem Adımları

1. **Dokümantasyon Envanteri**
   - CLAUDE.md'nin güncel ve eksiksiz olduğunu kontrol et
   - /docs klasöründeki tüm Türkçe dökümanları listele
   - README.md ve diğer root seviye dökümanları kontrol et
   - .claude/commands/ altındaki komut dökümanlarını gözden geçir
   - Her dökümanın CLAUDE.md ile tutarlı olduğunu doğrula
   - Eksik dokümantasyon alanlarını tespit et

2. **Tutarlılık ve Güncellik Kontrolü**
   - **Version Kontrolü**: Tüm dokümanlarda version bilgisi tutarlılığı
   - **Go Version**: go.mod'da belirtilen Go 1.22 requirement'ı
   - **MCP SDK**: mark3labs/mcp-go v0.6.0 versiyonu
   - **Tool Listesi**: 5 MCP tool'unun tamamının dokümante edildiğini doğrula
   - **Binary İsimleri**: gorev, gorev.exe tutarlılığı

3. **README.md Özel Kontrolleri**
   - Kurulum talimatlarının güncel ve çalışır durumda olduğunu doğrula
   - Docker komutlarının doğru olduğunu kontrol et
   - Claude Desktop/Code konfigürasyonlarının doğru olduğunu doğrula
   - GitHub URL'lerinin placeholder yerine gerçek değerler içermesi gerektiğini not et

4. **Docs Klasörü Kontrolleri**
   - docs/README.md - Ana dokümantasyon index'i
   - docs/kurulum.md - Detaylı kurulum rehberi
   - docs/kullanim.md - Kullanım kılavuzu
   - docs/mcp-araclari.md - MCP tool referansı
   - docs/mimari.md - Sistem mimarisi
   - Eksik dosyalar: api-referans.md, gelistirme.md, ornekler.md

5. **Çapraz Referans Doğrulama**
   - İlgili Dökümanlar bölümlerindeki tüm linkleri kontrol et
   - docs/ içindeki dosyaların birbirine olan referanslarını doğrula
   - Henüz yazılmamış dokümanlara olan linkleri işaretle

6. **İyileştirme Önerileri**
   - API referans dokümantasyonu oluşturulmalı
   - Geliştirici rehberi tamamlanmalı
   - Örnek kullanım senaryoları eklenmeli
   - CHANGELOG.md oluşturulmalı
   - LICENSE dosyası eklenmeli (MIT olarak belirtilmiş)

7. **Standardizasyon**
   - Türkçe terimler: görev, proje, öncelik (düşük/orta/yüksek)
   - Durum terimleri: beklemede, devam_ediyor, tamamlandı
   - Başlık hiyerarşisi: # > ## > ### tutarlılığı
   - Code block formatları: ```go, ```bash, ```json tutarlılığı
   - GitHub URL'leri: yourusername placeholder'ları gerçek değerlerle değiştirilmeli

## Özel Dikkat Edilecek Alanlar
- **MCP SDK Entegrasyonu**: Yeni tamamlandı, dokümantasyona yansıtılmalı
- **Arsiv-kotlin**: Eski kodların arşivlendiği belirtilmeli
- **Docker Build**: Multi-stage build process dokümante edilmeli
- **Test Coverage**: Integration test'lerin varlığı belirtilmeli
- **Binary Distribution**: GitHub releases üzerinden dağıtım planı

## Çıktı Formatı

1. **Durum Raporu**
   - Dokümantasyon tamlık skoru (1-10)
   - Tespit edilen tutarsızlıklar
   - Eksik dokümantasyon alanları

2. **Güncelleme Listesi**
   - Her dosya için yapılan güncellemeler
   - Eklenen yeni bölümler
   - Düzeltilen tutarsızlıklar

3. **Eylem Planı**
   - Öncelikli güncelleme gerektiren alanlar
   - Yeni oluşturulması gereken dökümanlar
   - GitHub repository setup önerileri