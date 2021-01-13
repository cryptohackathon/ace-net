export interface CypherAndDKRequest {
    /**
     * Sequential registration number of the client
     */
    clientSequenceId: number;

    /**
     * Label of the pool which is used in DMCFE encryption
     */
    poolLabel: string;

    /**
     * Base64 encoded client's cyphertext
     */
    cypherText: string[];

    /**
     * Base64 encoded client's decryption key share
     */
    decryptionKeyShare: string;
}