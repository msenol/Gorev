import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { TemplateTreeItem } from '../providers/templateTreeProvider';

export function registerTemplateCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  // Use Template command would be registered here
  // For now, we'll leave it empty as it's a more complex interaction
}