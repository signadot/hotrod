import { MainLayout } from "../components/layouts";
import {
    Badge,
    Box,
    Button,
    Flex,
    Heading,
    HStack,
    Stack,
    Text,
    useDisclosure,
} from "@chakra-ui/react";
import { motion } from "framer-motion";
import styles from "./home.module.css";

import { Logs } from "../components/features/logs/logs.tsx";
import { Map } from "../components/features/map/map.tsx";
import { RequestFlowDiagram } from "../components/features/requestFlow/RequestFlowDiagram.tsx";
import { useSession } from "../context/sessionContext/context.tsx";
import { useEffect, useRef, useState } from "react";
import { apiGet, apiPost } from "../services/http.ts";
import { Locations } from "../types/location.ts";
import { LocationSelect } from "../components/common/locationSelect/locationSelect.tsx";
import { useLogs } from "../hooks/useLogs.tsx";
import { NotificationResponse } from "../types/notifications.ts";
import { useGetRequestArrival } from "../hooks/useGetRequestArrival.tsx";
import Countdown, { CountdownRenderProps } from "react-countdown";

const countdownRenderer = ({ minutes, seconds, props }: CountdownRenderProps) => {
    if (minutes === 0 && seconds === 0) return <span style={{ color: '#48BB78' }}>Arrived</span>;
    return <span>{props.overtime ? "-" : ""}{minutes.toString().padStart(2, "0")}:{seconds.toString().padStart(2, "0")}</span>;
};

