const CONNECTION_STATES = Object.freeze({
    INITIALIZING: 'initializing',
    CONNECTING: 'connecting',
    CONNECTED: 'connected',
    RECONNECTING: 'reconnecting',
    DISCONNECTED: 'disconnected',
    CLOSED: 'closed'
});

export { CONNECTION_STATES };

export default class Pravah {
    constructor(options = {}) {
        this.baseUrl = options.baseUrl || '';
        this.streams = options.streams?.length ? [...options.streams].sort() : ['global'];
        this.subscriptionToken = options.subscriptionToken || '';
        this.getSubscriptionToken = options.getSubscriptionToken || null;
        this.id = null;
        this.eventSource = null;
        this.listeners = {};
        this.eventHandlers = {};
        this.lastEventId = null;
        this.connectionState = CONNECTION_STATES.INITIALIZING;
        this.channel = null;
        this.isLeader = false;
        this.closed = false;
        this.connecting = false;
        this.reconnectAttempt = 0;
        this.reconnectTimer = null;
        this.electionTimer = null;
        this.heartbeatTimer = null;
        this.electing = false;
        this.electionCandidates = new Set();
        this.lastLeaderSeenAt = 0;
        this.tabId = this.createTabId();
        this.channelName = '';
        this.handleBeforeUnload = () => this.releaseLeadership();

        this.init();
    }

    async init() {
        try {
            if (typeof fetch === 'undefined') {
                throw new Error('Fetch is required by the Pravah SDK');
            }
            const response = await fetch(`${this.baseUrl}/pravah/init`, { credentials: 'include' });
            if (!response.ok) {
                throw new Error(`Handshake failed with status ${response.status}`);
            }
            const data = await response.json();
            if (!data.success) {
                throw new Error('Handshake response was unsuccessful');
            }
            if (this.closed) return;
            if (typeof EventSource === 'undefined') {
                throw new Error('EventSource is required by the Pravah SDK');
            }

            this.id = data.pravah_id;
            await this.resolveSubscriptionToken(true);
            this.channelName = this.createChannelName();
            this.setupTabSync();
            this.manageConnection();
        } catch (error) {
            this.setConnectionState(CONNECTION_STATES.DISCONNECTED);
            this.emit('error', error);
            console.error('[Pravah SDK] Initialization failed:', error);
        }
    }

    createTabId() {
        if (globalThis.crypto?.randomUUID) return globalThis.crypto.randomUUID();
        return `${Date.now()}-${Math.random().toString(36).slice(2)}`;
    }

    createChannelName() {
        const scope = JSON.stringify([this.baseUrl, this.streams, this.subscriptionTenantId()]);
        let hash = 2166136261;
        for (const character of scope) {
            hash ^= character.charCodeAt(0);
            hash = Math.imul(hash, 16777619);
        }
        return `pravah_stream_${(hash >>> 0).toString(16)}`;
    }

    subscriptionTenantId() {
        try {
            const payload = this.subscriptionToken.split('.')[1];
            const base64 = payload.replace(/-/g, '+').replace(/_/g, '/');
            const decoded = globalThis.atob(base64.padEnd(Math.ceil(base64.length / 4) * 4, '='));
            return JSON.parse(decoded).appId || this.subscriptionToken;
        } catch {
            return this.subscriptionToken;
        }
    }

