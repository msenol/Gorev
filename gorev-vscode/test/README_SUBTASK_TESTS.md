# Subtask UI Test Documentation

This document describes the test coverage for the subtask UI functionality in the Gorev VS Code extension.

## Test Files

### 1. Unit Tests

#### `test/unit/subtaskUI.test.js`

Tests the basic subtask UI components:

- **Gorev Model Tests**: Verifies parent_id, alt_gorevler, and seviye fields
- **GorevHiyerarsi Structure**: Tests hierarchy information structure
- **MarkdownParser Hierarchy Tests**: Ensures correct parsing of hierarchical task structures
- **Command Registration**: Verifies CREATE_SUBTASK, CHANGE_PARENT, REMOVE_PARENT commands
- **Tree Item Context Values**: Tests task:parent, task:child context values

#### `test/unit/dragDropController.test.js`

Tests drag & drop functionality for parent changes:

- **Single Task Drag**: Tests dragging individual tasks
- **Multiple Tasks Drag**: Tests dragging multiple selected tasks
- **Drop on Task**: Tests parent change and dependency creation options
- **Drop on Empty Area**: Tests removing parent (making root task)
- **Error Handling**: Tests circular dependency and same-project validation
- **Configuration**: Tests allowParentChange configuration option

#### `test/unit/enhancedTreeViewHierarchy.test.js`

Tests hierarchical display in TreeView:

- **Root Task Display**: Only root tasks shown at top level
- **Task Expansion**: Parent tasks can be expanded to show children
- **Subtask Count**: Shows completion count (e.g., "📁 2/5")
- **Filtering with Hierarchy**: Maintains hierarchy during search/filter
- **Nested Subtasks**: Handles multiple levels of nesting
- **Context Menus**: Correct options for parent/child tasks

#### `test/unit/taskDetailPanelHierarchy.test.js`

Tests hierarchy display in task detail panel:

- **Hierarchy API Call**: Calls gorev_hiyerarsi_goster
- **Hierarchy Section**: Renders hierarchy statistics
- **Progress Bar**: Shows subtask completion progress
- **Parent Info**: Shows parent task information
- **Action Buttons**: Create subtask, change parent buttons
- **Message Handling**: Handles hierarchy-related messages

### 2. Integration Tests

#### `test/integration/subtaskCommands.test.js`

Tests command execution flow:

- **Create Subtask Command**: Full flow with user inputs
- **Change Parent Command**: Task selection and parent change
- **Remove Parent Command**: Making task root
- **Error Scenarios**: Circular dependency, different project errors
- **Tree Refresh**: Ensures UI updates after operations

## Running Tests

### Run All Tests

```bash
npm test
```

### Run Only Subtask Tests

```bash
# Using grep pattern
npm test -- --grep "Subtask|Hierarchy|Drag.*Drop.*parent"

# Or run specific test files
npm test test/unit/subtaskUI.test.js
npm test test/unit/dragDropController.test.js
npm test test/unit/enhancedTreeViewHierarchy.test.js
npm test test/unit/taskDetailPanelHierarchy.test.js
npm test test/integration/subtaskCommands.test.js
```

## Test Coverage Areas

### ✅ Covered

1. **Model Changes**: parent_id, alt_gorevler, seviye fields
2. **Hierarchical Display**: Tree structure with parent-child relationships
3. **Drag & Drop**: Parent changing via drag operations
4. **Commands**: Create subtask, change parent, remove parent
5. **UI Elements**: Buttons, context menus, progress bars
6. **Error Handling**: Circular dependencies, project constraints
7. **Markdown Parsing**: Hierarchical task structure parsing

### 🔄 Partially Covered

1. **Deep Nesting**: Tests up to 3 levels, but unlimited depth supported
2. **Performance**: Large hierarchy performance not tested
3. **Concurrent Updates**: Race conditions not tested

### ❌ Not Covered

1. **E2E Tests**: Full user workflow from UI to backend
2. **Visual Regression**: Screenshot comparisons
3. **Accessibility**: Keyboard navigation, screen readers
4. **Localization**: Turkish/English language switching

## Test Data Examples

### Simple Parent-Child

```javascript
const parentTask = {
    id: 'parent1',
    baslik: 'Parent Task',
    alt_gorevler: [
        { id: 'child1', baslik: 'Child 1', parent_id: 'parent1' },
        { id: 'child2', baslik: 'Child 2', parent_id: 'parent1' }
    ]
};
```

### Multi-Level Hierarchy

```javascript
const deepHierarchy = {
    id: 'root',
    baslik: 'Root Task',
    alt_gorevler: [{
        id: 'level1',
        baslik: 'Level 1',
        parent_id: 'root',
        alt_gorevler: [{
            id: 'level2',
            baslik: 'Level 2',
            parent_id: 'level1',
            alt_gorevler: []
        }]
    }]
};
```

### Hierarchy Statistics

```javascript
const hierarchyInfo = {
    gorev: task,
    ust_gorevler: [],
    toplam_alt_gorev: 10,
    tamamlanan_alt: 7,
    devam_eden_alt: 1,
    beklemede_alt: 2,
    ilerleme_yuzdesi: 70
};
```

## Mock Objects

### Mock MCP Client

```javascript
const mockMcpClient = {
    callTool: sinon.stub(),
    isConnected: sinon.stub().returns(true)
};
```

### Mock DataTransfer (Drag & Drop)

```javascript
const mockDataTransfer = {
    get: sinon.stub(),
    set: sinon.stub()
};
```

### Mock VS Code Window

```javascript
sinon.stub(vscode.window, 'showQuickPick');
sinon.stub(vscode.window, 'showInputBox');
sinon.stub(vscode.window, 'showErrorMessage');
sinon.stub(vscode.window, 'showInformationMessage');
```

## Common Test Patterns

### Testing Async Commands

```javascript
test('should handle async operation', async () => {
    mockMcpClient.callTool.resolves({ content: [{ text: 'response' }] });
    await vscode.commands.executeCommand('gorev.createSubtask', { task });
    assert.ok(mockMcpClient.callTool.called);
});
```

### Testing Error Scenarios

```javascript
test('should handle error', async () => {
    mockMcpClient.callTool.rejects(new Error('dairesel bağımlılık'));
    await command();
    assert.ok(vscode.window.showErrorMessage.called);
});
```

### Testing UI Updates

```javascript
test('should refresh tree', async () => {
    await command();
    assert.ok(mockTreeProvider.refresh.called);
});
```

## Future Improvements

1. **Performance Tests**: Test with 1000+ tasks in hierarchy
2. **E2E Tests**: Using VS Code's test API for real UI interaction
3. **Visual Tests**: Screenshot comparison for UI elements
4. **Stress Tests**: Concurrent operations, rapid parent changes
5. **Accessibility Tests**: Keyboard navigation through hierarchy
