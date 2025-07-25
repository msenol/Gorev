package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/msenol/gorev/internal/mcp"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
	langFlag  string
)

// getMigrationsPath returns the correct path to migrations folder
func getMigrationsPath() string {
	// First check if migrations exist in user's home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeMigrationsPath := filepath.Join(homeDir, ".gorev", "internal", "veri", "migrations")
		if _, err := os.Stat(homeMigrationsPath); err == nil {
			return "file://" + homeMigrationsPath
		}
	}

	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to relative path
		return "file://internal/veri/migrations"
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	// First, check if we can find the migrations in the standard project structure
	// This handles the case where the executable is in a temporary directory (go run)
	// or installed in a different location
	possiblePaths := []string{
		// Direct path from executable location
		filepath.Join(exeDir, "internal", "veri", "migrations"),
		// If in build directory
		filepath.Join(filepath.Dir(exeDir), "internal", "veri", "migrations"),
		// If in gorev-mcpserver directory
		filepath.Join(exeDir, "..", "internal", "veri", "migrations"),
		// Try to find gorev-mcpserver in parent directories
		filepath.Join(exeDir, "..", "..", "gorev-mcpserver", "internal", "veri", "migrations"),
		filepath.Join(exeDir, "..", "..", "..", "gorev-mcpserver", "internal", "veri", "migrations"),
	}

	// Also check GOREV_ROOT environment variable if set
	if gorevRoot := os.Getenv("GOREV_ROOT"); gorevRoot != "" {
		possiblePaths = append([]string{
			filepath.Join(gorevRoot, "internal", "veri", "migrations"),
		}, possiblePaths...)
	}

	// Try each possible path
	for _, path := range possiblePaths {
		// Clean the path to resolve .. and .
		cleanPath := filepath.Clean(path)
		// Check if the migrations directory exists
		if _, err := os.Stat(cleanPath); err == nil {
			// Found the migrations directory
			return "file://" + cleanPath
		}
	}

	// Fallback: assume migrations are relative to current working directory
	cwd, err := os.Getwd()
	if err == nil {
		migrationsPath := filepath.Join(cwd, "internal", "veri", "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return "file://" + migrationsPath
		}
	}

	// Last resort: return relative path and hope for the best
	return "file://internal/veri/migrations"
}

// getDatabasePath returns the correct path to database file
func getDatabasePath() string {
	// First check if database exists in user's home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeDBPath := filepath.Join(homeDir, ".gorev", "gorev.db")
		if _, err := os.Stat(homeDBPath); err == nil {
			return homeDBPath
		}
	}

	// Get the executable path
	exePath, err := os.Executable()
	if err != nil {
		// Fallback to current directory
		return "gorev.db"
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	// First, check if we can find an existing database in the standard locations
	possiblePaths := []string{
		// Direct path from executable location
		filepath.Join(exeDir, "gorev.db"),
		// If in build directory
		filepath.Join(filepath.Dir(exeDir), "gorev.db"),
		// If in gorev-mcpserver directory
		filepath.Join(exeDir, "..", "gorev.db"),
		// Try to find gorev-mcpserver in parent directories
		filepath.Join(exeDir, "..", "..", "gorev-mcpserver", "gorev.db"),
		filepath.Join(exeDir, "..", "..", "..", "gorev-mcpserver", "gorev.db"),
	}

	// Also check GOREV_ROOT environment variable if set
	if gorevRoot := os.Getenv("GOREV_ROOT"); gorevRoot != "" {
		possiblePaths = append([]string{
			filepath.Join(gorevRoot, "gorev.db"),
		}, possiblePaths...)
	}

	// Try each possible path for existing database
	for _, path := range possiblePaths {
		// Clean the path to resolve .. and .
		cleanPath := filepath.Clean(path)
		// Check if the database file exists
		if _, err := os.Stat(cleanPath); err == nil {
			// Found existing database
			return cleanPath
		}
	}

	// If no existing database found, determine where to create a new one
	// Priority: GOREV_ROOT > project root > current directory

	// 1. Try GOREV_ROOT if set
	if gorevRoot := os.Getenv("GOREV_ROOT"); gorevRoot != "" {
		return filepath.Join(gorevRoot, "gorev.db")
	}

	// 2. Try to find project root by looking for internal/veri/migrations
	for _, path := range possiblePaths {
		dir := filepath.Dir(path)
		migrationsPath := filepath.Join(dir, "internal", "veri", "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			// Found project root
			return filepath.Join(dir, "gorev.db")
		}
	}

	// 3. Fallback to current working directory
	cwd, err := os.Getwd()
	if err == nil {
		// Check if we're in the project root
		migrationsPath := filepath.Join(cwd, "internal", "veri", "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return filepath.Join(cwd, "gorev.db")
		}
	}

	// Last resort: use current directory
	return "gorev.db"
}

// detectLanguage detects the language preference from environment variables and CLI flags
func detectLanguage() string {
	// Priority 1: CLI flag (if we're processing it)
	if langFlag != "" {
		return langFlag
	}

	// Priority 2: GOREV_LANG environment variable
	if lang := os.Getenv("GOREV_LANG"); lang != "" {
		if lang == "en" || lang == "tr" {
			return lang
		}
	}

	// Priority 3: LANG environment variable (partial detection)
	if lang := os.Getenv("LANG"); lang != "" {
		if lang[:2] == "tr" {
			return "tr"
		}
		if lang[:2] == "en" {
			return "en"
		}
	}

	// Default to Turkish
	return "tr"
}

