package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/msenol/gorev/internal/api"
	"github.com/msenol/gorev/internal/config"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/msenol/gorev/internal/mcp"
	"github.com/spf13/cobra"
)

var (
	version       = "v0.17.0"
	buildTime     = "unknown"
	gitCommit     = "unknown"
	langFlag      string
	debugFlag     bool
	apiPortFlag   string
	noAPIFlag     bool
	serverModeFlag string
	dbPathFlag    string
)

// getMigrationsPath returns the correct path to migrations folder
// It first tries to use embedded migrations, then falls back to filesystem
func getMigrationsPath() string {
	// Priority 1: Use embedded migrations (always available in NPX package)
	return "embedded://migrations"
}

// getDatabasePath returns the correct path to database file
func getDatabasePath() string {
	// First priority: GOREV_DB_PATH environment variable (can be set by any MCP client)
	if dbPath := os.Getenv("GOREV_DB_PATH"); dbPath != "" {
		// Clean the path and make it absolute if not already
		cleanPath := filepath.Clean(dbPath)
		if !filepath.IsAbs(cleanPath) {
			if cwd, err := os.Getwd(); err == nil {
				cleanPath = filepath.Join(cwd, cleanPath)
			}
		}
		// Ensure the directory exists
		if dir := filepath.Dir(cleanPath); dir != "" {
			if err := os.MkdirAll(dir, 0755); err == nil {
				return cleanPath
			}
		}
		return cleanPath
	}

	// Second priority: check for workspace database in current working directory
	cwd, err := os.Getwd()
	if err == nil {
		workspaceDBPath := filepath.Join(cwd, ".gorev", "gorev.db")
		if _, err := os.Stat(workspaceDBPath); err == nil {
			return workspaceDBPath
		}

		// Also check parent directories for workspace database (monorepo support)
		parentPaths := []string{
			filepath.Join(filepath.Dir(cwd), ".gorev", "gorev.db"),
			filepath.Join(filepath.Dir(filepath.Dir(cwd)), ".gorev", "gorev.db"),
			filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(cwd))), ".gorev", "gorev.db"),
		}

		for _, path := range parentPaths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	// Third priority: check if database exists in user's home directory
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
	if cwd != "" {
		// Check if we're in the project root
		migrationsPath := filepath.Join(cwd, "internal", "veri", "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return filepath.Join(cwd, "gorev.db")
		}
	}

	// Last resort: use current directory
	return "gorev.db"
}

