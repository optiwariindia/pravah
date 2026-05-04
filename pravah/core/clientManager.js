/**
 * ClientManager Class (Singleton)
 * Manages active SSE connections and their stream subscriptions.
 */
class ClientManager {
    constructor() {
        // Map to store clients: Key is a unique connectionId (UUID), value is the client object
        this.clients = new Map();
    }

    /**
     * Add a new client connection to the registry
     * @param {string} connectionId - Unique ID for this specific SSE connection
     * @param {string} pravahId - The identifier from the cookie
     * @param {Set<string>} streams - Set of streams the client is subscribed to
     * @param {Object} res - The Express response object for the SSE stream
     */
    addClient(connectionId, pravahId, streams, res) {
        const client = {
            id: connectionId,
            pravahId,
            streams,
            connectedAt: new Date(),
            res
        };
        this.clients.set(connectionId, client);
        return client;
    }

    /**
     * Remove a client from the registry (on disconnect)
     * @param {string} connectionId
     */
    removeClient(connectionId) {
        if (this.clients.has(connectionId)) {
            this.clients.delete(connectionId);
            return true;
        }
        return false;
    }

    /**
     * Get all clients subscribed to a specific stream
     * @param {string} streamId
     * @returns {Array} List of matching client objects
     */
    getClientsByStream(streamId) {
        const matchingClients = [];
        for (const client of this.clients.values()) {
            if (client.streams.has(streamId)) {
                matchingClients.push(client);
            }
        }
        return matchingClients;
    }

    /**
     * Get total active connection count
     * @returns {number}
     */
    getConnectionCount() {
        return this.clients.size;
    }

    /**
     * Get all clients (useful for monitoring/debugging)
     */
    getAllClients() {
        return Array.from(this.clients.values()).map(c => ({
            id: c.id,
            pravahId: c.pravahId,
            streams: Array.from(c.streams),
            connectedAt: c.connectedAt
        }));
    }
}

const clientManager = new ClientManager();
export default clientManager;
