const assert = require('assert');
const { MarkdownParser } = require('../../dist/utils/markdownParser');

suite('MarkdownParser Test Suite', () => {
  
  suite('parseGorevListesi', () => {
    test('should parse simple task list', () => {
      const markdown = `## 📋 Görev Listesi

- [beklemede] Test görevi (orta öncelik)
  ID: 123e4567-e89b-12d3-a456-426614174000
  Proje: Test Projesi
  Test açıklaması`;

      const tasks = MarkdownParser.parseGorevListesi(markdown);
      
      assert.strictEqual(tasks.length, 1);
      assert.strictEqual(tasks[0].baslik, 'Test görevi');
      assert.strictEqual(tasks[0].id, '123e4567-e89b-12d3-a456-426614174000');
      assert.strictEqual(tasks[0].durum, 'beklemede');
      assert.strictEqual(tasks[0].oncelik, 'orta');
      assert.strictEqual(tasks[0].aciklama, 'Test açıklaması');
    });

    test('should parse tasks with tags and due dates', () => {
      const markdown = `- [devam_ediyor] Urgent task (yuksek öncelik)
  ID: 123e4567-e89b-12d3-a456-426614174001
  Son tarih: 2025-07-01
  Etiketler: bug, urgent`;

      const tasks = MarkdownParser.parseGorevListesi(markdown);
      
      assert.strictEqual(tasks.length, 1);
      assert.strictEqual(tasks[0].son_tarih, '2025-07-01');
      assert.deepStrictEqual(tasks[0].etiketler, ['bug', 'urgent']);
    });

    test('should parse multiple tasks', () => {
      const markdown = `- [tamamlandi] First task (dusuk öncelik)
  ID: task-1
  
- [beklemede] Second task (orta öncelik)
  ID: task-2
  Description line 1
  Description line 2`;

      const tasks = MarkdownParser.parseGorevListesi(markdown);
      
      assert.strictEqual(tasks.length, 2);
      assert.strictEqual(tasks[0].baslik, 'First task');
      assert.strictEqual(tasks[1].baslik, 'Second task');
      assert.strictEqual(tasks[1].aciklama, 'Description line 1 Description line 2');
    });
  });

  suite('parseProjeListesi', () => {
    test('should parse project list', () => {
      const markdown = `### 🔒 Test Project
**ID:** proj-123
**Tanım:** Test project description

### 📁 Another Project
**ID:** proj-456
**Tanım:** Another description`;

      const projects = MarkdownParser.parseProjeListesi(markdown);
      
      assert.strictEqual(projects.length, 2);
      assert.strictEqual(projects[0].isim, 'Test Project');
      assert.strictEqual(projects[0].id, 'proj-123');
      assert.strictEqual(projects[0].tanim, 'Test project description');
      assert.strictEqual(projects[1].isim, 'Another Project');
    });

    test('should handle projects without emoji', () => {
      const markdown = `### Simple Project
**ID:** proj-789
**Tanım:** Simple description`;

      const projects = MarkdownParser.parseProjeListesi(markdown);
      
      assert.strictEqual(projects.length, 1);
      assert.strictEqual(projects[0].isim, 'Simple Project');
    });
  });

  suite('parseTemplateListesi', () => {
    test('should parse template list', () => {
      const markdown = `### Araştırma

#### Test Template
- **ID:** \`template-123\`
- **Açıklama:** Test template description
- **Başlık Şablonu:** \`Test {{field}}\`
- **Alanlar:**
  - \`field1\` (text) *(zorunlu)*
  - \`field2\` (select) - varsayılan: opt1 - seçenekler: opt1, opt2, opt3`;

      const templates = MarkdownParser.parseTemplateListesi(markdown);
      
      assert.strictEqual(templates.length, 1);
      assert.strictEqual(templates[0].isim, 'Test Template');
      assert.strictEqual(templates[0].id, 'template-123');
      assert.strictEqual(templates[0].kategori, 'Araştırma');
      assert.strictEqual(templates[0].tanim, 'Test template description');
      assert.strictEqual(templates[0].varsayilan_baslik, 'Test {{field}}');
      assert.strictEqual(templates[0].alanlar.length, 2);
      
      // Check fields
      assert.strictEqual(templates[0].alanlar[0].isim, 'field1');
      assert.strictEqual(templates[0].alanlar[0].tur, 'metin');
      assert.strictEqual(templates[0].alanlar[0].zorunlu, true);
      
      assert.strictEqual(templates[0].alanlar[1].isim, 'field2');
      assert.strictEqual(templates[0].alanlar[1].tur, 'secim');
      assert.strictEqual(templates[0].alanlar[1].varsayilan, 'opt1');
      assert.deepStrictEqual(templates[0].alanlar[1].secenekler, ['opt1', 'opt2', 'opt3']);
    });

    test('should parse multiple templates', () => {
      const markdown = `### Teknik

#### Template 1
- **ID:** \`t1\`
- **Açıklama:** First template

### Özellik

#### Template 2
- **ID:** \`t2\`
- **Açıklama:** Second template`;

      const templates = MarkdownParser.parseTemplateListesi(markdown);
      
      assert.strictEqual(templates.length, 2);
      assert.strictEqual(templates[0].kategori, 'Teknik');
      assert.strictEqual(templates[1].kategori, 'Özellik');
    });
  });

  suite('parseGorevDetay', () => {
    test('should parse task detail', () => {
      const markdown = `# Test Task

**ID:** task-123
**Durum:** devam_ediyor
**Öncelik:** yuksek
Proje: Test Project (ID: proj-123)
Son Tarih: 2025-07-01
Etiketler: bug, urgent

## Açıklama

This is a detailed description
with multiple lines

## Bağımlılıklar

- Dependency Task (ID: dep-123) - tamamlandi
- Another Dependency (ID: dep-456) - beklemede`;

      const task = MarkdownParser.parseGorevDetay(markdown);
      
      assert.strictEqual(task.baslik, 'Test Task');
      assert.strictEqual(task.id, 'task-123');
      assert.strictEqual(task.durum, 'devam_ediyor');
      assert.strictEqual(task.oncelik, 'yuksek');
      assert.strictEqual(task.proje_id, 'proj-123');
      assert.strictEqual(task.son_tarih, '2025-07-01');
      assert.deepStrictEqual(task.etiketler, ['bug', 'urgent']);
      assert.strictEqual(task.aciklama, 'This is a detailed description\\nwith multiple lines');
      assert.strictEqual(task.bagimliliklar.length, 2);
      assert.strictEqual(task.bagimliliklar[0].hedef_baslik, 'Dependency Task');
      assert.strictEqual(task.bagimliliklar[0].hedef_durum, 'tamamlandi');
    });
  });

  suite('parseOzet', () => {
    test('should parse summary', () => {
      const markdown = `## 📊 Özet Bilgiler

Toplam görev sayısı: 25
Tamamlanan: 10
Devam eden: 5
Bekleyen: 10

Toplam proje sayısı: 3
Aktif proje: Test Project`;

      const summary = MarkdownParser.parseOzet(markdown);
      
      assert.strictEqual(summary.toplamGorev, 25);
      assert.strictEqual(summary.tamamlanan, 10);
      assert.strictEqual(summary.devamEden, 5);
      assert.strictEqual(summary.bekleyen, 10);
      assert.strictEqual(summary.toplamProje, 3);
      assert.strictEqual(summary.aktifProje, 'Test Project');
    });

    test('should handle no active project', () => {
      const markdown = `Toplam görev sayısı: 0
Aktif proje: Yok`;

      const summary = MarkdownParser.parseOzet(markdown);
      
      assert.strictEqual(summary.toplamGorev, 0);
      assert.strictEqual(summary.aktifProje, undefined);
    });
  });

  suite('markdownToHtml', () => {
    test('should convert headers', () => {
      const markdown = '# Header 1\\n## Header 2\\n### Header 3';
      const html = MarkdownParser.markdownToHtml(markdown);
      
      assert(html.includes('<h1>Header 1</h1>'));
      assert(html.includes('<h2>Header 2</h2>'));
      assert(html.includes('<h3>Header 3</h3>'));
    });

    test('should convert formatting', () => {
      const markdown = '**bold** and *italic*';
      const html = MarkdownParser.markdownToHtml(markdown);
      
      assert(html.includes('<strong>bold</strong>'));
      assert(html.includes('<em>italic</em>'));
    });

    test('should convert code and links', () => {
      const markdown = '`code` and [link](http://example.com)';
      const html = MarkdownParser.markdownToHtml(markdown);
      
      assert(html.includes('<code>code</code>'));
      assert(html.includes('<a href="http://example.com">link</a>'));
    });
  });
});