// initWorkspaceDatabase initializes a database in the specified location
func initWorkspaceDatabase(global bool) error {
	var dbPath string
	var dbDir string

	// Check GOREV_DB_PATH environment variable first
	if envDBPath := os.Getenv("GOREV_DB_PATH"); envDBPath != "" {
		dbPath = envDBPath
		dbDir = filepath.Dir(dbPath)
	} else if global {
		// Global database in user home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		dbDir = filepath.Join(homeDir, ".gorev")
		dbPath = filepath.Join(dbDir, "gorev.db")
	} else {
		// Workspace database in current directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		dbDir = filepath.Join(cwd, ".gorev")
		dbPath = filepath.Join(dbDir, "gorev.db")
	}

	// Check if database already exists
	if _, err := os.Stat(dbPath); err == nil {
		return fmt.Errorf("database already exists: %s", dbPath)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dbDir, err)
	}

	// Initialize database by creating a VeriYonetici instance
	migrationsPath := getMigrationsPath()
	var veriYonetici *gorev.VeriYonetici
	var err error

	if migrationsPath == "embedded://migrations" {
		// Use embedded migrations
		migrationsFS, fsErr := getEmbeddedMigrationsFS()
		if fsErr != nil {
			return fmt.Errorf("failed to get embedded migrations: %w", fsErr)
		}
		veriYonetici, err = gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
	} else {
		// Fallback to filesystem migrations
		veriYonetici, err = gorev.YeniVeriYonetici(dbPath, migrationsPath)
	}

	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() { _ = veriYonetici.Kapat() }()

	fmt.Printf("Database initialized successfully: %s\n", dbPath)
	return nil
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
		Use:     "gorev",
		Short:   i18n.T("cli.appDescription"),
		Long:    i18n.T("cli.appLongDescription"),
		Example: i18n.T("cli.appExamples"),
	}

	serveCmd := &cobra.Command{
		Use:     "serve",
		Short:   i18n.T("cli.serve"),
		Long:    i18n.T("cli.serveLongDescription"),
		Example: i18n.T("cli.serveExamples"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Re-initialize i18n if language flag was provided
			if langFlag != "" {
				if err := i18n.SetLanguage(langFlag); err != nil {
					return fmt.Errorf("failed to set language to %s: %w", langFlag, err)
				}
			}

			// Setup server configuration
			cfg := config.DefaultConfig()
			if serverModeFlag != "" {
				cfg.Mode = config.ServerMode(serverModeFlag)
			}
			if dbPathFlag != "" {
				cfg.CentralizedDBPath = dbPathFlag
			}
			cfg.Port = apiPortFlag
			cfg.AllowLocalPaths = cfg.Mode == config.ModeLocal
			config.SetGlobalConfig(cfg)

			log.Printf("🚀 Starting Gorev Server (mode: %s)...", cfg.Mode)

			return runServer()
		},
	}
	serveCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, i18n.T("cli.debug"))
	serveCmd.PersistentFlags().StringVar(&apiPortFlag, "api-port", "5082", "API server port")
	serveCmd.PersistentFlags().BoolVar(&noAPIFlag, "no-api", false, "Disable API server (MCP only)")
	serveCmd.PersistentFlags().StringVar(&serverModeFlag, "mode", "", "Server mode: local (default) or centralized")
	serveCmd.PersistentFlags().StringVar(&dbPathFlag, "db-path", "", "Database path (for centralized mode)")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("cli.version"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Gorev %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Git Commit: %s\n", gitCommit)
		},
	}

	// Init command
	var globalFlag bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize database",
		Long:  "Initialize a new Gorev database in the current workspace or global location.",
		Example: `  # Initialize workspace database in .gorev/gorev.db
  gorev init

  # Initialize global database in ~/.gorev/gorev.db
  gorev init --global`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initWorkspaceDatabase(globalFlag)
		},
	}
	initCmd.Flags().BoolVar(&globalFlag, "global", false, "Initialize global database in user home directory")

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

	templateAliasesCmd := &cobra.Command{
		Use:   "aliases",
		Short: i18n.T("cli.templateAliases"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTemplateAliases()
		},
	}

	templateCmd.AddCommand(templateListCmd, templateShowCmd, templateInitCmd, templateAliasesCmd)

	// MCP test commands
	mcpCmd := createMCPCommand()

	// IDE management commands
	ideCmd := createIDECommand()

	// Daemon management commands
	daemonCmd := createDaemonCommand()
	daemonStopCmd := createDaemonStopCommand()
	daemonStatusCmd := createDaemonStatusCommand()

	// MCP proxy command
	mcpProxyCmd := createMCPProxyCommand()

	// Seed test data command
	seedCmd := createSeedCommand()

	// Global flags
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", i18n.T("flags.language"))

	rootCmd.AddCommand(serveCmd, versionCmd, initCmd, templateCmd, mcpCmd, ideCmd, daemonCmd, daemonStopCmd, daemonStatusCmd, mcpProxyCmd, seedCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	// Veritabanını başlat
	dbPath := getDatabasePath()
	log.Printf("Using database: %s", dbPath)

	// Use embedded migrations if available, fallback to filesystem
	migrationsPath := getMigrationsPath()
	var veriYonetici *gorev.VeriYonetici
	var err error

	if migrationsPath == "embedded://migrations" {
		// Use embedded migrations
		migrationsFS, fsErr := getEmbeddedMigrationsFS()
		if fsErr != nil {
			return errors.New(i18n.T("error.dataManagerInit", map[string]interface{}{"Error": fsErr}))
		}
		veriYonetici, err = gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
	} else {
		// Fallback to filesystem migrations
		veriYonetici, err = gorev.YeniVeriYonetici(dbPath, migrationsPath)
	}

	if err != nil {
		return errors.New(i18n.T("error.dataManagerInit", map[string]interface{}{"Error": err}))
	}
	defer func() { _ = veriYonetici.Kapat() }()

	// İş mantığı servisini oluştur
	// Centralized modda: default workspace_id ile IsYonetici oluştur
	var isYonetici *gorev.IsYonetici
	cfg := config.GetGlobalConfig()
	if cfg.Mode == config.ModeCentralized {
		defaultWorkspaceID := os.Getenv("GOREV_WORKSPACE_ID")
		if defaultWorkspaceID == "" {
			defaultWorkspaceID = "default"
		}
		log.Printf("📦 Centralized mode: using workspace_id=%s", defaultWorkspaceID)
		isYonetici = gorev.YeniIsYoneticiWithWorkspaceID(veriYonetici, defaultWorkspaceID)
	} else {
		isYonetici = gorev.YeniIsYonetici(veriYonetici)
	}

	// IDE extension durumunu kontrol et (background'da)
	if !debugFlag {
		go checkAndPromptIDEExtensions()
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start API server (unless disabled)
	var apiServer *api.APIServer
	if !noAPIFlag {
		apiServer = api.NewAPIServer(apiPortFlag, isYonetici)

		// Set migrations FS for workspace manager
		if migrationsPath == "embedded://migrations" {
			migrationsFS, fsErr := getEmbeddedMigrationsFS()
			if fsErr != nil {
				log.Fatalf("Failed to get embedded migrations for workspace manager: %v", fsErr)
			}
			apiServer.SetMigrationsFS(migrationsFS)
		}

		// Serve embedded web UI static files
		if err := api.ServeStaticFiles(apiServer.App(), WebDistFS); err != nil {
			log.Printf("Warning: Failed to serve web UI: %v", err)
		}

		apiServer.StartAsync()
		log.Printf("🚀 Unified Gorev Server started:")
		log.Printf("📱 Web UI: http://localhost:%s", apiPortFlag)
		log.Printf("🔧 API: http://localhost:%s/api/v1", apiPortFlag)
		log.Printf("🔗 MCP: stdio (for AI assistants)")
	} else {
		log.Printf("🔗 Gorev MCP Server started (API disabled)")
	}

	// Debug mode mesajı
	if debugFlag {
		log.Printf("Debug mode enabled")
		if langFlag != "" {
			log.Printf("Language: %s", langFlag)
		}
	}

	// Start MCP server in a goroutine
	mcpErrChan := make(chan error, 1)
	go func() {
		sunucu, err := mcp.YeniMCPSunucuWithDebug(isYonetici, debugFlag)
		if err != nil {
			mcpErrChan <- fmt.Errorf("MCP server creation failed: %w", err)
			return
		}

		if err := mcp.ServeSunucu(sunucu); err != nil {
			mcpErrChan <- fmt.Errorf("MCP server failed: %w", err)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
	case err := <-mcpErrChan:
		log.Printf("MCP server error: %v", err)
		return err
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if apiServer != nil && !noAPIFlag {
		if err := apiServer.Shutdown(ctx); err != nil {
			log.Printf("API server shutdown error: %v", err)
		}
	}

	log.Println("🔽 Gorev server stopped")
	return nil
}

// createVeriYonetici creates a VeriYonetici instance using embedded migrations
func createVeriYonetici() (*gorev.VeriYonetici, error) {
	dbPath := getDatabasePath()
	migrationsPath := getMigrationsPath()

	if migrationsPath == "embedded://migrations" {
		// Use embedded migrations
		migrationsFS, err := getEmbeddedMigrationsFS()
		if err != nil {
			return nil, fmt.Errorf("failed to get embedded migrations: %w", err)
		}
		return gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
	} else {
		// Fallback to filesystem migrations
		return gorev.YeniVeriYonetici(dbPath, migrationsPath)
	}
}

// checkAndPromptIDEExtensions checks for IDEs and prompts for extension installation
func checkAndPromptIDEExtensions() {
	// 3 saniye bekle (server'ın tam olarak başlaması için)
	time.Sleep(3 * time.Second)

	detector := gorev.NewIDEDetector()
	detectedIDEs, err := detector.DetectAllIDEs()
	if err != nil {
		return // Sessizce çık, hata vermek server başlangıcını etkilemesin
	}

	if len(detectedIDEs) == 0 {
		return // Hiç IDE yok, çık
	}

	installer := gorev.NewExtensionInstaller(detector)
	defer func() {
		if err := installer.Cleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup error: %v\n", err)
		}
	}()

	// Kurulu olmayan IDE'leri bul
	var uninstalledIDEs []gorev.IDEType
	var outdatedIDEs []gorev.IDEType

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// GitHub'dan en son versiyon bilgisini al
	latestExtension, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	if err != nil {
		return // Network hatası varsa sessizce çık
	}

	for ideType, ide := range detectedIDEs {
		isInstalled, err := detector.IsExtensionInstalled(ideType, "mehmetsenol.gorev-vscode")
		if err != nil || !isInstalled {
			uninstalledIDEs = append(uninstalledIDEs, ideType)
		} else {
			// Version kontrolü
			installedVersion, err := detector.GetExtensionVersion(ideType, "mehmetsenol.gorev-vscode")
			if err == nil && installedVersion != latestExtension.Version {
				outdatedIDEs = append(outdatedIDEs, ideType)
			}
		}
		_ = ide // unused variable fix
	}

	// Kurulu olmayan IDE'ler için prompt
	if len(uninstalledIDEs) > 0 {
		fmt.Fprintf(os.Stderr, "\n🔍 %d IDE detected without Gorev extension:\n", len(uninstalledIDEs))
		for _, ideType := range uninstalledIDEs {
			ide := detectedIDEs[ideType]
			fmt.Fprintf(os.Stderr, "   %s %s\n", getIDEIconForCLI(ideType), ide.Name)
		}
		fmt.Fprintf(os.Stderr, "\n💡 Install with: gorev ide install\n\n")
	}

	// Güncellenebilir IDE'ler için prompt
	if len(outdatedIDEs) > 0 {
		fmt.Fprintf(os.Stderr, "\n🔄 %d IDE has outdated Gorev extension:\n", len(outdatedIDEs))
		for _, ideType := range outdatedIDEs {
			ide := detectedIDEs[ideType]
			fmt.Fprintf(os.Stderr, "   %s %s\n", getIDEIconForCLI(ideType), ide.Name)
		}
		fmt.Fprintf(os.Stderr, "\n💡 Update with: gorev ide update\n\n")
	}
}

func listTemplates(kategori string) error {
	// Veritabanını başlat
	veriYonetici, err := createVeriYonetici()
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer func() { _ = veriYonetici.Kapat() }()

	// Template'leri listele
	templates, err := veriYonetici.TemplateListele(context.Background(), kategori)
	if err != nil {
		return fmt.Errorf("template'ler listelenemedi: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println(i18n.T("display.noTemplates"))
		return nil
	}

	// Kategorilere göre grupla
	kategoriMap := make(map[string][]*gorev.GorevTemplate)
	for _, tmpl := range templates {
		kategoriMap[tmpl.Category] = append(kategoriMap[tmpl.Category], tmpl)
	}

	// Her kategoriyi göster
	for kat, tmpls := range kategoriMap {
		fmt.Printf("\n=== %s ===\n", kat)
		for _, tmpl := range tmpls {
			aliasInfo := ""
			if tmpl.Alias != "" {
				aliasInfo = fmt.Sprintf(" [%s: %s]", i18n.T("display.alias"), tmpl.Alias)
			}
			fmt.Printf("\n%s%s\n", tmpl.Name, aliasInfo)
			fmt.Printf("  %s\n", tmpl.Definition)
			fmt.Println("  " + i18n.T("display.title", map[string]interface{}{"Title": tmpl.DefaultTitle}))
			if tmpl.Alias != "" {
				fmt.Println("  " + i18n.T("display.quickUsage", map[string]interface{}{"Alias": tmpl.Alias}))
			}
		}
	}

	fmt.Println("\n" + i18n.T("display.detailInfo"))
	return nil
}

func showTemplate(templateID string) error {
	// Veritabanını başlat
	veriYonetici, err := createVeriYonetici()
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer func() { _ = veriYonetici.Kapat() }()

	// Template'i ID veya alias ile getir
	template, err := veriYonetici.TemplateIDVeyaAliasIleGetir(context.Background(), templateID)
	if err != nil {
		return fmt.Errorf("template bulunamadı: %w", err)
	}

	fmt.Printf("\n=== %s ===\n", template.Name)
	fmt.Println(i18n.T("display.id", map[string]interface{}{"ID": template.ID}))
	if template.Alias != "" {
		fmt.Printf("%s: %s\n", strings.Title(i18n.T("display.alias")), template.Alias)
	}
	fmt.Println(i18n.T("display.category", map[string]interface{}{"Category": template.Category}))
	fmt.Println(i18n.T("display.desc", map[string]interface{}{"Description": template.Definition}))
	fmt.Println("\n" + i18n.T("display.titleTemplate", map[string]interface{}{"Template": template.DefaultTitle}))

	fmt.Println("\n" + i18n.T("display.fields"))
	for _, alan := range template.Fields {
		zorunlu := ""
		if alan.Required {
			zorunlu = fmt.Sprintf(" (%s)", i18n.T("display.required"))
		}
		fmt.Printf("  - %s (%s)%s", alan.Name, alan.Type, zorunlu)
		if alan.Default != "" {
			fmt.Printf(" [%s: %s]", i18n.T("display.default"), alan.Default)
		}
		if len(alan.Options) > 0 {
			fmt.Printf("\n    %s", i18n.T("display.options", map[string]interface{}{"Options": fmt.Sprintf("%v", alan.Options)}))
		}
		fmt.Println()
	}

	fmt.Println("\n" + i18n.T("display.exampleDescription"))
	fmt.Println(template.DescriptionTemplate)

	return nil
}

func initTemplates() error {
	// Veritabanını başlat
	veriYonetici, err := createVeriYonetici()
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer func() { _ = veriYonetici.Kapat() }()

	// Varsayılan template'leri oluştur
	err = veriYonetici.VarsayilanTemplateleriOlustur(context.Background())
	if err != nil {
		return fmt.Errorf("template'ler oluşturulamadı: %w", err)
	}

	fmt.Println(i18n.T("success.defaultTemplatesCreated"))
	return nil
}

func listTemplateAliases() error {
	// Veritabanını başlat
	veriYonetici, err := createVeriYonetici()
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer func() { _ = veriYonetici.Kapat() }()

	// Template'leri listele
	templates, err := veriYonetici.TemplateListele(context.Background(), "")
	if err != nil {
		return fmt.Errorf("template'ler listelenemedi: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println(i18n.T("display.noTemplates"))
		return nil
	}

	fmt.Println(i18n.T("display.templateAliasesHeader"))
	fmt.Println("=" + strings.Repeat("=", 40))

	// Alias'leri topla
	aliases := make(map[string]*gorev.GorevTemplate)
	for _, tmpl := range templates {
		if tmpl.Alias != "" {
			aliases[tmpl.Alias] = tmpl
		}
	}

	if len(aliases) == 0 {
		fmt.Println(i18n.T("display.noAliasTemplates"))
		return nil
	}

	// Alias'leri alfabetik sırala ve göster
	for alias, tmpl := range aliases {
		fmt.Printf("\n🏷️  %s\n", alias)
		fmt.Printf("   → %s (%s)\n", tmpl.Name, tmpl.Category)
		fmt.Printf("   📝 %s\n", tmpl.Definition)
		fmt.Println("   " + i18n.T("display.templateAliasesInfo", map[string]interface{}{"ID": alias}))
	}

	fmt.Println("\n" + i18n.T("display.templatesListInfo"))
	fmt.Println(i18n.T("display.allTemplatesInfo"))

	return nil
}

// createIDECommand creates IDE management commands
func createIDECommand() *cobra.Command {
	ideCmd := &cobra.Command{
		Use:   "ide",
		Short: i18n.T("cli.ide"),
		Long:  i18n.T("cli.ideDescription"),
	}

	// IDE detect command
	ideDetectCmd := &cobra.Command{
		Use:   "detect",
		Short: i18n.T("cli.ideDetect"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIDEDetect()
		},
	}

	// IDE install command
	ideInstallCmd := &cobra.Command{
		Use:   "install [ide-type]",
		Short: i18n.T("cli.ideInstall"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ideType := "all"
			if len(args) > 0 {
				ideType = args[0]
			}
			return runIDEInstall(ideType)
		},
	}

	// IDE uninstall command
	ideUninstallCmd := &cobra.Command{
		Use:   "uninstall <ide-type>",
		Short: i18n.T("cli.ideUninstall"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIDEUninstall(args[0])
		},
	}

	// IDE status command
	ideStatusCmd := &cobra.Command{
		Use:   "status",
		Short: i18n.T("cli.ideStatus"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIDEStatus()
		},
	}

	// IDE update command
	ideUpdateCmd := &cobra.Command{
		Use:   "update [ide-type]",
		Short: i18n.T("cli.ideUpdate"),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ideType := "all"
			if len(args) > 0 {
				ideType = args[0]
			}
			return runIDEUpdate(ideType)
		},
	}

	// IDE config command
	ideConfigCmd := &cobra.Command{
		Use:   "config",
		Short: i18n.T("cli.ideConfig"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIDEConfig()
		},
	}

	// IDE config set command
	ideConfigSetCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: i18n.T("cli.ideConfigSet"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIDEConfigSet(args[0], args[1])
		},
	}

	ideConfigCmd.AddCommand(ideConfigSetCmd)

	ideCmd.AddCommand(ideDetectCmd, ideInstallCmd, ideUninstallCmd, ideStatusCmd, ideUpdateCmd, ideConfigCmd)
	return ideCmd
}

// runIDEDetect runs IDE detection
func runIDEDetect() error {
	detector := gorev.NewIDEDetector()
	detectedIDEs, err := detector.DetectAllIDEs()
	if err != nil {
		return err
	}

	if len(detectedIDEs) == 0 {
		fmt.Println(i18n.T("display.noIDEsDetected"))
		return nil
	}

	fmt.Println("🔍 Detected IDEs:")
	for ideType, ide := range detectedIDEs {
		icon := getIDEIconForCLI(ideType)
		fmt.Printf("\n%s %s\n", icon, ide.Name)
		fmt.Printf("   Path: %s\n", ide.ExecutablePath)
		if ide.Version != "unknown" && ide.Version != "" {
			fmt.Printf("   Version: %s\n", ide.Version)
		}
	}
	return nil
}

// runIDEInstall installs extension to IDE(s)
func runIDEInstall(ideType string) error {
	detector := gorev.NewIDEDetector()
	if _, err := detector.DetectAllIDEs(); err != nil {
		return err
	}

	installer := gorev.NewExtensionInstaller(detector)
	defer func() {
		if err := installer.Cleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup error: %v\n", err)
		}
	}()

	// Get latest extension info
	ctx := context.Background()
	extensionInfo, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	if err != nil {
		return fmt.Errorf("failed to get extension info: %w", err)
	}

	fmt.Printf("📦 Installing Gorev Extension v%s...\n\n", extensionInfo.Version)

	var results []gorev.InstallResult
	if ideType == "all" {
		results, err = installer.InstallToAllIDEs(ctx, extensionInfo)
	} else {
		result, installErr := installer.InstallExtension(ctx, gorev.IDEType(ideType), extensionInfo)
		if installErr != nil {
			err = installErr
		}
		if result != nil {
			results = append(results, *result)
		}
	}

	if err != nil && len(results) == 0 {
		return err
	}

	// Print results
	successCount := 0
	for _, result := range results {
		icon := "❌"
		if result.Success {
			icon = "✅"
			successCount++
		}
		fmt.Printf("%s %s: %s\n", icon, result.IDE, result.Message)
	}

	fmt.Printf("\nSummary: %d/%d installations successful\n", successCount, len(results))
	return nil
}

// runIDEUninstall uninstalls extension from IDE
func runIDEUninstall(ideType string) error {
	detector := gorev.NewIDEDetector()
	if _, err := detector.DetectAllIDEs(); err != nil {
		return err
	}

	installer := gorev.NewExtensionInstaller(detector)

	result, err := installer.UninstallExtension(gorev.IDEType(ideType), "mehmetsenol.gorev-vscode")
	if err != nil {
		return err
	}

	icon := "❌"
	if result.Success {
		icon = "✅"
	}
	fmt.Printf("%s %s: %s\n", icon, result.IDE, result.Message)
	return nil
}

// runIDEStatus shows extension status
func runIDEStatus() error {
	ctx := context.Background()
	detector := gorev.NewIDEDetector()
	detectedIDEs, err := detector.DetectAllIDEs()
	if err != nil {
		return err
	}

	installer := gorev.NewExtensionInstaller(detector)

	// Get latest version from GitHub
	extensionInfo, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	var latestVersion string
	if err == nil {
		latestVersion = extensionInfo.Version
	}

	fmt.Println("📊 Extension Status Report")
	if latestVersion != "" {
		fmt.Printf("Latest Available Version: v%s\n", latestVersion)
	}
	fmt.Println()

	for ideType, ide := range detectedIDEs {
		icon := getIDEIconForCLI(ideType)
		fmt.Printf("%s %s\n", icon, ide.Name)

		// Check if extension is installed
		isInstalled, err := detector.IsExtensionInstalled(ideType, "mehmetsenol.gorev-vscode")
		if err != nil {
			fmt.Printf("   Status: ❌ Error checking installation: %s\n", err)
		} else if !isInstalled {
			fmt.Printf("   Status: ❌ Not installed\n")
		} else {
			// Get installed version
			installedVersion, err := detector.GetExtensionVersion(ideType, "mehmetsenol.gorev-vscode")
			if err != nil {
				fmt.Printf("   Status: ✅ Installed (version unknown)\n")
			} else {
				statusIcon := "✅"
				updateStatus := ""

				if latestVersion != "" && installedVersion != latestVersion {
					statusIcon = "⚠️"
					updateStatus = fmt.Sprintf(" (Update available: v%s)", latestVersion)
				}

				fmt.Printf("   Status: %s Installed v%s%s\n", statusIcon, installedVersion, updateStatus)
			}
		}
		fmt.Println()
	}

	return nil
}

// runIDEUpdate updates extension to latest version
func runIDEUpdate(ideType string) error {
	detector := gorev.NewIDEDetector()
	if _, err := detector.DetectAllIDEs(); err != nil {
		return err
	}

	installer := gorev.NewExtensionInstaller(detector)
	defer func() {
		if err := installer.Cleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup error: %v\n", err)
		}
	}()

	// Get latest extension info
	ctx := context.Background()
	extensionInfo, err := installer.GetLatestExtensionInfo(ctx, "msenol", "Gorev")
	if err != nil {
		return fmt.Errorf("failed to get extension info: %w", err)
	}

	fmt.Printf("🔄 Updating Gorev Extension to v%s...\n\n", extensionInfo.Version)

	var results []gorev.InstallResult
	if ideType == "all" {
		// Update all detected IDEs
		allIDEs := detector.GetAllDetectedIDEs()
		for ideTypeKey := range allIDEs {
			result, err := installer.InstallExtension(ctx, ideTypeKey, extensionInfo)
			if result != nil {
				results = append(results, *result)
			}
			if err != nil {
				fmt.Printf("❌ Error updating %s: %s\n", ideTypeKey, err)
			}
		}
	} else {
		result, err := installer.InstallExtension(ctx, gorev.IDEType(ideType), extensionInfo)
		if result != nil {
			results = append(results, *result)
		}
		if err != nil {
			return err
		}
	}

	// Print results
	successCount := 0
	for _, result := range results {
		icon := "❌"
		if result.Success {
			icon = "✅"
			successCount++
		}
		fmt.Printf("%s %s: %s\n", icon, result.IDE, result.Message)
	}

	fmt.Printf("\nSummary: %d/%d updates successful\n", successCount, len(results))
	return nil
}

