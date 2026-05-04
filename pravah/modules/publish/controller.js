import clientManager from '../../core/clientManager.js';
import Event from '../../models/event.js';

/**
 * Event Publishing Controller
 * Receives events from parent app and broadcasts to SSE clients.
 */
export const publish = async (req, res) => {
    try {
        const { streamId, event, data } = req.body;

        // 1. Basic Validation
        if (!streamId || !event || !data || !data.id) {
            return res.status(400).json({ 
                error: 'Bad Request', 
                message: 'streamId, event, and data.id are required.' 
            });
        }

        // 2. Persistence (Signal-Only)
        // Storing the signal in DB for replay support
        const newEvent = new Event({
            streamId,
            event,
            data: {
                id: data.id,
                metadata: data.metadata || {}
            }
        });
        await newEvent.save();

        // 3. Broadcast to matching clients
        const clients = clientManager.getClientsByStream(streamId);
        
        const signal = {
            event,
            data: {
                id: data.id,
                timestamp: newEvent.createdAt
            }
        };

        const signalString = `data: ${JSON.stringify(signal)}\n\n`;

        clients.forEach(client => {
            try {
                client.res.write(signalString);
            } catch (err) {
                console.error(`[Pravah] Failed to send to client ${client.id}:`, err.message);
            }
        });

        return res.status(202).json({
            success: true,
            message: `Event '${event}' published to ${clients.length} clients on stream '${streamId}'.`,
            eventId: newEvent._id
        });

    } catch (error) {
        console.error('[Pravah] Publish Error:', error);
        res.status(500).json({
            error: 'Internal Server Error',
            message: error.message
        });
    }
};
