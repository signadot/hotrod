import {Log, LogEntry} from "./useLogs.tsx";
import {useEffect, useState} from "react";
import { faker } from "@faker-js/faker";

export const useGetRequestArrival = (logs: Log[]) => {
    const [driverArrival, setDriverArrival] = useState<number | undefined>();
    const [driverDetails, setDriverDetails] = useState({
        name: "",
        plate: "",
    });

    const parseDriverLogService = (entry: LogEntry) => {
        if (entry.service !== 'driver') return;

        const rawTimeMatch = entry.status.match(/(\d+m\d+s|\d+m|\d+s)/);

        if (!rawTimeMatch) return; // No substring that looks like time

        const rawTime = rawTimeMatch[0];

        /**
         * Possible cases
         * 3m
         * 3m2s
         * 2s
         * */
        const timeRegex = /^(?:(\d+)m)?(?:(\d+)s)?$/;
        const match = rawTime.match(timeRegex);

        if (!match) return;

        const minutes = parseInt(match[1] || '0', 10);
        const seconds = parseInt(match[2] || '0', 10);
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

        const plate = driverEntries.map(e => {
            return e.status.match(/Driver\s+(.*?)\s+arriving/)
        }).find(e => e !== null);

        setDriverDetails({
            name: faker.person.fullName(),
            plate: plate && plate.length > 1 ? plate[1] : "Unknown",
        });
        setDriverArrival(parsedTime);
    }, [logs]);


    return {
        driverArrival,
        driverDetails,
    };
};
