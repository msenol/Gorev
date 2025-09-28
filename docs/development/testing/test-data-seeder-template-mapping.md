# TestDataSeeder Template Mapping Strategy

## Overview

The TestDataSeeder currently uses the deprecated `gorev_olustur` command directly. We need to convert it to use the `templateden_gorev_olustur` command with appropriate templates.

## Available Templates

### Current Templates in Database

1. **Bug Raporu** (ID: 4dd56a2a-caf4-472c-8c0f-276bc8a1f880)
   - For bug reports and issues
   - Fields: baslik, aciklama, modul, ortam, adimlar, beklenen, mevcut, ekler, cozum, oncelik, etiketler

2. **Özellik İsteği** (ID: 6b083358-9c4d-4f4e-b041-9288c05a1bb7)
   - For feature requests
   - Fields: baslik, aciklama, amac, kullanicilar, kriterler, ui_ux, ilgili, efor, oncelik, etiketler

3. **Teknik Borç** (ID: 69e2b237-7c2e-4459-9d46-ea6c05aba39a)
   - For technical debt and refactoring
   - Fields: baslik, aciklama, alan, dosyalar, neden, analiz, cozum, riskler, iyilestirmeler, sure, oncelik, etiketler

4. **Araştırma Görevi** (ID: 13f04fe2-b5b6-4fd6-8684-5eca5dc2770d)
   - For research tasks
   - Fields: konu, amac, sorular, kaynaklar, alternatifler, kriterler, son_tarih, oncelik, etiketler

### New Templates (Need to be added to DB)

5. **Bug Raporu v2** - Enhanced bug report
6. **Spike Araştırma** - Time-boxed research
7. **Performans Sorunu** - Performance issues
8. **Güvenlik Düzeltmesi** - Security fixes
9. **Refactoring** - Code quality improvements

## Task Mapping Strategy

### 1. Bug/Issue Tasks → Bug Raporu Template

Tasks with tags like "bug", "critical", "urgent" or titles containing "hatası", "404", "error":

- "Login sayfası 404 hatası veriyor" → Bug Raporu
- "SSL sertifikası yenile" → Bug Raporu (infrastructure issue)

### 2. Feature Tasks → Özellik İsteği Template

Tasks with tags like "feature", "enhancement" or titles containing "implement", "ekle", "sistemi":

- "Kullanıcı giriş sistemi implement et" → Özellik İsteği
- "Push notification sistemi" → Özellik İsteği
- "Dark mode tema" → Özellik İsteği
- "Contact form entegrasyonu" → Özellik İsteği

### 3. Technical/Backend Tasks → Teknik Borç Template

Tasks with tags like "backend", "refactoring", "infrastructure", "performance":

- "Redis cache entegrasyonu" → Teknik Borç
- "Rate limiting implement et" → Teknik Borç
- "API dokümantasyonu güncelle" → Teknik Borç
- "ETL pipeline kurulumu" → Teknik Borç

### 4. Research/Analysis Tasks → Araştırma Görevi Template

Tasks with tags like "research", "araştırma", "analytics" or titles containing "araştırma", "analizi":

- "Chart library araştırması" → Araştırma Görevi
- "Makine öğrenmesi modeli" → Araştırma Görevi
- "SEO optimizasyonu" → Araştırma Görevi (can be research)

### 5. UI/Design Tasks → Özellik İsteği or Teknik Borç

Tasks with tags like "design", "ui", "frontend":

- "Ana sayfa tasarımını tamamla" → Özellik İsteği
- "Dashboard prototype hazırla" → Özellik İsteği
- "Responsive grid sistemi" → Teknik Borç

### 6. Security Tasks → Bug Raporu (until Güvenlik Düzeltmesi is available)

Tasks with tags like "security" or titles containing "güvenlik", "penetrasyon":

- "Penetrasyon testi yap" → Bug Raporu (security category)
- "2FA implementasyonu" → Özellik İsteği (security feature)

