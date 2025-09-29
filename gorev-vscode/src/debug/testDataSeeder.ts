import { ClientInterface } from '../interfaces/client';
import { t } from '../utils/l10n';
import { GorevDurum, GorevOncelik } from '../models/common';
import { Logger } from '../utils/logger';
import * as vscode from 'vscode';

/**
 * Debug için zengin test verileri oluşturur
 * Updated to use templates instead of deprecated gorev_olustur
 */
export class TestDataSeeder {
    // Template IDs from database
    private readonly TEMPLATE_IDS = {
        BUG_RAPORU: '4dd56a2a-caf4-472c-8c0f-276bc8a1f880',
        OZELLIK_ISTEGI: '6b083358-9c4d-4f4e-b041-9288c05a1bb7',
        TEKNIK_BORC: '69e2b237-7c2e-4459-9d46-ea6c05aba39a',
        ARASTIRMA_GOREVI: '13f04fe2-b5b6-4fd6-8684-5eca5dc2770d'
    };

    constructor(private mcpClient: ClientInterface) {}

    /**
     * Test verilerini oluştur
     */
    async seedTestData(): Promise<void> {
        const result = await vscode.window.showInformationMessage(
            t('testData.confirmSeed'),
            t('testData.yesCreate'),
            t('testData.no')
        );

        if (result !== t('testData.yesCreate')) {
            return;
        }

        try {
            await vscode.window.withProgress({
                location: vscode.ProgressLocation.Notification,
                title: t('testData.creating'),
                cancellable: false
            }, async (progress) => {
                // 1. Test projeleri oluştur
                progress.report({ increment: 10, message: t('testData.creatingProjects') });
                const projectIds = await this.createTestProjects();

                // 2. Test görevleri oluştur
                progress.report({ increment: 20, message: t('testData.creatingTasks') });
                const taskIds = await this.createTestTasks(projectIds);

                // 3. Pagination test için çok sayıda görev oluştur
                progress.report({ increment: 15, message: 'Creating tasks for pagination testing...' });
                const paginationTaskIds = await this.createPaginationTestTasks(projectIds);
                taskIds.push(...paginationTaskIds);

                // 4. Hierarchy test görevleri oluştur
                progress.report({ increment: 10, message: 'Creating hierarchy test tasks...' });
                const hierarchyTaskIds = await this.createHierarchyTestTasks(projectIds);
                taskIds.push(...hierarchyTaskIds);

                // 5. Bağımlılıklar oluştur
                progress.report({ increment: 10, message: t('testData.creatingDependencies') });
                await this.createTestDependencies(taskIds);

                // 6. Alt görevler oluştur
                progress.report({ increment: 10, message: t('testData.creatingSubtasks') });
                await this.createSubtasks(taskIds);

                // 7. Extra template görevler oluştur (örnekler için)
                progress.report({ increment: 10, message: t('testData.creatingExamples') });
                await this.createAdditionalTemplateExamples(projectIds);

                // 8. Bazı görevleri tamamla ve AI interaksiyonları ekle
                progress.report({ increment: 5, message: t('testData.updatingStatuses') });
                await this.updateSomeTaskStatuses(taskIds);

                // 9. AI context oluştur
                progress.report({ increment: 5, message: t('testData.creatingAIContext') });
                await this.setupAIContext(taskIds);

                progress.report({ increment: 10, message: t('testData.completed') });
            });

            vscode.window.showInformationMessage(t('testData.success'));
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : String(error);
            vscode.window.showErrorMessage(t('testData.failed', errorMessage));
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
                const result = await this.mcpClient.callTool('proje_olustur', project);
                // ID'yi response'tan çıkar - daha geniş bir regex kullan
                const responseText = result.content[0].text;
                Logger.debug('Project creation response:', responseText);
                
                // UUID formatında ID ara
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
            await this.mcpClient.callTool('proje_aktif_yap', { proje_id: projectIds[0] });
        }

        return projectIds;
    }

