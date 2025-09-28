# VS Code Extension Data Export/Import Guide

Complete guide for using the data export and import features in the Gorev VS Code Extension.

## Overview

The Gorev VS Code Extension provides a comprehensive visual interface for exporting and importing your task management data. This feature builds on top of the MCP server's `gorev_export` and `gorev_import` tools to provide an intuitive, multi-step user experience.

## Export Features

### 1. Export Data (Full Configuration)

**Command**: `Gorev: Export Data`  
**Access**: Command Palette → "Gorev: Export Data"

Features a 4-step wizard for comprehensive export configuration:

#### Step 1: Format Selection

- **JSON**: Structured data format, ideal for backup and migration
- **CSV**: Tabular format, perfect for analysis in spreadsheet applications

#### Step 2: Filtering Options

- **Include Completed Tasks**: Toggle to include/exclude completed tasks
- **Include Dependencies**: Export task dependency relationships
- **Include Templates**: Export task templates (optional)
- **Include AI Context**: Export AI interaction history (optional)
- **Project Filter**: Select specific projects to export
- **Tag Filter**: Filter tasks by specific tags
- **Date Range**: Export tasks within a specific date range

#### Step 3: Output Location

- **File Browser**: Choose where to save your export file
- **Format Detection**: Automatic file extension based on chosen format
- **Path Validation**: Ensures output directory is accessible

#### Step 4: Review and Export

- **Configuration Summary**: Review all export settings
- **Size Estimation**: Preview approximate export file size
- **Progress Tracking**: Real-time export progress with VS Code progress API

### 2. Export Current View (Quick Context Export)

**Command**: `Gorev: Export Current View`  
**Access**: TreeView title bar → Export icon, Context menu

Exports the currently displayed tasks based on your current filters and project selection:

- **Current Project Context**: Automatically uses active project if set
- **Applied Filters**: Respects current search and filter settings
- **Quick Configuration**: Minimal setup - choose format and location
- **Default Settings**: Uses sensible defaults for quick export

### 3. Quick Export (One-Click Export)

**Command**: `Gorev: Quick Export`  
**Access**: Command Palette → "Gorev: Quick Export"

Fastest export option with automatic settings:

- **Automatic Naming**: Uses timestamp format (`gorev-quick-export-YYYY-MM-DD-HHMMSS.json`)
- **Downloads Folder**: Automatically saves to user's Downloads directory
- **JSON Format**: Uses JSON format by default
- **Complete Data**: Includes tasks, projects, and dependencies (excludes templates and AI context)
- **Success Actions**: Option to open file or reveal in folder after export

## Import Features

### Import Data (Full Import Wizard)

**Command**: `Gorev: Import Data`  
**Access**: Command Palette → "Gorev: Import Data"

Features a 4-step wizard for safe data import:

#### Step 1: File Selection

- **File Browser**: Select JSON or CSV file to import
- **Format Detection**: Automatic format detection based on file extension
- **File Validation**: Basic file accessibility and format checks

#### Step 2: Project Mapping

- **Existing Projects**: View current projects in your system
- **Project Remapping**: Map imported project names to existing projects
- **New Project Creation**: Option to create new projects for unmatched imports
- **Conflict Preview**: See potential project conflicts before import

#### Step 3: Dry Run Preview

- **Safe Analysis**: Analyze import data without making changes
- **Conflict Detection**: Identify tasks that would conflict with existing data
- **Import Statistics**: Preview how many tasks, projects will be imported
- **Conflict Resolution Options**:
  - **Skip**: Skip conflicting items
  - **Overwrite**: Replace existing items with imported data
  - **Merge**: Attempt to merge conflicting data

#### Step 4: Import Execution

- **Final Confirmation**: Review all import settings
- **Progress Tracking**: Real-time import progress with detailed status
- **Success Summary**: Complete import statistics with success/failure counts
- **Error Reporting**: Detailed error messages for any failed imports

## Advanced Features

