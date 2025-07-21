import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { TemplateTreeItem } from '../providers/templateTreeProvider';
import { TemplateWizard } from '../ui/templateWizard';
import { GorevTemplate } from '../models/template';
import { Logger } from '../utils/logger';
import { COMMANDS } from '../utils/constants';
import { MarkdownParser } from '../utils/markdownParser';

export function registerTemplateCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  // Open template wizard
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.OPEN_TEMPLATE_WIZARD, async (templateId?: string) => {
      try {
        await TemplateWizard.show(mcpClient, context.extensionUri, templateId);
      } catch (error) {
        Logger.error('Failed to open template wizard:', error);
        vscode.window.showErrorMessage(vscode.l10n.t('template.wizardOpenFailed'));
      }
    })
  );

  // Create task from template (from tree view)
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_FROM_TEMPLATE, async (item?: TemplateTreeItem) => {
      try {
        if (!item || !item.template) {
          // If no template provided, open wizard
          await vscode.commands.executeCommand(COMMANDS.OPEN_TEMPLATE_WIZARD);
          return;
        }

        // Open wizard with selected template
        await TemplateWizard.show(mcpClient, context.extensionUri, item.template.id);
      } catch (error) {
        Logger.error('Failed to create task from template:', error);
        vscode.window.showErrorMessage(vscode.l10n.t('template.createFromFailed'));
      }
    })
  );

  // Quick create from template (with quick pick)
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.QUICK_CREATE_FROM_TEMPLATE, async () => {
      try {
        // Load all templates
        const result = await mcpClient.callTool('template_listele');
        const templates = MarkdownParser.parseTemplateListesi(result.content[0].text);

        if (templates.length === 0) {
          vscode.window.showInformationMessage(vscode.l10n.t('template.noTemplatesDefined'));
          return;
        }

        // Show quick pick
        const items = templates.map((t: GorevTemplate) => ({
          label: t.isim,
          description: t.kategori,
          detail: t.tanim,
          template: t
        }));

        const selected = await vscode.window.showQuickPick(items, {
          placeHolder: vscode.l10n.t('template.selectTemplate'),
          matchOnDescription: true,
          matchOnDetail: true
        });

        if (selected) {
          await TemplateWizard.show(mcpClient, context.extensionUri, selected.template.id);
        }
      } catch (error) {
        Logger.error('Failed to show template quick pick:', error);
        vscode.window.showErrorMessage(vscode.l10n.t('template.listLoadFailed'));
      }
    })
  );

  // Refresh templates
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.REFRESH_TEMPLATES, async () => {
      try {
        await providers.templateTreeProvider.refresh();
        vscode.window.showInformationMessage(vscode.l10n.t('template.refreshed'));
      } catch (error) {
        Logger.error('Failed to refresh templates:', error);
        vscode.window.showErrorMessage(vscode.l10n.t('template.refreshFailed'));
      }
    })
  );

  // Initialize default templates
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.INIT_DEFAULT_TEMPLATES, async () => {
      try {
        const answer = await vscode.window.showWarningMessage(
          vscode.l10n.t('template.initConfirm'),
          { modal: true },
          vscode.l10n.t('template.initConfirmYes')
        );

        if (answer !== vscode.l10n.t('template.initConfirmYes')) {
          return;
        }

        // Get the gorev server path from configuration
        const serverPath = vscode.workspace.getConfiguration('gorev').get<string>('serverPath');
        
        if (!serverPath) {
          vscode.window.showErrorMessage(vscode.l10n.t('template.serverPathNotConfigured'));
          return;
        }
        
        // Call gorev template init command
        const terminal = vscode.window.createTerminal('Gorev Template Init');
        terminal.sendText(`"${serverPath}" template init`);
        terminal.show();

        // Wait a bit and refresh
        setTimeout(async () => {
          await providers.templateTreeProvider.refresh();
          vscode.window.showInformationMessage(vscode.l10n.t('template.defaultsLoaded'));
        }, 2000);
      } catch (error) {
        Logger.error('Failed to initialize templates:', error);
        vscode.window.showErrorMessage(vscode.l10n.t('template.initFailed'));
      }
    })
  );

  // Show template details
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SHOW_TEMPLATE_DETAILS, async (item?: TemplateTreeItem) => {
      if (!item || !item.template) {
        return;
      }

      const template = item.template;
      const panel = vscode.window.createWebviewPanel(
        'templateDetails',
        template.isim,
        vscode.ViewColumn.One,
        {}
      );

      panel.webview.html = getTemplateDetailsHtml(template);
    })
  );

  // Export template as JSON
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.EXPORT_TEMPLATE, async (item?: TemplateTreeItem) => {
      if (!item || !item.template) {
        return;
      }

      const template = item.template;
      const uri = await vscode.window.showSaveDialog({
        defaultUri: vscode.Uri.file(`${template.isim.replace(/\s+/g, '_')}.json`),
        filters: {
          [vscode.l10n.t('template.jsonFilter')]: ['json']
        }
      });

      if (uri) {
        const content = JSON.stringify(template, null, 2);
        await vscode.workspace.fs.writeFile(uri, Buffer.from(content, 'utf-8'));
        vscode.window.showInformationMessage(vscode.l10n.t('template.exported'));
      }
    })
  );
}

