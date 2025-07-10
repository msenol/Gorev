import { Gorev, GorevDurum, GorevOncelik, Bagimlilik } from '../models/gorev';
import { Proje } from '../models/proje';
import { GorevTemplate } from '../models/template';
import { TemplateKategori } from '../models/common';
import { Logger } from '../utils/logger';

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
        
        // Check if this is the new compact format - check for indented lines and subtasks too
        const hasCompactFormat = lines.some(line => {
            const trimmed = line.trim();
            return /^\[⏳\]|\[🚀\]|\[✅\]|\[🔄\]|\[✓\]/.test(trimmed) || 
                   /^└─\s*\[⏳\]|^└─\s*\[🚀\]|^└─\s*\[✅\]|^└─\s*\[🔄\]|^└─\s*\[✓\]/.test(trimmed);
        });
        
        Logger.debug('[MarkdownParser] Starting to parse tasks...');
        Logger.debug('[MarkdownParser] Format detected:', hasCompactFormat ? 'Compact (v0.8.1+)' : 'Legacy');
        Logger.debug('[MarkdownParser] First 5 lines:', lines.slice(0, 5));
        Logger.debug('[MarkdownParser] Total lines:', lines.length);
        
        if (hasCompactFormat) {
            return this.parseCompactGorevListesi(markdown);
        }
        
        // Legacy parsing logic
        const gorevler: Gorev[] = [];
        const gorevMap: Map<string, Gorev> = new Map();
        const parentStack: { id: string, level: number }[] = [];
        let currentGorev: Partial<Gorev> | null = null;
        let descriptionLines: string[] = [];
        
        Logger.debug('[MarkdownParser] Total lines:', lines.length);
        
        for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            // Calculate indent level before trimming
            let indentLevel = 0;
            if (line.startsWith('  ')) {
                // Count spaces at the beginning
                const spaceCount = line.length - line.trimStart().length;
                indentLevel = Math.floor(spaceCount / 2);
            }
            const trimmedLine = line.trim();
            
            // Görev satırı formatları:
            // └─ [durum] başlık (öncelik) - alt görev
            // [durum] başlık (öncelik) - normal görev  
            // - [durum] başlık (öncelik) - liste formatında görev
            const taskRegex = /^(└─\s*|-\s*)?\[([^\]]+)\] (.+) \(([a-zA-ZğüşıöçĞÜŞİÖÇ]+) öncelik\)/;
            const taskMatch = trimmedLine.match(taskRegex);
            
            if (i < 10 && trimmedLine.includes('[')) {
                Logger.debug(`[MarkdownParser] Line ${i}: "${line}"`);
                Logger.debug(`[MarkdownParser] Indent level: ${indentLevel}, trimmed: "${trimmedLine}"`);
                Logger.debug(`[MarkdownParser] taskMatch:`, taskMatch);
            }
            
            if (taskMatch) {
                // Önceki görevi kaydet
                if (currentGorev && currentGorev.id) {
                    if (descriptionLines.length > 0) {
                        currentGorev.aciklama = descriptionLines.join(' ').trim();
                    }
                    const gorev = currentGorev as Gorev;
                    gorev.seviye = currentGorev.seviye || 0;
                    gorev.alt_gorevler = [];
                    
                    // Parent-child ilişkisini kur
                    while (parentStack.length > 0 && parentStack[parentStack.length - 1].level >= gorev.seviye) {
                        parentStack.pop();
                    }
                    
                    if (parentStack.length > 0) {
                        const parent = parentStack[parentStack.length - 1];
                        gorev.parent_id = parent.id;
                        const parentGorev = gorevMap.get(parent.id);
                        if (parentGorev && parentGorev.alt_gorevler) {
                            parentGorev.alt_gorevler.push(gorev);
                            Logger.debug(`[MarkdownParser] Added ${gorev.baslik} as child of ${parentGorev.baslik}`);
                        }
                    }
                    
                    // Add ALL tasks to the main array, regardless of hierarchy
                    gorevler.push(gorev);
                    Logger.debug(`[MarkdownParser] Added ${gorev.baslik} as task with proje_id: ${gorev.proje_id || 'NONE'}`); 
                    
                    gorevMap.set(gorev.id!, gorev);
                    parentStack.push({ id: gorev.id!, level: gorev.seviye });
                    Logger.debug(`[MarkdownParser] Parent stack after adding ${gorev.baslik}:`, parentStack.map(p => `${p.id}(L${p.level})`).join(' -> '));
                }
                
                // Yeni görev
                // taskMatch groups: [full match, prefix (└─ or -), durum, baslik, oncelik]
                const [, prefix, durum, baslik, oncelik] = taskMatch;
                currentGorev = {
                    baslik,
                    durum: this.parseDurum(durum),
                    oncelik: this.parseOncelik(oncelik),
                    etiketler: [],
                    alt_gorevler: [],
                    seviye: indentLevel
                };
                descriptionLines = [];
                continue;
            }
            
            // ID satırı
            if (currentGorev && trimmedLine.includes('ID:')) {
                const idMatch = trimmedLine.match(/ID:\s*([a-f0-9-]+)/);
                if (idMatch) {
                    currentGorev.id = idMatch[1];
                    Logger.debug(`[MarkdownParser] Parsed task ID: ${idMatch[1]} for task: ${currentGorev.baslik}`);
                }
                continue;
            }
            
            // Proje satırı
            if (currentGorev && trimmedLine.includes('Proje:')) {
                const projeMatch = trimmedLine.match(/Proje:\s*(.+)/);
                if (projeMatch) {
                    // Proje ismi var, sadece görsel için sakla
                    (currentGorev as any).proje_isim = projeMatch[1];
                }
                continue;
            }
            
            // ProjeID satırı
            if (currentGorev && trimmedLine.includes('ProjeID:')) {
                Logger.debug(`[MarkdownParser] Found ProjeID line: "${trimmedLine}"`);
                const projeIDMatch = trimmedLine.match(/ProjeID:\s*([a-f0-9-]+)/);
                if (projeIDMatch) {
                    currentGorev.proje_id = projeIDMatch[1];
                    Logger.debug(`[MarkdownParser] Parsed ProjeID: ${projeIDMatch[1]}`);
                } else {
                    Logger.debug(`[MarkdownParser] Failed to parse ProjeID from: "${trimmedLine}"`);
                }
                continue;
            }
            
            // Son Tarih
            if (currentGorev && trimmedLine.includes('Son tarih:')) {
                const tarihMatch = trimmedLine.match(/Son tarih:\s*(\d{4}-\d{2}-\d{2})/);
                if (tarihMatch) {
                    currentGorev.son_tarih = tarihMatch[1];
                }
                continue;
            }
            
            // Etiketler
            if (currentGorev && trimmedLine.includes('Etiketler:')) {
                const etiketMatch = trimmedLine.match(/Etiketler:\s*(.+)/);
                if (etiketMatch) {
                    currentGorev.etiketler = etiketMatch[1].split(',').map(e => e.trim());
                }
                continue;
            }
            
            // Bağımlı görev sayısı
            if (currentGorev && trimmedLine.includes('Bağımlı görev sayısı:')) {
                const match = trimmedLine.match(/Bağımlı görev sayısı:\s*(\d+)/);
                if (match) {
                    currentGorev.bagimli_gorev_sayisi = parseInt(match[1]);
                }
                continue;
            }
            
            // Tamamlanmamış bağımlılık sayısı
            if (currentGorev && trimmedLine.includes('Tamamlanmamış bağımlılık sayısı:')) {
                const match = trimmedLine.match(/Tamamlanmamış bağımlılık sayısı:\s*(\d+)/);
                if (match) {
                    currentGorev.tamamlanmamis_bagimlilik_sayisi = parseInt(match[1]);
                }
                continue;
            }
            
            // Bu göreve bağımlı sayısı
            if (currentGorev && trimmedLine.includes('Bu göreve bağımlı sayısı:')) {
                const match = trimmedLine.match(/Bu göreve bağımlı sayısı:\s*(\d+)/);
                if (match) {
                    currentGorev.bu_goreve_bagimli_sayisi = parseInt(match[1]);
                }
                continue;
            }
            
            // Açıklama satırı (ID, Proje, vs. değilse)
            if (currentGorev && trimmedLine && !trimmedLine.startsWith('[') && !trimmedLine.startsWith('##') && !trimmedLine.startsWith('└─')) {
                descriptionLines.push(trimmedLine);
            }
        }
        
        // Son görevi ekle
        if (currentGorev && currentGorev.id) {
            if (descriptionLines.length > 0) {
                currentGorev.aciklama = descriptionLines.join(' ').trim();
            }
            const gorev = currentGorev as Gorev;
            gorev.seviye = gorev.seviye || 0;
            gorev.alt_gorevler = [];
            
            // Parent-child ilişkisini kur
            while (parentStack.length > 0 && parentStack[parentStack.length - 1].level >= gorev.seviye) {
                parentStack.pop();
            }
            
            if (parentStack.length > 0) {
                const parent = parentStack[parentStack.length - 1];
                gorev.parent_id = parent.id;
                const parentGorev = gorevMap.get(parent.id);
                if (parentGorev && parentGorev.alt_gorevler) {
                    parentGorev.alt_gorevler.push(gorev);
                }
            }
            
            // Add ALL tasks to the main array, regardless of hierarchy
            gorevler.push(gorev);
            
            gorevMap.set(gorev.id, gorev);
        }
        
        Logger.debug('[MarkdownParser] Total tasks parsed:', gorevler.length);
        Logger.debug('[MarkdownParser] Tasks with subtasks:', gorevler.filter(g => g.alt_gorevler && g.alt_gorevler.length > 0).map(g => `${g.baslik} (${g.alt_gorevler!.length} subtasks)`));
        Logger.debug('[MarkdownParser] Tasks with proje_id:', gorevler.filter(g => g.proje_id).map(g => `${g.baslik} -> ${g.proje_id}`));
        Logger.debug('[MarkdownParser] Tasks WITHOUT proje_id:', gorevler.filter(g => !g.proje_id).map(g => g.baslik));
        return gorevler;
    }
    
    /**
     * Compact format görev listesi markdown'ını parse eder (v0.8.1+)
     * Format: [StatusIcon] Title (Priority)
     *         Description | Tarih: DD/MM | Etiket: tags | ID:uuid
     */
    static parseCompactGorevListesi(markdown: string): Gorev[] {
        const lines = markdown.split('\n');
        const gorevler: Gorev[] = [];
        const gorevMap = new Map<string, Gorev>(); // ID to Gorev mapping
        const parentStack: { id: string, indentLevel: number }[] = []; // Track parent hierarchy
        let i = 0;
        
        Logger.debug('[MarkdownParser] Parsing compact format with hierarchy...');
        Logger.debug('[MarkdownParser] Total lines to parse:', lines.length);
        
        while (i < lines.length) {
            const line = lines[i];
            
            // Calculate indent level before trimming
            let indentLevel = 0;
            if (line.startsWith('  ')) {
                const match = line.match(/^(\s*)/);
                if (match) {
                    indentLevel = Math.floor(match[1].length / 2);
                }
            }
            
            const trimmedLine = line.trim();
            
            // Skip empty lines and header lines
            if (!trimmedLine || trimmedLine.startsWith('#') || trimmedLine.startsWith('Görevler (') || trimmedLine.startsWith('Proje:')) {
                i++;
                continue;
            }
            
            // Match task header: [⏳|🚀|✅|✓|🔄] Title (Y|O|D|düşük|orta|yüksek)
            // Also handle subtasks with └─ prefix
            // Updated regex to handle titles with emojis and special characters
            const headerMatch = trimmedLine.match(/^(└─\s*)?\[(⏳|🚀|✅|✓|🔄)\]\s+(.+?)\s+\(([YOD]|düşük|orta|yüksek)\)$/);
            if (headerMatch) {
                const [_, prefix, statusIcon, title, priorityLetter] = headerMatch;
                Logger.debug(`[MarkdownParser] Found task header at line ${i}: "${title}", indent: ${indentLevel}, has prefix: ${!!prefix}`);
                
                // Map status icon to durum
                let durum = GorevDurum.Beklemede;
                if (statusIcon === '🚀' || statusIcon === '🔄') durum = GorevDurum.DevamEdiyor;
                else if (statusIcon === '✅' || statusIcon === '✓') durum = GorevDurum.Tamamlandi;
                
                // Map priority letter to oncelik
                let oncelik = GorevOncelik.Orta;
                if (priorityLetter === 'Y' || priorityLetter === 'yüksek') oncelik = GorevOncelik.Yuksek;
                else if (priorityLetter === 'D' || priorityLetter === 'düşük') oncelik = GorevOncelik.Dusuk;
                else if (priorityLetter === 'O' || priorityLetter === 'orta') oncelik = GorevOncelik.Orta;
                
                // Check next line(s) for details
                i++;
                let detailsLine = '';
                
                // Collect details from next line(s) - sometimes description spans multiple lines
                if (i < lines.length) {
                    const nextLine = lines[i];
                    detailsLine = nextLine.trim();
                    
                    // Check if this is a details line without ID (multiline description)
                    if (detailsLine && !detailsLine.includes('ID:')) {
                        // This might be a multiline description, keep reading until we find ID
                        let fullDetails = detailsLine;
                        i++;
                        while (i < lines.length) {
                            const continuationLine = lines[i].trim();
                            if (continuationLine.includes('ID:')) {
                                fullDetails += ' | ' + continuationLine;
                                break;
                            } else if (continuationLine && !continuationLine.startsWith('[')) {
                                fullDetails += ' ' + continuationLine;
                                i++;
                            } else {
                                // No ID found, this might not be a valid task
                                Logger.warn(`[MarkdownParser] No ID found for task: ${title}`);
                                break;
                            }
                        }
                        detailsLine = fullDetails;
                    }
                    
                    Logger.debug(`[MarkdownParser] Details line: ${detailsLine}`);
                    
                    // Parse different detail formats
                    let detailsMatch = null;
                    let description = '';
                    let dueDate = '';
                    let tags: string[] = [];
                    let projectId = '';
                    let projectName = '';
                    let id = '';
                    
                    // Format 1: Description | Bekleyen: N | ID:uuid
                    // Format 2: Description | Tarih: DD/MM | Bekleyen: N | ID:uuid
                    // Format 3: Description | Proje: Name | ID:uuid
                    // Format 4: - Description | ID:uuid (for subtasks)
                    
                    // Extract ID first (always at the end)
                    const idMatch = detailsLine.match(/ID:([a-f0-9-]+)$/);
                    if (idMatch) {
                        id = idMatch[1];
                        const detailsWithoutId = detailsLine.substring(0, detailsLine.lastIndexOf('|')).trim();
                        
                        // Extract description (everything before first |)
                        // Remove leading dash if present
                        const cleanDetails = detailsWithoutId.startsWith('- ') ? detailsWithoutId.substring(2) : detailsWithoutId;
                        description = cleanDetails.split('|')[0].trim();
                        
                        // Extract other fields
                        const parts = detailsWithoutId.split('|').slice(1);
                        for (const part of parts) {
                            const trimmedPart = part.trim();
                            if (trimmedPart.startsWith('Tarih:')) {
                                dueDate = trimmedPart.substring(6).trim();
                            } else if (trimmedPart.startsWith('Etiket:')) {
                                const etiketPart = trimmedPart.substring(7).trim();
                                // Handle "X adet" format
                                const adetMatch = etiketPart.match(/(\d+)\s+adet/);
                                if (adetMatch) {
                                    // Don't parse tags if it's just a count
                                    tags = [];
                                } else {
                                    tags = etiketPart.split(',').map(t => t.trim());
                                }
                            } else if (trimmedPart.startsWith('Proje:')) {
                                projectName = trimmedPart.substring(6).trim();
                                // Try to extract project ID from the project name if it contains ID
                                const projIdMatch = projectName.match(/\(ID:\s*([a-f0-9-]+)\)/);
                                if (projIdMatch) {
                                    projectId = projIdMatch[1];
                                    projectName = projectName.replace(/\s*\(ID:[^)]+\)/, '').trim();
                                }
                            }
                        }
                    }
                    
                    if (id) {
                        Logger.debug(`[MarkdownParser] Parsed details - ID: ${id}, Description: ${description}, ProjectID: ${projectId}, ProjectName: ${projectName}`);
                        
                        // Convert date from DD/MM to YYYY-MM-DD (assuming current year) if date exists
                        let formattedDate = '';
                        if (dueDate) {
                            const currentYear = new Date().getFullYear();
                            const [day, month] = dueDate.split('/');
                            if (day && month) {
                                formattedDate = `${currentYear}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
                            }
                        }
                        
                        const gorev: Gorev = {
                            id: id.trim(),
                            baslik: title,
                            aciklama: description ? description.trim() : '',
                            durum,
                            oncelik,
                            proje_id: projectId || '', // Use extracted project ID if available
                            son_tarih: formattedDate || '',
                            etiketler: tags || [],
                            olusturma_tarih: new Date().toISOString(), // Not provided in compact format
                            guncelleme_tarih: new Date().toISOString(), // Not provided in compact format
                            alt_gorevler: [],
                            seviye: indentLevel,
                            parent_id: '' // Will be set based on hierarchy
                        };
                        
                        // Update parent stack based on indent level
                        while (parentStack.length > 0 && parentStack[parentStack.length - 1].indentLevel >= indentLevel) {
                            parentStack.pop();
                        }
                        
                        // Set parent_id if this is a subtask
                        if (parentStack.length > 0 && indentLevel > 0) {
                            const parent = parentStack[parentStack.length - 1];
                            gorev.parent_id = parent.id;
                            // Add to parent's alt_gorevler
                            const parentGorev = gorevMap.get(parent.id);
                            if (parentGorev) {
                                parentGorev.alt_gorevler = parentGorev.alt_gorevler || [];
                                parentGorev.alt_gorevler.push(gorev);
                                Logger.debug(`[MarkdownParser] Added ${title} as child of ${parentGorev.baslik}`);
                            }
                        } else {
                            // Ensure root tasks have empty parent_id
                            gorev.parent_id = '';
                        }
                        
                        // Add to map and stack
                        gorevMap.set(gorev.id, gorev);
                        // Add ALL tasks to the main array, not just root tasks
                        gorevler.push(gorev);
                        
                        // Push to parent stack for potential children
                        // Only push if this task could potentially have children (not at max depth, etc.)
                        if (indentLevel === 0 || parentStack.some(p => p.indentLevel < indentLevel)) {
                            parentStack.push({ id: gorev.id, indentLevel });
                        }
                        
                        Logger.debug(`[MarkdownParser] Successfully parsed task: ${title} (${id}) at level ${indentLevel}, parent_id: ${gorev.parent_id || 'NONE'}`);
                    } else {
                        Logger.warn(`[MarkdownParser] Failed to parse details line: ${detailsLine}`);
                    }
                } else {
                    Logger.warn(`[MarkdownParser] No details line found for task: ${title}`);
                }
            } else if (trimmedLine.includes('[') && trimmedLine.includes(']')) {
                Logger.debug(`[MarkdownParser] Line ${i} contains brackets but didn't match pattern: "${trimmedLine}"`);
                Logger.debug(`[MarkdownParser] Regex test result:`, /^(└─\s*)?\[(⏳|🚀|✅|✓|🔄)\]\s+(.+?)\s+\(([YOD]|düşük|orta|yüksek)\)$/.test(trimmedLine));
            }
            i++;
        }
        
        Logger.debug('[MarkdownParser] Total tasks parsed:', gorevler.length);
        Logger.debug('[MarkdownParser] Tasks in map:', gorevMap.size);
        Logger.debug('[MarkdownParser] Task IDs parsed:', Array.from(gorevMap.keys()).sort());
        
        if (gorevler.length === 0) {
            Logger.warn('[MarkdownParser] No tasks were parsed! First 10 lines of markdown:');
            lines.slice(0, 10).forEach((line, idx) => Logger.debug(`  ${idx}: ${line}`));
        }
        
        // Return all tasks from the array (which should include all tasks regardless of hierarchy)
        return gorevler;
    }
    
    /**
     * Proje listesi markdown'ını parse eder
     */
    static parseProjeListesi(markdown: string): Proje[] {
        const lines = markdown.split('\n');
        const projeler: Proje[] = [];
        let currentProje: Partial<Proje> | null = null;
        
        Logger.debug('[MarkdownParser] Parsing project list, first few lines:', lines.slice(0, 10));
        
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
                    Logger.debug('[MarkdownParser] Found task count for project:', currentProje.isim, '=', currentProje.gorev_sayisi);
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
        
        // Handle emoji status indicators
        switch (durum) {
            case '✓':
            case '✅':
                return GorevDurum.Tamamlandi;
            case '🔄':
            case '⚡':
                return GorevDurum.DevamEdiyor;
            case '⏳':
            case '○':
                return GorevDurum.Beklemede;
        }
        
        // Handle text status
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