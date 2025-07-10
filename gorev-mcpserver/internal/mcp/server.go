package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/gorev"
)

// YeniMCPSunucu yeni bir MCP sunucusu oluşturur
func YeniMCPSunucu(isYonetici *gorev.IsYonetici) (*server.MCPServer, error) {
	// MCP sunucusunu oluştur
	s := server.NewMCPServer("gorev", "1.0.0")

	// Handler'ları oluştur ve kaydet
	handlers := YeniHandlers(isYonetici)
	handlers.RegisterTools(s)

	return s, nil
}

// ServeSunucu MCP sunucusunu stdio üzerinden çalıştırır
func ServeSunucu(s *server.MCPServer) error {
	return server.ServeStdio(s)
}

// NewServer creates a new MCP server with the given handlers
func NewServer(handlers *Handlers) *server.MCPServer {
	s := server.NewMCPServer("gorev", "1.0.0")
	handlers.RegisterTools(s)
	return s
}

// Tool represents an MCP tool
type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

// ListTools returns all registered tools
func ListTools() []Tool {
	// This is a placeholder - we need to implement this based on registered tools
	// For now, return a hardcoded list
	return []Tool{
		{Name: "gorev_olustur", Description: "Yeni bir görev oluştur"},
		{Name: "gorev_listele", Description: "Görevleri listele"},
		{Name: "gorev_detay", Description: "Görev detayını göster"},
		{Name: "gorev_guncelle", Description: "Görev durumunu güncelle"},
		{Name: "gorev_duzenle", Description: "Görev bilgilerini düzenle"},
		{Name: "gorev_sil", Description: "Görevi sil"},
		{Name: "gorev_bagimlilik_ekle", Description: "Göreve bağımlılık ekle"},
		{Name: "gorev_altgorev_olustur", Description: "Alt görev oluştur"},
		{Name: "gorev_ust_degistir", Description: "Görevin üst görevini değiştir"},
		{Name: "gorev_hiyerarsi_goster", Description: "Görev hiyerarşisini göster"},
		{Name: "proje_olustur", Description: "Yeni proje oluştur"},
		{Name: "proje_listele", Description: "Projeleri listele"},
		{Name: "proje_gorevleri", Description: "Proje görevlerini listele"},
		{Name: "proje_aktif_yap", Description: "Projeyi aktif yap"},
		{Name: "aktif_proje_goster", Description: "Aktif projeyi göster"},
		{Name: "aktif_proje_kaldir", Description: "Aktif proje ayarını kaldır"},
		{Name: "ozet_goster", Description: "Özet istatistikleri göster"},
		{Name: "template_listele", Description: "Görev şablonlarını listele"},
		{Name: "templateden_gorev_olustur", Description: "Şablondan görev oluştur"},
		{Name: "gorev_set_active", Description: "Aktif görevi ayarla"},
		{Name: "gorev_get_active", Description: "Aktif görevi getir"},
		{Name: "gorev_recent", Description: "Son görevleri listele"},
		{Name: "gorev_context_summary", Description: "AI bağlam özetini göster"},
		{Name: "gorev_batch_update", Description: "Toplu görev güncelleme"},
		{Name: "gorev_nlp_query", Description: "Doğal dil sorgusu ile görev ara"},
	}
}
