export interface RegistrationInfo {
    /**
     * Sequential registration number of the client
     */
    clientSequenceId: number;
    /**
     * Label of the pool which is used in DMCFE encryption
     */
    poolLabel: string;
    /**
     * UTC timestamp denoting user registration expiry time
     */
    registrationExpiry: string;

    /**
     * Labels for slots
     */
    slotLabels: string[];

    /**
     * Inner vector - vector 'y' that is used for inner product 
     */
    innerVector: number[];
    
    /**
     * Registration status in regard to relevant pool status.
     * - REGISTRATION - pool is accepting client registrations
     * - PK_COLLECTION - all users are registred, public key shares are being collected
     * - ENCRYPTION - public key is assembled, clients can send Cyphertexts and decryption key shares
     * - FINALIZED - all cyphertexts and decryption keys are collected
     * - EXPIRED - pool failed while PK_COLLECTION or ENCRYPTION phase before the the expiry time
     */
    status: 'REGISTRATION' | 'PK_COLLECTION' | 'ENCRYPTION' | 'FINALIZED' | 'EXPIRED'
}