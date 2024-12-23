import {useState} from "react";
import {Location} from "../types/location.ts";


export type Log = {
    requestID: number,
    entries: LogEntry[],
    pickupLocation: Location,
    dropoffLocation: Location,
}

export type LogEntry = {
    messageType: 'error' | 'info',
    date: Date,
    service: string,
    status: string,
    sandboxName?: string,
}

export const useLogs = () => {
    const [logs, setLogs] = useState<Log[]>([]);

    const addNewLog = (pickupLocation: Location, dropoffLocation: Location, requestID: number, log: LogEntry) => {
        setLogs(prev => [{ pickupLocation, dropoffLocation, requestID, entries: [log]}, ...prev])
    }

    const addEntry = (requestID: number, body: Omit<LogEntry, "messageType">, messageType: LogEntry['messageType']): boolean => {
        let updated = false;
        setLogs(prevLogs => {
            return prevLogs.map(log => {
                if (log.requestID === requestID) {
                    updated = true;
                    return {
                        ...log,
                        entries: [...log.entries, {...body, messageType}],
                    };
                }
                return log;
            });
        });
        return updated;
    };

    const addErrorEntry = (requestID: number, body: Omit<LogEntry, "messageType">) => {
        return addEntry(requestID, body, 'error');
    };

    const addInformationEntry = (requestID: number, body: Omit<LogEntry, "messageType">) => {
        return addEntry(requestID, body, 'info');
    };


    return {
        addErrorEntry,
        addInformationEntry,
        addNewLog,
        logs
    }
}