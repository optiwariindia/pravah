import mongoose from 'mongoose';
import clientManager from '../../core/clientManager.js';

/**
 * Health Check Controller
 * Returns service status and active connection metrics.
 */
export const getHealth = async (req, res) => {
    const dbStatus = mongoose.connection.readyState === 1 ? 'connected' : 'disconnected';
    
    const healthInfo = {
        status: 'healthy',
        service: 'Pravah',
        timestamp: new Date().toISOString(),
        uptime: process.uptime(),
        database: dbStatus,
        metrics: {
            activeConnections: clientManager.getConnectionCount(),
        }
    };

    // If DB is down, mark as unhealthy
    if (dbStatus !== 'connected') {
        healthInfo.status = 'degraded';
    }

    res.status(healthInfo.status === 'healthy' ? 200 : 503).json(healthInfo);
};
