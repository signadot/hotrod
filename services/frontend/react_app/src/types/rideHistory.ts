export type RideHistoryEntry = {
    sessionID: number,
    requestID: number,
    pickupLocation: string,
    dropoffLocation: string,
    requestedAt: string,
    driverPlate: string,
}

export type RideHistoryResponse = {
    totalCount: number,
    entries: RideHistoryEntry[],
}
