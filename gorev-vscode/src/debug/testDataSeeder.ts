import { MCPClient } from '../mcp/client';
import { GorevDurum, GorevOncelik } from '../models/common';
import { Logger } from '../utils/logger';
import * as vscode from 'vscode';

/**
 * Debug için zengin test verileri oluşturur
 */
export class TestDataSeeder {
    constructor(private mcpClient: MCPClient) {}

    /**
     * Test verilerini oluştur
     */
    async seedTestData(): Promise<void> {
        const result = await vscode.window.showInformationMessage(
            'Test verileri oluşturulacak. Mevcut veriler korunacak. Devam etmek istiyor musunuz?',
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

                // 2. Test görevleri oluştur
                progress.report({ increment: 30, message: 'Görevler oluşturuluyor...' });
                const taskIds = await this.createTestTasks(projectIds);

                // 3. Bağımlılıklar oluştur
                progress.report({ increment: 20, message: 'Bağımlılıklar oluşturuluyor...' });
                await this.createTestDependencies(taskIds);

                // 4. Alt görevler oluştur
                progress.report({ increment: 10, message: 'Alt görevler oluşturuluyor...' });
                await this.createSubtasks(taskIds);

                // 5. Template'lerden görevler oluştur
                progress.report({ increment: 10, message: 'Template görevleri oluşturuluyor...' });
                await this.createTasksFromTemplates(projectIds);

                // 6. Bazı görevleri tamamla ve AI interaksiyonları ekle
                progress.report({ increment: 10, message: 'Görev durumları güncelleniyor...' });
                await this.updateSomeTaskStatuses(taskIds);

                // 7. AI context oluştur
                progress.report({ increment: 10, message: 'AI context oluşturuluyor...' });
                await this.setupAIContext(taskIds);

                progress.report({ increment: 10, message: 'Tamamlandı!' });
            });

