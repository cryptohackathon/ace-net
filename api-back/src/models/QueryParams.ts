export interface QueryParams {
    /**
     * Sort order
     */
    sort?: 'ASC' | 'DESC';

    /**
     * Query limit
     */
    limit?: number;

    /**
     * Query offset
     */
    offset?: number;
}