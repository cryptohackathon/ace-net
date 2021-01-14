import { CypherAndDKRequest } from "../models/CypherAndDKRequest";
import { HistogramPayload } from "../models/HistogramPayload";
import { PoolDataPayload, PoolStatus } from "../models/PoolDataPayload";
import { PublicKeyShareRequest } from "../models/PublicKeyShareRequest";
import { RegistrationInfo } from "../models/RegistrationInfo";

export class ExposurePool {

    private _label: string;
    private _creationTime: Date;

    // all users have provided key shares
    private _registrationTime: Date;

    // all users have provided cyphertexts and decryption key shares
    private _finalizationTime: Date;

    // calculation time
    private _calculationTime: Date;

    private _status: PoolStatus;

    private _clientRegistrationExpiry: Date[]
    private _publicKeys: string[]
    private _cypherTexts: string[][]
    private _decryptionKeys: string[]
    private _histogram: number[]

    constructor(label?: string) {
        this._creationTime = new Date()
        this._label = label ? label : this._creationTime.toISOString()
        this._status = 'REGISTRATION'
        this._clientRegistrationExpiry = Array<Date>(this.size).fill(null)
        this._publicKeys = Array<string>(this.size).fill(null)
        this._cypherTexts = Array<string[]>(this.size).fill(null)
        this._decryptionKeys = Array<string>(this.size).fill(null)


    }

    get label(): string {
        return this._label
    }

    get size(): number {
        const tmp = parseInt(process.env.POOL_SIZE, 10)
        if (isNaN(tmp)) throw Error("Invalid POOL_SIZE");
        return tmp
    }

    get status() {
        return this._status
    }

    private addMiliseconds(time: Date, ms: string): Date {
        const milis = parseInt(ms, 10)
        if (isNaN(milis)) return time;
        const expiry = new Date(time)
        expiry.setMilliseconds(expiry.getMilliseconds() + milis)
        return expiry
    }

    get expires(): Date {
        switch (this._status) {
            case 'REGISTRATION':
                return this.addMiliseconds(this._creationTime, process.env.POOL_REGISTRATION_DURATION)
            case 'ENCRYPTION':
            case 'FINALIZED':
            case 'CALCULATED':
                return this.addMiliseconds(this._registrationTime, process.env.POOL_ENCRYPTION_DURATON)
            case 'EXPIRED': return this._creationTime
            default:
                throw Error("Unhandled case")
        }
    }

    get publicKeyAssebled() {
        return !this._publicKeys.some(x => x === null)
    }

    get cypherTextsAndDKsSubmitted() {
        return !this._cypherTexts.some(x => x === null)
    }

    get publicKeys(): string[] {
        if (['REGISTRATION', 'EXPIRED'].indexOf(this._status) >= 0) return null
        return this._publicKeys
    }

    get cypherTexts(): string[][] {
        if (['FINALIZED', 'CALCULATED'].indexOf(this._status) < 0) return null
        return this._cypherTexts
    }

    get decryptionKeys(): string[] {
        if (['FINALIZED', 'CALCULATED'].indexOf(this._status) < 0) return null
        return this._decryptionKeys
    }

    register() {
        const now = new Date()
        const exp = this.expires
        if (exp === null || this._status !== "REGISTRATION" || now > this.expires) {
            throw Error("Registrations closed")
        }
        console.log(this._clientRegistrationExpiry)
        const index = this._clientRegistrationExpiry.findIndex(x => x === null || x < now)
        console.log("N:", now, index, this._clientRegistrationExpiry[index], this._clientRegistrationExpiry[index] > now)
        if (index < 0) throw Error("Registrations temporary closed. Try later again.")
        console.log("INDEX:", index)
        this._clientRegistrationExpiry[index] = this.addMiliseconds(now, process.env.POOL_CLIENT_REGISTRATION_DURATION)
        this._publicKeys[index] = null

        return {
            clientSequenceId: index,
            poolLabel: this.label,
            registrationExpiry: this._clientRegistrationExpiry[index].toISOString(),
            status: this._status
        } as RegistrationInfo
    }

    submitPublicKey(payload: PublicKeyShareRequest): PoolDataPayload {
        if (this._status !== 'REGISTRATION') {
            throw Error("Registration period expired.")
        }
        if (this.label !== payload.poolLabel) {
            throw Error(`Wrong pool label. Provided: '${ payload.poolLabel }', should be '${ this.label }'.`)
        }
        const id = payload.clientSequenceId;
        if (id < 0 || id >= this.size) {
            throw Error(`Wrong client id. Provided ${ payload.clientSequenceId }, pool size: ${ this.size }.`)
        }
        const expiry = this._clientRegistrationExpiry[id]
        const now = new Date()
        if (now > expiry || payload.registrationExpiry !== expiry.toISOString()) {
            throw Error("Client registration expired.")
        }
        this._publicKeys[id] = payload.keyShare;
        console.log("ADDING PUBLIC KEY:", payload)
        if (!this.publicKeyAssebled) return this.info
        // public key is assembled, change status
        this._status = 'ENCRYPTION'
        this._registrationTime = now
        return this.info
    }

    postCypherTextAndDecryptionKeysShares(payload: CypherAndDKRequest): PoolDataPayload {
        if (this._status !== 'ENCRYPTION') {
            throw Error("Server not accepting cyphers.")
        }
        if (this.label !== payload.poolLabel) {
            throw Error(`Wrong pool label. Provided: '${ payload.poolLabel }', should be '${ this.label }'.`)
        }
        const id = payload.clientSequenceId;
        if (id < 0 || id >= this.size) {
            throw Error(`Wrong client id. Provided ${ payload.clientSequenceId }, pool size: ${ this.size }.`)
        }

        const now = new Date()
        if (now > this.expires) {
            this._status = 'EXPIRED'
            throw Error("Pool has expired.")
        }
        this._cypherTexts[id] = payload.cypherText
        this._decryptionKeys[id] = payload.decryptionKeyShare
        if (!this.cypherTextsAndDKsSubmitted) return this.info
        this._status = 'FINALIZED'
        this._finalizationTime = now
        return this.info
    }

    postHistogram(payload: HistogramPayload) {
        if (process.env.ANALYTICS_SECRET !== payload.secret) {
            throw Error('Authorization error')
        }
        if (this._status !== 'FINALIZED') {
            throw Error("Pool is not finalized.")
        }
        if (this.label !== payload.poolLabel) {
            throw Error(`Wrong pool label. Provided: '${ payload.poolLabel }', should be '${ this.label }'.`)
        }
        this._histogram = payload.histogram
        this._status = 'CALCULATED'
        this._calculationTime = new Date()
        return this.info
    }

    get info(): PoolDataPayload {
        const exp = this.expires
        return {
            status: this._status,
            poolLabel: this.label,
            poolExpiry: exp ? exp.toISOString() : null,
            creationTime: this._creationTime ? this._creationTime.toISOString() : null,
            registrationTime: this._registrationTime ? this._registrationTime.toISOString() : null,
            finalizationTime: this._finalizationTime ? this._finalizationTime.toISOString() : null,
            calculationTime: this._calculationTime ? this._calculationTime.toISOString() : null,
            publicKeys: this.publicKeys,
            cypherTexts: this.cypherTexts,
            decryptionKeys: this.decryptionKeys,
            histogram: this._histogram
        }
    }


}