# Son Değişikliklerin Dokümantasyonu ve Güncelleme

## Görev
Projede yapılan son değişiklikleri analiz et ve ilgili dokümantasyonu güncelle.

## İşlem Adımları

1. **Değişiklik Tespiti**
   - Son değişiklikleri git status ve git diff ile incele
   - Her değişikliğin türünü belirle (feature, fix, refactor, docs)
   - Etkilenen Go paketleri ve dosyaları listele
   - Breaking changes varsa özellikle vurgula

2. **Etki Analizi**
   - Değişikliklerin mevcut dokümanlara etkisini değerlendir
   - README.md'de güncellenmesi gereken bölümleri belirle
   - MCP tool değişikliklerinin docs/mcp-araclari.md'ye yansıması gerekiyorsa not al
   - Version numarasının güncellenmesi gerekip gerekmediğini değerlendir (Semantic Versioning)

3. **Dokümantasyon Güncellemeleri**
   - CLAUDE.md'yi gözden geçir ve gerekirse güncelle
   - README.md'deki version ve özellik bilgilerini güncelle
   - docs/ altındaki ilgili dokümantasyonu güncelle
   - CHANGELOG.md yoksa oluştur ve girişler ekle
   - Version bilgisi Makefile'da LDFLAGS ile yönetiliyor

4. **Referans Kontrolü**
   - Tüm dokümanlardaki çapraz referansların güncel olduğunu doğrula
   - Yeni eklenen MCP tool'ları için docs/mcp-araclari.md'ye ekle
   - Kaldırılan özellikler varsa ilgili referansları temizle

5. **Kalite Kontrolü**
   - Version numaralarının tutarlı olduğunu kontrol et
   - Go package import'larının doğru olduğunu kontrol et
   - TODO veya FIXME notlarının güncel olduğunu kontrol et

## Çıktı Formatı

### 📊 Değişiklik Özeti
```
Kategori: [Feature/Fix/Refactor/Docs]
Etkilenen Paketler:
- internal/mcp: [değişiklik açıklaması]
- internal/gorev: [değişiklik açıklaması]
- docs/: [değişiklik açıklaması]

Breaking Changes: [Var/Yok]
```

### 📝 Güncellenen Dosyalar
```
✅ README.md
   - Version: X.Y.Z
   - Özellikler güncellendi
   
✅ docs/mcp-araclari.md
   - Yeni tool eklendi: [tool adı]
   
✅ CHANGELOG.md
   - Yeni release notları eklendi
```

### ⚠️ Dikkat Edilecekler
```
- [ ] go.mod version sync kontrolü
- [ ] Docker image version tag'i
- [ ] GitHub release hazırlığı
```