// Package fixtures provides sample test data for the Gorev project.
// This data is used by the TestDataSeeder to populate databases for testing.
package fixtures

// SampleProject represents a sample project for testing
type SampleProject struct {
	NameTR       string
	NameEN       string
	DefinitionTR string
	DefinitionEN string
}

// SampleTask represents a sample task for testing
type SampleTask struct {
	TemplateAlias string            // Template alias: bug, feature, debt, research, test, doc
	ProjectIndex  int               // Index into SampleProjects array
	Values        map[string]string // Template field values
	Status        string            // beklemede, devam_ediyor, tamamlandi, iptal
	Priority      string            // dusuk, orta, yuksek
	DueDaysOffset int               // Days from now (negative = overdue)
	Tags          []string          // Tag names
	IsParent      bool              // If true, this task will have subtasks
}

// SampleSubtask represents a subtask in a hierarchy
type SampleSubtask struct {
	ParentIndex   int               // Index into parent task array
	TemplateAlias string            // Template alias
	Values        map[string]string // Template field values
	Status        string
	Priority      string
	Children      []SampleSubtask // Recursive children for deep hierarchies
}

// SampleDependency represents a task dependency
type SampleDependency struct {
	SourceIndex int    // Index of task that depends on target
	TargetIndex int    // Index of task that must be completed first
	Type        string // onceki, engelliyor
}

// SampleTag represents a tag/label
type SampleTag struct {
	NameTR string
	NameEN string
}

// SampleProjects contains 3 sample projects
var SampleProjects = []SampleProject{
	{
		NameTR:       "Mobil Uygulama",
		NameEN:       "Mobile App",
		DefinitionTR: "iOS ve Android uygulama geliştirme projesi",
		DefinitionEN: "iOS and Android app development project",
	},
	{
		NameTR:       "Backend API",
		NameEN:       "Backend API",
		DefinitionTR: "REST API geliştirme ve bakım projesi",
		DefinitionEN: "REST API development and maintenance project",
	},
	{
		NameTR:       "Web Dashboard",
		NameEN:       "Web Dashboard",
		DefinitionTR: "Admin paneli ve kullanıcı arayüzü projesi",
		DefinitionEN: "Admin panel and user interface project",
	},
}