function getTemplateDetailsHtml(template: GorevTemplate): string {
  return `
  <!DOCTYPE html>
  <html lang="tr">
  <head>
      <meta charset="UTF-8">
      <style>
          body {
              font-family: var(--vscode-font-family);
              color: var(--vscode-foreground);
              background-color: var(--vscode-editor-background);
              padding: 20px;
              line-height: 1.6;
          }
          h1 {
              color: var(--vscode-foreground);
              border-bottom: 1px solid var(--vscode-widget-border);
              padding-bottom: 10px;
          }
          .category {
              display: inline-block;
              padding: 4px 8px;
              background-color: var(--vscode-badge-background);
              color: var(--vscode-badge-foreground);
              border-radius: 3px;
              font-size: 12px;
              margin-bottom: 20px;
          }
          .description {
              color: var(--vscode-descriptionForeground);
              margin-bottom: 30px;
          }
          .field {
              margin-bottom: 20px;
              padding: 15px;
              background-color: var(--vscode-editor-background);
              border: 1px solid var(--vscode-widget-border);
              border-radius: 4px;
          }
          .field-name {
              font-weight: bold;
              color: var(--vscode-foreground);
          }
          .field-type {
              color: var(--vscode-textLink-foreground);
              font-size: 12px;
          }
          .field-required {
              color: var(--vscode-errorForeground);
              font-size: 12px;
          }
          .field-description {
              color: var(--vscode-descriptionForeground);
              font-size: 13px;
              margin-top: 5px;
          }
          code {
              background-color: var(--vscode-textCodeBlock-background);
              padding: 2px 4px;
              border-radius: 3px;
              font-family: var(--vscode-editor-font-family);
          }
      </style>
  </head>
  <body>
      <h1>${template.isim}</h1>
      <div class="category">${template.kategori || vscode.l10n.t('template.general')}</div>
      ${template.tanim ? `<div class="description">${template.tanim}</div>` : ''}
      
      <h2>${vscode.l10n.t('template.fields')}</h2>
      ${template.alanlar.map(field => `
          <div class="field">
              <div>
                  <span class="field-name">${field.isim}</span>
                  <span class="field-type">(${field.tur})</span>
                  ${field.zorunlu ? `<span class="field-required">${vscode.l10n.t('template.fieldRequired')}</span>` : ''}
              </div>
              ${field.varsayilan ? `<div class="field-description">${vscode.l10n.t('template.fieldDefault')}: <code>${field.varsayilan}</code></div>` : ''}
          </div>
      `).join('')}
      
      ${template.varsayilan_baslik ? `
          <h2>${vscode.l10n.t('template.defaultTitle')}</h2>
          <p><code>${template.varsayilan_baslik}</code></p>
      ` : ''}
      
      ${template.aciklama_template ? `
          <h2>${vscode.l10n.t('template.descriptionTemplate')}</h2>
          <pre>${template.aciklama_template}</pre>
      ` : ''}
  </body>
  </html>
  `;
}