// Gradient SVG icons
const ClockIcon = () => (
    <svg width="36" height="36" viewBox="0 0 24 24" fill="none" strokeLinecap="round" strokeLinejoin="round">
        <defs><linearGradient id="gc" x1="0" y1="0" x2="1" y2="1"><stop offset="0%" stopColor="#00B5D8"/><stop offset="100%" stopColor="#76E4F7"/></linearGradient></defs>
        <circle cx="12" cy="12" r="10" stroke="url(#gc)" strokeWidth="2"/><polyline points="12 6 12 12 16 14" stroke="url(#gc)" strokeWidth="2"/>
    </svg>
);
const UserIcon = () => (
    <svg width="36" height="36" viewBox="0 0 24 24" fill="none" strokeLinecap="round" strokeLinejoin="round">
        <defs><linearGradient id="gu" x1="0" y1="0" x2="1" y2="1"><stop offset="0%" stopColor="#9F7AEA"/><stop offset="100%" stopColor="#D6BCFA"/></linearGradient></defs>
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" stroke="url(#gu)" strokeWidth="2"/><circle cx="12" cy="7" r="4" stroke="url(#gu)" strokeWidth="2"/>
    </svg>
);
const CarIcon = () => (
    <svg width="36" height="36" viewBox="0 0 24 24" fill="none" strokeLinecap="round" strokeLinejoin="round">
        <defs><linearGradient id="gcar" x1="0" y1="0" x2="1" y2="1"><stop offset="0%" stopColor="#F6AD55"/><stop offset="100%" stopColor="#ECC94B"/></linearGradient></defs>
        <path d="M5 17h14M5 17a2 2 0 01-2-2V9a2 2 0 012-2h1l2-3h8l2 3h1a2 2 0 012 2v6a2 2 0 01-2 2M5 17a2 2 0 100 4 2 2 0 000-4zm14 0a2 2 0 100 4 2 2 0 000-4z" stroke="url(#gcar)" strokeWidth="1.8"/>
    </svg>
);

// Map notification body text to user-friendly status messages
function getStatusMessage(body: string): string {
    if (/Processing dispatch/i.test(body)) return 'Dispatching your request...';
    if (/Resolving locations/i.test(body)) return 'Resolving pickup & dropoff locations...';
    if (/Finding an available driver/i.test(body)) return 'Searching for nearby drivers...';
    if (/Resolving routes/i.test(body)) return 'Calculating fastest route...';
    if (/Driver .* arriving/i.test(body)) return 'Driver found!';
    return body;
}

export const HomePage = () => {
    const session = useSession();
    const notificationCursorRef = useRef(-1);
    const [locations, setLocations] = useState<Locations | undefined>();
    const [selectedLocations, setSelectedLocations] = useState({ pickupId: -1, dropoffId: -1 });
    const { logs, addNewLog, addErrorEntry, addInformationEntry } = useLogs();
    const lastRequestedDrive = useGetRequestArrival(logs);

    const [showArchitecture, setShowArchitecture] = useState(false);
    const [isRequesting, setIsRequesting] = useState(false);
    const [statusMessage, setStatusMessage] = useState('');
    const [, setStatusIdx] = useState(0);
    const requestStartRef = useRef<number>(0);
    const etaTargetRef = useRef<Date | null>(null);
    const lastEtaArrivalRef = useRef<number | undefined>(undefined);
    const logsPanel = useDisclosure();

    // Cycle through pre-set status messages for a realistic feel
    const STATUS_MESSAGES = [
        'Dispatching your request...',
        'Resolving pickup & dropoff...',
        'Searching for nearby drivers...',
        'Calculating fastest route...',
        'Matching with best driver...',
    ];

    useEffect(() => {
        if (!isRequesting) return;
        const interval = setInterval(() => {
            setStatusIdx(prev => {
                const next = prev + 1;
                if (next < STATUS_MESSAGES.length) {
                    setStatusMessage(STATUS_MESSAGES[next]);
                    return next;
                }
                return prev; // stay on last message
            });
        }, 1800);
        return () => clearInterval(interval);
    }, [isRequesting]);

    useEffect(() => { apiGet<Locations>('/splash').then(setLocations); }, []);

    useEffect(() => {
        const poll = async () => {
            const cursor = notificationCursorRef.current;
            const url = `/notifications?sessionID=${session.sessionID}&cursor=${cursor}&nonse=${Math.random()}`;
            const res = await apiGet<NotificationResponse>(url);
            notificationCursorRef.current = res.cursor;
            res.notifications.forEach(n => {
                addInformationEntry(n.context.request.id, {
                    date: new Date(n.timestamp),
                    service: n.context.baselineWorkload,
                    status: n.body,
                    sandboxName: n.context.sandboxName,
                });
                // Update the live status message
                if (isRequesting) {
                    setStatusMessage(getStatusMessage(n.body));
                }
            });
            if (res.notifications.some(n => /Driver .* arriving/i.test(n.body))) {
                // Minimum 4 seconds of spinner for realistic feel
                const elapsed = Date.now() - requestStartRef.current;
                const remaining = Math.max(0, 4000 - elapsed);
                setTimeout(() => {
                    setIsRequesting(false);
                    setStatusMessage('');
                    setStatusIdx(0);
                }, remaining);
            }
        };
        const id = setInterval(poll, 2000);
        return () => clearInterval(id);
    }, [isRequesting]);

    if (!locations) {
        return (
            <MainLayout>
                <Flex h="100%" alignItems="center" justifyContent="center">
                    <Heading size="md" color="whiteAlpha.600">Loading...</Heading>
                </Flex>
            </MainLayout>
        );
    }

    const handleRequestDrive = async () => {
        const { getLastRequestID, sessionID } = session;
        const requestID = getLastRequestID();
        const { pickupId, dropoffId } = selectedLocations;
        const pickupLocation = locations.Locations.find(l => l.id === pickupId);
        const dropoffLocation = locations.Locations.find(l => l.id === dropoffId);

        setIsRequesting(true);
        setStatusIdx(0);
        setStatusMessage(STATUS_MESSAGES[0]);
        requestStartRef.current = Date.now();

        addNewLog(pickupLocation!, dropoffLocation!, requestID, {
            messageType: "info", service: 'browser', date: new Date(), status: 'Requesting a ride.'
        });

        try {
            await apiPost<{}>('/dispatch', {
                sessionID, requestID,
                pickupLocationID: pickupId,
                dropoffLocationID: dropoffId,
            }, { "baggage": `session=${sessionID}, request=${requestID}` });
        } catch (e) {
            setIsRequesting(false);
            setStatusMessage('');
            addErrorEntry(requestID, {
                status: 'Error requesting a ride to frontend API',
                service: 'api', date: new Date(),
            });
        }
    };

    const hasResult = !!lastRequestedDrive.driverArrival;

    // Compute stable ETA target date — only recalculate when driverArrival changes
    if (lastRequestedDrive.driverArrival !== lastEtaArrivalRef.current) {
        lastEtaArrivalRef.current = lastRequestedDrive.driverArrival;
        if (lastRequestedDrive.driverArrival) {
            etaTargetRef.current = new Date(Date.now() + lastRequestedDrive.driverArrival * 1000);
        } else {
            etaTargetRef.current = null;
        }
    }
    const etaTarget = etaTargetRef.current;

    return (
        <MainLayout>
            <div className={styles.pageGrid}>
                {/* Left Panel */}
                <div className={styles.leftPanel}>
                    <Stack spacing={5}>
                        <Heading size="lg" color="whiteAlpha.900" textAlign="left" fontWeight={700}>Request a Ride</Heading>
                        <LocationSelect placeholder='Pickup location' locations={locations.Locations}
                            onSelect={id => setSelectedLocations(prev => ({ ...prev, pickupId: id }))}
                            selectedLocationID={selectedLocations.pickupId} />
                        <LocationSelect placeholder='Dropoff location' locations={locations.Locations}
                            onSelect={id => setSelectedLocations(prev => ({ ...prev, dropoffId: id }))}
                            selectedLocationID={selectedLocations.dropoffId} />
                        <Button variant='solid' colorScheme='cyan' size="lg" w="100%" fontSize="md" fontWeight={600} h="52px"
                            onClick={handleRequestDrive}
                            isLoading={isRequesting}
                            loadingText="Finding driver..."
                            isDisabled={
                                selectedLocations.pickupId === -1 || isNaN(selectedLocations.pickupId) ||
                                selectedLocations.dropoffId === -1 || isNaN(selectedLocations.dropoffId)
                            }>
                            Request Ride
                        </Button>
                    </Stack>

                    <Button variant='outline' size="md" colorScheme="gray" w="100%" fontSize="sm" fontWeight={600} h="44px"
                        onClick={() => setShowArchitecture(prev => !prev)}>
                        {showArchitecture ? 'Hide Architecture' : 'Show Architecture'}
                    </Button>

                    <Box borderRadius="md" overflow="hidden" h="260px" flexShrink={0} border="1px solid" borderColor="gray.700">
                        <Map dropoffLocationID={selectedLocations.dropoffId} pickupLocationID={selectedLocations.pickupId} />
                    </Box>
                </div>

                {/* Center Panel */}
                <div className={styles.centerPanel}>
                    {/* Result Card */}
                    <Box bg="gray.800" borderRadius="xl" border="1px solid" borderColor="gray.700" py={6} px={10} mb={4} flexShrink={0}>
                        {isRequesting && !hasResult ? (
                            <Flex direction="column" alignItems="center" gap={5} py={4}>
                                {/* Pulsing ring spinner — larger for demo */}
                                <Box position="relative" w="72px" h="72px">
                                    <motion.div
                                        animate={{ rotate: 360 }}
                                        transition={{ duration: 1.5, repeat: Infinity, ease: 'linear' }}
                                        style={{
                                            width: '72px', height: '72px', borderRadius: '50%',
                                            border: '4px solid transparent',
                                            borderTopColor: '#00B5D8', borderRightColor: '#9F7AEA',
                                            position: 'absolute',
                                        }}
                                    />
                                    <motion.div
                                        animate={{ scale: [1, 1.3, 1], opacity: [0.4, 0.1, 0.4] }}
                                        transition={{ duration: 2, repeat: Infinity, ease: 'easeInOut' }}
                                        style={{
                                            width: '72px', height: '72px', borderRadius: '50%',
                                            border: '2px solid #00B5D840',
                                            position: 'absolute',
                                        }}
                                    />
                                </Box>
                                <motion.div
                                    key={statusMessage}
                                    initial={{ opacity: 0, y: 5 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    transition={{ duration: 0.3 }}
                                >
                                    <Text color="cyan.300" fontSize="xl" fontWeight={600} textAlign="center">{statusMessage}</Text>
                                </motion.div>
                            </Flex>
                        ) : (
                            <HStack spacing={16} justifyContent="center">
                                <Box textAlign="center">
                                    <HStack spacing={3} justifyContent="center" mb={3}>
                                        <Box color={hasResult ? 'cyan.300' : 'whiteAlpha.300'}><ClockIcon /></Box>
                                        <Text fontSize="md" fontWeight={700} color={hasResult ? 'whiteAlpha.700' : 'whiteAlpha.500'} letterSpacing="0.5px">ETA</Text>
                                    </HStack>
                                    <Text fontSize="5xl" fontWeight={800} fontFamily="mono" color={hasResult ? 'cyan.300' : 'whiteAlpha.200'} lineHeight={1}>
                                        {hasResult ? (
                                            <Countdown date={etaTarget!}
                                                renderer={countdownRenderer} overtime={lastRequestedDrive.driverArrival! < 0} />
                                        ) : '--:--'}
                                    </Text>
                                </Box>
                                <Box textAlign="center">
                                    <HStack spacing={3} justifyContent="center" mb={3}>
                                        <Box color={hasResult ? 'purple.300' : 'whiteAlpha.300'}><UserIcon /></Box>
                                        <Text fontSize="md" fontWeight={700} color={hasResult ? 'whiteAlpha.700' : 'whiteAlpha.500'} letterSpacing="0.5px">Driver</Text>
                                    </HStack>
                                    <Text fontSize="2xl" fontWeight={700} color={hasResult ? 'whiteAlpha.900' : 'whiteAlpha.200'} lineHeight={1}>
                                        {hasResult ? lastRequestedDrive.driverDetails.name : '--'}
                                    </Text>
                                </Box>
                                <Box textAlign="center">
                                    <HStack spacing={3} justifyContent="center" mb={3}>
                                        <Box color={hasResult ? 'orange.300' : 'whiteAlpha.300'}><CarIcon /></Box>
                                        <Text fontSize="md" fontWeight={700} color={hasResult ? 'whiteAlpha.700' : 'whiteAlpha.500'} letterSpacing="0.5px">License #</Text>
                                    </HStack>
                                    <Text fontSize="2xl" fontWeight={700} fontFamily="mono" color={hasResult ? 'whiteAlpha.900' : 'whiteAlpha.200'} lineHeight={1}>
                                        {hasResult ? lastRequestedDrive.driverDetails.plate : '--'}
                                    </Text>
                                </Box>
                                {lastRequestedDrive.driverDistance != null && (
                                    <Box textAlign="center">
                                        <Text fontSize="md" fontWeight={700} color="whiteAlpha.700" mb={3} letterSpacing="0.5px">Distance</Text>
                                        <Text fontSize="2xl" fontWeight={700} color="whiteAlpha.900" lineHeight={1}>
                                            {lastRequestedDrive.driverDistance} mi
                                        </Text>
                                    </Box>
                                )}
                            </HStack>
                        )}
                    </Box>

                    {showArchitecture && (
                        <RequestFlowDiagram currentLog={logs.length > 0 ? logs[0] : undefined} showArchitecture={showArchitecture} />
                    )}
                </div>

                {/* Logs Panel */}
                {showArchitecture && (
                    <div className={`${styles.bottomPanel} ${logsPanel.isOpen ? styles.bottomPanelExpanded : styles.bottomPanelCollapsed}`}>
                        <div className={styles.bottomPanelHeader} onClick={logsPanel.onToggle}>
                            <HStack spacing={3}>
                                <Text fontSize="md" fontWeight={700} color="whiteAlpha.800">Request Logs</Text>
                                {logs.length > 0 && (
                                    <Badge variant="subtle" colorScheme="gray" fontSize="12px" px={2} py={1}>
                                        {logs.reduce((a, l) => a + l.entries.length, 0)} entries
                                    </Badge>
                                )}
                            </HStack>
                            <Text fontSize="md" color="whiteAlpha.500" transform={logsPanel.isOpen ? 'rotate(180deg)' : 'none'} transition="transform 0.2s">&#9660;</Text>
                        </div>
                        {logsPanel.isOpen && (
                            <div className={styles.bottomPanelBody}>
                                <Logs logs={logs} />
                            </div>
                        )}
                    </div>
                )}
            </div>
        </MainLayout>
    );
};
