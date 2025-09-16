package gorev

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/msenol/gorev/internal/i18n"
)

// SearchFilters represents search filter criteria
type SearchFilters struct {
	Status       []string `json:"status,omitempty"`
	Priority     []string `json:"priority,omitempty"`
	ProjectIDs   []string `json:"project_ids,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	EnableFuzzy  bool     `json:"enable_fuzzy,omitempty"`
	FuzzyThreshold int    `json:"fuzzy_threshold,omitempty"`
	CreatedAfter  string  `json:"created_after,omitempty"`
	CreatedBefore string  `json:"created_before,omitempty"`
	DueAfter     string   `json:"due_after,omitempty"`
	DueBefore    string   `json:"due_before,omitempty"`
}

// FilterProfile represents a saved filter configuration
type FilterProfile struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Filters     SearchFilters `json:"filters"`
	SearchQuery string        `json:"search_query,omitempty"`
	IsDefault   bool          `json:"is_default,omitempty"`
	UseCount    int           `json:"use_count,omitempty"`
	LastUsedAt  *time.Time    `json:"last_used_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// FilterProfileManager handles filter profile operations
type FilterProfileManager struct {
	db *sql.DB
}

// NewFilterProfileManager creates a new filter profile manager
func NewFilterProfileManager(db *sql.DB) *FilterProfileManager {
	return &FilterProfileManager{
		db: db,
	}
}

// SaveFilterProfile saves a new filter profile or updates existing one
func (fpm *FilterProfileManager) SaveFilterProfile(profile *FilterProfile) error {
	if profile.ID == "" {
		// Create new profile
		return fpm.createFilterProfile(profile)
	} else {
		// Update existing profile
		return fpm.updateFilterProfile(profile)
	}
}

// CreateProfile creates a new filter profile and returns it
func (fpm *FilterProfileManager) CreateProfile(profile FilterProfile) (*FilterProfile, error) {
	// Generate new ID
	profile.ID = fmt.Sprintf("fp_%d", time.Now().UnixNano())
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	err := fpm.createFilterProfile(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// GetProfile retrieves a filter profile by ID
func (fpm *FilterProfileManager) GetProfile(id string) (*FilterProfile, error) {
	if id == "" {
		return nil, fmt.Errorf(i18n.T("error.filterProfileIdRequired"))
	}

	query := `
		SELECT id, name, description, filters, search_query, is_default, created_at, updated_at
		FROM filter_profiles
		WHERE id = ?
	`

	var profile FilterProfile
	var filtersJSON string
	var searchQuery sql.NullString

	err := fpm.db.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Description,
		&filtersJSON,
		&searchQuery,
		&profile.IsDefault,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(i18n.T("filter_profile_not_found", map[string]interface{}{"id": id}))
		}
		return nil, fmt.Errorf(i18n.T("error.filterProfileGetFailed", map[string]interface{}{"Error": err}))
	}

	// Unmarshal filters
	err = json.Unmarshal([]byte(filtersJSON), &profile.Filters)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileUnmarshalFailed", map[string]interface{}{"Error": err}))
	}

	if searchQuery.Valid {
		profile.SearchQuery = searchQuery.String
	}

	return &profile, nil
}

// ListProfiles retrieves all filter profiles
func (fpm *FilterProfileManager) ListProfiles() ([]FilterProfile, error) {
	query := `
		SELECT id, name, description, filters, search_query, is_default, use_count, last_used_at, created_at, updated_at
		FROM filter_profiles
		ORDER BY name
	`

	rows, err := fpm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileListFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var profiles []FilterProfile
	for rows.Next() {
		var profile FilterProfile
		var filtersJSON string
		var searchQuery sql.NullString
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Description,
			&filtersJSON,
			&searchQuery,
			&profile.IsDefault,
			&profile.UseCount,
			&lastUsedAt,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.filterProfileScanFailed", map[string]interface{}{"Error": err}))
		}

		// Unmarshal filters
		err = json.Unmarshal([]byte(filtersJSON), &profile.Filters)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.filterProfileUnmarshalFailed", map[string]interface{}{"Error": err}))
		}

		if searchQuery.Valid {
			profile.SearchQuery = searchQuery.String
		}

		if lastUsedAt.Valid {
			profile.LastUsedAt = &lastUsedAt.Time
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// UpdateProfile updates an existing filter profile
func (fpm *FilterProfileManager) UpdateProfile(profile FilterProfile) error {
	profile.UpdatedAt = time.Now()
	return fpm.updateFilterProfile(&profile)
}

// DeleteProfile deletes a filter profile by ID
func (fpm *FilterProfileManager) DeleteProfile(id string) error {
	if id == "" {
		return fmt.Errorf(i18n.T("error.filterProfileIdRequired"))
	}

	// Check if profile exists
	existingProfile, err := fpm.GetProfile(id)
	if err != nil {
		return err
	}

	// Don't allow deletion of default profiles
	if existingProfile.IsDefault {
		return fmt.Errorf(i18n.T("cannot_delete_default_profile"))
	}

	query := "DELETE FROM filter_profiles WHERE id = ?"
	result, err := fpm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileDeleteFailed", map[string]interface{}{"Error": err}))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileDeleteFailed", map[string]interface{}{"Error": err}))
	}

	if rowsAffected == 0 {
		return fmt.Errorf(i18n.T("filter_profile_not_found", map[string]interface{}{"id": id}))
	}

	log.Printf("%s", i18n.T("success.filterProfileDeleted", map[string]interface{}{
		"Name": existingProfile.Name,
	}))

	return nil
}

