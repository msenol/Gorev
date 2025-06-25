# Haftalık Bakım Görevleri

## Bakım Kontrol Listesi

1. **Bağımlılık Yönetimi**
   - Go modül güncellemelerini kontrol et (`go list -u -m all`)
   - Güvenlik güncellemelerini incele
   - Patch versiyonları güncelle (`go get -u`)
   - CLAUDE.md'deki komutları kullan: `make deps`

2. **Kod Kalitesi**
   - CLAUDE.md'deki Development Commands kullan:
   - Tüm testleri çalıştır (`make test`)
   - Test coverage kontrol et (`make test-coverage`)
   - Race condition testi (`go test -race ./...`)
   - Format kontrolü (`make fmt`)
   - Lint kontrolü (`make lint`)
   - go vet analizi (`go vet ./...`)
   - TODO/FIXME yorumlarını gözden geçir

3. **Dokümantasyon**
   - CLAUDE.md'nin güncel olduğunu kontrol et
   - CHANGELOG.md güncelle
   - README.md doğruluğunu kontrol et
   - docs/ altındaki dökümanları gözden geçir
   - Eksik dökümanları tamamla (CLAUDE.md'deki "Important Files" bölümüne göre)
   - GitHub URL placeholder'larını güncelle

4. **Build ve Dağıtım**
   - Binary build test et (`make build`)
   - Docker image build et (`make docker-build`)
   - Cross-platform build kontrol et (`make build-all`)
   - GitHub releases hazırlığı

5. **Güvenlik ve Temizlik**
   - Geçici dosyaları temizle (`make clean`)
   - .gitignore kontrolü
   - Hassas bilgi taraması
   - Binary dosyaların commit edilmediğini doğrula

6. **Performans**
   - SQLite veritabanı optimizasyonu
   - Benchmark testleri çalıştır
   - Memory profiling (pprof)
   - Goroutine leak kontrolü

## Haftalık Rapor Formatı

```markdown
## Haftalık Bakım Raporu - [Tarih]

### ✅ Tamamlanan Görevler
- [ ] Bağımlılık güncellemeleri
- [ ] Test suite başarılı
- [ ] Dokümantasyon güncel
- [ ] Build kontrolleri

### 🔄 Güncellenen Bağımlılıklar
- package_name: v1.2.3 → v1.2.4

### 🐛 Tespit Edilen Sorunlar
- [Sorun açıklaması ve çözüm planı]

### 📈 Metrikler
- Test Coverage: %XX
- Build Süresi: XX saniye
- Binary Boyutu: XX MB

### 📋 Gelecek Hafta Planı
- [Planlanan iyileştirmeler]
```