func main() {
	// Initialize i18n system
	lang := detectLanguage()
	if err := i18n.Initialize(lang); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize i18n: %v\n", err)
	}

	rootCmd := &cobra.Command{
		Use:   "gorev",
		Short: i18n.T("cli.appDescription"),
		Long: `Gorev, Model Context Protocol (MCP) üzerinden AI asistanlarına
görev yönetimi yetenekleri sağlayan modern bir sunucudur.`,
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: i18n.T("cli.serve"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Re-initialize i18n if language flag was provided
			if langFlag != "" {
				if err := i18n.SetLanguage(langFlag); err != nil {
					return fmt.Errorf("failed to set language to %s: %w", langFlag, err)
				}
			}
			return runServer()
		},
	}
	serveCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "Language preference (tr, en)")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("cli.version"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Gorev %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Git Commit: %s\n", gitCommit)
		},
	}

	// Template komutları
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: i18n.T("cli.template"),
		Long:  i18n.T("cli.templateDescription"),
	}

	templateListCmd := &cobra.Command{
		Use:   "list [kategori]",
		Short: i18n.T("cli.templateList"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var kategori string
			if len(args) > 0 {
				kategori = args[0]
			}
			return listTemplates(kategori)
		},
	}

	templateShowCmd := &cobra.Command{
		Use:   "show <template-id>",
		Short: i18n.T("cli.templateShow"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showTemplate(args[0])
		},
	}

	templateInitCmd := &cobra.Command{
		Use:   "init",
		Short: i18n.T("cli.templateInit"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return initTemplates()
		},
	}

	templateCmd.AddCommand(templateListCmd, templateShowCmd, templateInitCmd)

	// MCP test commands
	mcpCmd := createMCPCommand()

	rootCmd.AddCommand(serveCmd, versionCmd, templateCmd, mcpCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	// Veritabanını başlat
	veriYonetici, err := gorev.YeniVeriYonetici(getDatabasePath(), getMigrationsPath())
	if err != nil {
		return fmt.Errorf(i18n.T("error.dataManagerInit", map[string]interface{}{"Error": err}))
	}
	defer veriYonetici.Kapat()

	// İş mantığı servisini oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// MCP sunucusunu başlat
	sunucu, err := mcp.YeniMCPSunucu(isYonetici)
	if err != nil {
		return fmt.Errorf(i18n.T("error.mcpServerCreate", map[string]interface{}{"Error": err}))
	}

	// Sunucuyu çalıştır
	// fmt.Fprintln(os.Stderr, "Gorev MCP sunucusu başlatılıyor...")
	if err := mcp.ServeSunucu(sunucu); err != nil {
		return fmt.Errorf(i18n.T("error.serverStart", map[string]interface{}{"Error": err}))
	}

	return nil
}

func listTemplates(kategori string) error {
	// Veritabanını başlat
	veriYonetici, err := gorev.YeniVeriYonetici(getDatabasePath(), getMigrationsPath())
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer veriYonetici.Kapat()

	// Template'leri listele
	templates, err := veriYonetici.TemplateListele(kategori)
	if err != nil {
		return fmt.Errorf("template'ler listelenemedi: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("Henüz template bulunmuyor.")
		return nil
	}

	// Kategorilere göre grupla
	kategoriMap := make(map[string][]*gorev.GorevTemplate)
	for _, tmpl := range templates {
		kategoriMap[tmpl.Kategori] = append(kategoriMap[tmpl.Kategori], tmpl)
	}

	// Her kategoriyi göster
	for kat, tmpls := range kategoriMap {
		fmt.Printf("\n=== %s ===\n", kat)
		for _, tmpl := range tmpls {
			fmt.Printf("\n%s (ID: %s)\n", tmpl.Isim, tmpl.ID)
			fmt.Printf("  %s\n", tmpl.Tanim)
			fmt.Printf("  Başlık: %s\n", tmpl.VarsayilanBaslik)
		}
	}

	fmt.Println("\nDetaylı bilgi için: gorev template show <template-id>")
	return nil
}

func showTemplate(templateID string) error {
	// Veritabanını başlat
	veriYonetici, err := gorev.YeniVeriYonetici(getDatabasePath(), getMigrationsPath())
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer veriYonetici.Kapat()

	// Template'i getir
	template, err := veriYonetici.TemplateGetir(templateID)
	if err != nil {
		return fmt.Errorf("template bulunamadı: %w", err)
	}

	fmt.Printf("\n=== %s ===\n", template.Isim)
	fmt.Printf("ID: %s\n", template.ID)
	fmt.Printf("Kategori: %s\n", template.Kategori)
	fmt.Printf("Açıklama: %s\n", template.Tanim)
	fmt.Printf("\nBaşlık Şablonu: %s\n", template.VarsayilanBaslik)

	fmt.Println("\nAlanlar:")
	for _, alan := range template.Alanlar {
		zorunlu := ""
		if alan.Zorunlu {
			zorunlu = " (zorunlu)"
		}
		fmt.Printf("  - %s (%s)%s", alan.Isim, alan.Tip, zorunlu)
		if alan.Varsayilan != "" {
			fmt.Printf(" [varsayılan: %s]", alan.Varsayilan)
		}
		if len(alan.Secenekler) > 0 {
			fmt.Printf("\n    Seçenekler: %v", alan.Secenekler)
		}
		fmt.Println()
	}

	fmt.Println("\nÖrnek Açıklama:")
	fmt.Println(template.AciklamaTemplate)

	return nil
}

func initTemplates() error {
	// Veritabanını başlat
	veriYonetici, err := gorev.YeniVeriYonetici(getDatabasePath(), getMigrationsPath())
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer veriYonetici.Kapat()

	// Varsayılan template'leri oluştur
	err = veriYonetici.VarsayilanTemplateleriOlustur()
	if err != nil {
		return fmt.Errorf("template'ler oluşturulamadı: %w", err)
	}

	fmt.Println("✓ Varsayılan template'ler başarıyla oluşturuldu.")
	return nil
}
