import mongoose from 'mongoose';
import { v4 as uuidv4 } from 'uuid';

const appSchema = new mongoose.Schema({
    name: { type: String, required: true, unique: true },
    appId: { type: String, required: true, unique: true, default: uuidv4 },
    apiKey: { type: String, required: true, unique: true, default: () => uuidv4().replace(/-/g, '') },
    status: { type: String, enum: ['active', 'revoked'], default: 'active' },
    createdAt: { type: Date, default: Date.now }
});

const App = mongoose.model('App', appSchema);

export default App;
