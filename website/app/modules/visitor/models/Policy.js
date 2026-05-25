import { MongooseModel } from 'express-web-tools';
import { Schema } from 'mongoose';

const policySchema = new Schema({
    name: { type: String, required: true, unique: true }, // e.g., 'privacy-policy', 'terms-of-service'
    title: { type: String, required: true },
    content: { type: String, required: true }
});

const config = {
    softDelete: true,
    auditEnforce: false, // Set to true if we have admin auth context
    timestamps: true
};

const policyModelWrapper = new MongooseModel('Policy', policySchema, null, null, config);
export const Policy = policyModelWrapper.model();
