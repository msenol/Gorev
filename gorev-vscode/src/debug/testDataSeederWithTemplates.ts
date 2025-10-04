import { ApiClient } from '../api/client';
import { GorevDurum, GorevOncelik } from '../models/common';
import { Logger } from '../utils/logger';
import * as vscode from 'vscode';

/**
 * Template-based test data seeder for Gorev
 * Uses templateden_gorev_olustur instead of deprecated gorev_olustur
 */
export class TestDataSeederWithTemplates {
    // Template IDs from database
    private readonly TEMPLATE_IDS = {
        BUG_RAPORU: '4dd56a2a-caf4-472c-8c0f-276bc8a1f880',
        OZELLIK_ISTEGI: '6b083358-9c4d-4f4e-b041-9288c05a1bb7',
        TEKNIK_BORC: '69e2b237-7c2e-4459-9d46-ea6c05aba39a',
        ARASTIRMA_GOREVI: '13f04fe2-b5b6-4fd6-8684-5eca5dc2770d'
    };

    constructor(private apiClient: ApiClient) {}

    /**
     * Test verilerini oluştur
     */
    async seedTestData(): Promise<void> {
        const result = await vscode.window.showInformationMessage(
            'Template-based test verileri oluşturulacak. Mevcut veriler korunacak. Devam etmek istiyor musunuz?',
            'Evet, Oluştur',
            'Hayır'
        );

        if (result !== 'Evet, Oluştur') {
            return;
        }

        try {
            await vscode.window.withProgress({
                location: vscode.ProgressLocation.Notification,
                title: 'Test verileri oluşturuluyor...',
                cancellable: false
            }, async (progress) => {
                // 1. Test projeleri oluştur
                progress.report({ increment: 10, message: 'Projeler oluşturuluyor...' });
                const projectIds = await this.createTestProjects();

                // 2. Test görevleri oluştur (template-based)
                progress.report({ increment: 30, message: 'Template görevleri oluşturuluyor...' });
                const taskIds = await this.createTemplateBasedTasks(projectIds);

                // 3. Bağımlılıklar oluştur
                progress.report({ increment: 20, message: 'Bağımlılıklar oluşturuluyor...' });
                await this.createTestDependencies(taskIds);

                // 4. Alt görevler oluştur
                progress.report({ increment: 10, message: 'Alt görevler oluşturuluyor...' });
                await this.createSubtasks(taskIds);

                // 5. Bazı görevleri tamamla ve AI interaksiyonları ekle
                progress.report({ increment: 20, message: 'Görev durumları güncelleniyor...' });
                await this.updateSomeTaskStatuses(taskIds);

                // 6. AI context oluştur
                progress.report({ increment: 10, message: 'AI context oluşturuluyor...' });
                await this.setupAIContext(taskIds);

                progress.report({ increment: 10, message: 'Tamamlandı!' });
            });

            vscode.window.showInformationMessage('✅ Template-based test verileri başarıyla oluşturuldu!');
        } catch (error) {
            vscode.window.showErrorMessage(`Test verileri oluşturulamadı: ${error}`);
            Logger.error('Test data seeding failed:', error);
        }
    }

    /**
     * Test projeleri oluştur
     */
    private async createTestProjects(): Promise<string[]> {
        const projects = [
            {
                isim: '🚀 Yeni Web Sitesi',
                tanim: 'Şirket web sitesinin yeniden tasarımı ve geliştirilmesi'
            },
            {
                isim: '📱 Mobil Uygulama',
                tanim: 'iOS ve Android için mobil uygulama geliştirme projesi'
            },
            {
                isim: '🔧 Backend API',
                tanim: 'RESTful API ve mikroservis mimarisi geliştirme'
            },
            {
                isim: '📊 Veri Analitiği',
                tanim: 'Müşteri davranış analizi ve raporlama sistemi'
            },
            {
                isim: '🔒 Güvenlik Güncellemeleri',
                tanim: 'Sistem güvenliği ve penetrasyon testi projesi'
            }
        ];

        const projectIds: string[] = [];

        for (const project of projects) {
            try {
                const result = await this.apiClient.callTool('proje_olustur', project);
                const responseText = result.content[0].text;
                Logger.debug('Project creation response:', responseText);
                
                const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                if (idMatch) {
                    projectIds.push(idMatch[1]);
                    Logger.info(`Created project: ${project.isim} with ID: ${idMatch[1]}`);
                } else {
                    Logger.warn(`Could not parse project ID for: ${project.isim}`);
                }
            } catch (error) {
                Logger.error('Failed to create project:', error);
            }
        }

        // İlk projeyi aktif yap
        if (projectIds.length > 0) {
            await this.apiClient.callTool('proje_aktif_yap', { proje_id: projectIds[0] });
        }

        return projectIds;
    }

