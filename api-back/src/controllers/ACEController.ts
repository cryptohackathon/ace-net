// import express from "express";
// import { Hello } from "src/models/Hello";
// import { Hello } from "src/models/Hello";

// import { EnergyEngine } from "src/engines/energy";
import { Body, Controller, Get, Path, Post, Query, Route, Tags } from "tsoa";
import { Inject } from "typescript-ioc";
import { ACEEngine } from "../engines/aceEngine";
import { ApiResponse, handleApiResponse } from "../models/ApiResponse";
import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
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
    @Get("register/{date}")
    public async register(
        @Path() date: string,
    ): Promise<ApiResponse<RegistrationInfo>> {
        return handleApiResponse(
            this.aceEngine.register()
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

    @Post("cyphertext-and-dk/upload")
    public async postCypherTextAndDecryptionKeysShares(
        @Body() requestBody: CypherAndDKRequest): Promise<ApiResponse<PoolDataPayload>> {
        return handleApiResponse(
            this.aceEngine.postCypherTextAndDecryptionKeysShares(requestBody)
        )
    }

    @Get("pools/list/{date}")
    public async getPools(
        @Path() date: string,
        @Query() sort?: 'ASC' | 'DESC',
        @Query() limit?: number,
        @Query() offset?: number,
    ): Promise<ApiResponse<PaginatedList<PoolDataPayload>>> {
        return handleApiResponse(
            this.aceEngine.listPools(date, {sort, limit, offset})
        )
    }





}