import { Gorev, GorevDurum, GorevOncelik, Bagimlilik } from '../models/gorev';
import { Proje } from '../models/proje';
import { GorevTemplate } from '../models/template';
import { TemplateKategori } from '../models/common';

/**
 * MCP response'larını parse eden utility sınıfı
 */
export class MarkdownParser {
    
    /**
     * Görev detay markdown'ını parse eder
     */
    static parseGorevDetay(markdown: string): Partial<Gorev> {
        const lines = markdown.split('\n');
        const gorev: Partial<Gorev> = {
            bagimliliklar: []
        };
        
        let inBagimlilikSection = false;
        
        for (const line of lines) {
            // Başlık
            if (line.startsWith('# ')) {
                gorev.baslik = line.substring(2).trim();
            }
            
            // ID
            const idMatch = line.match(/\*\*ID:\*\*\s*([a-f0-9-]+)/);
            if (idMatch) {
                gorev.id = idMatch[1];
            }
            
            // Durum
            if (line.includes('Durum:')) {
                const durumMatch = line.match(/\*?\*?Durum:?\*?\*?\s*([\w_]+)/);
                if (durumMatch) {
                    gorev.durum = this.parseDurum(durumMatch[1]);
                }
            }
            
            // Öncelik
            if (line.includes('Öncelik:')) {
                const oncelikMatch = line.match(/\*?\*?Öncelik:?\*?\*?\s*(\w+)/);
                if (oncelikMatch) {
                    gorev.oncelik = this.parseOncelik(oncelikMatch[1]);
                }
            }
            
            // Proje
            if (line.includes('Proje:')) {
                const projeMatch = line.match(/Proje:\s*(.+)/);
                if (projeMatch) {
                    const projeIdMatch = projeMatch[1].match(/\(ID:\s*([^)]+)\)/);
                    if (projeIdMatch) {
                        gorev.proje_id = projeIdMatch[1];
                    }
                }
            }
            
            // Son Tarih
            if (line.includes('Son Tarih:')) {
                const tarihMatch = line.match(/Son Tarih:\s*(\d{4}-\d{2}-\d{2})/);
                if (tarihMatch) {
                    gorev.son_tarih = tarihMatch[1];
                }
            }
            
            // Etiketler
            if (line.includes('Etiketler:')) {
                const etiketMatch = line.match(/Etiketler:\s*(.+)/);
                if (etiketMatch) {
                    gorev.etiketler = etiketMatch[1].split(',').map(e => e.trim());
                }
            }
            
            // Açıklama bölümü başlangıcı
            if (line === '## Açıklama') {
                // Sonraki satırları açıklama olarak topla
                const acikamaLines: string[] = [];
                let i = lines.indexOf(line) + 1;
                while (i < lines.length && !lines[i].startsWith('#')) {
                    if (lines[i].trim()) {
                        acikamaLines.push(lines[i]);
                    }
                    i++;
                }
                gorev.aciklama = acikamaLines.join('\n').trim();
            }
            
            // Bağımlılıklar bölümü
            if (line === '## Bağımlılıklar') {
                inBagimlilikSection = true;
                continue;
            }
            
            if (inBagimlilikSection && line.startsWith('- ')) {
                const bagimlilik = this.parseBagimlilik(line);
                if (bagimlilik) {
                    gorev.bagimliliklar!.push(bagimlilik);
                }
            }
            
            if (inBagimlilikSection && line.startsWith('#')) {
                inBagimlilikSection = false;
            }
        }
        
