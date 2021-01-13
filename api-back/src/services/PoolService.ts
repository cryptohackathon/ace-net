import { Factory, Singleton } from "typescript-ioc";
import { ExposurePool } from "../classes/ExposurePool";
import { PoolDataPayload } from "../models/PoolDataPayload";
import { PublicKeyShareRequest } from "../models/PublicKeyShareRequest";
import { RegistrationInfo } from "../models/RegistrationInfo";

@Singleton
@Factory(() => new PoolService())
export class PoolService {
    poolMap = new Map<string, ExposurePool>()

    register(poolLabel?: string): RegistrationInfo {   
        let pool: ExposurePool = null;     
        if(!poolLabel) {   // register to any pool
            let sortedKeys = [...this.poolMap.keys()].sort()
            for(let label of sortedKeys) {
                let tmpPool = this.poolMap.get(label)
                if(tmpPool.status != 'REGISTRATION') continue
                pool = tmpPool
            }
            if(pool === null) {
                pool = new ExposurePool()
                this.poolMap.set(pool.label, pool)
            }
        } else {   // label provided, must exist
            let pool = this.poolMap.get(poolLabel)
            if(!pool) {
                throw Error(`Pool '${poolLabel}' does not exist.`)
            }            
        }
        return pool.register()
    }

    status(poolLabel: string): PoolDataPayload {
        let pool = this.poolMap.get(poolLabel)
        if(pool === null) {
            throw Error(`Wrong pool label '${poolLabel}'.`)
        }
        return pool.info
    }

    postPublicKeyShare(payload: PublicKeyShareRequest) {
        let poolLabel = payload.poolLabel
        let pool = this.poolMap.get(poolLabel)
        if(pool === null) {
            throw Error(`Wrong pool label '${poolLabel}'.`)
        }
        return pool.submitPublicKey(payload)
    }

}