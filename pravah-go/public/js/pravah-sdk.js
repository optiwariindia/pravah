/**
 * Pravah SDK
 * Managed real-time event streaming with Multi-Tab Synchronization.
 * Version: 1.0.1 (Generic Signal Support)
 */
class Pravah {
    constructor(options = {}) {
        this.baseUrl = options.baseUrl || ''; // Pravah service URL
        this.streams = options.streams || ['global'];
        this.subscriptionToken = options.subscriptionToken || '';
        this.id = null;
        this.eventSource = null;
        this.listeners = {};
        this.eventHandlers = {};
        this.lastEventId = null;
        
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
                if (payload.id) this.lastEventId = payload.id;
                this.emit(payload.event, payload.data);
            } else if (type === 'LEADER_PING' && this.isLeader) {
                // I am the leader, notify the new tab
                this.channel.postMessage({ type: 'LEADER_ALIVE' });
            } else if (type === 'LEADER_ALIVE') {
                // A leader exists, I will be a follower
                this.isLeader = false;
                this.closeEventSource();
            } else if (type === 'LEADER_CLOSED') {
                if (payload?.id) this.lastEventId = payload.id;
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
                this.channel.postMessage({ type: 'LEADER_CLOSED', payload: { id: this.lastEventId } });
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

        const streams = this.streams.map(encodeURIComponent).join(',');
        const lastEventId = this.lastEventId ? `&lastEventId=${encodeURIComponent(this.lastEventId)}` : '';
        const token = `&token=${encodeURIComponent(this.subscriptionToken)}`;
        const url = `${this.baseUrl}/pravah/subscribe?streams=${streams}${token}${lastEventId}`;
        this.eventSource = new EventSource(url, { withCredentials: true });

        this.attachSignalListener();
        Object.keys(this.listeners)
            .filter(event => event === 'connected' || event === 'heartbeat')
            .forEach(event => this.attachEventListener(event));

        this.eventSource.onerror = () => {
            console.error('[Pravah SDK] SSE connection interrupted. The browser is reconnecting...');
            this.emit('disconnected', { timestamp: new Date().toISOString() });
        };
    }

    attachEventListener(event) {
        if (!this.eventSource || this.eventHandlers[event]) return;

        const handler = (e) => {
            try {
                const data = JSON.parse(e.data);
                if (e.lastEventId) this.lastEventId = e.lastEventId;
                this.channel.postMessage({
                    type: 'EVENT_SIGNAL',
                    payload: { event, data, id: e.lastEventId || null }
                });
                this.emit(event, data);
            } catch (err) {
                console.warn('[Pravah SDK] Received non-JSON message:', e.data);
            }
        };

        this.eventHandlers[event] = handler;
        this.eventSource.addEventListener(event, handler);
    }

    attachSignalListener() {
        if (!this.eventSource || this.eventHandlers.signal) return;

        const handler = (e) => {
            try {
                const payload = JSON.parse(e.data);
                if (e.lastEventId) this.lastEventId = e.lastEventId;
                this.channel.postMessage({
                    type: 'EVENT_SIGNAL',
                    payload: { event: payload.event, data: payload.data, id: e.lastEventId || null }
                });
                this.emit(payload.event, payload.data);
            } catch (err) {
                console.warn('[Pravah SDK] Received invalid signal:', e.data);
            }
        };

        this.eventHandlers.signal = handler;
        this.eventSource.addEventListener('signal', handler);
    }

    closeEventSource() {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
            this.eventHandlers = {};
        }
    }

    /**
     * Public API: Listen for signals
     * Usage: pravah.on('bid-updated', (data) => { ... })
     */
    on(event, callback) {
        if (!this.listeners[event]) this.listeners[event] = [];
        this.listeners[event].push(callback);
        if (event === 'connected' || event === 'heartbeat') {
            this.attachEventListener(event);
        }
        console.log(`[Pravah SDK] Listener registered for: ${event}`);
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
