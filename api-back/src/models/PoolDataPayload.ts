
/**
 * Registration status in regard to relevant pool status.
 * - REGISTRATION - pool is accepting client registrations and after them uploads of public key shares
 * - ENCRYPTION - public key is assembled, clients can send Cyphertexts and decryption key shares
 * - FINALIZED - all cyphertexts and decryption keys are collected
 * - CALCULATED - histograms are calculated
 * - EXPIRED - pool failed while PK_COLLECTION or ENCRYPTION phase before the the expiry time
 */
export type PoolStatus = 'REGISTRATION' | 'ENCRYPTION' | 'FINALIZED' | 'CALCULATED' | 'EXPIRED';

export interface PoolDataPayload {
    /**
     * Status
     */
    status: PoolStatus;

    /**
     * Creation time
     */
    creationTime?: string;

    /**
     * End of registration time
     */
    registrationTime?: string;

    /**
     * End of finalization time
     */
    finalizationTime?: string;

    /**
     * Time of calculation of histograms
     */
    calculationTime?: string;

    /**
     * Label of the pool which is used in DMCFE encryption
     */
    poolLabel: string;
    /**
     * UTC timestamp denoting pool expiry time
     */
    poolExpiry: string;

    /**
     * In case of status 'ENCRYPTION' or 'FINALIZED' contains Base64 encoded list of public key shares of each client
     */
    publicKeys?: string[];

    /**
     * In case of status 'FINALIZED' contains Base64 encoded list of user cyphertexts
     */
    cypherTexts?: string[][];

    /**
     * In case of status 'FINALIZED' contains Base64 encoded list of user decryption key shares together forming the
     * decryption key.
     */
    decryptionKeys?: string[][];

    /**
     * Labels for slots
     */
    slotLabels: string[];

    /**
     * Inner vector - vector 'y' that is used for inner product 
     */
    innerVector: number[];

    /**
     * Calculated histogram, submitted by certified analytics server
     */
    histogram?: number[]
}