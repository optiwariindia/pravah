import App from '../models/app.js';

/**
 * Dual-Mode Publishing Authentication Middleware
 * 1. Checks for 'pravah_id' cookie (Client-side trigger)
 * 2. Checks for 'X-API-KEY' in headers (Server-side/Internal trigger)
 */
export const validatePublishAuth = async (req, res, next) => {
    // Mode 1: Cookie-based (Handshake ID)
    const pravahId = req.cookies.pravah_id;
    if (pravahId) {
        // If we have a valid session, we trust it for publishing
        return next();
    }

    // Mode 2: API Key-based (DB Lookup)
    const apiKey = req.headers['x-api-key'];
    if (apiKey) {
        try {
            const app = await App.findOne({ apiKey, status: 'active' });
            if (app) {
                req.appContext = app; // Attach app info for potential logging
                return next();
            }
        } catch (error) {
            console.error('[Pravah] DB Error during API key validation:', error);
            return res.status(500).json({ error: 'Internal Server Error', message: 'Auth check failed.' });
        }
    }

    // If neither is present or valid
    return res.status(403).json({ 
        error: 'Forbidden', 
        message: 'Unauthorized source. Missing pravah_id cookie or valid X-API-KEY.' 
    });
};