// runIDEConfig shows current IDE configuration
func runIDEConfig() error {
	configManager := gorev.NewIDEConfigManager()
	if err := configManager.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	config := configManager.GetConfig()

	fmt.Println("🔧 IDE Configuration")
	fmt.Printf("Config file: %s\n\n", configManager.GetConfigPath())

	fmt.Printf("Auto Install:    %t\n", config.AutoInstall)
	fmt.Printf("Auto Update:     %t\n", config.AutoUpdate)
	fmt.Printf("Check Interval:  %v\n", config.CheckInterval)
	fmt.Printf("Disable Prompts: %t\n", config.DisablePrompts)
	fmt.Printf("Extension ID:    %s\n", config.ExtensionID)
	fmt.Printf("Supported IDEs:  %v\n", config.SupportedIDEs)
	fmt.Printf("Last Check:      %v\n", config.LastUpdateCheck.Format("2006-01-02 15:04:05"))

	fmt.Println("\n💡 Configure with: gorev ide config set <key> <value>")
	fmt.Println("Available keys: auto_install, auto_update, disable_prompts")

	return nil
}

// runIDEConfigSet sets a configuration value
func runIDEConfigSet(key, value string) error {
	configManager := gorev.NewIDEConfigManager()
	if err := configManager.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var err error
	switch key {
	case "auto_install":
		boolValue := value == "true" || value == "1" || value == "yes"
		err = configManager.SetAutoInstall(boolValue)
		fmt.Printf("✅ Auto install set to: %t\n", boolValue)
	case "auto_update":
		boolValue := value == "true" || value == "1" || value == "yes"
		err = configManager.SetAutoUpdate(boolValue)
		fmt.Printf("✅ Auto update set to: %t\n", boolValue)
	case "disable_prompts":
		boolValue := value == "true" || value == "1" || value == "yes"
		err = configManager.SetDisablePrompts(boolValue)
		fmt.Printf("✅ Disable prompts set to: %t\n", boolValue)
	case "check_interval":
		duration, parseErr := time.ParseDuration(value)
		if parseErr != nil {
			return fmt.Errorf("invalid duration format: %s (example: 24h, 30m)", value)
		}
		err = configManager.SetCheckInterval(duration)
		fmt.Printf("✅ Check interval set to: %v\n", duration)
	default:
		return fmt.Errorf("unknown config key: %s. Available: auto_install, auto_update, disable_prompts, check_interval", key)
	}

	return err
}

// getIDEIconForCLI returns CLI-friendly icons for IDE types
func getIDEIconForCLI(ideType gorev.IDEType) string {
	switch ideType {
	case gorev.IDETypeVSCode:
		return "🔵"
	case gorev.IDETypeCursor:
		return "🖱️"
	case gorev.IDETypeWindsurf:
		return "🌊"
	default:
		return "💻"
	}
}
