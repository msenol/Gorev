import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ClientInterface } from '../interfaces/client';
import { ApiClient, ApiError } from '../api/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { ProjeTreeItem } from '../providers/projeTreeProvider';
import { Proje } from '../models/proje';
import { Logger } from '../utils/logger';

async function getProjectList(apiClient: ApiClient): Promise<Proje[]> {
  try {
    const response = await apiClient.getProjects();
    if (!response.success || !response.data) {
      return [];
    }

    // Convert API Project[] to internal Proje[] model
    return response.data.map(project => ({
      id: project.id,
      isim: project.isim,
      tanim: project.tanim || '',
      gorev_sayisi: project.gorev_sayisi,
      olusturma_tarih: project.olusturma_tarihi,
      guncelleme_tarih: project.olusturma_tarihi,
    }));
  } catch (error) {
    if (error instanceof ApiError) {
      Logger.error(`[getProjectList] API Error ${error.statusCode}:`, error.apiError);
    } else {
      Logger.error('[getProjectList] Failed to get projects:', error);
    }
    return [];
  }
}

export function registerProjeCommands(
  context: vscode.ExtensionContext,
  mcpClient: ClientInterface,
  providers: CommandContext
): void {
  // Initialize API client
  const apiClient = mcpClient instanceof ApiClient ? mcpClient : new ApiClient();

  // Create Project
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_PROJECT, async () => {
      try {
        const isim = await vscode.window.showInputBox({
          prompt: t('project.namePrompt'),
          placeHolder: t('project.namePlaceholder'),
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return t('project.nameRequired');
            }
            return null;
          },
        });

        if (!isim) return;

        const tanim = await vscode.window.showInputBox({
          prompt: t('project.descriptionPrompt'),
          placeHolder: t('project.descriptionPlaceholder'),
        });

        // Use REST API to create project
        await apiClient.createProject({
          isim,
          tanim: tanim || '',
        });

        vscode.window.showInformationMessage(t('project.createdSuccess'));
        await providers.projeTreeProvider.refresh();
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[CreateProject] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('project.createFailed', error.apiError));
        } else {
          vscode.window.showErrorMessage(t('project.createFailed', String(error)));
        }
      }
    })
  );

  // Set Active Project
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SET_ACTIVE_PROJECT, async (item?: ProjeTreeItem) => {
      try {
        // If no item provided (e.g., command palette), show project picker
        if (!item) {
          const projects = await getProjectList(apiClient);
          if (projects.length === 0) {
            vscode.window.showWarningMessage(t('project.noProjectsFound'));
            return;
          }

          const selected = await vscode.window.showQuickPick(
            projects.map(p => ({
              label: p.isim,
              description: t('project.taskCount', (p.gorev_sayisi || 0).toString()),
              project: p
            })),
            {
              placeHolder: t('project.selectToActivate')
            }
          );

          if (!selected) return;

          // Use REST API to activate project
          await apiClient.activateProject(selected.project.id);
          vscode.window.showInformationMessage(t('project.nowActive', selected.project.isim));
        } else {
          // Item provided from tree view
          if (item.isActive) {
            const deactivate = await vscode.window.showQuickPick(
              [t('project.deactivateOption'), t('project.cancelOption')],
              {
                placeHolder: t('project.alreadyActivePrompt'),
              }
            );

            if (deactivate === t('project.deactivateOption')) {
              // Use REST API to remove active project
              await apiClient.removeActiveProject();
              vscode.window.showInformationMessage(t('project.deactivated'));
            }
          } else {
            // Use REST API to activate project
            await apiClient.activateProject(item.project.id);
            vscode.window.showInformationMessage(t('project.nowActive', item.project.isim));
          }
        }

        await providers.projeTreeProvider.refresh();
        await providers.gorevTreeProvider.refresh();
        providers.statusBarManager.update();
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[SetActiveProject] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('project.updateActiveFailed', error.apiError));
        } else {
          vscode.window.showErrorMessage(t('project.updateActiveFailed', String(error)));
        }
      }
    })
  );
}