    /**
     * Test görevleri oluştur - template-based approach
     */
    private async createTestTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];

        // Create Bug Report tasks using template
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
                projectId: projectIds[1],
                degerler: {
                    baslik: 'Push notification Android\'de çalışmıyor',
                    aciklama: 'Firebase Cloud Messaging entegrasyonu iOS\'ta çalışıyor ama Android\'de bildirimler gelmiyor',
                    modul: 'Mobile/Notifications',
                    ortam: 'production',
                    adimlar: '1. Android cihazda uygulamayı aç\n2. Bildirim izni ver\n3. Test bildirimi gönder',
                    beklenen: 'Push notification alınmalı',
                    mevcut: 'Bildirimler Android\'de alınmıyor',
                    oncelik: 'yuksek',
                    etiketler: 'mobile,bug,firebase'
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
                    adimlar: '1. HTTPS ile siteye git\n2. Sertifika uyarısı görünüyor',
                    beklenen: 'Valid SSL sertifikası',
                    mevcut: 'NET::ERR_CERT_DATE_INVALID',
                    oncelik: 'yuksek',
                    etiketler: 'security,infrastructure,urgent'
                }
            }
        ];

        // Create Feature Request tasks using template
        const featureTasks = [
            {
                templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                projectId: projectIds[0],
                degerler: {
                    baslik: 'Ana sayfa tasarımını tamamla',
                    aciklama: 'Modern ve responsive ana sayfa tasarımı yapılacak',
                    amac: 'Kullanıcı deneyimini iyileştirmek ve dönüşüm oranını artırmak',
                    kullanicilar: 'Tüm web sitesi ziyaretçileri',
                    kriterler: '1. Hero section\n2. Özellikler bölümü\n3. Footer\n4. Mobile responsive',
                    ui_ux: 'Modern, minimal tasarım. Hızlı yükleme.',
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
                    amac: 'Güvenli kullanıcı kimlik doğrulama sistemi',
                    kullanicilar: 'Kayıtlı kullanıcılar',
                    kriterler: '1. JWT token\n2. Login/Register/Forgot password\n3. Session management',
                    ui_ux: 'Clean login forms, social login options',
                    efor: 'büyük',
                    oncelik: 'yuksek',
                    etiketler: 'backend,security,feature'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                projectId: projectIds[1],
                degerler: {
                    baslik: 'Dark mode tema',
                    aciklama: 'Sistem ayarlarına göre otomatik tema değişimi',
                    amac: 'Kullanıcı deneyimini iyileştirmek ve göz yorgunluğunu azaltmak',
                    kullanicilar: 'Tüm mobil uygulama kullanıcıları',
                    kriterler: '1. Sistem temasını takip et\n2. Manuel toggle\n3. Tercih kaydetme',
                    ui_ux: 'Settings sayfasında toggle switch',
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
                    aciklama: 'Figma\'da interaktif dashboard prototipi',
                    amac: 'Veri görselleştirme için kullanıcı dostu arayüz',
                    kullanicilar: 'Data analysts, managers',
                    kriterler: '1. Real-time updates\n2. Customizable widgets\n3. Export',
                    ui_ux: 'Clean, modern dashboard with drag-drop',
                    efor: 'büyük',
                    oncelik: 'yuksek',
                    etiketler: 'design,analytics,prototype'
                }
            }
        ];

        // Create Technical Debt tasks using template
        const techDebtTasks = [
            {
                templateId: this.TEMPLATE_IDS.TEKNIK_BORC,
                projectId: projectIds[2],
                degerler: {
                    baslik: 'Redis cache entegrasyonu',
                    aciklama: 'Performans artışı için Redis cache katmanı',
                    alan: 'Backend/Performance',
                    neden: 'API response time\'ları yüksek',
                    analiz: 'Ortalama response time 800ms',
                    cozum: 'Redis ile cache layer ekle',
                    iyilestirmeler: '%70 performans artışı bekleniyor',
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
                    aciklama: 'API güvenliği için rate limiting',
                    alan: 'Backend/Security',
                    neden: 'DDoS ve abuse riski var',
                    analiz: 'Rate limiting yok',
                    cozum: 'Redis-based rate limiting ekle',
                    iyilestirmeler: 'Güvenlik artışı',
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
                    aciklama: 'Apache Airflow ile veri pipeline',
                    alan: 'Data Engineering',
                    neden: 'Mevcut cron jobs unreliable',
                    analiz: 'Manuel müdahale gerekiyor',
                    cozum: 'Airflow DAGs ile otomatize et',
                    iyilestirmeler: 'Better reliability ve monitoring',
                    sure: '2+ hafta',
                    oncelik: 'yuksek',
                    etiketler: 'data,backend,infrastructure'
                }
            }
        ];

        // Create Research tasks using template
        const researchTasks = [
            {
                templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                projectId: projectIds[0],
                degerler: {
                    konu: 'SEO optimizasyonu',
                    amac: 'Web sitesi SEO performansını artırmak',
                    sorular: '1. Core Web Vitals?\n2. Meta tags?\n3. Sitemap?',
                    kriterler: 'Google PageSpeed, Lighthouse scores',
                    son_tarih: this.getDateString(14),
                    oncelik: 'orta',
                    etiketler: 'seo,performance,research'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                projectId: projectIds[4],
                degerler: {
                    konu: 'Penetrasyon testi metodolojisi',
                    amac: 'OWASP Top 10 güvenlik testleri',
                    sorular: '1. Hangi açıklar var?\n2. Risk seviyeleri?\n3. Çözüm önerileri?',
                    kriterler: 'OWASP guidelines, security best practices',
                    son_tarih: this.getDateString(7),
                    oncelik: 'yuksek',
                    etiketler: 'security,testing,critical'
                }
            }
        ];

        // Additional tasks using general templates
        const additionalTasks = [
            {
                templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                projectId: projectIds[0],
                degerler: {
                    konu: 'Team meeting hazırlığı',
                    amac: 'Haftalık geliştirici toplantısı için sunum hazırla',
                    sorular: '1. Hangi konular ele alınacak?\n2. Sprint progress nasıl?\n3. Blockerlar neler?',
                    kriterler: 'Toplantı sunumu hazır ve agenda belirlenmiş',
                    son_tarih: this.getDateString(1),
                    oncelik: 'orta',
                    etiketler: 'meeting,planning'
                }
            },
            {
                templateId: this.TEMPLATE_IDS.TEKNIK_BORC,
                projectId: projectIds[0],
                degerler: {
                    baslik: 'Code review yapılacak PR\'lar',
                    aciklama: '5 adet bekleyen pull request incelenecek',
                    alan: 'Code Quality',
                    neden: 'Kod kalitesi ve takım standartları',
                    analiz: 'Bekleyen PR\'ların durumu',
                    cozum: 'Systematic code review process',
                    iyilestirmeler: 'Better code quality and team alignment',
                    sure: '1 gün',
                    oncelik: 'yuksek',
                    etiketler: 'review,git,urgent'
                }
            }
        ];

        // Process template-based tasks
        for (const taskGroup of [bugTasks, featureTasks, techDebtTasks, researchTasks, additionalTasks]) {
            for (const task of taskGroup) {
                try {
                    const result = await this.mcpClient.callTool('templateden_gorev_olustur', {
                        template_id: task.templateId,
                        degerler: task.degerler
                    });
                    
                    const responseText = result.content[0].text;
                    Logger.debug(`Template task creation response:`, responseText);
                    
                    const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                    if (idMatch) {
                        const taskId = idMatch[1];
                        taskIds.push(taskId);
                        
                        // Assign to project if specified
                        if (task.projectId) {
                            try {
                                await this.mcpClient.callTool('gorev_duzenle', {
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
        }

        // All tasks now use templates - no direct creation needed

        return taskIds;
    }

    /**
     * Test bağımlılıkları oluştur
     */
    private async createTestDependencies(taskIds: string[]): Promise<void> {
        // Örnek bağımlılıklar
        const dependencies = [
            { kaynak: 0, hedef: 4, tip: 'blocks' }, // Login bug -> User auth feature
            { kaynak: 4, hedef: 5, tip: 'blocks' }, // User auth -> Other features
            { kaynak: 7, hedef: 8, tip: 'depends_on' }, // Rate limiting depends on Redis
        ];

        for (const dep of dependencies) {
            if (taskIds[dep.kaynak] && taskIds[dep.hedef]) {
                try {
                    await this.mcpClient.callTool('gorev_bagimlilik_ekle', {
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
     * Bazı görevlerin durumlarını güncelle
     */
    private async updateSomeTaskStatuses(taskIds: string[]): Promise<void> {
        // Bazı görevleri "devam ediyor" yap
        const inProgressTasks = [1, 4, 7];
        for (const index of inProgressTasks) {
            if (taskIds[index]) {
                try {
                    await this.mcpClient.callTool('gorev_guncelle', {
                        id: taskIds[index],
                        durum: GorevDurum.DevamEdiyor
                    });
                } catch (error) {
                    Logger.error('Failed to update task status:', error);
                }
            }
        }

        // Bazı görevleri tamamla
        const completedTasks = [2, 10];
        for (const index of completedTasks) {
            if (taskIds[index]) {
                try {
                    await this.mcpClient.callTool('gorev_guncelle', {
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
     * Alt görevler oluştur
     */
    private async createSubtasks(parentTaskIds: string[]): Promise<void> {
        // Ana sayfa tasarımı için alt görevler (assuming it's the 3rd feature task)
        if (parentTaskIds[3]) {
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
                },
                {
                    parent_id: parentTaskIds[3],
                    baslik: 'Animation ve transitions',
                    aciklama: 'Smooth scroll ve hover effect animasyonları',
                    oncelik: GorevOncelik.Dusuk,
                    etiketler: 'frontend,animation,ux'
                }
            ];

            for (const subtask of subtasks) {
                try {
                    const result = await this.mcpClient.callTool('gorev_altgorev_olustur', subtask);
                    
                    // İkinci seviye alt görev ekle
                    if (subtask.baslik === 'Hero section mockup') {
                        const responseText = result.content[0].text;
                        const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                        if (idMatch) {
                            await this.mcpClient.callTool('gorev_altgorev_olustur', {
                                parent_id: idMatch[1],
                                baslik: 'Color palette seçimi',
                                aciklama: 'Brand guidelines\'a uygun renk paleti',
                                oncelik: GorevOncelik.Yuksek,
                                etiketler: 'design,branding'
                            });
                        }
                    }
                } catch (error) {
                    Logger.error('Failed to create subtask:', error);
                }
            }
        }

        // Login sistemi için alt görevler (assuming it's the 4th feature task)
        if (parentTaskIds[4]) {
            const subtasks = [
                {
                    parent_id: parentTaskIds[4],
                    baslik: 'JWT token implementasyonu',
                    aciklama: 'Access ve refresh token yönetimi',
                    oncelik: GorevOncelik.Yuksek,
                    son_tarih: this.getDateString(3),
                    etiketler: 'backend,security,jwt'
                },
                {
                    parent_id: parentTaskIds[4],
                    baslik: 'Password reset flow',
                    aciklama: 'Email ile şifre sıfırlama akışı',
                    oncelik: GorevOncelik.Orta,
                    etiketler: 'backend,email,feature'
                }
            ];

            for (const subtask of subtasks) {
                try {
                    await this.mcpClient.callTool('gorev_altgorev_olustur', subtask);
                } catch (error) {
                    Logger.error('Failed to create subtask:', error);
                }
            }
        }
    }

    /**
     * Template'lerden extra görevler oluştur
     */
    private async createAdditionalTemplateExamples(projectIds: string[]): Promise<void> {
        try {
            // Example of using different template variations
            const additionalTasks = [
                {
                    templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                    degerler: {
                        baslik: 'GraphQL API endpoint',
                        aciklama: 'REST API yanında GraphQL desteği',
                        amac: 'Flexible data fetching için GraphQL layer',
                        kullanicilar: 'Frontend developers, API consumers',
                        kriterler: '1. Schema definition\n2. Resolvers\n3. Subscriptions',
                        efor: 'büyük',
                        oncelik: 'orta',
                        etiketler: 'backend,api,feature'
                    }
                },
                {
                    templateId: this.TEMPLATE_IDS.ARASTIRMA_GOREVI,
                    degerler: {
                        konu: 'Makine öğrenmesi için framework',
                        amac: 'Churn prediction modeli için ML framework seçimi',
                        sorular: '1. TensorFlow vs PyTorch?\n2. Deployment options?\n3. Performance?',
                        kriterler: 'Ease of use, community, deployment',
                        oncelik: 'orta',
                        etiketler: 'ml,data-science,research'
                    }
                }
            ];

            for (const task of additionalTasks) {
                try {
                    const result = await this.mcpClient.callTool('templateden_gorev_olustur', {
                        template_id: task.templateId,
                        degerler: task.degerler
                    });
                    
                    // Assign to appropriate project if needed
                    if (projectIds.length > 2) {
                        const responseText = result.content[0].text;
                        const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                        if (idMatch) {
                            await this.mcpClient.callTool('gorev_duzenle', {
                                id: idMatch[1],
                                proje_id: projectIds[2] // Backend API project
                            });
                        }
                    }
                } catch (error) {
                    Logger.error('Failed to create additional template task:', error);
                }
            }
        } catch (error) {
            Logger.error('Failed in createAdditionalTemplateExamples:', error);
        }
    }

    /**
     * AI context ve interaksiyonları oluştur
     */
    private async setupAIContext(taskIds: string[]): Promise<void> {
        // Bazı görevleri AI için aktif yap
        if (taskIds.length > 0) {
            try {
                // İlk görevi aktif yap
                await this.mcpClient.callTool('gorev_set_active', {
                    task_id: taskIds[0]
                });
                Logger.info('Set active task for AI context');

                // Doğal dil sorgusu test et
                const nlpResults = [
                    await this.mcpClient.callTool('gorev_nlp_query', { query: 'bugün yapılacak görevler' }),
                    await this.mcpClient.callTool('gorev_nlp_query', { query: 'yüksek öncelikli görevler' }),
                    await this.mcpClient.callTool('gorev_nlp_query', { query: 'etiket:bug' })
                ];

                Logger.info('Tested NLP queries');

                // Context summary al
                const contextSummary = await this.mcpClient.callTool('gorev_context_summary', {});
                Logger.info('Generated AI context summary');

                // Batch update test et - bazı görevlerin durumunu toplu güncelle
                if (taskIds.length > 5) {
                    await this.mcpClient.callTool('gorev_batch_update', {
                        updates: [
                            { id: taskIds[2], updates: { durum: 'devam_ediyor' } },
                            { id: taskIds[3], updates: { durum: 'devam_ediyor' } }
                        ]
                    });
                    Logger.info('Performed batch update');
                }
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

    /**
     * Test verilerini temizle
     */
    async clearTestData(): Promise<void> {
        const result = await vscode.window.showWarningMessage(
            t('testData.confirmClear'),
            t('testData.yesDelete'),
            t('testData.no')
        );

        if (result !== t('testData.yesDelete')) {
            return;
        }

        try {
            // Önce tüm görevleri listele ve sil
            const tasksResult = await this.mcpClient.callTool('gorev_listele', {
                tum_projeler: true
            });

            // Parse task IDs from response
            const taskIdMatches = tasksResult.content[0].text.matchAll(/ID: ([a-f0-9-]+)/g);
            const taskIds: string[] = [];
            for (const match of taskIdMatches) {
                if (match[1]) {
                    taskIds.push(match[1]);
                }
            }

            for (const taskId of taskIds) {
                try {
                    await this.mcpClient.callTool('gorev_sil', {
                        id: taskId,
                        onay: true
                    });
                } catch (error) {
                    Logger.error('Failed to delete task:', error);
                }
            }

            vscode.window.showInformationMessage(t('testData.cleared'));
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : String(error);
            vscode.window.showErrorMessage(t('testData.clearFailed', errorMessage));
            Logger.error('Failed to clear test data:', error);
        }
    }

    /**
     * Create large number of tasks for pagination testing
     */
    private async createPaginationTestTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];
        const targetTaskCount = 150; // To test pagination (pageSize is 100)

        Logger.info(`[TestDataSeeder] Creating ${targetTaskCount} tasks for pagination testing...`);

        for (let i = 0; i < targetTaskCount; i++) {
            try {
                const taskData = {
                    templateId: this.TEMPLATE_IDS.BUG_RAPORU,
                    projectId: projectIds[i % projectIds.length],
                    degerler: {
                        baslik: `Pagination Test Task ${i + 1}`,
                        aciklama: `This is test task #${i + 1} created for pagination testing. It helps verify that the VS Code extension can handle large numbers of tasks properly.`,
                        oncelik: ['dusuk', 'orta', 'yuksek'][i % 3],
                        etiketler: [`pagination-test`, `batch-${Math.floor(i / 50) + 1}`]
                    }
                };

                const result = await this.mcpClient.callTool('templateden_gorev_olustur', taskData);
                const taskId = this.extractTaskId(result);
                if (taskId) {
                    taskIds.push(taskId);
                }

                // Log progress every 25 tasks
                if ((i + 1) % 25 === 0) {
                    Logger.info(`[TestDataSeeder] Created ${i + 1}/${targetTaskCount} pagination test tasks`);
                }
            } catch (error) {
                Logger.error(`Failed to create pagination test task ${i + 1}:`, error);
            }
        }

        Logger.info(`[TestDataSeeder] Successfully created ${taskIds.length} pagination test tasks`);
        return taskIds;
    }

    /**
     * Create tasks with complex hierarchy for hierarchy testing
     */
    private async createHierarchyTestTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];

        Logger.info('[TestDataSeeder] Creating hierarchy test tasks...');

        // Create parent tasks first
        const parentTasks = [
            {
                baslik: 'Feature: User Management System',
                aciklama: 'Complete user management system with authentication and authorization',
                oncelik: 'yuksek'
            },
            {
                baslik: 'Feature: Reporting Dashboard',
                aciklama: 'Advanced reporting dashboard with charts and export functionality',
                oncelik: 'orta'
            },
            {
                baslik: 'Infrastructure: Database Migration',
                aciklama: 'Migrate from old database schema to new optimized structure',
                oncelik: 'yuksek'
            }
        ];

        const parentTaskIds: string[] = [];

        for (const parentTaskData of parentTasks) {
            try {
                const taskData = {
                    templateId: this.TEMPLATE_IDS.OZELLIK_ISTEGI,
                    projectId: projectIds[0],
                    degerler: {
                        ...parentTaskData,
                        etiketler: ['hierarchy-test', 'parent-task']
                    }
                };

                const result = await this.mcpClient.callTool('templateden_gorev_olustur', taskData);
                const taskId = this.extractTaskId(result);
                if (taskId) {
                    parentTaskIds.push(taskId);
                    taskIds.push(taskId);
                    Logger.info(`Created parent task: ${parentTaskData.baslik} (${taskId})`);
                }
            } catch (error) {
                Logger.error(`Failed to create parent task: ${parentTaskData.baslik}`, error);
            }
        }

        // Create subtasks for each parent
        const subtaskTemplates = [
            { baslik: 'Design Phase', aciklama: 'Design user interface and user experience' },
            { baslik: 'Backend Implementation', aciklama: 'Implement backend API and business logic' },
            { baslik: 'Frontend Implementation', aciklama: 'Create frontend components and pages' },
            { baslik: 'Testing Phase', aciklama: 'Write and execute comprehensive tests' },
            { baslik: 'Documentation', aciklama: 'Create user and technical documentation' }
        ];

        for (let i = 0; i < parentTaskIds.length; i++) {
            const parentId = parentTaskIds[i];

            for (const subtaskTemplate of subtaskTemplates) {
                try {
                    const taskData = {
                        templateId: this.TEMPLATE_IDS.TEKNIK_BORC,
                        projectId: projectIds[0],
                        degerler: {
                            baslik: `${subtaskTemplate.baslik} (Parent ${i + 1})`,
                            aciklama: subtaskTemplate.aciklama,
                            oncelik: 'orta',
                            etiketler: ['hierarchy-test', 'subtask', `parent-${i + 1}`]
                        }
                    };

                    const result = await this.mcpClient.callTool('templateden_gorev_olustur', taskData);
                    const taskId = this.extractTaskId(result);
                    if (taskId) {
                        taskIds.push(taskId);

                        // Set parent relationship
                        await this.mcpClient.callTool('gorev_ust_gorev_degistir', {
                            gorev_id: taskId,
                            ust_gorev_id: parentId
                        });

                        Logger.debug(`Created subtask: ${taskData.degerler.baslik} under parent ${parentId}`);
                    }
                } catch (error) {
                    Logger.error(`Failed to create subtask: ${subtaskTemplate.baslik}`, error);
                }
            }
        }

        Logger.info(`[TestDataSeeder] Successfully created ${taskIds.length} hierarchy test tasks (${parentTaskIds.length} parents)`);
        return taskIds;
    }

    /**
     * Extract task ID from MCP response
     */
    private extractTaskId(result: any): string | null {
        try {
            const responseText = result.content[0].text;
            const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
            return idMatch ? idMatch[1] : null;
        } catch (error) {
            Logger.error('Failed to extract task ID from response:', error);
            return null;
        }
    }
}