            vscode.window.showInformationMessage('✅ Test verileri başarıyla oluşturuldu!');
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
     * Test görevleri oluştur
     */
    private async createTestTasks(projectIds: string[]): Promise<string[]> {
        const taskTemplates = [
            // Web Sitesi görevleri
            {
                baslik: 'Ana sayfa tasarımını tamamla',
                aciklama: 'Modern ve responsive ana sayfa tasarımı yapılacak. Hero section, özellikler bölümü ve footer dahil.',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 0 ? projectIds[0] : undefined,
                son_tarih: this.getDateString(2),
                etiketler: 'design,frontend,urgent'
            },
            {
                baslik: 'Kullanıcı giriş sistemi implement et',
                aciklama: 'JWT tabanlı authentication sistemi kurulacak. Login, register, forgot password sayfaları dahil.',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 0 ? projectIds[0] : undefined,
                son_tarih: this.getDateString(5),
                etiketler: 'backend,security,feature'
            },
            {
                baslik: 'Ürün kataloğu sayfası',
                aciklama: 'Filtreleme, sıralama ve pagination özellikleri ile ürün listeleme sayfası',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 0 ? projectIds[0] : undefined,
                son_tarih: this.getDateString(7),
                etiketler: 'frontend,feature'
            },
            {
                baslik: 'SEO optimizasyonu',
                aciklama: 'Meta taglar, sitemap, robots.txt ve sayfa hızı optimizasyonu',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 0 ? projectIds[0] : undefined,
                son_tarih: this.getDateString(14),
                etiketler: 'seo,performance'
            },
            {
                baslik: 'Contact form entegrasyonu',
                aciklama: 'Email gönderimi ile iletişim formu. Spam koruması dahil.',
                oncelik: GorevOncelik.Dusuk,
                proje_id: projectIds.length > 0 ? projectIds[0] : undefined,
                etiketler: 'frontend,feature'
            },

            // Mobil Uygulama görevleri
            {
                baslik: 'Push notification sistemi',
                aciklama: 'Firebase Cloud Messaging entegrasyonu ile bildirim sistemi',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 1 ? projectIds[1] : undefined,
                son_tarih: this.getDateString(-2), // Gecikmiş
                etiketler: 'mobile,feature,firebase'
            },
            {
                baslik: 'Offline mode desteği',
                aciklama: 'SQLite ile local veri saklama ve senkronizasyon',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 1 ? projectIds[1] : undefined,
                son_tarih: this.getDateString(10),
                etiketler: 'mobile,feature,database'
            },
            {
                baslik: 'Dark mode tema',
                aciklama: 'Sistem ayarlarına göre otomatik tema değişimi',
                oncelik: GorevOncelik.Dusuk,
                proje_id: projectIds.length > 1 ? projectIds[1] : undefined,
                etiketler: 'mobile,ui,enhancement'
            },
            {
                baslik: 'App Store deployment',
                aciklama: 'iOS App Store submission hazırlıkları ve yayınlama',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 1 ? projectIds[1] : undefined,
                son_tarih: this.getDateString(0), // Bugün
                etiketler: 'deployment,ios,critical'
            },

            // Backend API görevleri
            {
                baslik: 'GraphQL endpoint ekle',
                aciklama: 'REST API yanında GraphQL desteği eklenecek',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 2 ? projectIds[2] : undefined,
                son_tarih: this.getDateString(21),
                etiketler: 'backend,api,feature'
            },
            {
                baslik: 'Rate limiting implement et',
                aciklama: 'API güvenliği için rate limiting ve throttling',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 2 ? projectIds[2] : undefined,
                son_tarih: this.getDateString(3),
                etiketler: 'backend,security,performance'
            },
            {
                baslik: 'Redis cache entegrasyonu',
                aciklama: 'Performans artışı için Redis cache katmanı',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 2 ? projectIds[2] : undefined,
                etiketler: 'backend,performance,redis'
            },
            {
                baslik: 'API dokümantasyonu güncelle',
                aciklama: 'Swagger/OpenAPI dokümantasyonu güncellenecek',
                oncelik: GorevOncelik.Dusuk,
                proje_id: projectIds.length > 2 ? projectIds[2] : undefined,
                son_tarih: this.getDateString(30),
                etiketler: 'documentation,api'
            },

            // Veri Analitiği görevleri
            {
                baslik: 'Dashboard prototype hazırla',
                aciklama: 'Figma\'da interaktif dashboard prototipi',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 3 ? projectIds[3] : undefined,
                son_tarih: this.getDateString(1),
                etiketler: 'design,analytics,prototype'
            },
            {
                baslik: 'ETL pipeline kurulumu',
                aciklama: 'Apache Airflow ile veri pipeline\'ı kurulacak',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 3 ? projectIds[3] : undefined,
                son_tarih: this.getDateString(7),
                etiketler: 'data,backend,infrastructure'
            },
            {
                baslik: 'Makine öğrenmesi modeli',
                aciklama: 'Müşteri churn prediction modeli geliştirilecek',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 3 ? projectIds[3] : undefined,
                etiketler: 'ml,data-science,python'
            },

            // Güvenlik görevleri
            {
                baslik: 'Penetrasyon testi yap',
                aciklama: 'OWASP Top 10 güvenlik açıklarını test et',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 4 ? projectIds[4] : undefined,
                son_tarih: this.getDateString(-5), // Gecikmiş
                etiketler: 'security,testing,critical'
            },
            {
                baslik: 'SSL sertifikası yenile',
                aciklama: 'Tüm subdomain\'ler için wildcard SSL sertifikası',
                oncelik: GorevOncelik.Yuksek,
                proje_id: projectIds.length > 4 ? projectIds[4] : undefined,
                son_tarih: this.getDateString(-1), // Gecikmiş
                etiketler: 'security,infrastructure,urgent'
            },
            {
                baslik: '2FA implementasyonu',
                aciklama: 'Google Authenticator ile iki faktörlü doğrulama',
                oncelik: GorevOncelik.Orta,
                proje_id: projectIds.length > 4 ? projectIds[4] : undefined,
                son_tarih: this.getDateString(14),
                etiketler: 'security,feature,backend'
            },

            // Projesiz görevler
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
                son_tarih: this.getDateString(0), // Bugün
                etiketler: 'review,git,urgent'
            },
            {
                baslik: 'Teknik blog yazısı',
                aciklama: 'Microservices best practices hakkında blog yazısı',
                oncelik: GorevOncelik.Dusuk,
                etiketler: 'writing,documentation'
            },
            {
                baslik: 'Yeni developer onboarding',
                aciklama: 'Yeni başlayan developer için environment setup',
                oncelik: GorevOncelik.Orta,
                son_tarih: this.getDateString(2),
                etiketler: 'hr,setup,training'
            }
        ];

        const taskIds: string[] = [];

        // Proje ID'lerini logla
        Logger.info(`Available project IDs: ${projectIds.join(', ')}`);

