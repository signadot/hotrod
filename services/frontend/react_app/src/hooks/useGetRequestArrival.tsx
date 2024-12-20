import {Log, LogEntry} from "./useLogs.tsx";
import {useEffect, useState} from "react";
import { faker } from "@faker-js/faker";

export const useGetRequestArrival = (logs: Log[]) => {
    const [driverArrival, setDriverArrival] = useState<number | undefined>();
    const [driverName, setDriverName] = useState("");

    const parseDriverLogService = (entry: LogEntry) => {
        if (entry.service !== 'driver') return;

        const timeRegex = /(\d+)m(\d+)s/;
        const match = entry.status.match(timeRegex);

        if (!match) return;

        const minutes = parseInt(match[1], 10);
        const seconds = parseInt(match[2], 10);

        return minutes * 60 + seconds;
    };

    useEffect(() => {
        if (logs.length === 0) return;

        const lastRequestDrive = logs[0];

        setDriverArrival(undefined);

        if (!lastRequestDrive) {
            setDriverArrival(undefined);
            return;
        }

        const driverEntries = lastRequestDrive.entries.filter((e) => e.service === 'driver');
        const parsedTime = driverEntries
            .map((e) => parseDriverLogService(e))
            .find((e) => e !== undefined);

        setDriverName(faker.person.fullName());
        setDriverArrival(parsedTime);
    }, [logs]);


    return {
        driverArrival,
        driverName,
    };
};
