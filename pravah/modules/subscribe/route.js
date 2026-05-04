import express from 'express';
import { subscribe } from './controller.js';

const router = express.Router();

/**
 * GET /pravah/subscribe
 * SSE stream endpoint.
 */
router.get('/subscribe', subscribe);

export default router;