    setupTabSync() {
        if (typeof BroadcastChannel === 'undefined') {
            this.becomeLeader();
            return;
        }

        this.channel = new BroadcastChannel(this.channelName);
        this.channel.onmessage = ({ data }) => {
            const { type, payload = {} } = data || {};
            if (type === 'EVENT_SIGNAL') {
                if (payload.id) this.lastEventId = payload.id;
                this.emit('signal', { event: payload.event, data: payload.data });
                this.emit(payload.event, payload.data);
            } else if (type === 'LEADER_PING' && this.isLeader && payload.target) {
                this.postMessage('LEADER_ALIVE', { target: payload.target, id: this.lastEventId });
            } else if (type === 'LEADER_ALIVE' && payload.target === this.tabId) {
                this.lastLeaderSeenAt = Date.now();
                this.cancelElection();
                this.isLeader = false;
                this.closeEventSource();
            } else if (type === 'LEADER_HEARTBEAT') {
                this.lastLeaderSeenAt = Date.now();
                if (this.isLeader && payload.sender !== this.tabId) {
                    this.compareLeadership(payload.sender);
                }
            } else if (type === 'LEADER_CLOSED') {
                if (payload.id) this.lastEventId = payload.id;
                this.lastLeaderSeenAt = 0;
                this.scheduleElection(50);
            } else if (type === 'LEADER_CLAIM') {
                if (this.isLeader) {
                    this.postMessage('LEADER_ALIVE', { target: payload.sender, id: this.lastEventId });
                } else {
                    this.startElection(payload.sender);
                }
            }
        };

        globalThis.window?.addEventListener('beforeunload', this.handleBeforeUnload);
        this.postMessage('LEADER_PING', { target: this.tabId });
        this.scheduleElection(400);
        this.heartbeatTimer = setInterval(() => {
            if (this.closed) return;
            if (this.isLeader) {
                this.postMessage('LEADER_HEARTBEAT', { sender: this.tabId, id: this.lastEventId });
            } else if (Date.now() - this.lastLeaderSeenAt > 6000) {
                this.startElection();
            }
        }, 2000);
    }

    postMessage(type, payload = {}) {
        this.channel?.postMessage({ type, payload });
    }

    manageConnection() {
        if (this.closed) return;
        if (!this.channel) return;
        this.postMessage('LEADER_PING', { target: this.tabId });
        this.scheduleElection(400);
    }

    scheduleElection(delay) {
        if (this.closed || this.electionTimer) return;
        this.electionTimer = setTimeout(() => {
            this.electionTimer = null;
            if (!this.lastLeaderSeenAt && !this.isLeader) this.startElection();
        }, delay);
    }

    startElection(peerTabId = null) {
        if (this.closed || this.isLeader) return;
        if (this.electing) {
            if (peerTabId) this.electionCandidates.add(peerTabId);
            return;
        }
        this.cancelElection();
        this.electing = true;
        this.electionCandidates = new Set([this.tabId, ...(peerTabId ? [peerTabId] : [])]);
        this.postMessage('LEADER_CLAIM', { sender: this.tabId });
        this.electionTimer = setTimeout(() => {
            this.electionTimer = null;
            this.electing = false;
            const winner = [...this.electionCandidates].sort()[0];
            if (winner === this.tabId) {
                this.becomeLeader();
            }
        }, 250);
    }

    cancelElection() {
        if (this.electionTimer) clearTimeout(this.electionTimer);
        this.electionTimer = null;
        this.electing = false;
    }

    compareLeadership(otherTabId) {
        if (!otherTabId || this.tabId <= otherTabId) return;
        this.releaseLeadership();
        this.lastLeaderSeenAt = Date.now();
    }

    async becomeLeader() {
        if (this.closed || this.isLeader) return;
        this.cancelElection();
        this.isLeader = true;
        this.postMessage('LEADER_HEARTBEAT', { sender: this.tabId, id: this.lastEventId });
        await this.connectSSE();
    }

    async resolveSubscriptionToken(refresh = false) {
        if (this.getSubscriptionToken && (refresh || !this.subscriptionToken)) {
            const token = await this.getSubscriptionToken();
            if (typeof token !== 'string' || !token) {
                throw new Error('getSubscriptionToken did not return a token');
            }
            this.subscriptionToken = token;
        }
        if (!this.subscriptionToken) {
            throw new Error('A subscription token is required');
        }
        return this.subscriptionToken;
    }

