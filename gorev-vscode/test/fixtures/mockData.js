module.exports = {
  mockTasks: [
    {
      id: '123e4567-e89b-12d3-a456-426614174000',
      baslik: 'Test Task 1',
      aciklama: 'Test task description',
      durum: 'beklemede',
      oncelik: 'orta',
      proje_id: 'proj-123',
      son_tarih: '2025-07-01',
      etiketler: ['test', 'mock'],
      olusturma_tarih: '2025-06-28T10:00:00Z',
      guncelleme_tarih: '2025-06-28T10:00:00Z'
    },
    {
      id: '123e4567-e89b-12d3-a456-426614174001',
      baslik: 'Test Task 2',
      aciklama: 'Another test task',
      durum: 'devam_ediyor',
      oncelik: 'yuksek',
      proje_id: 'proj-123',
      etiketler: ['urgent'],
      olusturma_tarih: '2025-06-28T11:00:00Z',
      guncelleme_tarih: '2025-06-28T11:00:00Z'
    },
    {
      id: '123e4567-e89b-12d3-a456-426614174002',
      baslik: 'Completed Task',
      aciklama: 'This task is done',
      durum: 'tamamlandi',
      oncelik: 'dusuk',
      proje_id: 'proj-456',
      olusturma_tarih: '2025-06-28T09:00:00Z',
      guncelleme_tarih: '2025-06-28T12:00:00Z'
    }
  ],

  mockProjects: [
    {
      id: 'proj-123',
      isim: 'Test Project 1',
      tanim: 'First test project',
      olusturma_tarih: '2025-06-01T10:00:00Z',
      guncelleme_tarih: '2025-06-01T10:00:00Z'
    },
    {
      id: 'proj-456',
      isim: 'Test Project 2',
      tanim: 'Second test project',
      olusturma_tarih: '2025-06-01T11:00:00Z',
      guncelleme_tarih: '2025-06-01T11:00:00Z'
    }
  ],

  mockTemplates: [
    {
      id: 'template-123',
      isim: 'Bug Report',
      tanim: 'Template for bug reports',
      varsayilan_baslik: '🐛 [{{module}}] {{title}}',
      aciklama_template: '## Bug Description\\n{{description}}\\n\\n## Steps to Reproduce\\n{{steps}}',
      alanlar: [
        {
          isim: 'title',
          tur: 'metin',
          zorunlu: true
        },
        {
          isim: 'module',
          tur: 'metin',
          zorunlu: true
        },
        {
          isim: 'priority',
          tur: 'secim',
          zorunlu: true,
          varsayilan: 'orta',
          secenekler: ['dusuk', 'orta', 'yuksek']
        }
      ],
      kategori: 'Teknik',
      aktif: true
    },
    {
      id: 'template-456',
      isim: 'Feature Request',
      tanim: 'Template for feature requests',
      varsayilan_baslik: '✨ {{title}}',
      aciklama_template: '## Feature Description\\n{{description}}\\n\\n## User Story\\n{{userStory}}',
      alanlar: [
        {
          isim: 'title',
          tur: 'metin',
          zorunlu: true
        },
        {
          isim: 'description',
          tur: 'metin',
          zorunlu: true
        },
        {
          isim: 'effort',
          tur: 'secim',
          secenekler: ['small', 'medium', 'large']
        }
      ],
      kategori: 'Özellik',
      aktif: true
    }
  ],

  mockDependencies: [
    {
      kaynak_id: '123e4567-e89b-12d3-a456-426614174001',
      hedef_id: '123e4567-e89b-12d3-a456-426614174000',
      hedef_baslik: 'Test Task 1',
      hedef_durum: 'beklemede',
      baglanti_tip: 'engelliyor'
    }
  ],

  mockSummary: {
    toplamGorev: 25,
    tamamlanan: 10,
    devamEden: 5,
    bekleyen: 10,
    toplamProje: 3,
    aktifProje: 'Test Project 1'
  },

  // Mock MCP responses
  mockMCPResponses: {
    gorev_listele: {
      content: [{
        type: 'text',
        text: `## 📋 Görev Listesi

### Bekleyen Görevler (2)

- [beklemede] Test Task 1 (orta öncelik)
  ID: 123e4567-e89b-12d3-a456-426614174000
  Proje: Test Project 1
  Son tarih: 2025-07-01
  Etiketler: test, mock
  Test task description

### Devam Eden Görevler (1)

- [devam_ediyor] Test Task 2 (yuksek öncelik)
  ID: 123e4567-e89b-12d3-a456-426614174001
  Proje: Test Project 1
  Etiketler: urgent
  Another test task`
      }]
    },

    proje_listele: {
      content: [{
        type: 'text',
        text: `## 📁 Proje Listesi

### 🔒 Test Project 1
**ID:** proj-123
**Tanım:** First test project
**Görev Sayısı:** Toplam: 2, Tamamlanan: 0, Devam Eden: 1, Bekleyen: 1

### 📁 Test Project 2
**ID:** proj-456
**Tanım:** Second test project
**Görev Sayısı:** Toplam: 1, Tamamlanan: 1, Devam Eden: 0, Bekleyen: 0`
      }]
    },

    template_listele: {
      content: [{
        type: 'text',
        text: `## 📋 Görev Template'leri

### Teknik

#### Bug Report
- **ID:** \`template-123\`
- **Açıklama:** Template for bug reports
- **Başlık Şablonu:** \`🐛 [{{module}}] {{title}}\`
- **Alanlar:**
  - \`title\` (text) *(zorunlu)*
  - \`module\` (text) *(zorunlu)*
  - \`priority\` (select) *(zorunlu)* - varsayılan: orta - seçenekler: dusuk, orta, yuksek

### Özellik

#### Feature Request
- **ID:** \`template-456\`
- **Açıklama:** Template for feature requests
- **Başlık Şablonu:** \`✨ {{title}}\`
- **Alanlar:**
  - \`title\` (text) *(zorunlu)*
  - \`description\` (text) *(zorunlu)*
  - \`effort\` (select) - seçenekler: small, medium, large`
      }]
    },

    ozet_goster: {
      content: [{
        type: 'text',
        text: `## 📊 Özet Bilgiler

Toplam görev sayısı: 25
Tamamlanan: 10
Devam eden: 5
Bekleyen: 10

Toplam proje sayısı: 3
Aktif proje: Test Project 1

### 📈 İstatistikler
- Tamamlanma oranı: %40
- Ortalama görev süresi: 3.5 gün
- Bu hafta tamamlanan: 7`
      }]
    }
  }
};