import express from 'express';
import { getHealth } from './controller.js';

const router = express.Router();

/**
 * GET /pravah/health
 * Public health check endpoint.
 */
router.get('/health', getHealth);

export default router;