### Progress Tracking

All export/import operations provide real-time progress feedback:

- **VS Code Progress API**: Native progress bars in notification area
- **Detailed Status**: Step-by-step progress messages
- **Cancellation Support**: Ability to cancel long-running operations (where supported)
- **Error Recovery**: Graceful handling of interrupted operations

### Conflict Resolution

Advanced conflict handling for import operations:

- **Smart Detection**: Identifies conflicts by task title, project, and metadata
- **Resolution Strategies**:
  - **Skip Conflicts**: Keep existing data, skip conflicting imports
  - **Overwrite Existing**: Replace existing data with imported data
  - **Interactive Resolution**: Prompt for each conflict (future enhancement)
- **Dry Run Mode**: Test import operations without making changes
- **Rollback Support**: Future enhancement for undoing imports

### Localization

Complete bilingual support:

- **Turkish Interface**: Native Turkish language support for all UI elements
- **English Interface**: Full English localization
- **Context-Aware Messages**: Dynamic message formatting based on operation results
- **Consistent Terminology**: Unified terms across all export/import features

## Error Handling

Comprehensive error handling throughout the export/import process:

### Common Error Scenarios

1. **File Access Errors**:
   - Insufficient permissions for output directory
   - Invalid file paths or non-existent directories
   - File already exists (with overwrite protection)

2. **Format Errors**:
   - Invalid JSON structure in import files
   - Malformed CSV data
   - Unsupported file formats

3. **Data Validation Errors**:
   - Missing required fields
   - Invalid task relationships
   - Circular dependency detection

4. **MCP Connection Errors**:
   - Server not connected
   - Timeout during large operations
   - Network connectivity issues

### Error Recovery

- **Graceful Degradation**: Partial success handling for bulk operations
- **User-Friendly Messages**: Clear, actionable error descriptions
- **Recovery Suggestions**: Specific steps to resolve common issues
- **Support Information**: Links to troubleshooting resources

## Best Practices

### Export Best Practices

1. **Regular Backups**: Use Quick Export for daily backups
2. **Full Exports**: Use Export Data for complete system backups
3. **Selective Exports**: Use Current View for sharing specific project data
4. **Format Selection**: Use JSON for backups, CSV for analysis
5. **Date Range Filtering**: Export specific time periods for archival

### Import Best Practices

1. **Always Use Dry Run**: Preview imports before executing
2. **Backup Before Import**: Export current data before importing
3. **Project Mapping**: Carefully map projects to avoid duplicates
4. **Conflict Resolution**: Choose appropriate resolution strategy
5. **Verify Results**: Check imported data after successful import

### Performance Considerations

1. **Large Datasets**: Use filtering to reduce export size
2. **Network Operations**: Ensure stable connection for large imports
3. **Disk Space**: Verify sufficient disk space for exports
4. **Memory Usage**: Consider system memory for large operations

## Troubleshooting

### Common Issues

1. **Export Fails**:
   - Check output directory permissions
   - Verify sufficient disk space
   - Try smaller data sets with filtering

2. **Import Fails**:
   - Validate file format and structure
   - Check for corrupted files
   - Verify MCP server connection

3. **UI Not Responsive**:
   - Close and reopen export/import dialogs
   - Restart VS Code extension
   - Check MCP server status

### Getting Help

- **Extension Logs**: Check VS Code Output panel → "Gorev Extension"
- **MCP Server Logs**: Use debug mode for detailed logging
- **GitHub Issues**: Report bugs at project repository
- **Documentation**: Reference MCP tools guide for advanced usage

## Integration with MCP Tools

The VS Code extension seamlessly integrates with the underlying MCP server tools:

- **gorev_export**: Powers all export functionality
- **gorev_import**: Handles all import operations
- **Real-time Validation**: Uses server-side validation for data integrity
- **Consistent Behavior**: Same results whether using CLI or VS Code interface

This ensures that exports created through VS Code are fully compatible with CLI operations and vice versa.
