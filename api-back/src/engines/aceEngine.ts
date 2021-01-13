import { Singleton, Factory, Inject } from "typescript-ioc";
import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
import { PaginatedList } from "../models/PaginatedList";
import { PoolDataPayload } from "../models/PoolDataPayload";
import { PublicKeyShareRequest } from "../models/PublicKeyShareRequest";
import { QueryParams } from "../models/QueryParams";
import { RegistrationInfo } from "../models/RegistrationInfo";
import { PoolService } from "../services/PoolService";

@Singleton
@Factory(() => new ACEEngine())
export class ACEEngine {

    @Inject
    private poolService: PoolService;

    public async register(poolLabel?: string): Promise<RegistrationInfo> {
        return this.poolService.register(poolLabel)
        // return {
        //     clientSequenceId: 0,
        //     poolLabel: "today_1",
        //     registrationExpiry: "2021-01-20T11:15:46.163Z",
        //     status: 'REGISTRATION'
        // } as RegistrationInfo
    }

    public async status(poolLabel: string): Promise<PoolDataPayload> {
        return this.poolService.status(poolLabel)
        // return {
        //     poolLabel: "today_1",
        //     poolExpiry: "2021-01-20T11:15:46.163Z",
        //     status: 'PK_COLLECTION'
        // } as PoolDataPayload
    }

    public async postPublicKeyShare(payload: PublicKeyShareRequest): Promise<PoolDataPayload> {
        return this.poolService.postPublicKeyShare(payload)
        // return {
        //     poolLabel: "today_1",
        //     poolExpiry: "2021-01-20T11:15:46.163Z",
        //     status: 'ENCRYPTION'
        // } as PoolDataPayload
    }

    public async postCypherTextAndDecryptionKeysShares(payload: CypherAndDKRequest): Promise<PoolDataPayload> {
        return this.poolService.postCypherTextAndDecryptionKeysShares(payload)
        // return {
        //     status: 'FINALIZED',
        //     poolLabel: "today_1",
        //     poolExpiry: "2021-01-20T11:15:46.163Z",
        //     publicKeys: ["1", "2", "3"],
        //     cypherTexts: [["4", "5", "6"], ["4", "5", "6"]],
        //     decryptionKeys: ["7", "8", "9"]
        // } as PoolDataPayload
    }

    public async listPools(queryParams: QueryParams): Promise<PaginatedList<PoolDataPayload>> {
        const items = this.poolService.listPools()
        return {
            count: items.length,
            items,
            // limit: 100,
            offset: 0
        }
    }




}