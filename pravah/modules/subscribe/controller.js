import { v4 as uuidv4 } from 'uuid';
import clientManager from '../../core/clientManager.js';

/**
 * SSE Subscription Controller
 * Handles persistent connection and stream registration.
 */
export const subscribe = (req, res) => {
    // 1. Identify User (from Handshake cookie)
    const pravahId = req.cookies.pravah_id;

    if (!pravahId) {
        return res.status(401).json({ error: 'Unauthorized', message: 'No Pravah ID found. Please init first.' });
    }

    const connectionId = uuidv4();

    // 2. Parse Subscribed Streams
    const streams = req.query.streams 
        ? new Set(req.query.streams.split(',').map(s => s.trim()))
        : new Set(['global']);

    // 3. Set SSE Headers
    res.writeHead(200, {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
        'X-Accel-Buffering': 'no' // Important for Nginx
    });

    // 4. Send Initial Connection Data (Signal Format)
    const initData = {
        event: 'connected',
        data: {
            connectionId,
            pravahId,
            streams: Array.from(streams),
            timestamp: new Date().toISOString()
        }
    };
    res.write(`data: ${JSON.stringify(initData)}\n\n`);

    // 5. Register in ClientManager
    clientManager.addClient(connectionId, pravahId, streams, res);

    // 6. Heartbeat (Every 30 seconds)
    const heartbeatInterval = setInterval(() => {
        const heartbeat = {
            event: 'heartbeat',
            data: { timestamp: new Date().toISOString() }
        };
        res.write(`data: ${JSON.stringify(heartbeat)}\n\n`);
    }, 30000);

    // 7. Handle Disconnect
    req.on('close', () => {
        clearInterval(heartbeatInterval);
        clientManager.removeClient(connectionId);
        console.log(`[Pravah] Stream Closed: ${connectionId} (User: ${pravahId})`);
    });

    console.log(`[Pravah] Stream Opened: ${connectionId} (User: ${pravahId}) Streams: ${Array.from(streams).join(', ')}`);
};
