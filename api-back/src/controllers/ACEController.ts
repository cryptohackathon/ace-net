// import express from "express";
// import { Hello } from "src/models/Hello";
// import { Hello } from "src/models/Hello";

// import { EnergyEngine } from "src/engines/energy";
import { Body, Controller, Get, Path, Post, Query, Route, Tags } from "tsoa";
import { Inject } from "typescript-ioc";
import { ACEEngine } from "../engines/aceEngine";
import { ApiResponse, handleApiResponse } from "../models/ApiResponse";
import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
import { HistogramPayload } from "../models/HistogramPayload";
import { PaginatedList } from "../models/PaginatedList";
import { PoolDataPayload } from "../models/PoolDataPayload";
import { PublicKeyShareRequest } from "../models/PublicKeyShareRequest";
import { RegistrationInfo } from "../models/RegistrationInfo";

@Tags('ACE')
@Route("ace")
export class ACEController extends Controller {

    @Inject
    private aceEngine: ACEEngine;

    /**
     * Registers a client to the pool
     * @param date date of registration in the form of `YYYY-MM-dd`
     */
    @Post("register")
    public async register(
        @Query() label?: string,
    ): Promise<ApiResponse<RegistrationInfo>> {
        return handleApiResponse(
            this.aceEngine.register(label)
        )
    }

    /**
     * Returns status info about the pool specific to the client
     * @param clientSequenceId - clients id
     * @param poolLabel - pool
     */
    @Get("status/{poolLabel}")
    public async status(
        @Path() poolLabel: string,
    ): Promise<ApiResponse<PoolDataPayload>> {
        return handleApiResponse(
            this.aceEngine.status(poolLabel)
        )
    }

    /**
     * Posts specific user's public key share to form common master key
     * @param requestBody
     */
    @Post("public-key-share")
    public async postPublicKeyShare(
        @Body() requestBody: PublicKeyShareRequest): Promise<ApiResponse<PoolDataPayload>> {
        return handleApiResponse(
            this.aceEngine.postPublicKeyShare(requestBody)
        )
    }

    @Post("cyphertext-and-dk")
    public async postCypherTextAndDecryptionKeysShares(
        @Body() requestBody: CypherAndDKRequest): Promise<ApiResponse<PoolDataPayload>> {
        return handleApiResponse(
            this.aceEngine.postCypherTextAndDecryptionKeysShares(requestBody)
        )
    }

    @Get("pools")
    public async getPools(
        @Query() status?: string
    ): Promise<ApiResponse<PoolDataPayload[]>> {
        return handleApiResponse(
            this.aceEngine.listPools(status)
        )
    }

    @Post("histogram")
    public async postHistogram(
        @Body() requestBody: HistogramPayload): Promise<ApiResponse<PoolDataPayload>> {
        return handleApiResponse(
            this.aceEngine.postHistogram(requestBody)
        )
    }

    @Post("reset")
    public async reset(): Promise<ApiResponse<any>> {
        return handleApiResponse(
            this.aceEngine.reset()
        )
    }

}