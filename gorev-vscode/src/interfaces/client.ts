import { EventEmitter } from 'events';

export interface ClientInterface extends EventEmitter {
  connect(serverPath?: string): Promise<void>;
  disconnect(): void;
  isConnected(): boolean;
  callTool(name: string, params?: any): Promise<any>;
  getTools(): any[];
}

export type ClientType = ClientInterface;