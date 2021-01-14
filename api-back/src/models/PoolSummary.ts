import { PoolStatus } from "./PoolDataPayload";


export interface PoolSummary {
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
     * Size of a pool
     */
    poolSize: number;

    /**
     * Number of public key shares
     */
    noPublicKeys: number;

    /**
     * Nubmer of cyphertexts
     */
    noCyphertexts: number;

    /**
     * Labels for slots
     */
    slotLabels: string[];
    
    /**
     * Calculated histogram, submitted by certified analytics server
     */
    histogram?: number[]
}