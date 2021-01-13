
export interface HistogramPayload {
    /**
     * Secret known only to analytics server
     */
    secret: string;

    /**
     * Label of the pool which is used in DMCFE encryption
     */
    poolLabel: string;

    /**
     * Calculated histogram, submitted by certified analytics server
     */
    histogram?: number[]
}