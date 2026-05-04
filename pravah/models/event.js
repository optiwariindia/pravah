import mongoose from 'mongoose';

const eventSchema = new mongoose.Schema({
    streamId: { type: String, required: true, index: true },
    event: { type: String, required: true },
    data: { 
        id: { type: String, required: true },
        metadata: { type: Object, default: {} }
    },
    createdAt: { type: Date, default: Date.now, expires: 86400 } // 24 hours TTL
});

// Compound index for replay efficiency
eventSchema.index({ streamId: 1, _id: 1 });

const Event = mongoose.model('Event', eventSchema);

export default Event;
