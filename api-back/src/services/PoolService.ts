import { Factory, Singleton } from "typescript-ioc";
import { ExposurePool } from "../classes/ExposurePool";
import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
import { HistogramPayload } from "../models/HistogramPayload";
import { PoolDataPayload } from "../models/PoolDataPayload";
import { PublicKeyShareRequest } from "../models/PublicKeyShareRequest";
import { RegistrationInfo } from "../models/RegistrationInfo";
import { wsServer } from "../server";

@Singleton
@Factory(() => new PoolService())
export class PoolService {

    constructor() {
        setInterval(() => {
            console.log("WORKING")
            wsServer.clients.forEach(client => {
                client.send(JSON.stringify({
                    message: 'POOLS_SUMMARY',
                    data: [...this.poolMap.values()].map(pool => pool.summary)
                }));
            });
        }, 1000)
    }

    poolMap = new Map<string, ExposurePool>()

    private _innerVector: number[] = Array<number>(this.size).fill(1)  // all ones vector
    private _slotLabels: string[] = ["L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9", "L10"]

    get size(): number {
        const tmp = parseInt(process.env.POOL_SIZE, 10)
        if (isNaN(tmp)) throw Error("Invalid POOL_SIZE");
        return tmp
    }

    wsServer = wsServer

    register(poolLabel?: string): RegistrationInfo {
        let pool: ExposurePool = null;
        if (!poolLabel) {   // register to any pool
            const sortedKeys = [...this.poolMap.keys()].sort()
            for (const label of sortedKeys) {  // find first pool where you can register
                const tmpPool = this.poolMap.get(label)
                if (tmpPool.status !== 'REGISTRATION') continue
                try {
                    const res = tmpPool.register()
                    if (res) {
                        // wsServer.clients.forEach(client => {
                        //     client.send(JSON.stringify({
                        //         message: `REGISTRATION: -> pool: ${ res.poolLabel }, id: ${ res.clientSequenceId }`
                        //     }));
                        // });
                        return res
                    }
                } catch (e) {
                    console.log("ERR:", e)
                    continue
                }
            }
            if (pool === null) {  // need a new pool
                console.log("CREATING NEW POOL")
                pool = new ExposurePool(this._slotLabels, this._innerVector)
                this.poolMap.set(pool.label, pool)
                const res = pool.register()
                // wsServer.clients.forEach(client => {
                //     client.send(JSON.stringify({
                //         message: `REGISTRATION: -> pool: ${ res.poolLabel }, id: ${ res.clientSequenceId }`
                //     }));
                // });
                return res
            }
        } else {   // label provided, must exist, try to register to specific pool
            pool = this.poolMap.get(poolLabel)
            if (!pool) {
                throw Error(`Pool '${ poolLabel }' does not exist.`)
            }
            const res = pool.register()
            // wsServer.clients.forEach(client => {
            //     client.send(JSON.stringify({
            //         message: `REGISTRATION: -> pool: ${ res.poolLabel }, id: ${ res.clientSequenceId }`
            //     }));
            // });
            return res
        }
    }

    status(poolLabel: string): PoolDataPayload {
        const pool = this.poolMap.get(poolLabel)
        if (pool === null) {
            throw Error(`Wrong pool label '${ poolLabel }'.`)
        }
        return pool.info
    }

    postPublicKeyShare(payload: PublicKeyShareRequest): PoolDataPayload {
        const poolLabel = payload.poolLabel
        const pool = this.poolMap.get(poolLabel)
        if (pool === null) {
            throw Error(`Wrong pool label '${ poolLabel }'.`)
        }
        const res = pool.submitPublicKey(payload)
        // wsServer.clients.forEach(client => {
        //     client.send(JSON.stringify({
        //         message: `SHARE SENT: -> pool: ${ res.poolLabel }, id: ${ payload.clientSequenceId }`
        //     }));
        //     if (res.status === 'ENCRYPTION') {
        //         client.send(JSON.stringify({
        //             message: `ENCRYPTION: -> pool: ${ res.poolLabel }`
        //         }));
        //     }

        // });
        return res
    }

    postCypherTextAndDecryptionKeysShares(payload: CypherAndDKRequest): PoolDataPayload {
        const poolLabel = payload.poolLabel
        const pool = this.poolMap.get(poolLabel)
        if (pool === null) {
            throw Error(`Wrong pool label '${ poolLabel }'.`)
        }
        const res = pool.postCypherTextAndDecryptionKeysShares(payload)
        // wsServer.clients.forEach(client => {
        //     client.send(JSON.stringify({
        //         message: `CYPHER AND DK SENT: -> pool: ${ res.poolLabel }, id: ${ payload.clientSequenceId }`
        //     }));
        //     if (res.status === 'FINALIZED') {
        //         client.send(JSON.stringify({
        //             message: `FINALIZED: -> pool: ${ res.poolLabel }`
        //         }));
        //     }
        // });
        return res
    }

    listPools(status?: string): PoolDataPayload[] {
        const res = [...this.poolMap.values()].map(x => x.info)
        if (status) {
            return res.filter(x => x.status === status)
        }
        return res
    }

    postHistogram(payload: HistogramPayload) {
        if (process.env.ANALYTICS_SECRET !== payload.secret) {
            throw Error('Authorization error')
        }
        const poolLabel = payload.poolLabel
        const pool = this.poolMap.get(poolLabel)
        if (pool === null) {
            throw Error(`Wrong pool label '${ poolLabel }'.`)
        }
        return pool.postHistogram(payload)
    }

    reset() {
        this.poolMap.clear()
    }

}