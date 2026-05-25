import { Router } from 'express';
import visitorController from '../controllers/VisitorController.js';

const router = Router();

router.get('/', visitorController.index);
router.get('/pricing', visitorController.pricing);

// Dynamic Policy Routes
router.get('/privacy-policy', (req, res, next) => {
    req.params.policyName = 'privacy-policy';
    visitorController.showPolicy(req, res, next);
});

router.get('/terms-of-service', (req, res, next) => {
    req.params.policyName = 'terms-of-service';
    visitorController.showPolicy(req, res, next);
});

export default router;