// createFilterProfile creates a new filter profile
func (fpm *FilterProfileManager) createFilterProfile(profile *FilterProfile) error {
	// Validate input
	if profile.Name == "" {
		return fmt.Errorf(i18n.T("errors.filter_profile_name_required", nil))
	}

	// Check if name already exists
	exists, err := fpm.filterProfileNameExists(profile.Name, "")
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileGetFailed", map[string]interface{}{"Error": err}))
	}
	if exists {
		return fmt.Errorf(i18n.T("errors.filter_profile_name_exists", map[string]interface{}{
			"name": profile.Name,
		}))
	}

	// Marshal filters to JSON
	filtersJSON, err := json.Marshal(profile.Filters)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileMarshalFailed", map[string]interface{}{"Error": err}))
	}

	// Insert into database
	query := `
		INSERT INTO filter_profiles (name, description, filters, search_query, is_default, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := fpm.db.Exec(query,
		profile.Name,
		profile.Description,
		string(filtersJSON),
		profile.SearchQuery,
		profile.IsDefault,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileCreateFailed", map[string]interface{}{"Error": err}))
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileCreateFailed", map[string]interface{}{"Error": err}))
	}

	profile.ID = fmt.Sprintf("%d", id)
	profile.CreatedAt = time.Now()

	log.Printf("%s", i18n.T("success.filterProfileCreated", map[string]interface{}{
		"Name": profile.Name,
		"ID":   profile.ID,
	}))
	return nil
}

// updateFilterProfile updates an existing filter profile
func (fpm *FilterProfileManager) updateFilterProfile(profile *FilterProfile) error {
	// Validate input
	if profile.Name == "" {
		return fmt.Errorf(i18n.T("errors.filter_profile_name_required", nil))
	}

	// Check if name already exists (excluding current profile)
	exists, err := fpm.filterProfileNameExists(profile.Name, profile.ID)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileGetFailed", map[string]interface{}{"Error": err}))
	}
	if exists {
		return fmt.Errorf(i18n.T("errors.filter_profile_name_exists", map[string]interface{}{
			"name": profile.Name,
		}))
	}

	// Marshal filters to JSON
	filtersJSON, err := json.Marshal(profile.Filters)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileMarshalFailed", map[string]interface{}{"Error": err}))
	}

	// Update in database
	query := `
		UPDATE filter_profiles
		SET name = ?, description = ?, filters = ?, search_query = ?, is_default = ?
		WHERE id = ?
	`

	result, err := fpm.db.Exec(query,
		profile.Name,
		profile.Description,
		string(filtersJSON),
		profile.SearchQuery,
		profile.IsDefault,
		profile.ID,
	)

	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileUpdateFailed", map[string]interface{}{"Error": err}))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileUpdateFailed", map[string]interface{}{"Error": err}))
	}

	if rowsAffected == 0 {
		return fmt.Errorf(i18n.T("errors.filter_profile_not_found", map[string]interface{}{
			"id": profile.ID,
		}))
	}

	log.Printf("%s", i18n.T("success.filterProfileUpdated", map[string]interface{}{
		"Name": profile.Name,
	}))
	return nil
}

// GetFilterProfile retrieves a filter profile by ID
func (fpm *FilterProfileManager) GetFilterProfile(id int) (*FilterProfile, error) {
	query := `
		SELECT id, name, description, filters, search_query, is_default,
		       created_at, last_used_at, use_count
		FROM filter_profiles
		WHERE id = ?
	`

	var profile FilterProfile
	var filtersJSON string
	var lastUsedAt sql.NullTime

	err := fpm.db.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Description,
		&filtersJSON,
		&profile.SearchQuery,
		&profile.IsDefault,
		&profile.CreatedAt,
		&lastUsedAt,
		&profile.UseCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(i18n.T("errors.filter_profile_not_found", map[string]interface{}{
				"id": id,
			}))
		}
		return nil, fmt.Errorf(i18n.T("error.filterProfileGetFailed", map[string]interface{}{"Error": err}))
	}

	// Parse filters JSON
	err = json.Unmarshal([]byte(filtersJSON), &profile.Filters)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileParseFailed", map[string]interface{}{"Error": err}))
	}

	// Set last used time if available
	if lastUsedAt.Valid {
		profile.LastUsedAt = &lastUsedAt.Time
	}

	return &profile, nil
}

// ListFilterProfiles retrieves all filter profiles, optionally filtering by defaults
func (fpm *FilterProfileManager) ListFilterProfiles(defaultsOnly bool) ([]*FilterProfile, error) {
	query := `
		SELECT id, name, description, filters, search_query, is_default,
		       created_at, last_used_at, use_count
		FROM filter_profiles
	`

	args := []interface{}{}
	if defaultsOnly {
		query += " WHERE is_default = TRUE"
	}

	query += " ORDER BY is_default DESC, use_count DESC, name ASC"

	rows, err := fpm.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileListFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var profiles []*FilterProfile

	for rows.Next() {
		var profile FilterProfile
		var filtersJSON string
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Description,
			&filtersJSON,
			&profile.SearchQuery,
			&profile.IsDefault,
			&profile.CreatedAt,
			&lastUsedAt,
			&profile.UseCount,
		)

		if err != nil {
			log.Printf("%s", i18n.T("error.scanResultFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Parse filters JSON
		err = json.Unmarshal([]byte(filtersJSON), &profile.Filters)
		if err != nil {
			log.Printf("%s", i18n.T("error.filterProfileParseFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Set last used time if available
		if lastUsedAt.Valid {
			profile.LastUsedAt = &lastUsedAt.Time
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

// DeleteFilterProfile deletes a filter profile by ID
func (fpm *FilterProfileManager) DeleteFilterProfile(id int) error {
	// Don't allow deletion of default profiles
	profile, err := fpm.GetFilterProfile(id)
	if err != nil {
		return err
	}

	if profile.IsDefault {
		return fmt.Errorf(i18n.T("errors.cannot_delete_default_profile", nil))
	}

	query := "DELETE FROM filter_profiles WHERE id = ?"
	result, err := fpm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileDeleteFailed", map[string]interface{}{"Error": err}))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileUpdateFailed", map[string]interface{}{"Error": err}))
	}

	if rowsAffected == 0 {
		return fmt.Errorf(i18n.T("errors.filter_profile_not_found", map[string]interface{}{
			"id": id,
		}))
	}

	log.Printf("%s", i18n.T("success.filterProfileDeleted", map[string]interface{}{
		"Name": profile.Name,
	}))
	return nil
}

// MarkProfileUsed updates the usage statistics for a profile
func (fpm *FilterProfileManager) MarkProfileUsed(id int) error {
	query := `
		UPDATE filter_profiles
		SET last_used_at = ?, use_count = use_count + 1
		WHERE id = ?
	`

	result, err := fpm.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileUseFailed", map[string]interface{}{"Error": err}))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileUpdateFailed", map[string]interface{}{"Error": err}))
	}

	if rowsAffected == 0 {
		return fmt.Errorf(i18n.T("errors.filter_profile_not_found", map[string]interface{}{
			"id": id,
		}))
	}

	return nil
}

// GetDefaultProfiles returns all default filter profiles
func (fpm *FilterProfileManager) GetDefaultProfiles() ([]*FilterProfile, error) {
	return fpm.ListFilterProfiles(true)
}

// GetMostUsedProfiles returns the most frequently used filter profiles
func (fpm *FilterProfileManager) GetMostUsedProfiles(limit int) ([]*FilterProfile, error) {
	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT id, name, description, filters, search_query, is_default,
		       created_at, last_used_at, use_count
		FROM filter_profiles
		WHERE use_count > 0
		ORDER BY use_count DESC, last_used_at DESC
		LIMIT ?
	`

	rows, err := fpm.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileGetFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var profiles []*FilterProfile

	for rows.Next() {
		var profile FilterProfile
		var filtersJSON string
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Description,
			&filtersJSON,
			&profile.SearchQuery,
			&profile.IsDefault,
			&profile.CreatedAt,
			&lastUsedAt,
			&profile.UseCount,
		)

		if err != nil {
			log.Printf("%s", i18n.T("error.scanResultFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Parse filters JSON
		err = json.Unmarshal([]byte(filtersJSON), &profile.Filters)
		if err != nil {
			log.Printf("%s", i18n.T("error.filterProfileParseFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Set last used time if available
		if lastUsedAt.Valid {
			profile.LastUsedAt = &lastUsedAt.Time
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

// SearchFilterProfiles searches filter profiles by name or description
func (fpm *FilterProfileManager) SearchFilterProfiles(query string) ([]*FilterProfile, error) {
	searchQuery := `
		SELECT id, name, description, filters, search_query, is_default,
		       created_at, last_used_at, use_count
		FROM filter_profiles
		WHERE name LIKE ? OR description LIKE ?
		ORDER BY use_count DESC, name ASC
	`

	searchTerm := "%" + query + "%"
	rows, err := fpm.db.Query(searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileSearchFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var profiles []*FilterProfile

	for rows.Next() {
		var profile FilterProfile
		var filtersJSON string
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Description,
			&filtersJSON,
			&profile.SearchQuery,
			&profile.IsDefault,
			&profile.CreatedAt,
			&lastUsedAt,
			&profile.UseCount,
		)

		if err != nil {
			log.Printf("%s", i18n.T("error.scanResultFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Parse filters JSON
		err = json.Unmarshal([]byte(filtersJSON), &profile.Filters)
		if err != nil {
			log.Printf("%s", i18n.T("error.filterProfileParseFailed", map[string]interface{}{"Error": err}))
			continue
		}

		// Set last used time if available
		if lastUsedAt.Valid {
			profile.LastUsedAt = &lastUsedAt.Time
		}

		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

// filterProfileNameExists checks if a profile name already exists
func (fpm *FilterProfileManager) filterProfileNameExists(name string, excludeID string) (bool, error) {
	query := "SELECT COUNT(*) FROM filter_profiles WHERE name = ?"
	args := []interface{}{name}

	if excludeID != "" {
		query += " AND id != ?"
		args = append(args, excludeID)
	}

	var count int
	err := fpm.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetSearchHistory returns recent search history
func (fpm *FilterProfileManager) GetSearchHistory(limit int) ([]*SearchHistoryEntry, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, query, filters, result_count, execution_time_ms, created_at
		FROM search_history
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := fpm.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.filterProfileGetFailed", map[string]interface{}{"Error": err}))
	}
	defer rows.Close()

	var history []*SearchHistoryEntry

	for rows.Next() {
		var entry SearchHistoryEntry
		var filters sql.NullString

		err := rows.Scan(
			&entry.ID,
			&entry.Query,
			&filters,
			&entry.ResultCount,
			&entry.ExecutionTimeMs,
			&entry.CreatedAt,
		)

		if err != nil {
			log.Printf("%s", i18n.T("error.scanResultFailed", map[string]interface{}{"Error": err}))
			continue
		}

		if filters.Valid {
			entry.Filters = filters.String
		}

		history = append(history, &entry)
	}

	return history, nil
}

// CleanOldSearchHistory removes search history older than specified days
func (fpm *FilterProfileManager) CleanOldSearchHistory(daysOld int) error {
	if daysOld <= 0 {
		daysOld = 30 // Default to 30 days
	}

	cutoffDate := time.Now().AddDate(0, 0, -daysOld)

	query := "DELETE FROM search_history WHERE created_at < ?"
	result, err := fpm.db.Exec(query, cutoffDate)
	if err != nil {
		return fmt.Errorf(i18n.T("error.searchHistoryCleanFailed", map[string]interface{}{"Error": err}))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(i18n.T("error.filterProfileUpdateFailed", map[string]interface{}{"Error": err}))
	}

	log.Printf("%s", i18n.T("success.searchHistoryCleaned", map[string]interface{}{
		"Count": rowsAffected,
		"Days":  daysOld,
	}))
	return nil
}
