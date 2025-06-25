package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/gorev/internal/mcp"
	"github.com/yourusername/gorev/internal/gorev"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gorev",
		Short: "Güçlü ve esnek görev yönetimi MCP sunucusu",
		Long: `Gorev, Model Context Protocol (MCP) üzerinden AI asistanlarına
görev yönetimi yetenekleri sağlayan modern bir sunucudur.`,
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "MCP sunucusunu başlat",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer()
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Versiyon bilgisini göster",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Gorev %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Git Commit: %s\n", gitCommit)
		},
	}

	rootCmd.AddCommand(serveCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	// Veritabanını başlat
	veriYonetici, err := gorev.YeniVeriYonetici("gorev.db")
	if err != nil {
		return fmt.Errorf("veri yönetici başlatılamadı: %w", err)
	}
	defer veriYonetici.Kapat()

	// İş mantığı servisini oluştur
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// MCP sunucusunu başlat
	sunucu, err := mcp.YeniMCPSunucu(isYonetici)
	if err != nil {
		return fmt.Errorf("MCP sunucusu oluşturulamadı: %w", err)
	}

	// Sunucuyu çalıştır
	fmt.Fprintln(os.Stderr, "Gorev MCP sunucusu başlatılıyor...")
	if err := mcp.ServeSunucu(sunucu); err != nil {
		return fmt.Errorf("sunucu başlatılamadı: %w", err)
	}

	return nil
}