    /**
     * Template-based görevler oluştur
     */
    private async createTemplateBasedTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];

        // Bug Raporu örnekleri
        const bugTasks = [
            {
                templateId: this.TEMPLATE_IDS.BUG_RAPORU,
                projectId: projectIds[0],
                degerler: {
                    baslik: 'Login sayfası 404 hatası veriyor',
                    aciklama: 'Production ortamında /login URL\'ine gittiğimizde 404 hatası alıyoruz',
                    modul: 'Authentication',
                    ortam: 'production',
                    adimlar: '1. Production URL\'ine git\n2. /login sayfasına git\n3. 404 hatası görünüyor',
                    beklenen: 'Login sayfası açılmalı',
                    mevcut: '404 Not Found hatası',
                    cozum: 'Routing konfigürasyonu kontrol edilmeli',
                    oncelik: 'yuksek',
                    etiketler: 'bug,critical,production'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.BUG_RAPORU,
                projectId: projectIds[4],
                degerler: {
                    baslik: 'SSL sertifikası expire olmuş',
                    aciklama: 'Ana domain ve subdomain\'lerde SSL sertifikası süresi dolmuş',
                    modul: 'Infrastructure',
                    ortam: 'production',
                    adimlar: '1. Herhangi bir subdomain\'e HTTPS ile git\n2. Sertifika uyarısı görünüyor',
                    beklenen: 'Valid SSL sertifikası ile güvenli bağlantı',
                    mevcut: 'NET::ERR_CERT_DATE_INVALID hatası',
                    cozum: 'Wildcard SSL sertifikası yenilenmeli',
                    oncelik: 'yuksek',
                    etiketler: 'security,infrastructure,urgent'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.BUG_RAPORU,
                projectId: projectIds[1],
                degerler: {
                    baslik: 'Push notification Android\'de çalışmıyor',
                    aciklama: 'Firebase Cloud Messaging entegrasyonu iOS\'ta çalışıyor ama Android\'de bildirimler gelmiyor',
                    modul: 'Mobile/Notifications',
                    ortam: 'production',
                    adimlar: '1. Android cihazda uygulamayı aç\n2. Bildirim izni ver\n3. Test bildirimi gönder\n4. Bildirim gelmiyor',
                    beklenen: 'Push notification alınmalı',
                    mevcut: 'Bildirimler Android\'de alınmıyor',
                    oncelik: 'yuksek',
                    etiketler: 'mobile,bug,firebase'
                }
            }
        ];

        // Özellik İsteği örnekleri
        const featureTasks = [
            {
                templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                projectId: projectIds[0],
                degerler: {
                    baslik: 'Ana sayfa hero section tasarımı',
                    aciklama: 'Modern ve responsive ana sayfa tasarımı yapılacak',
                    amac: 'Kullanıcıların ilk izlenimini güçlendirmek ve dönüşüm oranını artırmak',
                    kullanicilar: 'Tüm web sitesi ziyaretçileri',
                    kriterler: '1. Mobile responsive\n2. Hızlı yükleme\n3. A/B test ready\n4. SEO optimized',
                    ui_ux: 'Hero section, özellikler bölümü, testimonials ve CTA butonları',
                    ilgili: 'Landing page, conversion optimization',
                    efor: 'orta',
                    oncelik: 'yuksek',
                    etiketler: 'design,frontend,urgent'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                projectId: projectIds[0],
                degerler: {
                    baslik: 'Kullanıcı giriş sistemi',
                    aciklama: 'JWT tabanlı authentication sistemi kurulacak',
                    amac: 'Güvenli kullanıcı kimlik doğrulama ve yetkilendirme sistemi',
                    kullanicilar: 'Kayıtlı kullanıcılar ve yöneticiler',
                    kriterler: '1. JWT token\n2. Refresh token\n3. Remember me\n4. Social login (Google, GitHub)',
                    ui_ux: 'Login, register, forgot password sayfaları',
                    ilgili: 'User management, session handling',
                    efor: 'büyük',
                    oncelik: 'yuksek',
                    etiketler: 'backend,security,feature'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                projectId: projectIds[1],
                degerler: {
                    baslik: 'Dark mode desteği',
                    aciklama: 'Sistem ayarlarına göre otomatik tema değişimi',
                    amac: 'Kullanıcı deneyimini iyileştirmek ve göz yorgunluğunu azaltmak',
                    kullanicilar: 'Tüm mobil uygulama kullanıcıları',
                    kriterler: '1. Sistem temasını takip et\n2. Manuel toggle\n3. Tercih kaydetme\n4. Smooth transitions',
                    ui_ux: 'Settings sayfasında toggle switch, tüm ekranlarda tema desteği',
                    efor: 'küçük',
                    oncelik: 'dusuk',
                    etiketler: 'mobile,ui,enhancement'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                projectId: projectIds[3],
                degerler: {
                    baslik: 'Dashboard prototype',
                    aciklama: 'Interaktif analytics dashboard prototipi',
                    amac: 'Veri görselleştirme ve raporlama için kullanıcı dostu arayüz',
                    kullanicilar: 'Data analysts, managers, executives',
                    kriterler: '1. Real-time updates\n2. Customizable widgets\n3. Export functionality\n4. Mobile responsive',
                    ui_ux: 'Figma\'da interaktif prototype, drag-drop widget support',
                    ilgili: 'Data visualization, reporting module',
                    efor: 'büyük',
                    oncelik: 'yuksek',
                    etiketler: 'design,analytics,prototype'
                }
            }
        ];

        // Teknik Borç örnekleri
        const techDebtTasks = [
            {
                templateId: this.TEMPLATE_IDS.TEKNIK_BORC,
                projectId: projectIds[2],
                degerler: {
                    baslik: 'Redis cache layer implementasyonu',
                    aciklama: 'API performansını artırmak için Redis cache katmanı eklenecek',
                    alan: 'Backend/Cache',
                    dosyalar: 'src/services/*, src/middleware/cache.js',
                    neden: 'Database query\'leri yavaş, response time\'lar yüksek',
                    analiz: 'Ortalama response time 800ms, hedef 200ms altı',
                    cozum: 'Redis ile frequently accessed data cache\'lenmeli',
                    riskler: 'Cache invalidation complexity, memory usage',
                    iyilestirmeler: '%70 performance improvement, reduced DB load',
                    sure: '1 hafta',
                    oncelik: 'orta',
                    etiketler: 'backend,performance,redis'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.TEKNIK_BORC,
                projectId: projectIds[2],
                degerler: {
                    baslik: 'API rate limiting',
                    aciklama: 'DDoS ve abuse prevention için rate limiting',
                    alan: 'Backend/Security',
                    dosyalar: 'src/middleware/rateLimiter.js',
                    neden: 'API güvenliği ve resource protection gerekli',
                    analiz: 'Mevcut durumda rate limiting yok, abuse riski var',
                    cozum: 'Redis-based rate limiting with sliding window',
                    riskler: 'Legitimate user impact, configuration complexity',
                    iyilestirmeler: 'Enhanced security, predictable resource usage',
                    sure: '2-3 gün',
                    oncelik: 'yuksek',
                    etiketler: 'backend,security,performance'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.TEKNIK_BORC,
                projectId: projectIds[3],
                degerler: {
                    baslik: 'ETL pipeline modernizasyonu',
                    aciklama: 'Legacy ETL scripts\'leri Apache Airflow\'a taşınacak',
                    alan: 'Data Engineering',
                    dosyalar: 'etl/*, airflow/dags/*',
                    neden: 'Mevcut cron-based system unreliable ve monitoring yok',
                    analiz: 'Manuel interventions, no retry logic, poor visibility',
                    cozum: 'Airflow DAGs with proper monitoring ve alerting',
                    riskler: 'Data pipeline downtime during migration',
                    iyilestirmeler: 'Better reliability, monitoring, scalability',
                    sure: '2+ hafta',
                    oncelik: 'yuksek',
                    etiketler: 'data,backend,infrastructure'
                }
            }
        ];

        // Araştırma Görevi örnekleri
        const researchTasks = [
            {
                templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                projectId: projectIds[0],
                degerler: {
                    konu: 'Next.js 14 App Router migration',
                    amac: 'Mevcut Pages Router\'dan App Router\'a geçiş feasibility',
                    sorular: '1. Performance improvements?\n2. Migration effort?\n3. Breaking changes?\n4. Team learning curve?',
                    kaynaklar: 'Next.js docs, Vercel blog, migration guides',
                    alternatifler: 'Stay with Pages Router, Remix, SvelteKit',
                    kriterler: 'Performance, DX, SEO, bundle size, community support',
                    son_tarih: this.getDateString(14),
                    oncelik: 'orta',
                    etiketler: 'research,frontend,nextjs'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                projectId: projectIds[3],
                degerler: {
                    konu: 'Chart library comparison',
                    amac: 'Dashboard için en uygun chart library seçimi',
                    sorular: '1. Performance with large datasets?\n2. Customization options?\n3. Mobile support?\n4. Bundle size?',
                    kaynaklar: 'Chart.js, D3.js, ApexCharts, Recharts docs',
                    alternatifler: 'Chart.js, D3.js, ApexCharts, Recharts, Victory',
                    kriterler: 'Performance, features, size, ease of use, community',
                    oncelik: 'orta',
                    etiketler: 'research,frontend,visualization'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                projectId: projectIds[4],
                degerler: {
                    konu: 'OWASP Top 10 security audit',
                    amac: 'Mevcut sistemin güvenlik açıklarını tespit etmek',
                    sorular: '1. Which vulnerabilities exist?\n2. Risk levels?\n3. Remediation effort?\n4. Priority order?',
                    kaynaklar: 'OWASP guides, security scanners, pen test tools',
                    alternatifler: 'Manual audit, automated scanning, external pen test',
                    kriterler: 'Coverage, accuracy, actionability, cost',
                    son_tarih: this.getDateString(7),
                    oncelik: 'yuksek',
                    etiketler: 'security,testing,critical'
                }
            }
        ];

        // Non-template tasks (meetings, reviews, etc.)
        const directTasks = [
            {
                baslik: 'Team meeting hazırlığı',
                aciklama: 'Haftalık geliştirici toplantısı için sunum hazırla',
                oncelik: GorevOncelik.Orta,
                son_tarih: this.getDateString(1),
                etiketler: 'meeting,planning'
            },
            {
                baslik: 'Code review yapılacak PR\'lar',
                aciklama: '5 adet bekleyen pull request incelenecek',
                oncelik: GorevOncelik.Yuksek,
                son_tarih: this.getDateString(0),
                etiketler: 'review,git,urgent'
            },
            {
                baslik: 'Teknik blog yazısı',
                aciklama: 'Microservices best practices hakkında blog yazısı',
                oncelik: GorevOncelik.Dusuk,
                etiketler: 'writing,documentation'
            }
        ];

        // Create tasks from templates
        for (const task of [...bugTasks, ...featureTasks, ...techDebtTasks, ...researchTasks]) {
            try {
                const result = await this.apiClient.callTool('templateden_gorev_olustur', {
                    template_id: task.templateId,
                    degerler: task.degerler
                });
                
                const responseText = result.content[0].text;
                Logger.debug(`Template task creation response:`, responseText);
                
                // Extract task ID
                const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                if (idMatch) {
                    const taskId = idMatch[1];
                    taskIds.push(taskId);
                    
                    // Assign to project if specified
                    if (task.projectId) {
                        try {
                            await this.apiClient.callTool('gorev_duzenle', {
                                id: taskId,
                                proje_id: task.projectId
                            });
                            Logger.info(`Assigned task to project: ${task.projectId}`);
                        } catch (error) {
                            Logger.error(`Failed to assign task to project:`, error);
                        }
                    }
                }
            } catch (error) {
                Logger.error(`Failed to create template task:`, error);
            }
        }

        // Create direct tasks (without templates)
        for (const task of directTasks) {
            try {
                const result = await this.apiClient.callTool('gorev_olustur', task);
                const responseText = result.content[0].text;
                
                const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                if (idMatch) {
                    taskIds.push(idMatch[1]);
                }
            } catch (error) {
                Logger.error(`Failed to create direct task:`, error);
            }
        }

        return taskIds;
    }

    /**
     * Test bağımlılıkları oluştur
     */
    private async createTestDependencies(taskIds: string[]): Promise<void> {
        const dependencies = [
            { kaynak: 0, hedef: 1, tip: 'blocks' },
            { kaynak: 1, hedef: 2, tip: 'blocks' },
            { kaynak: 5, hedef: 6, tip: 'depends_on' },
            { kaynak: 8, hedef: 9, tip: 'blocks' },
        ];

        for (const dep of dependencies) {
            if (taskIds[dep.kaynak] && taskIds[dep.hedef]) {
                try {
                    await this.apiClient.callTool('gorev_bagimlilik_ekle', {
                        kaynak_id: taskIds[dep.kaynak],
                        hedef_id: taskIds[dep.hedef],
                        baglanti_tipi: dep.tip
                    });
                } catch (error) {
                    Logger.error('Failed to create dependency:', error);
                }
            }
        }
    }

    /**
     * Alt görevler oluştur
     */
    private async createSubtasks(parentTaskIds: string[]): Promise<void> {
        // Ana sayfa tasarımı için alt görevler
        if (parentTaskIds[3]) { // Feature task for ana sayfa
            const subtasks = [
                {
                    parent_id: parentTaskIds[3],
                    baslik: 'Hero section mockup',
                    aciklama: 'Ana sayfa hero bölümü için Figma mockup hazırla',
                    oncelik: GorevOncelik.Yuksek,
                    etiketler: 'design,ui,mockup'
                },
                {
                    parent_id: parentTaskIds[3],
                    baslik: 'Responsive grid sistemi',
                    aciklama: 'Bootstrap 5 veya Tailwind CSS ile responsive grid',
                    oncelik: GorevOncelik.Orta,
                    etiketler: 'frontend,css,responsive'
                }
            ];

            for (const subtask of subtasks) {
                try {
                    await this.apiClient.callTool('gorev_altgorev_olustur', subtask);
                } catch (error) {
                    Logger.error('Failed to create subtask:', error);
                }
            }
        }
    }

    /**
     * Bazı görevlerin durumlarını güncelle
     */
    private async updateSomeTaskStatuses(taskIds: string[]): Promise<void> {
        // Bazı görevleri "devam ediyor" yap
        const inProgressTasks = [1, 4, 7];
        for (const index of inProgressTasks) {
            if (taskIds[index]) {
                try {
                    await this.apiClient.callTool('gorev_guncelle', {
                        id: taskIds[index],
                        durum: GorevDurum.DevamEdiyor
                    });
                } catch (error) {
                    Logger.error('Failed to update task status:', error);
                }
            }
        }

        // Bazı görevleri tamamla
        const completedTasks = [2, 5];
        for (const index of completedTasks) {
            if (taskIds[index]) {
                try {
                    await this.apiClient.callTool('gorev_guncelle', {
                        id: taskIds[index],
                        durum: GorevDurum.Tamamlandi
                    });
                } catch (error) {
                    Logger.error('Failed to update task status:', error);
                }
            }
        }
    }

    /**
     * AI context ve interaksiyonları oluştur
     */
    private async setupAIContext(taskIds: string[]): Promise<void> {
        if (taskIds.length > 0) {
            try {
                // İlk görevi aktif yap
                await this.apiClient.callTool('gorev_set_active', {
                    task_id: taskIds[0]
                });
                Logger.info('Set active task for AI context');

                // Context summary al
                const contextSummary = await this.apiClient.callTool('gorev_context_summary', {});
                Logger.info('Generated AI context summary');
            } catch (error) {
                Logger.error('Failed to setup AI context:', error);
            }
        }
    }

    /**
     * Bugünden itibaren belirtilen gün sayısı kadar sonraki tarihi döndür
     */
    private getDateString(daysFromNow: number): string {
        const date = new Date();
        date.setDate(date.getDate() + daysFromNow);
        return date.toISOString().split('T')[0];
    }
}