        for (const task of taskTemplates) {
            try {
                // Proje ID'sini kontrol et ve logla
                if (task.proje_id) {
                    if (!projectIds.includes(task.proje_id)) {
                        Logger.warn(`Invalid project ID in task "${task.baslik}": ${task.proje_id}`);
                        Logger.warn(`Available project IDs: ${projectIds.join(', ')}`);
                    } else {
                        Logger.debug(`Task "${task.baslik}" assigned to project ID: ${task.proje_id}`);
                    }
                } else {
                    Logger.debug(`Task "${task.baslik}" has no project (projesiz)`);
                }
                
                const result = await this.mcpClient.callTool('gorev_olustur', task);
                const responseText = result.content[0].text;
                Logger.debug(`Task creation response for "${task.baslik}":`, responseText);
                
                // UUID formatında ID ara
                const idMatch = responseText.match(/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})/i);
                if (idMatch) {
                    taskIds.push(idMatch[1]);
                    Logger.info(`Created task: ${task.baslik} with project: ${task.proje_id || 'none'}`);
                }
            } catch (error) {
                Logger.error(`Failed to create task "${task.baslik}":`, error);
            }
        }

        return taskIds;
    }

    /**
     * Test bağımlılıkları oluştur
     */
    private async createTestDependencies(taskIds: string[]): Promise<void> {
        // Örnek bağımlılıklar
        const dependencies = [
            { kaynak: 0, hedef: 1, tip: 'blocks' }, // Ana sayfa tasarımı -> Login sistemi'ni bloklar
            { kaynak: 1, hedef: 2, tip: 'blocks' }, // Login sistemi -> Ürün kataloğu'nu bloklar
            { kaynak: 11, hedef: 12, tip: 'depends_on' }, // Redis cache -> Rate limiting'e bağlı
            { kaynak: 5, hedef: 8, tip: 'blocks' }, // Push notification -> App Store deployment'ı bloklar
            { kaynak: 14, hedef: 15, tip: 'depends_on' }, // ETL pipeline -> Dashboard prototype'a bağlı
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
        const inProgressTasks = [1, 5, 9, 14, 20];
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
        const completedTasks = [4, 7, 12, 15];
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
        // Ana sayfa tasarımı için alt görevler
        if (parentTaskIds[0]) {
            const subtasks = [
                {
                    parent_id: parentTaskIds[0],
                    baslik: 'Hero section mockup',
                    aciklama: 'Ana sayfa hero bölümü için Figma mockup hazırla',
                    oncelik: GorevOncelik.Yuksek,
                    etiketler: 'design,ui,mockup'
                },
                {
                    parent_id: parentTaskIds[0],
                    baslik: 'Responsive grid sistemi',
                    aciklama: 'Bootstrap 5 veya Tailwind CSS ile responsive grid',
                    oncelik: GorevOncelik.Orta,
                    etiketler: 'frontend,css,responsive'
                },
                {
                    parent_id: parentTaskIds[0],
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

        // Login sistemi için alt görevler
        if (parentTaskIds[1]) {
            const subtasks = [
                {
                    parent_id: parentTaskIds[1],
                    baslik: 'JWT token implementasyonu',
                    aciklama: 'Access ve refresh token yönetimi',
                    oncelik: GorevOncelik.Yuksek,
                    son_tarih: this.getDateString(3),
                    etiketler: 'backend,security,jwt'
                },
                {
                    parent_id: parentTaskIds[1],
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

        // Dashboard prototype için alt görevler
        if (parentTaskIds[13]) {
            const subtasks = [
                {
                    parent_id: parentTaskIds[13],
                    baslik: 'KPI cards tasarımı',
                    aciklama: 'Ana metrikleri gösteren kart componentleri',
                    oncelik: GorevOncelik.Yuksek,
                    etiketler: 'design,dashboard,component'
                },
                {
                    parent_id: parentTaskIds[13],
                    baslik: 'Chart library araştırması',
                    aciklama: 'Chart.js vs D3.js vs ApexCharts karşılaştırması',
                    oncelik: GorevOncelik.Orta,
                    etiketler: 'research,frontend,visualization'
                },
                {
                    parent_id: parentTaskIds[13],
                    baslik: 'Real-time data updates',
                    aciklama: 'WebSocket ile canlı veri güncellemeleri',
                    oncelik: GorevOncelik.Orta,
                    etiketler: 'frontend,websocket,realtime'
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
     * Template'lerden görevler oluştur
     */
    private async createTasksFromTemplates(projectIds: string[]): Promise<void> {
        try {
            // Önce template'leri listele
            const templatesResult = await this.mcpClient.callTool('template_listele', {});
            Logger.debug('Available templates:', templatesResult.content[0].text);

            // Bug raporu template'inden görev oluştur
            try {
                await this.mcpClient.callTool('templateden_gorev_olustur', {
                    template_id: '311422f9-51ad-4678-8631-e0f7957aae47', // Bug Raporu template ID
                    degerler: {
                        baslik: 'Login sayfası 404 hatası veriyor',
                        aciklama: 'Production ortamında /login URL\'ine gittiğimizde 404 hatası alıyoruz',
                        modul: 'Authentication',
                        ortam: 'production',
                        adimlar: '1. Production URL\'ine git\n2. /login sayfasına git\n3. 404 hatası görünüyor',
                        beklenen: 'Login sayfası açılmalı',
                        mevcut: '404 Not Found hatası',
                        oncelik: 'yuksek',
                        etiketler: 'bug,critical,production'
                    }
                });
            } catch (error) {
                Logger.error('Failed to create task from bug template:', error);
            }

            // Araştırma görevi template'inden oluştur
            try {
                await this.mcpClient.callTool('templateden_gorev_olustur', {
                    template_id: '146837f2-bd50-4a88-9d93-38da1d7c09d6', // Araştırma Görevi template ID
                    degerler: {
                        konu: 'Next.js 14 App Router',
                        amac: 'Yeni projede kullanmak için Next.js 14 App Router özelliklerini araştırmak',
                        sorular: '1. Performance karşılaştırması?\n2. Migration süreci?\n3. Edge runtime desteği?',
                        kaynaklar: 'Next.js dokümantasyonu, Vercel blog, YouTube tutorialları',
                        alternatifler: 'Pages Router, Remix, SvelteKit',
                        kriterler: 'Performance, Developer Experience, SEO, Bundle Size',
                        son_tarih: this.getDateString(14),
                        oncelik: 'orta',
                        etiketler: 'araştırma,nextjs,frontend'
                    }
                });
            } catch (error) {
                Logger.error('Failed to create task from research template:', error);
            }

            // Özellik isteği template'inden oluştur
            try {
                await this.mcpClient.callTool('templateden_gorev_olustur', {
                    template_id: '430d308c-440d-49cd-a307-9db78f8608bf', // Özellik İsteği template ID
                    degerler: {
                        baslik: 'Dark mode toggle özelliği',
                        aciklama: 'Kullanıcılar tema tercihlerini kaydedebilmeli',
                        amac: 'Kullanıcı deneyimini iyileştirmek ve göz yorgunluğunu azaltmak',
                        kullanicilar: 'Tüm kullanıcılar, özellikle gece çalışanlar',
                        kriterler: '1. Sistem temasına uyum\n2. Manuel toggle\n3. Tercih kaydetme\n4. Smooth transition',
                        ui_ux: 'Header\'da toggle switch, sistem temasını takip et opsiyonu',
                        efor: 'orta',
                        oncelik: 'orta',
                        etiketler: 'özellik,ui,enhancement'
                    }
                });
            } catch (error) {
                Logger.error('Failed to create task from feature template:', error);
            }
        } catch (error) {
            Logger.error('Failed to list templates:', error);
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
                            { id: taskIds[3], updates: { durum: 'devam_ediyor' } },
                            { id: taskIds[4], updates: { durum: 'tamamlandi' } }
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
            '⚠️ DİKKAT: Tüm görevler ve projeler silinecek! Emin misiniz?',
            'Evet, Sil',
            'Hayır'
        );

        if (result !== 'Evet, Sil') {
            return;
        }

        try {
            // Önce tüm görevleri listele ve sil
            const tasksResult = await this.mcpClient.callTool('gorev_listele', {
                tum_projeler: true
            });

            // Parse task IDs from response
            const taskIdMatches = tasksResult.content[0].text.matchAll(/ID: ([a-f0-9-]+)/g);
            const taskIds = Array.from(taskIdMatches).map(match => match[1]);

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

            vscode.window.showInformationMessage('✅ Test verileri temizlendi!');
        } catch (error) {
            vscode.window.showErrorMessage(`Test verileri temizlenemedi: ${error}`);
            Logger.error('Failed to clear test data:', error);
        }
    }
}