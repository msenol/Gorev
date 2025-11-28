import { ApiClient, MCPToolResult } from '../api/client';
import { t } from '../utils/l10n';
import { GorevDurum, GorevOncelik } from '../models/common';
import { Logger } from '../utils/logger';
import * as vscode from 'vscode';

/**
 * Debug için zengin test verileri oluşturur
 * Updated to use templates instead of deprecated gorev_olustur
 * Template IDs are now dynamically fetched by alias
 */
export class TestDataSeeder {
    // Template IDs will be fetched dynamically
    private templateIds: Record<string, string> = {};

    constructor(private apiClient: ApiClient) {}

    /**
     * Fetch template IDs by alias from the API
     */
    private async fetchTemplateIds(): Promise<void> {
        try {
            const result = await this.apiClient.getTemplates();
            if (result.success && result.data) {
                for (const template of result.data) {
                    if (template.alias) {
                        this.templateIds[template.alias] = template.id;
                    }
                }
            }
            Logger.info('[TestDataSeeder] Fetched template IDs:', this.templateIds);
        } catch (error) {
            Logger.error('[TestDataSeeder] Failed to fetch template IDs:', error);
            throw new Error('Failed to fetch templates. Make sure the server is running.');
        }
    }

    /**
     * Get template ID by alias, throws if not found
     */
    private getTemplateId(alias: string): string {
        const id = this.templateIds[alias];
        if (!id) {
            throw new Error(`Template with alias '${alias}' not found. Available: ${Object.keys(this.templateIds).join(', ')}`);
        }
        return id;
    }

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
                // 0. Fetch template IDs first
                progress.report({ increment: 5, message: 'Fetching templates...' });
                await this.fetchTemplateIds();

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

                // 8. Bazı görevleri tamamla
                progress.report({ increment: 10, message: t('testData.updatingStatuses') });
                await this.updateSomeTaskStatuses(taskIds);

