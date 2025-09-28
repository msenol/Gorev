# VS Code Extension Testing Guide

Bu rehber, Gorev VS Code extension'ının bağımsız olarak test edilmesi için kapsamlı bir kılavuz sağlar.

## 🚨 Task Display Issue Solutions

### Major Issues Fixed

1. **Configuration Mismatch**: pageSize fallback artık package.json ile uyumlu (100)
2. **Hierarchy Filtering**: Orphaned subtask'ları göstermek için geliştirildi
3. **Enhanced Logging**: Task visibility sorunları için detaylı log sistemi

## 🧪 Test Scenarios

### 1. Pagination Testing

```bash
# VS Code Command Palette'te:
> Gorev: Seed Test Data
```

Bu komut şunları oluşturur:

- **150 pagination test task** (pageSize=100'ü test etmek için)
- **Hierarchy test tasks** (parent-child ilişkileri)
- **Normal test tasks** (çeşitli durumlar)

### 2. VS Code Extension Development Testing

#### Method 1: Extension Development Host

1. VS Code'da projeyi aç: `/home/msenol/Projects/Gorev/gorev-vscode`
2. `F5` tuşuna bas (Run and Debug)
3. Yeni VS Code penceresi açılır (Extension Development Host)
4. Yeni pencerede extension test et

#### Method 2: Terminal Testing

```bash
cd /home/msenol/Projects/Gorev/gorev-vscode
npm run compile
npm test
```

#### Method 3: Manual VSIX Testing

```bash
cd /home/msenol/Projects/Gorev/gorev-vscode
npm run package
code --install-extension gorev-vscode-*.vsix
```

### 3. MCP Server Connection Testing

#### Debug Commands (Command Palette)

- `> Gorev Debug: Test MCP Connection`
- `> Gorev Debug: Seed Test Data`
- `> Gorev Debug: Clear Test Data`
- `> Gorev Debug: Show Debug Logs`
- `> Gorev Debug: Toggle Debug Mode`

#### Configuration Testing

VS Code Settings (`Ctrl+,` → Extensions → Gorev):

```json
{
    "gorev.pagination.pageSize": 100,
    "gorev.debug.useWrapper": true,
    "gorev.debug.logPath": "/tmp/gorev-debug",
    "gorev.databaseMode": "auto"
}
```

## 📊 Test Data Overview

### Standard Test Data (gorev.debug.seedTestData)

- **3 Projects**: Web App, Mobile App, Security Updates
- **~200+ Tasks**: Various priorities, statuses, templates
- **Hierarchy**: Parent-child relationships
- **Dependencies**: Task dependencies and blocking relationships
- **Templates**: All major template types used

### Pagination Test Data

- **150 Tasks**: Specifically for pagination testing
- **Batch Tagging**: `pagination-test`, `batch-1`, `batch-2`, `batch-3`
- **Multiple Projects**: Distributed across all test projects

### Hierarchy Test Data

- **3 Parent Tasks**: Complex features requiring subtasks
- **15 Subtasks**: 5 subtasks per parent (Design, Backend, Frontend, Testing, Docs)
- **Proper Relationships**: Parent-child relationships properly set

## 🔍 Debugging Task Visibility Issues

### 1. Enable Debug Logging

```bash
# Command Palette
> Gorev Debug: Toggle Debug Mode
```

### 2. Check Logs

```bash
# Command Palette
> Gorev Debug: Show Debug Logs

# Or check console
> Developer: Toggle Developer Tools
```

### 3. Key Log Messages to Look For

```
[EnhancedGorevTreeProvider] Task loading summary:
  - Total tasks fetched: X
  - Expected total: Y
  - Page size used: 100
  - Show all projects: true/false

[EnhancedGorevTreeProvider] Group groupName: X total, Y root, Z orphaned

[EnhancedGorevTreeProvider] TASK COUNT MISMATCH: Expected X, got Y
```

### 4. Configuration Verification

```bash
# Check current config
> Preferences: Open Settings (JSON)
```

## 🚀 Manual Testing Procedures

### Test Case 1: Basic Task Display

1. Start MCP server: `./gorev serve --debug`
2. Open VS Code Extension Development Host (F5)
3. Verify connection in status bar
4. Check if tasks are displayed in tree view
5. Count visible tasks vs. expected

### Test Case 2: Pagination

1. Clear existing data: `> Gorev Debug: Clear Test Data`
2. Seed test data: `> Gorev Debug: Seed Test Data`
3. Verify 200+ tasks are visible
4. Check for pagination test tasks in tree view
5. Verify no tasks are missing

### Test Case 3: Hierarchy Display

1. Look for "Feature:" parent tasks
2. Expand parent tasks
3. Verify subtasks are shown
4. Check for orphaned subtasks (should be visible)

### Test Case 4: Filtering and Grouping

1. Try different grouping options (status, priority, project)
2. Apply filters (high priority, overdue, etc.)
3. Verify tasks don't disappear unexpectedly
4. Clear filters and verify all tasks return

### Test Case 5: Performance with Large Data

1. Create 500+ tasks using multiple seed operations
2. Monitor loading times
3. Check for safety limit warnings in logs
4. Verify UI remains responsive

## 🛠️ Troubleshooting Common Issues

### Issue: No tasks displayed

**Check:**

1. MCP server running? `./gorev serve --debug`
2. Connection status in VS Code status bar
3. Database has tasks? `./gorev gorev_listele`
4. pageSize configuration (should be 100)

### Issue: Some tasks missing

**Check:**

1. Debug logs for "TASK COUNT MISMATCH"
2. Safety limits (1000 task cap)
3. Hierarchy filtering (orphaned subtasks)
4. Project filters (showAllProjects setting)

### Issue: Extension not loading

**Check:**

1. Compilation errors: `npm run compile`
2. Package.json validity
3. Extension activation events
4. VS Code version compatibility

## 📝 Creating Custom Test Scenarios

### Add More Test Tasks

```typescript
// In testDataSeeder.ts, modify createPaginationTestTasks
const targetTaskCount = 500; // Increase for stress testing
```

### Test Specific Edge Cases

```typescript
// Create tasks with specific properties
const edgeCaseTask = {
    templateId: this.TEMPLATE_IDS.BUG_RAPORU,
    projectId: projectIds[0],
    degerler: {
        baslik: 'Edge Case Task',
        parent_id: 'non-existent-parent-id', // Test orphaned subtask handling
        etiketler: ['edge-case']
    }
};
```

## 📊 Performance Benchmarks

### Expected Performance

- **Task Loading**: <2 seconds for 200 tasks
- **UI Response**: <500ms for filter changes
- **Pagination**: <1 second per page load
- **Memory Usage**: <50MB for extension

### Monitoring Tools

```bash
# VS Code Performance
> Developer: Reload Window With Extensions Disabled
> Developer: Show Running Extensions

# System monitoring
htop
```

## 🎯 Success Criteria

### ✅ All Tests Should Pass

1. **All tasks visible**: No missing tasks due to pagination/hierarchy issues
2. **Proper grouping**: Tasks correctly grouped and filterable
3. **Performance**: Responsive UI with large datasets
4. **Error handling**: Graceful handling of server disconnection
5. **Configuration**: Settings work as expected

### ✅ Debug Information Available

1. **Comprehensive logs**: Detailed logging for troubleshooting
2. **Count verification**: Task count matches expected
3. **Error reporting**: Clear error messages for issues
4. **Status visibility**: Connection and loading status clear

---

**Bu kılavuz ile VS Code extension'ının task display sorunlarını tespit edip çözebilirsiniz. Test data generation sistemi ile çeşitli senaryoları test edebilir, debug logging ile sorunları izleyebilirsiniz.**
