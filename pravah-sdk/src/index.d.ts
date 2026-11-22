export type ConnectionState =
    | 'initializing'
    | 'connecting'
    | 'connected'
    | 'reconnecting'
    | 'disconnected'
    | 'closed';

export interface PravahOptions {
    baseUrl?: string;
    streams?: string[];
    subscriptionToken?: string;
    getSubscriptionToken?: () => Promise<string>;
}

export const CONNECTION_STATES: Readonly<Record<string, ConnectionState>>;

export default class Pravah {
    constructor(options?: PravahOptions);
    connectionState: ConnectionState;
    isLeader: boolean;
    lastEventId: string | null;
    on(event: string, callback: (data: any) => void): void;
    off(event: string, callback: (data: any) => void): void;
    updateSubscriptionToken(token: string): void;
    close(): void;
}

export function usePravah(options: PravahOptions): Pravah | null;
export function usePravahListener(
    pravah: Pravah | null,
    event: string,
    callback: (data: any) => void
): void;
export function usePravahState<T = any>(pravah: Pravah | null, event: string): T | null;
