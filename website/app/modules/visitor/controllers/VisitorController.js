import { CrudController, asyncHandler, API } from 'express-web-tools';
import { Policy } from '../models/Policy.js';

class VisitorController extends CrudController {
    constructor() {
        super(Policy);
    }

    /**
     * Render Home Page
     */
    index = asyncHandler(async (req, res) => {
        return res.render('pages/index');
    });

    /**
     * Render Pricing Page
     */
    pricing = asyncHandler(async (req, res) => {
        return res.render('pages/pricing');
    });

    /**
     * Render Dynamic Policy
     */
    showPolicy = asyncHandler(async (req, res) => {
        const name = req.params.policyName;
        const policy = await Policy.findOne({ name, isDeleted: false });
        
        if (!policy) {
            return res.status(404).render('pages/404', { message: 'Policy not found' });
        }

        return res.render('pages/policy', { policy });
    });
}

export default new VisitorController();
