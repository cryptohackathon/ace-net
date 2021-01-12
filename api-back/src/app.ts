import bodyParser from "body-parser";
import dotenv from "dotenv";
import express, { NextFunction, Request as ExRequest, Response as ExResponse } from "express";
import swaggerUi from 'swagger-ui-express';
import { RegisterRoutes } from "../build/routes";
import swaggerDocument from "../static/swagger.json";
import { ValidateError } from "tsoa";
import { ApiResponse } from "./models/ApiResponse";
import compression from "compression";
import helmet from "helmet";
import cookieParser from "cookie-parser";

// initialize configuration
dotenv.config();

export const app = express();

app.use(helmet());
app.use(compression()); // Compress all routes

app.use(cookieParser());
app.use(bodyParser.json({ limit: '50mb' }));
// Use body parser to read sent json payloads
app.use(
    bodyParser.urlencoded({
        limit: '50mb',
        extended: true,
        parameterLimit: 50000
    })
);

app.use('/api-doc', swaggerUi.serve, swaggerUi.setup(swaggerDocument));
app.use('/static', express.static('static'));


RegisterRoutes(app);

app.use(function notFoundHandler(_req, res: ExResponse) {
    res.status(404).send({
        message: "Not Found",
    });
});

app.use(function errorHandler(
    err: unknown,
    req: ExRequest,
    res: ExResponse,
    next: NextFunction
): ExResponse | void {
    if (err instanceof ApiResponse) {
        return res.status(400).json(err);
    }
    if (err instanceof ValidateError) {
        console.warn(`Caught Validation Error for ${ req.path }:`, err.fields);
        return res.status(422).json(
            new ApiResponse<any>(undefined, 'VALIDATION_ERROR', undefined, err ? err.fields : null)
        );
    }
    if (err instanceof Error) {
        if((err as any).status) {
            return res.status((err as any).status).json(
                new ApiResponse<any>(undefined, 'ERROR', err.message)
            );
        }
        return res.status(500).json(
            new ApiResponse<any>(undefined, 'ERROR', err.message)
        );
    }
    next();
});
