// import express from "express";
// import { Hello } from "src/models/Hello";
// import { Hello } from "src/models/Hello";

// import { EnergyEngine } from "src/engines/energy";
import { Controller, Get, Route, Tags } from "tsoa";
import { Inject } from "typescript-ioc";
import { EnergyEngine } from "../engines/energy";
import { ApiResponse, handleApiResponse } from "../models/ApiResponse";
import { Hello } from "../models/Hello";

@Tags('Demo')
@Route("demo")
export class OrganizationController extends Controller {

    @Inject
    private energyEngine: EnergyEngine;

    @Get("test")
    public async test(
    ): Promise<ApiResponse<Hello>> {
        return handleApiResponse(
            this.energyEngine.test()
        )
    }
}