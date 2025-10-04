-- Add search support for advanced search functionality
-- Using regular SQL for maximum compatibility (no FTS5 dependency)

-- Create search index table for better performance
CREATE TABLE IF NOT EXISTS gorevler_search (
    id TEXT PRIMARY KEY,
    baslik TEXT,
    aciklama TEXT,
    etiketler TEXT,
    proje_adi TEXT,
    search_text TEXT, -- Combined searchable text
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster searching
CREATE INDEX IF NOT EXISTS idx_gorevler_search_baslik ON gorevler_search(baslik);
CREATE INDEX IF NOT EXISTS idx_gorevler_search_text ON gorevler_search(search_text);
CREATE INDEX IF NOT EXISTS idx_gorevler_search_combined ON gorevler_search(baslik, aciklama);

-- Create filter profiles table for saved search configurations
CREATE TABLE IF NOT EXISTS filter_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    filters TEXT NOT NULL, -- JSON field containing filter configuration
    search_query TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    use_count INTEGER DEFAULT 0,
    last_used_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create search history table for tracking search patterns
CREATE TABLE IF NOT EXISTS search_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    filters TEXT, -- JSON field for filter state
    result_count INTEGER DEFAULT 0,
    execution_time_ms INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better search performance
CREATE INDEX IF NOT EXISTS idx_filter_profiles_name ON filter_profiles(name);
CREATE INDEX IF NOT EXISTS idx_filter_profiles_default ON filter_profiles(is_default);
CREATE INDEX IF NOT EXISTS idx_filter_profiles_last_used ON filter_profiles(last_used_at);
CREATE INDEX IF NOT EXISTS idx_search_history_query ON search_history(query);
CREATE INDEX IF NOT EXISTS idx_search_history_created ON search_history(created_at);

-- Insert some default filter profiles
INSERT INTO filter_profiles (name, description, filters, is_default) VALUES
('Yüksek Öncelik', 'Yüksek öncelikli görevler', '{"priority":["yuksek"]}', TRUE),
('Devam Ediyor', 'Şu anda üzerinde çalışılan görevler', '{"status":["devam_ediyor"]}', TRUE),
('Bekleyen', 'Beklemede olan görevler', '{"status":["beklemede"]}', TRUE),
('Acil Görevler', 'Yüksek öncelikli bekleyen görevler', '{"status":["beklemede"],"priority":["yuksek"]}', TRUE);