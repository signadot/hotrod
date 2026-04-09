export type RideHistoryEntry = {
    requestID: number,
    pickupLocation: string,
    dropoffLocation: string,
    requestedAt: string,
    driverPlate: string,
}

export type RideHistoryResponse = {
    total: number,
    rides: RideHistoryEntry[],
}
