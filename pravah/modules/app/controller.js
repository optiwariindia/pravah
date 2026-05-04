import App from '../../models/app.js';

/**
 * Register a new application
 * Generates an appId and apiKey for the application.
 */
export const registerApp = async (req, res) => {
    try {
        const { name } = req.body;

        if (!name) {
            return res.status(400).json({ 
                error: 'Bad Request', 
                message: 'Application name is required.' 
            });
        }

        // Check if app already exists
        const existingApp = await App.findOne({ name });
        if (existingApp) {
            return res.status(409).json({ 
                error: 'Conflict', 
                message: 'An application with this name already exists.' 
            });
        }

        const newApp = new App({ name });
        await newApp.save();

        return res.status(201).json({
            success: true,
            message: 'Application registered successfully.',
            data: {
                name: newApp.name,
                appId: newApp.appId,
                apiKey: newApp.apiKey,
                status: newApp.status
            }
        });

    } catch (error) {
        console.error('[Pravah] App Registration Error:', error);
        res.status(500).json({
            error: 'Internal Server Error',
            message: error.message
        });
    }
};

/**
 * Get all registered apps (Internal/Admin use)
 */
export const listApps = async (req, res) => {
    try {
        const apps = await App.find({}, 'name appId status createdAt');
        res.status(200).json({ success: true, data: apps });
    } catch (error) {
        res.status(500).json({ error: 'Internal Server Error', message: error.message });
    }
};
