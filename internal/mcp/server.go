package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/yourusername/gorev/internal/gorev"
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