# Migration Guide: v0.16.x → v0.17.0

**Release Date**: October 11, 2025
**Migration Type**: Breaking Changes - Automatic Database Migration

## 🚨 Overview

Version 0.17.0 is a **major refactoring release** that renames all database columns and API fields from Turkish to English. This change improves international accessibility and code maintainability while preserving the project's Turkish domain terminology (`gorevler`, `projeler`).

**Key Points:**
- ✅ **Automatic Migration**: Database schema updates automatically on first run
- ✅ **Data Preserved**: All your tasks, projects, and templates remain intact
- ⚠️ **Irreversible**: Once migrated, you cannot downgrade to v0.16.x
- ⚠️ **Breaking Changes**: Custom scripts and integrations need updates

---

## 📋 Pre-Migration Checklist

### 1. Backup Your Database

**CRITICAL**: Always backup before upgrading.

```bash
# Backup your Gorev database
cp ~/.gorev/gorev.db ~/.gorev/gorev.db.backup-$(date +%Y%m%d)

# Verify backup exists
ls -lh ~/.gorev/gorev.db.backup-*
```

### 2. Check Current Version

```bash
gorev version
# Should show: Gorev 0.16.x
```

### 3. Stop Running Instances

```bash
# Stop Gorev daemon if running
pkill -f "gorev daemon"

# Verify no processes running
ps aux | grep gorev
```

### 4. Document Custom Integrations

Review and document any custom scripts, tools, or integrations that:
- Use Gorev REST API
- Query the SQLite database directly
- Reference Turkish field names in code

---

## 🔄 Migration Steps

### Step 1: Install v0.17.0

#### NPM Users

```bash
npm update -g gorev-mcp-server

# Verify new version
gorev version
# Should show: Gorev 0.17.0
```

#### VS Code Extension Users

1. Open VS Code Extensions panel
2. Search for "Gorev"
3. Update to v0.17.0
4. Reload VS Code

#### Building from Source

```bash
cd gorev-mcpserver
git pull origin main
git checkout v0.17.0
make build

./gorev version
# Should show: Gorev 0.17.0
```

### Step 2: First Run (Automatic Migration)

```bash
# Start Gorev - migration runs automatically
gorev serve --debug

# You'll see migration logs:
# DEBUG: Migration 11 already applied, skipping
# SUCCESS: Database migrated successfully
```

**Expected Output:**
```
2025/10/11 10:00:00 DEBUG: Schema migrations table ready
2025/10/11 10:00:00 DEBUG: Found 11 migration files
2025/10/11 10:00:00 DEBUG: Migration 11 applied: rename_fields_to_english
2025/10/11 10:00:00 SUCCESS: Database migrated successfully
```

### Step 3: Verify Migration

```bash
# Check database schema
sqlite3 ~/.gorev/gorev.db "PRAGMA table_info(gorevler);" | grep -E "(title|description|status)"

# Should see:
# 1|title|TEXT|1||0
# 2|description|TEXT|0||0
# 3|status|TEXT|1||0
```

### Step 4: Test Basic Operations

```bash
# List tasks (should work with new schema)
gorev mcp list-tasks --all-projects

# Create a test task from template
gorev template-task bug \
  --title "Test Migration" \
  --description "Testing v0.17.0 migration" \
  --modul "test" \
  --priority "orta"

# Verify task created
gorev task-detail <task-id>
```

---

## 🔀 Field Name Changes

### Complete Mapping Table

| Table | Old Field Name (TR) | New Field Name (EN) |
|-------|-------------------|-------------------|
| **gorevler** | baslik | title |
| | aciklama | description |
| | durum | status |
| | oncelik | priority |
| | proje_id | project_id |
| | olusturma_tarih | created_at |
| | guncelleme_tarih | updated_at |
| | son_tarih | due_date |
| **projeler** | isim | name |
| | tanim | definition |
| | olusturma_tarih | created_at |
| | guncelleme_tarih | updated_at |
| **etiketler** | isim | name |
| **gorev_etiketleri** | gorev_id | task_id |
| | etiket_id | tag_id |
| **baglantilar** | kaynak_id | source_id |
| | hedef_id | target_id |
| | baglanti_tip | connection_type |
| **gorev_templateleri** | isim | name |
| | tanim | definition |
| | varsayilan_baslik | default_title |
| | aciklama_template | description_template |
| | ornek_degerler | sample_values |
| | kategori | category |
| | aktif | active |
| | alanlar | fields |
| **ai_interactions** | gorev_id | task_id |
| **aktif_proje** | proje_id | project_id |

---

## 🛠️ Updating Custom Integrations

### REST API Changes

**Before (v0.16.x):**
```typescript
// Old field names
const task = {
  baslik: "Task Title",
  aciklama: "Description",
  durum: "beklemede",
  oncelik: "yuksek",
  proje_id: "proj-123"
};

fetch('http://localhost:5082/api/gorevler', {
  method: 'POST',
  body: JSON.stringify(task)
});
```

**After (v0.17.0):**
```typescript
// New field names
const task = {
  title: "Task Title",
  description: "Description",
  status: "beklemede",
  priority: "yuksek",
  project_id: "proj-123"
};

fetch('http://localhost:5082/api/gorevler', {
  method: 'POST',
  body: JSON.stringify(task)
});
```

### Direct SQL Queries

**Before (v0.16.x):**
```sql
SELECT baslik, aciklama, durum, oncelik
FROM gorevler
WHERE proje_id = 'proj-123'
  AND durum = 'beklemede'
ORDER BY olusturma_tarih DESC;
```

