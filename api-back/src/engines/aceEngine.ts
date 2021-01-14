import { Singleton, Factory, Inject } from "typescript-ioc";
import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
import { HistogramPayload } from "../models/HistogramPayload";
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
    }

    public async status(poolLabel: string): Promise<PoolDataPayload> {
        return this.poolService.status(poolLabel)
    }

    public async postPublicKeyShare(payload: PublicKeyShareRequest): Promise<PoolDataPayload> {
        return this.poolService.postPublicKeyShare(payload)
    }

    public async postCypherTextAndDecryptionKeysShares(payload: CypherAndDKRequest): Promise<PoolDataPayload> {
        return this.poolService.postCypherTextAndDecryptionKeysShares(payload)
    }

    public async listPools(status?: string): Promise<PoolDataPayload[]> {
        return this.poolService.listPools(status)
    }

    public async postHistogram(payload: HistogramPayload) {
        return this.poolService.postHistogram(payload)
    }

    public async reset() {
        return this.poolService.reset()
    }


}