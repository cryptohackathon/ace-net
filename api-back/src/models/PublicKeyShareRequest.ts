export interface PublicKeyShareRequest {
    /**
     * Sequential registration number of the client
     */
    clientSequenceId: number;

    /**
     * Label of the pool which is used in DMCFE encryption
     */
    poolLabel: string;

    /**
     * UTC timestamp denoting user registration expiry time. Should be provided exactly as obtained 
     * at registration! // TODO: currently it is abused as registration token
     */
    registrationExpiry: string;
    
    /**
     * Base64 encoded key share provided by the client
     */
    keyShare: string;
}