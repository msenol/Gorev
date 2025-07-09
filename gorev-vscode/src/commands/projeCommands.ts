import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { ProjeTreeItem } from '../providers/projeTreeProvider';
import { MarkdownParser } from '../utils/markdownParser';
import { Proje } from '../models/proje';

async function getProjectList(mcpClient: MCPClient): Promise<Proje[]> {
  try {
    const result = await mcpClient.callTool('proje_listele');
    return MarkdownParser.parseProjeListesi(result.content[0].text);
  } catch (error) {
    console.error('Failed to get project list:', error);
    return [];
  }
}

export function registerProjeCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  // Create Project
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_PROJECT, async () => {
      try {
        const isim = await vscode.window.showInputBox({
          prompt: 'Project name',
          placeHolder: 'Enter project name',
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return 'Project name is required';
            }
            return null;
          },
        });

        if (!isim) return;

        const tanim = await vscode.window.showInputBox({
          prompt: 'Project description (optional)',
          placeHolder: 'Enter project description',
        });

        await mcpClient.callTool('proje_olustur', {
          isim,
          tanim: tanim || '',
        });

        vscode.window.showInformationMessage('Project created successfully');
        await providers.projeTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to create project: ${error}`);
      }
    })
  );

  // Set Active Project
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SET_ACTIVE_PROJECT, async (item?: ProjeTreeItem) => {
      try {
        // If no item provided (e.g., command palette), show project picker
        if (!item) {
          const projects = await getProjectList(mcpClient);
          if (projects.length === 0) {
            vscode.window.showWarningMessage('No projects found. Create a project first.');
            return;
          }

          const selected = await vscode.window.showQuickPick(
            projects.map(p => ({
              label: p.isim,
              description: `${p.gorev_sayisi || 0} tasks`,
              project: p
            })),
            {
              placeHolder: 'Select a project to activate'
            }
          );

          if (!selected) return;

          await mcpClient.callTool('proje_aktif_yap', {
            proje_id: selected.project.id,
          });
          vscode.window.showInformationMessage(`"${selected.project.isim}" is now the active project`);
        } else {
          // Item provided from tree view
          if (item.isActive) {
            const deactivate = await vscode.window.showQuickPick(
              ['Deactivate', 'Cancel'],
              {
                placeHolder: 'This project is already active. Do you want to deactivate it?',
              }
            );

            if (deactivate === 'Deactivate') {
              await mcpClient.callTool('aktif_proje_kaldir');
              vscode.window.showInformationMessage('Project deactivated');
            }
          } else {
            await mcpClient.callTool('proje_aktif_yap', {
              proje_id: item.project.id,
            });
            vscode.window.showInformationMessage(`"${item.project.isim}" is now the active project`);
          }
        }

        await providers.projeTreeProvider.refresh();
        await providers.gorevTreeProvider.refresh();
        providers.statusBarManager.update();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to update active project: ${error}`);
      }
    })
  );
}