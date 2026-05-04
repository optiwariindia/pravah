/**
 * Pravah SDK
 * Managed real-time event streaming with Multi-Tab Synchronization.
 * Version: 1.0.1 (Generic Signal Support)
 */
export default class Pravah {
    constructor(options = {}) {
        this.baseUrl = options.baseUrl || ''; // Pravah service URL
        this.streams = options.streams || ['global'];
        this.id = null;
        this.eventSource = null;
        this.listeners = {};
        
        // Multi-tab synchronization using BroadcastChannel
        this.channel = new BroadcastChannel('pravah_stream');
        this.isLeader = false;

        this.init();
    }

    /**
     * Initialization: Handshake -> Leader Election -> Connection
     */
    async init() {
        try {
            // 1. Handshake to ensure pravah_id cookie is set
            const response = await fetch(`${this.baseUrl}/pravah/init`, { credentials: 'include' });
            const data = await response.json();
            
            if (data.success) {
                console.log('[Pravah SDK] Handshake successful:', data.message);
                this.id = data.pravah_id;
                
                // 2. Setup tab synchronization
                this.setupTabSync();
                
                // 3. Start Connection logic
                this.manageConnection();
            }
        } catch (error) {
            console.error('[Pravah SDK] Initialization failed:', error);
        }
    }

    /**
     * Leader Election & Message Handling
     */
    setupTabSync() {
        this.channel.onmessage = (msg) => {
            const { type, payload } = msg.data;
            
            if (type === 'EVENT_SIGNAL') {
                // Received from Leader tab
                this.emit(payload.event, payload.data);
            } else if (type === 'LEADER_PING' && this.isLeader) {
                // I am the leader, notify the new tab
                this.channel.postMessage({ type: 'LEADER_ALIVE' });
            } else if (type === 'LEADER_ALIVE') {
                // A leader exists, I will be a follower
                this.isLeader = false;
                this.closeEventSource();
            } else if (type === 'LEADER_CLOSED') {
                // Leader closed, try to become leader
                this.manageConnection();
            }
        };

        // Check if a leader already exists
        this.channel.postMessage({ type: 'LEADER_PING' });
        
        // If no leader responds in 400ms, claim leadership
        setTimeout(() => {
            if (!this.eventSource && !this.isLeader) {
                this.becomeLeader();
            }
        }, 400);

        // Handle page close
        window.addEventListener('beforeunload', () => {
            if (this.isLeader) {
                this.channel.postMessage({ type: 'LEADER_CLOSED' });
            }
        });
    }

    becomeLeader() {
        console.log('[Pravah SDK] Connection active in this tab (Leader).');
        this.isLeader = true;
        this.connectSSE();
    }

    manageConnection() {
        // Simple logic to ensure only one connection
        if (!this.isLeader) {
            this.channel.postMessage({ type: 'LEADER_PING' });
        }
    }

    /**
     * Establish SSE Connection (Only for Leader)
     */
    connectSSE() {
        if (this.eventSource || !this.isLeader) return;

        const url = `${this.baseUrl}/pravah/subscribe?streams=${this.streams.join(',')}`;
        this.eventSource = new EventSource(url, { withCredentials: true });

        /**
         * Generic Message Handler
         * Expects data format: { "event": "event-name", "data": { "id": "xyz" } }
         */
        this.eventSource.onmessage = (e) => {
            try {
                const payload = JSON.parse(e.data);
                const { event, data } = payload;
                
                if (event) {
                    // 1. Broadcast to other tabs
                    this.channel.postMessage({ type: 'EVENT_SIGNAL', payload: { event, data } });
                    // 2. Emit locally
                    this.emit(event, data);
                }
            } catch (err) {
                console.warn('[Pravah SDK] Received non-JSON message:', e.data);
            }
        };

        this.eventSource.addEventListener('connected', (e) => {
            console.log('[Pravah SDK] SSE Stream connected.');
        });

        this.eventSource.onerror = (err) => {
            console.error('[Pravah SDK] SSE Error. Reconnecting in 5s...');
            this.emit('disconnected', { timestamp: new Date().toISOString() });
            this.closeEventSource();
            setTimeout(() => this.connectSSE(), 5000);
        };
    }

    closeEventSource() {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }
    }

    /**
     * Public API: Listen for signals
     * Usage: pravah.on('bid-updated', (data) => { ... })
     */
    on(event, callback) {
        if (!this.listeners[event]) this.listeners[event] = [];
        this.listeners[event].push(callback);
        console.log(`[Pravah SDK] Listener registered for: ${event}`);
    }

    /**
     * Public API: Remove listener
     */
    off(event, callback) {
        if (!this.listeners[event]) return;
        this.listeners[event] = this.listeners[event].filter(cb => cb !== callback);
        console.log(`[Pravah SDK] Listener removed for: ${event}`);
    }

    /**
     * Internal: Trigger local listeners
     */
    emit(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(cb => cb(data));
        }
    }
}
export * from './hooks';
