import { Factory, Singleton } from "typescript-ioc";
import { ExposurePool } from "../classes/ExposurePool";
import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
import { PoolDataPayload } from "../models/PoolDataPayload";
import { PublicKeyShareRequest } from "../models/PublicKeyShareRequest";
import { QueryParams } from "../models/QueryParams";
import { RegistrationInfo } from "../models/RegistrationInfo";

@Singleton
@Factory(() => new PoolService())
export class PoolService {
    poolMap = new Map<string, ExposurePool>()

    register(poolLabel?: string): RegistrationInfo {
        let pool: ExposurePool = null;
        console.log("LAB:", poolLabel)
        if(!poolLabel) {   // register to any pool
            const sortedKeys = [...this.poolMap.keys()].sort()
            console.log("SK:", sortedKeys)
            for(const label of sortedKeys) {  // find first pool where you can register
                const tmpPool = this.poolMap.get(label)
                if(tmpPool.status !== 'REGISTRATION') continue
                try {
                    console.log("Trying")
                    const res = tmpPool.register()
                    if(res) return res
                } catch(e) {
                    console.log("ERR:", e)
                    continue
                }
            }
            if(pool === null) {  // need a new pool
                console.log("CREATING NEW POOL")
                pool = new ExposurePool()
                this.poolMap.set(pool.label, pool)
                return pool.register()
            }
        } else {   // label provided, must exist, try to register to specific pool
            pool = this.poolMap.get(poolLabel)
            if(!pool) {
                throw Error(`Pool '${poolLabel}' does not exist.`)
            }
            return pool.register()
        }
    }

    status(poolLabel: string): PoolDataPayload {
        const pool = this.poolMap.get(poolLabel)
        if(pool === null) {
            throw Error(`Wrong pool label '${poolLabel}'.`)
        }
        return pool.info
    }

    postPublicKeyShare(payload: PublicKeyShareRequest): PoolDataPayload {
        const poolLabel = payload.poolLabel
        const pool = this.poolMap.get(poolLabel)
        if(pool === null) {
            throw Error(`Wrong pool label '${poolLabel}'.`)
        }
        return pool.submitPublicKey(payload)
    }

    postCypherTextAndDecryptionKeysShares(payload: CypherAndDKRequest): PoolDataPayload {
        const poolLabel = payload.poolLabel
        const pool = this.poolMap.get(poolLabel)
        if(pool === null) {
            throw Error(`Wrong pool label '${poolLabel}'.`)
        }
        return pool.postCypherTextAndDecryptionKeysShares(payload)
    }

    // TODO: now it returns all, fix for queryParameters
    listPools(queryParams?: QueryParams): PoolDataPayload[] {
        return [...this.poolMap.values()].map(x => x.info)
    }

}