// SampleTasks contains 15 sample tasks with various states
var SampleTasks = []SampleTask{
	// Mobil Uygulama Project (Index 0) - 5 tasks
	{
		TemplateAlias: "bug",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Login sayfası 404 hatası",
			"description": "Production ortamında login sayfasına giderken 404 hatası alınıyor",
			"module":      "auth",
			"environment": "production",
			"steps":       "1. Uygulamayı aç\n2. Login butonuna tıkla\n3. 404 hatası görülür",
		},
		Status:        "beklemede",
		Priority:      "yuksek",
		DueDaysOffset: 3,
		Tags:          []string{"bug", "kritik"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Push notification sistemi",
			"description": "Firebase Cloud Messaging ile push notification implementasyonu",
			"scope":       "Tüm kullanıcılar için bildirim gönderme özelliği",
		},
		Status:        "devam_ediyor",
		Priority:      "orta",
		DueDaysOffset: 14,
		Tags:          []string{"ozellik", "mobil"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Dark mode tema",
			"description": "Karanlık mod tema desteği eklenmesi",
			"scope":       "Sistem ayarlarına göre otomatik tema değişimi",
		},
		Status:        "beklemede",
		Priority:      "dusuk",
		DueDaysOffset: 30,
		Tags:          []string{"ozellik", "frontend"},
	},
	{
		TemplateAlias: "bug",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Bellek sızıntısı",
			"description": "Uzun süreli kullanımda bellek sızıntısı tespit edildi",
			"module":      "core",
			"environment": "production",
			"steps":       "1. Uygulamayı uzun süre açık tut\n2. Memory profiler ile izle",
		},
		Status:        "tamamlandi",
		Priority:      "yuksek",
		DueDaysOffset: -5, // Completed 5 days ago
		Tags:          []string{"bug", "performans"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Offline mod desteği",
			"description": "İnternet bağlantısı olmadan çalışma özelliği",
			"scope":       "Temel özelliklerin offline çalışması",
		},
		Status:        "iptal",
		Priority:      "orta",
		DueDaysOffset: 0,
		Tags:          []string{"ozellik"},
	},

	// Backend API Project (Index 1) - 5 tasks
	{
		TemplateAlias: "debt",
		ProjectIndex:  1,
		Values: map[string]string{
			"title":       "Redis cache entegrasyonu",
			"description": "API response caching için Redis implementasyonu",
			"impact":      "Yüksek - %50 performans artışı bekleniyor",
			"solution":    "Redis cluster kurulumu ve cache layer implementasyonu",
		},
		Status:        "tamamlandi",
		Priority:      "yuksek",
		DueDaysOffset: -10,
		Tags:          []string{"backend", "performans"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  1,
		Values: map[string]string{
			"title":       "API rate limiting",
			"description": "DDoS koruması için rate limiting implementasyonu",
			"scope":       "IP bazlı ve kullanıcı bazlı limit",
		},
		Status:        "beklemede",
		Priority:      "orta",
		DueDaysOffset: 7,
		Tags:          []string{"backend", "guvenlik"},
	},
	{
		TemplateAlias: "bug",
		ProjectIndex:  1,
		Values: map[string]string{
			"title":       "Timeout hatası düzeltme",
			"description": "Büyük veri setlerinde timeout hatası alınıyor",
			"module":      "data",
			"environment": "production",
			"steps":       "1. 10000+ kayıt içeren sorgu yap\n2. Timeout hatası görülür",
		},
		Status:        "devam_ediyor",
		Priority:      "orta",
		DueDaysOffset: 2,
		Tags:          []string{"bug", "backend"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  1,
		Values: map[string]string{
			"title":       "GraphQL API desteği",
			"description": "REST API yanında GraphQL endpoint eklenmesi",
			"scope":       "Temel CRUD operasyonları için GraphQL",
		},
		Status:        "beklemede",
		Priority:      "dusuk",
		DueDaysOffset: 60,
		Tags:          []string{"ozellik", "backend"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  1,
		Values: map[string]string{
			"title":       "API dokümantasyonu güncelleme",
			"description": "Swagger/OpenAPI dokümantasyonunun güncellenmesi",
			"scope":       "Tüm endpoint'ler için güncel dokümantasyon",
		},
		Status:        "tamamlandi",
		Priority:      "dusuk",
		DueDaysOffset: -3,
		Tags:          []string{"dokumantasyon"},
	},

	// Web Dashboard Project (Index 2) - 5 tasks
	{
		TemplateAlias: "feature",
		ProjectIndex:  2,
		Values: map[string]string{
			"title":       "Dashboard ana sayfa tasarımı",
			"description": "Modern ve kullanıcı dostu dashboard tasarımı",
			"scope":       "Ana metrikler, grafikler ve widget'lar",
		},
		Status:        "devam_ediyor",
		Priority:      "orta",
		DueDaysOffset: 10,
		Tags:          []string{"ozellik", "frontend"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  2,
		Values: map[string]string{
			"title":       "Kullanıcı yönetimi modülü",
			"description": "Admin panelinde kullanıcı CRUD işlemleri",
			"scope":       "Kullanıcı listesi, ekleme, düzenleme, silme, rol atama",
		},
		Status:        "devam_ediyor",
		Priority:      "yuksek",
		DueDaysOffset: 20,
		Tags:          []string{"ozellik", "admin"},
		IsParent:      true, // This task will have subtasks
	},
	{
		TemplateAlias: "bug",
		ProjectIndex:  2,
		Values: map[string]string{
			"title":       "Tablo sıralama hatası",
			"description": "Tablo sütunlarında sıralama düzgün çalışmıyor",
			"module":      "table",
			"environment": "staging",
			"steps":       "1. Herhangi bir tabloyu aç\n2. Sütun başlığına tıkla\n3. Sıralama yanlış",
		},
		Status:        "tamamlandi",
		Priority:      "orta",
		DueDaysOffset: -1,
		Tags:          []string{"bug", "frontend"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  2,
		Values: map[string]string{
			"title":       "Raporlama modülü",
			"description": "Detaylı raporlar ve analitik dashboard",
			"scope":       "PDF export, grafik çeşitleri, tarih filtreleri",
		},
		Status:        "beklemede",
		Priority:      "orta",
		DueDaysOffset: 45,
		Tags:          []string{"ozellik", "analitik"},
	},
	{
		TemplateAlias: "debt",
		ProjectIndex:  2,
		Values: map[string]string{
			"title":       "React güncellemesi",
			"description": "React 17'den React 18'e güncelleme",
			"impact":      "Orta - Concurrent features kullanılabilir",
			"solution":    "Aşamalı güncelleme ve test",
		},
		Status:        "tamamlandi",
		Priority:      "dusuk",
		DueDaysOffset: -15,
		Tags:          []string{"teknik-borc", "frontend"},
	},
}

// SampleSubtasks contains subtask hierarchies (3 levels deep)
// These will be attached to the "Kullanıcı yönetimi modülü" task (index 11)
var SampleSubtasks = []SampleSubtask{
	{
		ParentIndex:   11, // "Kullanıcı yönetimi modülü"
		TemplateAlias: "feature",
		Values: map[string]string{
			"title":       "Kayıt formu",
			"description": "Yeni kullanıcı kayıt formu implementasyonu",
			"scope":       "Email, şifre, ad-soyad alanları",
		},
		Status:   "devam_ediyor",
		Priority: "orta",
		Children: []SampleSubtask{
			{
				TemplateAlias: "feature",
				Values: map[string]string{
					"title":       "Form validasyonu",
					"description": "Frontend ve backend form validasyonu",
					"scope":       "Email format, şifre güçlülüğü kontrolü",
				},
				Status:   "devam_ediyor",
				Priority: "orta",
				Children: []SampleSubtask{
					{
						TemplateAlias: "feature",
						Values: map[string]string{
							"title":       "Email doğrulama",
							"description": "Email doğrulama linki gönderimi",
							"scope":       "Doğrulama email'i ve link sistemi",
						},
						Status:   "beklemede",
						Priority: "orta",
					},
				},
			},
			{
				TemplateAlias: "feature",
				Values: map[string]string{
					"title":       "Captcha entegrasyonu",
					"description": "Bot koruması için reCAPTCHA",
					"scope":       "Google reCAPTCHA v3 entegrasyonu",
				},
				Status:   "beklemede",
				Priority: "dusuk",
			},
		},
	},
	{
		ParentIndex:   11, // "Kullanıcı yönetimi modülü"
		TemplateAlias: "feature",
		Values: map[string]string{
			"title":       "Giriş sistemi",
			"description": "Kullanıcı giriş ve oturum yönetimi",
			"scope":       "Login formu, session management, remember me",
		},
		Status:   "beklemede",
		Priority: "yuksek",
		Children: []SampleSubtask{
			{
				TemplateAlias: "feature",
				Values: map[string]string{
					"title":       "2FA implementasyonu",
					"description": "İki faktörlü kimlik doğrulama",
					"scope":       "TOTP (Google Authenticator) desteği",
				},
				Status:   "beklemede",
				Priority: "orta",
			},
		},
	},
}

// SampleDependencies contains task dependencies
var SampleDependencies = []SampleDependency{
	{
		SourceIndex: 6, // "API rate limiting" depends on
		TargetIndex: 5, // "Redis cache entegrasyonu"
		Type:        "onceki",
	},
	{
		SourceIndex: 10, // "Dashboard ana sayfa" depends on
		TargetIndex: 5,  // "Redis cache entegrasyonu" (backend ready)
		Type:        "engelliyor",
	},
	{
		SourceIndex: 13, // "Raporlama modülü" depends on
		TargetIndex: 10, // "Dashboard ana sayfa"
		Type:        "onceki",
	},
}

// SampleTags contains all available tags
var SampleTags = []SampleTag{
	{NameTR: "bug", NameEN: "bug"},
	{NameTR: "kritik", NameEN: "critical"},
	{NameTR: "ozellik", NameEN: "feature"},
	{NameTR: "backend", NameEN: "backend"},
	{NameTR: "frontend", NameEN: "frontend"},
	{NameTR: "mobil", NameEN: "mobile"},
	{NameTR: "arastirma", NameEN: "research"},
	{NameTR: "dokumantasyon", NameEN: "documentation"},
	{NameTR: "performans", NameEN: "performance"},
	{NameTR: "guvenlik", NameEN: "security"},
	{NameTR: "admin", NameEN: "admin"},
	{NameTR: "analitik", NameEN: "analytics"},
	{NameTR: "teknik-borc", NameEN: "tech-debt"},
}

// MinimalSampleProjects contains minimal data for quick tests
var MinimalSampleProjects = []SampleProject{
	{
		NameTR:       "Test Projesi",
		NameEN:       "Test Project",
		DefinitionTR: "Test amaçlı proje",
		DefinitionEN: "Project for testing purposes",
	},
}

// MinimalSampleTasks contains minimal tasks for quick tests
var MinimalSampleTasks = []SampleTask{
	{
		TemplateAlias: "bug",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Test Bug",
			"description": "Test bug açıklaması",
			"module":      "test",
			"environment": "development",
			"steps":       "Test adımları",
		},
		Status:        "beklemede",
		Priority:      "orta",
		DueDaysOffset: 7,
		Tags:          []string{"bug"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Test Feature",
			"description": "Test özellik açıklaması",
			"scope":       "Test kapsamı",
		},
		Status:        "devam_ediyor",
		Priority:      "yuksek",
		DueDaysOffset: 14,
		Tags:          []string{"ozellik"},
	},
	{
		TemplateAlias: "feature",
		ProjectIndex:  0,
		Values: map[string]string{
			"title":       "Completed Task",
			"description": "Tamamlanmış görev",
			"scope":       "Test",
		},
		Status:        "tamamlandi",
		Priority:      "dusuk",
		DueDaysOffset: -3,
		Tags:          []string{"ozellik"},
	},
}
