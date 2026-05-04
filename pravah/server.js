import cookieParser from "cookie-parser";
import cors from "cors";
import express from "express";
import expressPureIP from "express-pureip";

import mongoose from "mongoose";
import { v4 as uuid } from "uuid";

import { checkEnvironmentVariables } from "./utils/index.js";
import handshakeRouter from "./core/handshake.js";
import subscribeRouter from "./modules/subscribe/route.js";
import publishRouter from "./modules/publish/route.js";
import healthRouter from "./modules/health/route.js";
import appRouter from "./modules/app/route.js";

(async () => {
    await checkEnvironmentVariables();
    await mongoose.connect(process.env.MONGODB)
    const app = express();
    app
        .set("view-engine", "twig")
        .use(express.static("public"))
        .use([
            cors({
                origin: true, 
                credentials: true
            }),
            expressPureIP,
            express.json(),
            express.urlencoded({
                extended: true
            }),
            cookieParser()
        ])
        .use("/pravah", handshakeRouter)
        .use("/pravah", subscribeRouter)
        .use("/pravah", publishRouter)
        .use("/pravah", healthRouter)
        .use("/api/v1", appRouter)
        // .use() //main router

        .use((req, res) => {
            res.status(404).json({
                error: "Not Found",
                message: `Cannot ${req.method} ${req.originalUrl}`
            });
        })
        .use((err, req, res, next) => {
            try {
                if (res.headersSent) {
                    return next(err);
                }

                const status = err.status || err.statusCode || 500;
                res.status(status).json({
                    error: err.name || "Internal Server Error",
                    message: err.message || "Internal Server Error"
                });
            } catch (error) {
                console.log(error)
            }
        })
        .listen(process.env.PORT ?? 3000, (error) => {
            if (error) console.log(error.message)
            else
                console.log(`Server started on port ${process.env.PORT ?? 3000}`)
        })
})();
