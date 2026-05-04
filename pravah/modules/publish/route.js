import express from 'express';
import { publish } from './controller.js';
import { validatePublishAuth } from '../../middleware/auth.js';

const router = express.Router();

/**
 * POST /pravah/publish
 * Internal API to ingest and broadcast events.
 * Protected by Dual-Mode Auth (Cookie or DB-based API Key).
 */
router.post('/publish', validatePublishAuth, publish);

export default router;