                // Note: AI context tools are MCP-only, skipped for VS Code extension testing

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
                const result = await this.apiClient.callTool('proje_olustur', project);
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
            await this.apiClient.callTool('aktif_proje_ayarla', { proje_id: projectIds[0] });
        }

        return projectIds;
    }

    /**
     * Test görevleri oluştur - template-based approach
     */
    private async createTestTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];

        // Create Bug Report tasks using template (English field names)
        const bugTasks = [
            {
                templateId: this.getTemplateId('bug'),
                projectId: projectIds[0],
                degerler: {
                    title: 'Login sayfası 404 hatası veriyor',
                    description: 'Production ortamında /login URL\'ine gittiğimizde 404 hatası alıyoruz',
                    module: 'Authentication',
                    environment: 'production',
                    steps: '1. Production URL\'ine git\n2. /login sayfasına git\n3. 404 hatası görünüyor',
                    expected: 'Login sayfası açılmalı',
                    actual: '404 Not Found hatası',
                    solution: 'Routing konfigürasyonu kontrol edilmeli',
                    priority: 'yuksek',
                    tags: 'bug,critical,production'
                }
            },
            {
                templateId: this.getTemplateId('bug'),
                projectId: projectIds[1],
                degerler: {
                    title: 'Push notification Android\'de çalışmıyor',
                    description: 'Firebase Cloud Messaging entegrasyonu iOS\'ta çalışıyor ama Android\'de bildirimler gelmiyor',
                    module: 'Mobile/Notifications',
                    environment: 'production',
                    steps: '1. Android cihazda uygulamayı aç\n2. Bildirim izni ver\n3. Test bildirimi gönder',
                    expected: 'Push notification alınmalı',
                    actual: 'Bildirimler Android\'de alınmıyor',
                    priority: 'yuksek',
                    tags: 'mobile,bug,firebase'
                }
            },
            {
                templateId: this.getTemplateId('bug'),
                projectId: projectIds[4],
                degerler: {
                    title: 'SSL sertifikası expire olmuş',
                    description: 'Ana domain ve subdomain\'lerde SSL sertifikası süresi dolmuş',
                    module: 'Infrastructure',
                    environment: 'production',
                    steps: '1. HTTPS ile siteye git\n2. Sertifika uyarısı görünüyor',
                    expected: 'Valid SSL sertifikası',
                    actual: 'NET::ERR_CERT_DATE_INVALID',
                    priority: 'yuksek',
                    tags: 'security,infrastructure,urgent'
                }
            }
        ];

        // Create Feature Request tasks using template (English field names)
        const featureTasks = [
            {
                templateId: this.getTemplateId('feature'),
                projectId: projectIds[0],
                degerler: {
                    title: 'Ana sayfa tasarımını tamamla',
                    description: 'Modern ve responsive ana sayfa tasarımı yapılacak',
                    purpose: 'Kullanıcı deneyimini iyileştirmek ve dönüşüm oranını artırmak',
                    users: 'Tüm web sitesi ziyaretçileri',
                    criteria: '1. Hero section\n2. Özellikler bölümü\n3. Footer\n4. Mobile responsive',
                    ui_ux: 'Modern, minimal tasarım. Hızlı yükleme.',
                    effort: 'orta',
                    priority: 'yuksek',
                    tags: 'design,frontend,urgent'
                }
            },
            {
                templateId: this.getTemplateId('feature'),
                projectId: projectIds[0],
                degerler: {
                    title: 'Kullanıcı giriş sistemi',
                    description: 'JWT tabanlı authentication sistemi kurulacak',
                    purpose: 'Güvenli kullanıcı kimlik doğrulama sistemi',
                    users: 'Kayıtlı kullanıcılar',
                    criteria: '1. JWT token\n2. Login/Register/Forgot password\n3. Session management',
                    ui_ux: 'Clean login forms, social login options',
                    effort: 'büyük',
                    priority: 'yuksek',
                    tags: 'backend,security,feature'
                }
            },
            {
                templateId: this.getTemplateId('feature'),
                projectId: projectIds[1],
                degerler: {
                    title: 'Dark mode tema',
                    description: 'Sistem ayarlarına göre otomatik tema değişimi',
                    purpose: 'Kullanıcı deneyimini iyileştirmek ve göz yorgunluğunu azaltmak',
                    users: 'Tüm mobil uygulama kullanıcıları',
                    criteria: '1. Sistem temasını takip et\n2. Manuel toggle\n3. Tercih kaydetme',
                    ui_ux: 'Settings sayfasında toggle switch',
                    effort: 'küçük',
                    priority: 'dusuk',
                    tags: 'mobile,ui,enhancement'
                }
            },
            {
                templateId: this.getTemplateId('feature'),
                projectId: projectIds[3],
                degerler: {
                    title: 'Dashboard prototype',
                    description: 'Figma\'da interaktif dashboard prototipi',
                    purpose: 'Veri görselleştirme için kullanıcı dostu arayüz',
                    users: 'Data analysts, managers',
                    criteria: '1. Real-time updates\n2. Customizable widgets\n3. Export',
                    ui_ux: 'Clean, modern dashboard with drag-drop',
                    effort: 'büyük',
                    priority: 'yuksek',
                    tags: 'design,analytics,prototype'
                }
            }
        ];

        // Create Technical Debt tasks using template (mixed English/Turkish field names per template)
        const techDebtTasks = [
            {
                templateId: this.getTemplateId('debt'),
                projectId: projectIds[2],
                degerler: {
                    title: 'Redis cache entegrasyonu',
                    description: 'Performans artışı için Redis cache katmanı',
                    alan: 'Backend/Performance',
                    neden: 'API response time\'ları yüksek',
                    analiz: 'Ortalama response time 800ms',
                    cozum: 'Redis ile cache layer ekle',
                    iyilestirmeler: '%70 performans artışı bekleniyor',
                    sure: '1 hafta',
                    priority: 'orta',
                    tags: 'backend,performance,redis'
                }
            },
            {
                templateId: this.getTemplateId('debt'),
                projectId: projectIds[2],
                degerler: {
                    title: 'API rate limiting',
                    description: 'API güvenliği için rate limiting',
                    alan: 'Backend/Security',
                    neden: 'DDoS ve abuse riski var',
                    analiz: 'Rate limiting yok',
                    cozum: 'Redis-based rate limiting ekle',
                    iyilestirmeler: 'Güvenlik artışı',
                    sure: '2-3 gün',
                    priority: 'yuksek',
                    tags: 'backend,security,performance'
                }
            },
            {
                templateId: this.getTemplateId('debt'),
                projectId: projectIds[3],
                degerler: {
                    title: 'ETL pipeline modernizasyonu',
                    description: 'Apache Airflow ile veri pipeline',
                    alan: 'Data Engineering',
                    neden: 'Mevcut cron jobs unreliable',
                    analiz: 'Manuel müdahale gerekiyor',
                    cozum: 'Airflow DAGs ile otomatize et',
                    iyilestirmeler: 'Better reliability ve monitoring',
                    sure: '2+ hafta',
                    priority: 'yuksek',
                    tags: 'data,backend,infrastructure'
                }
            }
        ];

        // Create Research tasks using template (English field names)
        const researchTasks = [
            {
                templateId: this.getTemplateId('research'),
                projectId: projectIds[0],
                degerler: {
                    topic: 'SEO optimizasyonu',
                    purpose: 'Web sitesi SEO performansını artırmak',
                    questions: '1. Core Web Vitals?\n2. Meta tags?\n3. Sitemap?',
                    criteria: 'Google PageSpeed, Lighthouse scores',
                    due_date: this.getDateString(14),
                    priority: 'orta',
                    tags: 'seo,performance,research'
                }
            },
            {
                templateId: this.getTemplateId('research'),
                projectId: projectIds[4],
                degerler: {
                    topic: 'Penetrasyon testi metodolojisi',
                    purpose: 'OWASP Top 10 güvenlik testleri',
                    questions: '1. Hangi açıklar var?\n2. Risk seviyeleri?\n3. Çözüm önerileri?',
                    criteria: 'OWASP guidelines, security best practices',
                    due_date: this.getDateString(7),
                    priority: 'yuksek',
                    tags: 'security,testing,critical'
                }
            }
        ];

        // Additional tasks using general templates (English field names)
        const additionalTasks = [
            {
                templateId: this.getTemplateId('research'),
                projectId: projectIds[0],
                degerler: {
                    topic: 'Team meeting hazırlığı',
                    purpose: 'Haftalık geliştirici toplantısı için sunum hazırla',
                    questions: '1. Hangi konular ele alınacak?\n2. Sprint progress nasıl?\n3. Blockerlar neler?',
                    criteria: 'Toplantı sunumu hazır ve agenda belirlenmiş',
                    due_date: this.getDateString(1),
                    priority: 'orta',
                    tags: 'meeting,planning'
                }
            },
            {
                templateId: this.getTemplateId('debt'),
                projectId: projectIds[0],
                degerler: {
                    title: 'Code review yapılacak PR\'lar',
                    description: '5 adet bekleyen pull request incelenecek',
                    alan: 'Code Quality',
                    neden: 'Kod kalitesi ve takım standartları',
                    analiz: 'Bekleyen PR\'ların durumu',
                    cozum: 'Systematic code review process',
                    iyilestirmeler: 'Better code quality and team alignment',
                    sure: '1 gün',
                    priority: 'yuksek',
                    tags: 'review,git,urgent'
                }
            }
        ];

        // Process template-based tasks
        for (const taskGroup of [bugTasks, featureTasks, techDebtTasks, researchTasks, additionalTasks]) {
            for (const task of taskGroup) {
                try {
                    const result = await this.apiClient.callTool('templateden_gorev_olustur', {
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
        const completedTasks = [2, 10];
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
                    tags: 'design,ui,mockup'
                },
                {
                    parent_id: parentTaskIds[3],
                    baslik: 'Responsive grid sistemi',
                    aciklama: 'Bootstrap 5 veya Tailwind CSS ile responsive grid',
                    oncelik: GorevOncelik.Orta,
                    tags: 'frontend,css,responsive'
                },
                {
                    parent_id: parentTaskIds[3],
                    baslik: 'Animation ve transitions',
                    aciklama: 'Smooth scroll ve hover effect animasyonları',
                    oncelik: GorevOncelik.Dusuk,
                    tags: 'frontend,animation,ux'
                }
            ];

            for (const subtask of subtasks) {
                try {
                    const result = await this.apiClient.callTool('gorev_altgorev_olustur', subtask);
                    
                    // İkinci seviye alt görev ekle
                    if (subtask.baslik === 'Hero section mockup') {
                        const responseText = result.content[0].text;
                        const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                        if (idMatch) {
                            await this.apiClient.callTool('gorev_altgorev_olustur', {
                                parent_id: idMatch[1],
                                baslik: 'Color palette seçimi',
                                aciklama: 'Brand guidelines\'a uygun renk paleti',
                                oncelik: GorevOncelik.Yuksek,
                                tags: 'design,branding'
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
                    tags: 'backend,security,jwt'
                },
                {
                    parent_id: parentTaskIds[4],
                    baslik: 'Password reset flow',
                    aciklama: 'Email ile şifre sıfırlama akışı',
                    oncelik: GorevOncelik.Orta,
                    tags: 'backend,email,feature'
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
     * Template'lerden extra görevler oluştur (English field names)
     */
    private async createAdditionalTemplateExamples(projectIds: string[]): Promise<void> {
        try {
            // Example of using different template variations
            const additionalTasks = [
                {
                    templateId: this.getTemplateId('feature'),
                    degerler: {
                        title: 'GraphQL API endpoint',
                        description: 'REST API yanında GraphQL desteği',
                        purpose: 'Flexible data fetching için GraphQL layer',
                        users: 'Frontend developers, API consumers',
                        criteria: '1. Schema definition\n2. Resolvers\n3. Subscriptions',
                        effort: 'büyük',
                        priority: 'orta',
                        tags: 'backend,api,feature'
                    }
                },
                {
                    templateId: this.getTemplateId('research'),
                    degerler: {
                        topic: 'Makine öğrenmesi için framework',
                        purpose: 'Churn prediction modeli için ML framework seçimi',
                        questions: '1. TensorFlow vs PyTorch?\n2. Deployment options?\n3. Performance?',
                        criteria: 'Ease of use, community, deployment',
                        priority: 'orta',
                        tags: 'ml,data-science,research'
                    }
                }
            ];

            for (const task of additionalTasks) {
                try {
                    const result = await this.apiClient.callTool('templateden_gorev_olustur', {
                        template_id: task.templateId,
                        degerler: task.degerler
                    });

                    // Assign to appropriate project if needed
                    if (projectIds.length > 2) {
                        const responseText = result.content[0].text;
                        const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                        if (idMatch) {
                            await this.apiClient.callTool('gorev_duzenle', {
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

    // Note: AI context tools (gorev_set_active, gorev_nlp_query, gorev_context_summary,
    // gorev_batch_update) are MCP-only and not available via REST API.
    // They are used by AI assistants (Claude, Cursor, etc.), not by VS Code extension.

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
            const tasksResult = await this.apiClient.callTool('gorev_listele', {
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
                    await this.apiClient.callTool('gorev_sil', {
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
     * Uses feature template with English field names
     */
    private async createPaginationTestTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];
        const targetTaskCount = 150; // To test pagination (pageSize is 100)
        const priorities = ['dusuk', 'orta', 'yuksek'];

        Logger.info(`[TestDataSeeder] Creating ${targetTaskCount} tasks for pagination testing...`);

        for (let i = 0; i < targetTaskCount; i++) {
            try {
                const taskData = {
                    template_id: this.getTemplateId('feature'),
                    degerler: {
                        title: `Pagination Test Task ${i + 1}`,
                        description: `This is test task #${i + 1} created for pagination testing. It helps verify that the VS Code extension can handle large numbers of tasks properly.`,
                        purpose: 'Pagination testing for VS Code extension',
                        users: 'Extension developers and testers',
                        criteria: '1. Task loads correctly\n2. Pagination works\n3. No performance issues',
                        priority: priorities[i % 3],
                        tags: `pagination-test,batch-${Math.floor(i / 50) + 1}`
                    }
                };

                const result = await this.apiClient.callTool('templateden_gorev_olustur', taskData);
                const taskId = this.extractTaskId(result);
                if (taskId) {
                    taskIds.push(taskId);

                    // Assign to project
                    const projectId = projectIds[i % projectIds.length];
                    if (projectId) {
                        try {
                            await this.apiClient.callTool('gorev_duzenle', {
                                id: taskId,
                                proje_id: projectId
                            });
                        } catch {
                            // Ignore project assignment errors
                        }
                    }
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
     * Uses feature template with English field names
     */
    private async createHierarchyTestTasks(projectIds: string[]): Promise<string[]> {
        const taskIds: string[] = [];

        Logger.info('[TestDataSeeder] Creating hierarchy test tasks...');

        // Create parent tasks first (using feature template)
        const parentTasks = [
            {
                title: 'Feature: User Management System',
                description: 'Complete user management system with authentication and authorization',
                purpose: 'Enable secure user authentication and role-based access',
                users: 'All application users',
                criteria: '1. Login/Register\n2. Role management\n3. Session handling',
                priority: 'yuksek'
            },
            {
                title: 'Feature: Reporting Dashboard',
                description: 'Advanced reporting dashboard with charts and export functionality',
                purpose: 'Provide business insights through visual analytics',
                users: 'Managers and analysts',
                criteria: '1. Charts\n2. Filters\n3. Export to PDF/Excel',
                priority: 'orta'
            },
            {
                title: 'Infrastructure: Database Migration',
                description: 'Migrate from old database schema to new optimized structure',
                purpose: 'Improve database performance and maintainability',
                users: 'Development team',
                criteria: '1. Zero downtime\n2. Data integrity\n3. Rollback plan',
                priority: 'yuksek'
            }
        ];

        const parentTaskIds: string[] = [];

        for (const parentTaskData of parentTasks) {
            try {
                const taskData = {
                    template_id: this.getTemplateId('feature'),
                    degerler: {
                        ...parentTaskData,
                        tags: 'hierarchy-test,parent-task'
                    }
                };

                const result = await this.apiClient.callTool('templateden_gorev_olustur', taskData);
                const taskId = this.extractTaskId(result);
                if (taskId) {
                    parentTaskIds.push(taskId);
                    taskIds.push(taskId);

                    // Assign to project
                    if (projectIds[0]) {
                        try {
                            await this.apiClient.callTool('gorev_duzenle', {
                                id: taskId,
                                proje_id: projectIds[0]
                            });
                        } catch {
                            // Ignore project assignment errors
                        }
                    }
                    Logger.info(`Created parent task: ${parentTaskData.title} (${taskId})`);
                }
            } catch (error) {
                Logger.error(`Failed to create parent task: ${parentTaskData.title}`, error);
            }
        }

        // Create subtasks for each parent (using feature template)
        const subtaskTemplates = [
            { title: 'Design Phase', description: 'Design user interface and user experience' },
            { title: 'Backend Implementation', description: 'Implement backend API and business logic' },
            { title: 'Frontend Implementation', description: 'Create frontend components and pages' },
            { title: 'Testing Phase', description: 'Write and execute comprehensive tests' },
            { title: 'Documentation', description: 'Create user and technical documentation' }
        ];

        for (let i = 0; i < parentTaskIds.length; i++) {
            const parentId = parentTaskIds[i];

            for (const subtaskTemplate of subtaskTemplates) {
                try {
                    const taskData = {
                        template_id: this.getTemplateId('feature'),
                        degerler: {
                            title: `${subtaskTemplate.title} (Parent ${i + 1})`,
                            description: subtaskTemplate.description,
                            purpose: 'Part of parent feature implementation',
                            users: 'Development team',
                            criteria: '1. Implementation complete\n2. Code reviewed\n3. Tests passing',
                            priority: 'orta',
                            tags: `hierarchy-test,subtask,parent-${i + 1}`
                        }
                    };

                    const result = await this.apiClient.callTool('templateden_gorev_olustur', taskData);
                    const taskId = this.extractTaskId(result);
                    if (taskId) {
                        taskIds.push(taskId);

                        // Assign to project first
                        if (projectIds[0]) {
                            try {
                                await this.apiClient.callTool('gorev_duzenle', {
                                    id: taskId,
                                    proje_id: projectIds[0]
                                });
                            } catch {
                                // Ignore project assignment errors
                            }
                        }

                        // Set parent relationship
                        await this.apiClient.callTool('gorev_ust_gorev_degistir', {
                            gorev_id: taskId,
                            ust_gorev_id: parentId
                        });

                        Logger.debug(`Created subtask: ${subtaskTemplate.title} under parent ${parentId}`);
                    }
                } catch (error) {
                    Logger.error(`Failed to create subtask: ${subtaskTemplate.title}`, error);
                }
            }
        }

        Logger.info(`[TestDataSeeder] Successfully created ${taskIds.length} hierarchy test tasks (${parentTaskIds.length} parents)`);
        return taskIds;
    }

    /**
     * Extract task ID from MCP response
     */
    private extractTaskId(result: MCPToolResult): string | null {
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