    async connectSSE() {
        if (this.closed || !this.isLeader || this.eventSource || this.connecting) return;
        this.connecting = true;
        this.setConnectionState(
            this.reconnectAttempt ? CONNECTION_STATES.RECONNECTING : CONNECTION_STATES.CONNECTING
        );

        try {
            const token = await this.resolveSubscriptionToken(this.reconnectAttempt > 0);
            if (this.closed || !this.isLeader) return;

            const streams = this.streams.map(encodeURIComponent).join(',');
            const lastEventId = this.lastEventId ? `&lastEventId=${encodeURIComponent(this.lastEventId)}` : '';
            const url = `${this.baseUrl}/pravah/subscribe?streams=${streams}&token=${encodeURIComponent(token)}${lastEventId}`;
            this.eventSource = new EventSource(url, { withCredentials: true });
            this.attachSignalListener();
            this.attachEventListener('connected');
            this.attachEventListener('heartbeat');

            this.eventSource.onerror = () => {
                if (this.closed || !this.isLeader) return;
                this.setConnectionState(CONNECTION_STATES.RECONNECTING);
                this.emit('disconnected', { timestamp: new Date().toISOString() });
                this.closeEventSource();
                this.scheduleReconnect();
            };
        } catch (error) {
            if (!this.closed) {
                this.setConnectionState(CONNECTION_STATES.DISCONNECTED);
                this.emit('error', error);
                this.scheduleReconnect();
            }
        } finally {
            this.connecting = false;
        }
    }

    scheduleReconnect() {
        if (this.closed || this.reconnectTimer || !this.isLeader) return;
        const delay = Math.min(1000 * (2 ** this.reconnectAttempt), 15000);
        this.reconnectAttempt += 1;
        this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null;
            this.connectSSE();
        }, delay);
    }

    attachEventListener(event) {
        if (!this.eventSource || this.eventHandlers[event]) return;
        const handler = (messageEvent) => {
            try {
                const data = JSON.parse(messageEvent.data);
                if (event === 'connected') {
                    this.reconnectAttempt = 0;
                    this.setConnectionState(CONNECTION_STATES.CONNECTED);
                }
                this.broadcastEvent(event, data, messageEvent.lastEventId);
            } catch (error) {
                console.warn('[Pravah SDK] Received non-JSON message:', messageEvent.data);
            }
        };
        this.eventHandlers[event] = handler;
        this.eventSource.addEventListener(event, handler);
    }

    attachSignalListener() {
        if (!this.eventSource || this.eventHandlers.signal) return;
        const handler = (messageEvent) => {
            try {
                const payload = JSON.parse(messageEvent.data);
                this.broadcastEvent(payload.event, payload.data, messageEvent.lastEventId);
            } catch (error) {
                console.warn('[Pravah SDK] Received invalid signal:', messageEvent.data);
            }
        };
        this.eventHandlers.signal = handler;
        this.eventSource.addEventListener('signal', handler);
    }

    broadcastEvent(event, data, lastEventId) {
        if (lastEventId) this.lastEventId = lastEventId;
        this.postMessage('EVENT_SIGNAL', { event, data, id: lastEventId || null });
        this.emit('signal', { event, data });
        this.emit(event, data);
    }

    closeEventSource() {
        if (!this.eventSource) return;
        this.eventSource.close();
        this.eventSource = null;
        this.eventHandlers = {};
    }

    releaseLeadership() {
        if (!this.isLeader) return;
        this.postMessage('LEADER_CLOSED', { id: this.lastEventId });
        this.isLeader = false;
    }

    close() {
        if (this.closed) return;
        this.setConnectionState(CONNECTION_STATES.CLOSED);
        this.closed = true;
        this.cancelElection();
        if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
        if (this.heartbeatTimer) clearInterval(this.heartbeatTimer);
        this.releaseLeadership();
        this.closeEventSource();
        globalThis.window?.removeEventListener('beforeunload', this.handleBeforeUnload);
        this.channel?.close();
        this.channel = null;
    }

    setConnectionState(state) {
        if (this.connectionState === state) return;
        this.connectionState = state;
        this.emit('statechange', state);
    }

    updateSubscriptionToken(token) {
        if (typeof token !== 'string' || !token) {
            throw new Error('A non-empty subscription token is required');
        }
        this.subscriptionToken = token;
        if (this.isLeader && !this.eventSource) {
            if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
            this.reconnectAttempt = 0;
            this.connectSSE();
        }
    }

    on(event, callback) {
        if (typeof callback !== 'function') return;
        if (!this.listeners[event]) this.listeners[event] = [];
        this.listeners[event].push(callback);
    }

    off(event, callback) {
        if (!this.listeners[event]) return;
        this.listeners[event] = this.listeners[event].filter(listener => listener !== callback);
    }

    emit(event, data) {
        for (const listener of [...(this.listeners[event] || [])]) {
            listener(data);
        }
    }
}