        return gorev;
    }
    
    /**
     * Görev listesi markdown'ını parse eder
     */
    static parseGorevListesi(markdown: string): Gorev[] {
        const lines = markdown.split('\n');
        const gorevler: Gorev[] = [];
        let currentGorev: Partial<Gorev> | null = null;
        let descriptionLines: string[] = [];
        
        for (let i = 0; i < lines.length; i++) {
            const line = lines[i].trim();
            
            // Görev satırı: - [durum] başlık (öncelik)
            const taskMatch = line.match(/^- \[([^\]]+)\] (.+) \((\w+) öncelik\)/);
            if (taskMatch) {
                // Önceki görevi kaydet
                if (currentGorev && currentGorev.id) {
                    if (descriptionLines.length > 0) {
                        currentGorev.aciklama = descriptionLines.join(' ').trim();
                    }
                    gorevler.push(currentGorev as Gorev);
                }
                
                // Yeni görev
                const [, durum, baslik, oncelik] = taskMatch;
                currentGorev = {
                    baslik,
                    durum: this.parseDurum(durum),
                    oncelik: this.parseOncelik(oncelik),
                    etiketler: []
                };
                descriptionLines = [];
                continue;
            }
            
            // ID satırı
            if (currentGorev && line.includes('ID:')) {
                const idMatch = line.match(/ID:\s*([a-f0-9-]+)/);
                if (idMatch) {
                    currentGorev.id = idMatch[1];
                }
                continue;
            }
            
            // Proje satırı
            if (currentGorev && line.includes('Proje:')) {
                const projeMatch = line.match(/Proje:\s*(.+)/);
                if (projeMatch) {
                    // Proje ismi var, sadece görsel için sakla
                    currentGorev.proje_ismi = projeMatch[1];
                }
                continue;
            }
            
            // ProjeID satırı
            if (currentGorev && line.includes('ProjeID:')) {
                const projeIDMatch = line.match(/ProjeID:\s*([a-f0-9-]+)/);
                if (projeIDMatch) {
                    currentGorev.proje_id = projeIDMatch[1];
                }
                continue;
            }
            
            // Son Tarih
            if (currentGorev && line.includes('Son tarih:')) {
                const tarihMatch = line.match(/Son tarih:\s*(\d{4}-\d{2}-\d{2})/);
                if (tarihMatch) {
                    currentGorev.son_tarih = tarihMatch[1];
                }
                continue;
            }
            
            // Etiketler
            if (currentGorev && line.includes('Etiketler:')) {
                const etiketMatch = line.match(/Etiketler:\s*(.+)/);
                if (etiketMatch) {
                    currentGorev.etiketler = etiketMatch[1].split(',').map(e => e.trim());
                }
                continue;
            }
            
            // Açıklama satırı (ID, Proje, vs. değilse)
            if (currentGorev && line && !line.startsWith('- [') && !line.startsWith('##')) {
                descriptionLines.push(line);
            }
        }
        
        // Son görevi ekle
        if (currentGorev && currentGorev.id) {
            if (descriptionLines.length > 0) {
                currentGorev.aciklama = descriptionLines.join(' ').trim();
            }
            gorevler.push(currentGorev as Gorev);
        }
        
        return gorevler;
    }
    
    /**
     * Proje listesi markdown'ını parse eder
     */
    static parseProjeListesi(markdown: string): Proje[] {
        const lines = markdown.split('\n');
        const projeler: Proje[] = [];
        let currentProje: Partial<Proje> | null = null;
        
        console.log('[MarkdownParser] Parsing project list, first few lines:', lines.slice(0, 10));
        
        for (const line of lines) {
            // Proje başlığı: ### 🔒 Proje İsmi
            if (line.startsWith('###')) {
                if (currentProje && currentProje.id) {
                    projeler.push(currentProje as Proje);
                }
                
                // Emoji ve proje ismini ayıkla
                const projeName = line.replace(/^###\s*/, '').replace(/^[\u{1F300}-\u{1F9FF}\u{1F600}-\u{1F64F}\u{1F680}-\u{1F6FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}]\s*/u, '').trim();
                currentProje = {
                    isim: projeName
                };
            }
            
            if (!currentProje) continue;
            
            // ID satırı
            if (line.includes('**ID:**')) {
                const idMatch = line.match(/\*\*ID:\*\*\s*([a-f0-9-]+)/);
                if (idMatch) {
                    currentProje.id = idMatch[1];
                }
            }
            
            // Tanım satırı
            if (line.includes('**Tanım:**')) {
                const tanimMatch = line.match(/\*\*Tanım:\*\*\s*(.+)/);
                if (tanimMatch) {
                    currentProje.tanim = tanimMatch[1].trim();
                }
            }
            
            // Görev Sayısı satırı
            if (line.includes('**Görev Sayısı:**')) {
                const sayiMatch = line.match(/\*\*Görev Sayısı:\*\*\s*(\d+)/);
                if (sayiMatch) {
                    currentProje.gorev_sayisi = parseInt(sayiMatch[1]);
                    console.log('[MarkdownParser] Found task count for project:', currentProje.isim, '=', currentProje.gorev_sayisi);
                }
            }
        }
        
        // Son projeyi ekle
        if (currentProje && currentProje.id) {
            projeler.push(currentProje as Proje);
        }
        
        return projeler;
    }
    
    /**
     * Template listesi markdown'ını parse eder
     */
    static parseTemplateListesi(markdown: string): GorevTemplate[] {
        const lines = markdown.split('\n');
        const templates: GorevTemplate[] = [];
        let currentTemplate: Partial<GorevTemplate> | null = null;
        let currentCategory: string | null = null;
        let inAlanlarSection = false;
        let alanlar: any[] = [];
        
        for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            
            // Skip empty lines
            if (!line.trim()) {
                inAlanlarSection = false;
                continue;
            }
            
            // Category header: ### Category
            if (line.startsWith('### ') && !line.includes('ID:')) {
                currentCategory = line.substring(4).trim();
                continue;
            }
            
            // Template name header: #### Template Name
            if (line.startsWith('#### ')) {
                // Save previous template
                if (currentTemplate && currentTemplate.id) {
                    currentTemplate.alanlar = alanlar;
                    templates.push(currentTemplate as GorevTemplate);
                }
                
                // Start new template
                const templateName = line.substring(5).trim();
                currentTemplate = {
                    isim: templateName,
                    kategori: currentCategory as TemplateKategori,
                    alanlar: [],
                    aktif: true,
                    aciklama_template: '',
                    ornek_degerler: {}
                };
                alanlar = [];
                inAlanlarSection = false;
                continue;
            }
            
            if (!currentTemplate) continue;
            
            // ID line: - **ID:** `xxxx`
            if (line.includes('**ID:**')) {
                const idMatch = line.match(/\*\*ID:\*\*\s*`?([a-f0-9-]+)`?/);
                if (idMatch) {
                    currentTemplate.id = idMatch[1];
                }
                continue;
            }
            
            // Description line: - **Açıklama:** xxx
            if (line.includes('**Açıklama:**')) {
                const descMatch = line.match(/\*\*Açıklama:\*\*\s*(.+)/);
                if (descMatch) {
                    currentTemplate.tanim = descMatch[1].trim();
                }
                continue;
            }
            
            // Default title line: - **Başlık Şablonu:** `xxx`
            if (line.includes('**Başlık Şablonu:**')) {
                const titleMatch = line.match(/\*\*Başlık Şablonu:\*\*\s*`?([^`]+)`?/);
                if (titleMatch) {
                    currentTemplate.varsayilan_baslik = titleMatch[1].trim();
                }
                continue;
            }
            
            // Fields section start: - **Alanlar:**
            if (line.includes('**Alanlar:**')) {
                inAlanlarSection = true;
                continue;
            }
            
            // Field line:   - `fieldname` (type) *(zorunlu)* - varsayılan: value - seçenekler: opt1, opt2
            if (inAlanlarSection && line.trim().startsWith('- `')) {
                const fieldMatch = line.match(/- `(\w+)` \((\w+)\)(.*)$/);
                if (fieldMatch) {
                    const [, fieldName, fieldType, rest] = fieldMatch;
                    const field: any = {
                        isim: fieldName,
                        tur: this.mapFieldType(fieldType),
                        zorunlu: rest.includes('*(zorunlu)*'),
                        varsayilan: '',
                        secenekler: []
                    };
                    
                    // Extract default value
                    const defaultMatch = rest.match(/varsayılan:\s*([^-]+)/);
                    if (defaultMatch) {
                        field.varsayilan = defaultMatch[1].trim();
                    }
                    
                    // Extract options
                    const optionsMatch = rest.match(/seçenekler:\s*(.+)$/);
                    if (optionsMatch) {
                        field.secenekler = optionsMatch[1].split(',').map(opt => opt.trim());
                    }
                    
                    alanlar.push(field);
                }
                continue;
            }
        }
        
        // Don't forget the last template
        if (currentTemplate && currentTemplate.id) {
            currentTemplate.alanlar = alanlar;
            templates.push(currentTemplate as GorevTemplate);
        }
        
        return templates;
    }
    
    /**
     * Özet bilgilerini parse eder
     */
    static parseOzet(markdown: string): {
        toplamGorev: number;
        tamamlanan: number;
        devamEden: number;
        bekleyen: number;
        toplamProje: number;
        aktifProje?: string;
    } {
        const ozet: {
            toplamGorev: number;
            tamamlanan: number;
            devamEden: number;
            bekleyen: number;
            toplamProje: number;
            aktifProje?: string;
        } = {
            toplamGorev: 0,
            tamamlanan: 0,
            devamEden: 0,
            bekleyen: 0,
            toplamProje: 0
        };
        
        const lines = markdown.split('\n');
        
        for (const line of lines) {
            const toplamMatch = line.match(/Toplam görev sayısı:\s*(\d+)/);
            if (toplamMatch) {
                ozet.toplamGorev = parseInt(toplamMatch[1]);
            }
            
            const tamamMatch = line.match(/Tamamlanan:\s*(\d+)/);
            if (tamamMatch) {
                ozet.tamamlanan = parseInt(tamamMatch[1]);
            }
            
            const devamMatch = line.match(/Devam eden:\s*(\d+)/);
            if (devamMatch) {
                ozet.devamEden = parseInt(devamMatch[1]);
            }
            
            const bekleyenMatch = line.match(/Bekleyen:\s*(\d+)/);
            if (bekleyenMatch) {
                ozet.bekleyen = parseInt(bekleyenMatch[1]);
            }
            
            const projeMatch = line.match(/Toplam proje sayısı:\s*(\d+)/);
            if (projeMatch) {
                ozet.toplamProje = parseInt(projeMatch[1]);
            }
            
            const aktifMatch = line.match(/Aktif proje:\s*(.+)/);
            if (aktifMatch && !aktifMatch[1].includes('Yok')) {
                ozet.aktifProje = aktifMatch[1];
            }
        }
        
        return ozet;
    }
    
    /**
     * Durum string'ini enum'a çevirir
     */
    private static parseDurum(durum: string): GorevDurum {
        const normalizedDurum = durum.toLowerCase().replace(/\s/g, '_');
        
        switch (normalizedDurum) {
            case 'beklemede':
            case 'bekleyen':
            case 'pending':
                return GorevDurum.Beklemede;
            case 'devam_ediyor':
            case 'devam':
            case 'in_progress':
                return GorevDurum.DevamEdiyor;
            case 'tamamlandi':
            case 'tamamlandı':
            case 'completed':
                return GorevDurum.Tamamlandi;
            default:
                return GorevDurum.Beklemede;
        }
    }
    
    /**
     * Öncelik string'ini enum'a çevirir
     */
    private static parseOncelik(oncelik: string): GorevOncelik {
        const normalizedOncelik = oncelik.toLowerCase();
        
        switch (normalizedOncelik) {
            case 'yuksek':
            case 'yüksek':
            case 'high':
                return GorevOncelik.Yuksek;
            case 'orta':
            case 'medium':
                return GorevOncelik.Orta;
            case 'dusuk':
            case 'düşük':
            case 'low':
                return GorevOncelik.Dusuk;
            default:
                return GorevOncelik.Orta;
        }
    }
    
    /**
     * Bağımlılık satırını parse eder
     */
    private static parseBagimlilik(line: string): Bagimlilik | null {
        // Format: - Görev başlığı (ID: xxx) - Durum
        const match = line.match(/- (.+) \(ID: ([^)]+)\) - (.+)/);
        if (match) {
            return {
                hedef_baslik: match[1],
                hedef_id: match[2],
                hedef_durum: this.parseDurum(match[3]),
                kaynak_id: '', // Bilinmiyor
                baglanti_tip: 'engelliyor' // Varsayılan
            };
        }
        
        return null;
    }
    
    /**
     * Markdown'ı HTML'e çevirir (basit versiyon)
     */
    static markdownToHtml(markdown: string): string {
        let html = markdown;
        
        // Headers
        html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>');
        html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>');
        html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>');
        
        // Bold
        html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
        
        // Italic
        html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');
        
        // Code
        html = html.replace(/`(.+?)`/g, '<code>$1</code>');
        
        // Links
        html = html.replace(/\[(.+?)\]\((.+?)\)/g, '<a href="$2">$1</a>');
        
        // Lists
        html = html.replace(/^- (.+)$/gm, '<li>$1</li>');
        html = html.replace(/(<li>.*<\/li>)/s, '<ul>$1</ul>');
        
        // Blockquotes
        html = html.replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>');
        
        // Line breaks
        html = html.replace(/\n/g, '<br>');
        
        return html;
    }
    
    /**
     * HTML'den güvenli metin çıkarır
     */
    static extractTextFromHtml(html: string): string {
        return html
            .replace(/<[^>]*>/g, '') // HTML etiketlerini kaldır
            .replace(/&amp;/g, '&')
            .replace(/&lt;/g, '<')
            .replace(/&gt;/g, '>')
            .replace(/&quot;/g, '"')
            .replace(/&#039;/g, "'");
    }
    
    /**
     * Field type'ı map eder
     */
    private static mapFieldType(type: string): 'metin' | 'sayi' | 'tarih' | 'secim' {
        switch (type.toLowerCase()) {
            case 'text':
            case 'metin':
                return 'metin';
            case 'number':
            case 'sayi':
                return 'sayi';
            case 'date':
            case 'tarih':
                return 'tarih';
            case 'select':
            case 'secim':
                return 'secim';
            default:
                return 'metin';
        }
    }
}