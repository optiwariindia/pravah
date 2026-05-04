import express from 'express';
import { v4 as uuidv4 } from 'uuid';

const router = express.Router();

/**
 * Handshake Logic
 * Manages the 'pravah_id' identity.
 * Sets a persistent HttpOnly cookie if it doesn't exist.
 */
router.get('/init', (req, res) => {
    try {
        let pravahId = req.cookies.pravah_id;
        let isNew = false;

        if (!pravahId) {
            pravahId = uuidv4();
            isNew = true;

            // Set the cookie
            // Note: HttpOnly prevents JS access (Security), Lax prevents CSRF risk
            res.cookie('pravah_id', pravahId, {
                httpOnly: true,
                secure: process.env.NODE_ENV === 'production',
                sameSite: 'Lax',
                maxAge: 365 * 24 * 60 * 60 * 1000, // 1 year
                path: '/'
            });
        }

        return res.status(200).json({
            success: true,
            pravah_id: pravahId,
            new: isNew,
            message: isNew ? 'Pravah ID initialized' : 'Pravah ID recognized'
        });
    } catch (error) {
        res.status(500).json({
            error: 'Handshake Error',
            message: error.message
        });
    }
});

export default router;
