import express from 'express';
import { registerApp, listApps } from './controller.js';

const router = express.Router();

/**
 * PUT /register
 * Register a new application to get an API Key.
 */
router.put('/register', registerApp);

/**
 * GET /pravah/apps
 * List all registered applications.
 */
router.get('/apps', listApps);

export default router;