**After (v0.17.0):**
```sql
SELECT title, description, status, priority
FROM gorevler
WHERE project_id = 'proj-123'
  AND status = 'beklemede'
ORDER BY created_at DESC;
```

### Template Values

**Before (v0.16.x):**
```json
{
  "template_id": "bug",
  "values": {
    "baslik": "Bug Title",
    "aciklama": "Bug Description",
    "oncelik": "yuksek"
  }
}
```

**After (v0.17.0):**
```json
{
  "template_id": "bug",
  "values": {
    "title": "Bug Title",
    "description": "Bug Description",
    "priority": "yuksek"
  }
}
```

---

## 🐛 Troubleshooting

### Issue: Migration Fails

**Symptoms:**
```
ERROR: Migration 11 failed: SQL error
```

**Solutions:**
1. Restore from backup:
   ```bash
   cp ~/.gorev/gorev.db.backup-YYYYMMDD ~/.gorev/gorev.db
   ```

2. Check database integrity:
   ```bash
   sqlite3 ~/.gorev/gorev.db "PRAGMA integrity_check;"
   ```

3. Review migration logs:
   ```bash
   gorev serve --debug 2>&1 | grep -A 10 "Migration 11"
   ```

### Issue: "no such column: baslik" Error

**Cause:** You're using old field names after migration.

**Solution:** Update your code/scripts to use new English field names (see mapping table above).

### Issue: Template Placeholders Not Working

**Cause:** Templates created before v0.17.0 may use Turkish placeholders (`{{baslik}}`).

**Solution:** Built-in templates are auto-updated. For custom templates:
```bash
# List your templates
gorev template-list

# Edit custom template to use {{title}} instead of {{baslik}}
sqlite3 ~/.gorev/gorev.db "UPDATE gorev_templateleri
SET varsayilan_baslik = REPLACE(varsayilan_baslik, '{{baslik}}', '{{title}}'),
    aciklama_template = REPLACE(aciklama_template, '{{aciklama}}', '{{description}}')
WHERE id = 'your-template-id';"
```

### Issue: VS Code Extension Shows Old Field Names

**Cause:** Extension not updated to v0.17.0.

**Solution:**
1. Update VS Code extension to v0.17.0
2. Reload VS Code window (Cmd/Ctrl + Shift + P → "Developer: Reload Window")
3. Restart Gorev daemon

---

## 🔙 Rollback (Emergency Only)

**⚠️ WARNING:** Rollback is only possible if you have a v0.16.x database backup BEFORE migration.

```bash
# Stop Gorev
pkill -f "gorev"

# Restore old database
cp ~/.gorev/gorev.db.backup-YYYYMMDD ~/.gorev/gorev.db

# Downgrade to v0.16.3
npm install -g gorev-mcp-server@0.16.3

# Start Gorev
gorev serve
```

**Note:** Any tasks created after migration to v0.17.0 will be LOST during rollback.

---

## ✅ Post-Migration Validation

### 1. Database Schema Check

```bash
sqlite3 ~/.gorev/gorev.db << EOF
.mode column
.headers on

-- Check gorevler table has new columns
PRAGMA table_info(gorevler);

-- Verify data integrity
SELECT COUNT(*) as total_tasks FROM gorevler;
SELECT COUNT(DISTINCT project_id) as total_projects FROM gorevler;

-- Test FTS search
SELECT task_id, title FROM gorevler_fts WHERE gorevler_fts MATCH 'test' LIMIT 5;
EOF
```

### 2. API Endpoint Test

```bash
# Health check
curl http://localhost:5082/health

# List tasks with new field names
curl http://localhost:5082/api/gorevler | jq '.[] | {title, status, priority}'

# Create test task
curl -X POST http://localhost:5082/api/gorevler \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Migration Test Task",
    "description": "Testing v0.17.0 API",
    "status": "beklemede",
    "priority": "orta"
  }' | jq
```

### 3. MCP Tools Test

```bash
# Test MCP tools still work
gorev mcp list-tasks --project "your-project" --status "beklemede"

# Test template creation
gorev mcp create-from-template bug \
  --values '{"title":"Test Bug","description":"Testing","modul":"test","priority":"orta"}'

# Test advanced search
gorev mcp search --query "status:beklemede priority:yuksek"
```

---

## 📚 Additional Resources

- **Full CHANGELOG**: [CHANGELOG.md](../CHANGELOG.md)
- **API Documentation**: [docs/API_REFERENCE.md](API_REFERENCE.md)
- **MCP Tools Reference**: [docs/MCP_TOOLS_REFERENCE.md](MCP_TOOLS_REFERENCE.md)
- **Database Schema**: [docs/architecture/database-schema.md](architecture/database-schema.md)

---

## 🆘 Getting Help

If you encounter issues during migration:

1. **Check Logs**: Run with `--debug` flag to see detailed migration logs
2. **GitHub Issues**: https://github.com/msenol/gorev/issues
3. **Backup First**: Always ensure you have a working backup before attempting fixes

**Common Questions:**
- **Q: Will my data be lost?**
  A: No, migration preserves all data. Only column names change.

- **Q: Can I use v0.16.x and v0.17.0 simultaneously?**
  A: No, once migrated to v0.17.0, the database is incompatible with v0.16.x.

- **Q: Do I need to update VS Code extension?**
  A: Yes, VS Code extension must be v0.17.0 to work with v0.17.0 server.

---

**Last Updated**: October 11, 2025
**Version**: v0.17.0
