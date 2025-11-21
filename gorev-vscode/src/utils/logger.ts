import * as vscode from 'vscode';

export enum LogLevel {
  Debug = 0,
  Info = 1,
  Warning = 2,
  Error = 3,
}

export class Logger {
  private static outputChannel: vscode.OutputChannel;
  private static logLevel: LogLevel = LogLevel.Info;

  static {
    this.outputChannel = vscode.window.createOutputChannel('Gorev');
  }

  static setLogLevel(level: LogLevel): void {
    this.logLevel = level;
  }

  static debug(...args: unknown[]): void {
    if (this.logLevel <= LogLevel.Debug) {
      this.log('DEBUG', ...args);
    }
  }

  static info(...args: unknown[]): void {
    if (this.logLevel <= LogLevel.Info) {
      this.log('INFO', ...args);
    }
  }

  static warn(...args: unknown[]): void {
    if (this.logLevel <= LogLevel.Warning) {
      this.log('WARN', ...args);
    }
  }

  static error(...args: unknown[]): void {
    if (this.logLevel <= LogLevel.Error) {
      this.log('ERROR', ...args);
    }
  }

  static show(): void {
    this.outputChannel.show();
  }

  private static log(level: string, ...args: unknown[]): void {
    const timestamp = new Date().toISOString();
    const message = args
      .map((arg) => {
        if (typeof arg === 'object') {
          try {
            return JSON.stringify(arg, null, 2);
          } catch {
            return String(arg);
          }
        }
        return String(arg);
      })
      .join(' ');

    this.outputChannel.appendLine(`[${timestamp}] [${level}] ${message}`);
  }
}