# API Changes - Version 1.1.0

## Date: 25 June 2025

### New MCP Tools Added

#### 1. gorev_detay
- **Purpose**: Display detailed task information in markdown format
- **Parameters**: 
  - `id` (string, required): Task ID
- **Returns**: Markdown formatted task details including status, priority, dates, project, and description

#### 2. gorev_duzenle
- **Purpose**: Edit task properties (title, description, priority, or project)
- **Parameters**:
  - `id` (string, required): Task ID
  - `baslik` (string, optional): New title
  - `aciklama` (string, optional): New description (supports markdown)
  - `oncelik` (string, optional): New priority (dusuk, orta, yuksek)
  - `proje_id` (string, optional): New project ID
- **Returns**: Success confirmation with task ID
- **Note**: At least one optional field must be provided

#### 3. gorev_sil
- **Purpose**: Delete a task permanently
- **Parameters**:
  - `id` (string, required): Task ID
  - `onay` (boolean, required): Must be true to confirm deletion
- **Returns**: Success confirmation with deleted task name and ID

#### 4. proje_listele
- **Purpose**: List all projects with their task counts
- **Parameters**: None
- **Returns**: Markdown formatted list of projects with metadata and task counts

#### 5. proje_gorevleri
- **Purpose**: List all tasks for a specific project, grouped by status
- **Parameters**:
  - `proje_id` (string, required): Project ID
- **Returns**: Markdown formatted task list grouped by status (Devam Ediyor, Beklemede, Tamamlandı)

### Enhanced Features

1. **Markdown Support**
   - Task descriptions now support full markdown formatting
   - Markdown is preserved in `gorev_detay` output
   - `gorev_duzenle` accepts markdown in description field

2. **Partial Updates**
   - `gorev_duzenle` only updates specified fields
   - Unspecified fields remain unchanged
   - Validation ensures at least one field is provided

3. **Safety Features**
   - `gorev_sil` requires explicit confirmation parameter
   - Task name is shown before deletion for verification
   - Error messages are more descriptive and user-friendly

### Technical Changes

1. **Business Logic Layer**
   - Added: `GorevDetayAl(id string) (*Gorev, error)`
   - Added: `ProjeDetayAl(id string) (*Proje, error)`
   - Added: `GorevDuzenle(id, baslik, aciklama, oncelik, projeID string, flags...bool) error`
   - Added: `GorevSil(id string) error`
   - Added: `ProjeListele() ([]*Proje, error)`
   - Added: `ProjeGorevleri(projeID string) ([]*Gorev, error)`
   - Added: `ProjeGorevSayisi(projeID string) (int, error)`

2. **Data Access Layer**
   - Added: `ProjeGetir(id string) (*Proje, error)`
   - Added: `GorevSil(id string) error`
   - Added: `ProjeGorevleriGetir(projeID string) ([]*Gorev, error)`

3. **MCP Handlers**
   - All handlers updated to match mark3labs/mcp-go v0.6.0 API
   - Removed context parameter from handler functions
   - Updated tool registration to use new API

### Breaking Changes
None - All existing tools maintain backward compatibility

### Migration Notes
- Existing tasks will work with new tools without modification
- Task descriptions can be gradually updated to use markdown
- No database schema changes required