### 7. General/Meeting Tasks → Create without template

Tasks like meetings, reviews, documentation:

- "Team meeting hazırlığı" → Use gorev_olustur directly
- "Code review yapılacak PR'lar" → Use gorev_olustur directly

## Implementation Steps

1. **Add helper function to determine template**:

   ```typescript
   private getTemplateForTask(task: any): { templateId: string; fields: any } | null
   ```

2. **Map task fields to template fields**:
   - baslik → baslik (most templates)
   - aciklama → aciklama (most templates)
   - oncelik → oncelik (all templates)
   - etiketler → etiketler (all templates)
   - son_tarih → son_tarih (some templates)
   - proje_id → Not part of template, handled separately

3. **Handle special cases**:
   - Tasks without matching template → fallback to direct creation
   - Subtasks → use gorev_altgorev_olustur (no template needed)
   - Tasks with dependencies → create task first, then add dependencies

## Sample Conversions

### Example 1: Bug Task

```typescript
// Old way
await this.mcpClient.callTool('gorev_olustur', {
    baslik: 'Login sayfası 404 hatası veriyor',
    aciklama: 'Production ortamında /login URL\'ine gittiğimizde 404 hatası alıyoruz',
    oncelik: GorevOncelik.Yuksek,
    proje_id: projectIds[0],
    son_tarih: this.getDateString(5),
    etiketler: 'bug,critical,production'
});

// New way
await this.mcpClient.callTool('templateden_gorev_olustur', {
    template_id: '4dd56a2a-caf4-472c-8c0f-276bc8a1f880', // Bug Raporu
    degerler: {
        baslik: 'Login sayfası 404 hatası veriyor',
        aciklama: 'Production ortamında /login URL\'ine gittiğimizde 404 hatası alıyoruz',
        modul: 'Authentication',
        ortam: 'production',
        adimlar: '1. Production URL\'ine git\\n2. /login sayfasına git\\n3. 404 hatası görünüyor',
        beklenen: 'Login sayfası açılmalı',
        mevcut: '404 Not Found hatası',
        oncelik: 'yuksek',
        etiketler: 'bug,critical,production'
    }
});
```

### Example 2: Feature Task

```typescript
// Old way
await this.mcpClient.callTool('gorev_olustur', {
    baslik: 'Dark mode tema',
    aciklama: 'Sistem ayarlarına göre otomatik tema değişimi',
    oncelik: GorevOncelik.Dusuk,
    proje_id: projectIds[1],
    etiketler: 'mobile,ui,enhancement'
});

// New way
await this.mcpClient.callTool('templateden_gorev_olustur', {
    template_id: '6b083358-9c4d-4f4e-b041-9288c05a1bb7', // Özellik İsteği
    degerler: {
        baslik: 'Dark mode tema',
        aciklama: 'Sistem ayarlarına göre otomatik tema değişimi',
        amac: 'Kullanıcı deneyimini iyileştirmek ve göz yorgunluğunu azaltmak',
        kullanicilar: 'Tüm mobil uygulama kullanıcıları',
        kriterler: '1. Sistem temasına uyum\\n2. Manuel toggle\\n3. Tercih kaydetme',
        ui_ux: 'Settings sayfasında toggle switch',
        efor: 'orta',
        oncelik: 'dusuk',
        etiketler: 'mobile,ui,enhancement'
    }
});
```

## Notes

1. **Project Assignment**: Templates don't include proje_id. After creating a task from template, use `gorev_duzenle` to assign to project.
2. **Due Dates**: Only some templates have son_tarih field. For others, use `gorev_duzenle` after creation.
3. **Priority Mapping**: Convert GorevOncelik enum values to lowercase strings (yuksek, orta, dusuk).
4. **Field Defaults**: Provide sensible defaults for required template fields that aren't in original task data.
5. **Validation**: Ensure all required fields for chosen template are provided.
