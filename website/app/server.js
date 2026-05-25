import express from 'express';
import { ExpressServer, configTest } from 'express-web-tools';
import path from 'path';
import { fileURLToPath } from 'url';

// Resolving __dirname for ES Modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Ensure required environment variables are present
configTest(['MONGODB', 'PORT']);

const server = new ExpressServer();

// View Engine Setup
server.views('twig', path.join(__dirname, 'views'));

// Static Files
server.addMiddleware(express.static(path.join(__dirname, 'public')));

// Standard Middleware
server.addMiddleware(express.json());
server.addMiddleware(express.urlencoded({ extended: true }));

// Routing
import visitorRouter from './modules/visitor/routes/index.js';
server.addRoute('/', visitorRouter);

// Start Server
const port = process.env.PORT || 3000;
server.start(port);
