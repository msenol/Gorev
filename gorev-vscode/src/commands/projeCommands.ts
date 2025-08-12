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
    const errorMessage = error instanceof Error ? error.message : String(error);
    console.error(vscode.l10n.t('project.getListFailed', errorMessage));
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
          prompt: vscode.l10n.t('project.namePrompt'),
          placeHolder: vscode.l10n.t('project.namePlaceholder'),
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return vscode.l10n.t('project.nameRequired');
            }
            return null;
          },
        });

        if (!isim) return;

        const tanim = await vscode.window.showInputBox({
          prompt: vscode.l10n.t('project.descriptionPrompt'),
          placeHolder: vscode.l10n.t('project.descriptionPlaceholder'),
        });

        await mcpClient.callTool('proje_olustur', {
          isim,
          tanim: tanim || '',
        });

        vscode.window.showInformationMessage(vscode.l10n.t('project.createdSuccess'));
        await providers.projeTreeProvider.refresh();
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : String(error);
        vscode.window.showErrorMessage(vscode.l10n.t('project.createFailed', errorMessage));
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
            vscode.window.showWarningMessage(vscode.l10n.t('project.noProjectsFound'));
            return;
          }

          const selected = await vscode.window.showQuickPick(
            projects.map(p => ({
              label: p.isim,
              description: vscode.l10n.t('project.taskCount', (p.gorev_sayisi || 0).toString()),
              project: p
            })),
            {
              placeHolder: vscode.l10n.t('project.selectToActivate')
            }
          );

          if (!selected) return;

          await mcpClient.callTool('proje_aktif_yap', {
            proje_id: selected.project.id,
          });
          vscode.window.showInformationMessage(vscode.l10n.t('project.nowActive', selected.project.isim));
        } else {
          // Item provided from tree view
          if (item.isActive) {
            const deactivate = await vscode.window.showQuickPick(
              [vscode.l10n.t('project.deactivateOption'), vscode.l10n.t('project.cancelOption')],
              {
                placeHolder: vscode.l10n.t('project.alreadyActivePrompt'),
              }
            );

            if (deactivate === vscode.l10n.t('project.deactivateOption')) {
              await mcpClient.callTool('aktif_proje_kaldir');
              vscode.window.showInformationMessage(vscode.l10n.t('project.deactivated'));
            }
          } else {
            await mcpClient.callTool('proje_aktif_yap', {
              proje_id: item.project.id,
            });
            vscode.window.showInformationMessage(vscode.l10n.t('project.nowActive', item.project.isim));
          }
        }

        await providers.projeTreeProvider.refresh();
        await providers.gorevTreeProvider.refresh();
        providers.statusBarManager.update();
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : String(error);
        vscode.window.showErrorMessage(vscode.l10n.t('project.updateActiveFailed', errorMessage));
      }
    })
  );
}