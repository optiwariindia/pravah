import * as React from 'react';

export type LogLevel = 'none' | 'error' | 'warn' | 'info';

export interface PravahOptions {
  baseUrl: string;
  streams?: string[];
  logLevel?: LogLevel;
}

export default class Pravah {
  constructor(options?: PravahOptions);
  baseUrl: string;
  streams: string[];
  logLevel: LogLevel;
  id: string | null;
  isLeader: boolean;
  on(event: string, callback: (data: any) => void): void;
  off(event: string, callback: (data: any) => void): void;
  closeEventSource(): void;
}

export interface PravahProviderProps {
  options: PravahOptions;
  children: React.ReactNode;
}

export const PravahProvider: React.FC<PravahProviderProps>;

export function usePravah(options?: PravahOptions): Pravah | null;

export function usePravahContext(): Pravah | null;

export function usePravahListener(
  pravah: Pravah | null | undefined,
  event: string,
  callback: (data: any) => void
): void;

export function usePravahState<T = any>(
  pravah: Pravah | null | undefined,
  event: string
): T | null;
