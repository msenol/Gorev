package gorev

import (
	"database/sql"
	"fmt"

	"github.com/msenol/gorev/internal/i18n"
)

// AktifProjeAyarla aktif projeyi ayarlar
func (vy *VeriYonetici) AktifProjeAyarla(projeID string) error {
	// Önce projenin var olduğunu kontrol et
	var count int
	err := vy.db.QueryRow("SELECT COUNT(*) FROM projeler WHERE id = ?", projeID).Scan(&count)
	if err != nil {
		return fmt.Errorf(i18n.T("error.check_failed", map[string]interface{}{"Entity": "proje", "Error": err}))
	}
	if count == 0 {
		return fmt.Errorf(i18n.T("error.projectNotFoundId", map[string]interface{}{"Id": projeID}))
	}

	// Aktif proje tablosunu güncelle (INSERT OR REPLACE)
	sorgu := `INSERT OR REPLACE INTO aktif_proje (id, project_id) VALUES (1, ?)`
	_, err = vy.db.Exec(sorgu, projeID)
	if err != nil {
		return fmt.Errorf(i18n.T("error.activeProjectSetFailed", map[string]interface{}{"Error": err}))
	}

	return nil
}

// AktifProjeGetir aktif projeyi getirir
func (vy *VeriYonetici) AktifProjeGetir() (string, error) {
	var projeID string
	err := vy.db.QueryRow("SELECT project_id FROM aktif_proje WHERE id = 1").Scan(&projeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Aktif proje yok
		}
		return "", fmt.Errorf(i18n.T("error.activeProjectGetFailed", map[string]interface{}{"Error": err}))
	}
	return projeID, nil
}

// AktifProjeKaldir aktif proje ayarını kaldırır
func (vy *VeriYonetici) AktifProjeKaldir() error {
	_, err := vy.db.Exec("DELETE FROM aktif_proje WHERE id = 1")
	if err != nil {
		return fmt.Errorf(i18n.T("error.activeProjectRemoveFailed", map[string]interface{}{"Error": err}))
	}
	return nil
}
