import {MainLayout} from "../components/layouts";
import {
    Box,
    Button,
    Card,
    CardBody,
    CardHeader,
    Heading,
    HStack,
    Stack,
    StackDivider,
    Text,
    useDisclosure,
    useToast,
} from "@chakra-ui/react";
import styles from "./home.module.css";

import {Logs} from "../components/features/logs/logs.tsx";
import {Map} from "../components/features/map/map.tsx";
import {useSession} from "../context/sessionContext/context.tsx";
import {useEffect, useRef, useState} from "react";
import {apiGet, apiPost} from "../services/http.ts";
import {Locations} from "../types/location.ts";
import {LocationSelect} from "../components/common/locationSelect/locationSelect.tsx";
import {useLogs} from "../hooks/useLogs.tsx";
import {NotificationResponse} from "../types/notifications.ts";
import {useGetRequestArrival} from "../hooks/useGetRequestArrival.tsx";
import Countdown, {CountdownRenderProps} from "react-countdown";

const countdownRenderer = ({minutes, seconds, completed, driverName, driverPlate}: CountdownRenderProps & { driverName: string, driverPlate: string }) => {
    if (completed) {
        return <Text>{driverName} arrived</Text>;
    } else {
        return <Text as="b">The driver {driverName} ({driverPlate}) will arrive in {minutes.toString().padStart(2, "0")}:{seconds.toString().padStart(2, "0")}</Text>;
    }
};

export const HomePage = () => {
    const session = useSession();
    const notificationCursorRef = useRef(-1);
    const [locations, setLocations] = useState<Locations | undefined>();

    const [selectedLocations, setSelectedLocations] = useState({pickupId: -1, dropoffId: -1});
    const {logs, addNewLog, addErrorEntry, addInformationEntry} = useLogs();
    const lastRequestedDrive = useGetRequestArrival(logs);

    const logsModal = useDisclosure();
    const toast = useToast();

    useEffect(() => {
        const fetchLocations = async () => {
            const locations = await apiGet<Locations>('/splash');
            setLocations(locations)
        }

        fetchLocations();
    }, []);

    useEffect(() => {
        const pollNotifications = async () => {
            const cursor = notificationCursorRef.current;
            const nonse = Math.random();
            const url = `/notifications?sessionID=${session.sessionID}&cursor=${cursor}&nonse=${nonse}`
            const notification = await apiGet<NotificationResponse>(url);

            notificationCursorRef.current = notification.cursor;

            notification.notifications.forEach(notification => {
                addInformationEntry(notification.context.request.id, {
                    date: new Date(notification.timestamp),
                    service: notification.context.baselineWorkload,
                    status: notification.body,
                    sandboxName: notification.context.sandboxName,
                })
            });
        }

        const intervalID = setInterval(pollNotifications, 2000);

        return () => {
            clearInterval(intervalID);
        }

    }, []);

    if (!locations) {
        return (
            <MainLayout titleSuffix="">
                <Heading>Loading</Heading>
            </MainLayout>
        )
    }

    const handleRequestDrive = async () => {
        const {getLastRequestID, sessionID} = session;
        const requestID = getLastRequestID();

        const {pickupId, dropoffId} = selectedLocations;
        const data = {
            sessionID: sessionID,
            requestID: requestID,
            pickupLocationID: pickupId,
            dropoffLocationID: dropoffId
        }

        const pickupLocation = locations?.Locations.find(l => l.id === pickupId);
        const dropoffLocation = locations?.Locations.find(l => l.id === dropoffId);

        toast({
            title: 'Hotrod::Ride requested.',
            description: `Ride requested to ${dropoffLocation?.name}`,
            status: 'success',
            duration: 9000,
            isClosable: true,
            position: "top"
        });

        // Reset locations
        setSelectedLocations({pickupId: -1, dropoffId: -1});

        addNewLog(pickupLocation!, dropoffLocation!, requestID, {
            messageType: "info",
            service: 'browser',
            date: new Date(),
            status: 'Requesting a ride.'
        });

        try {
            await apiPost<{}>('/dispatch', data, {"baggage": `session=${sessionID}, request=${requestID}`});
        } catch (e) {
            addErrorEntry(requestID, {
                status: 'Error requesting a ride to frontend API',
                service: 'api',
                date: new Date(),
            })
        }
    }

    return (
        <MainLayout titleSuffix={locations.TitleSuffix}>
            <HStack alignItems='flex-start' p={4} gap={8} justifyContent='space-between' h='100%'>
                <Stack flexGrow={1} w='50%'>
                    <Card border={12} maxW={600}>
                        <CardHeader>
                            <Heading size='lg' textAlign='left'>
                                Go anywhere with HotROD
                            </Heading>
                            <Heading size='xs' textAlign='left'>
                                Request a ride, hop in, and go.
                            </Heading>
                        </CardHeader>
                        <CardBody>
                            <Stack divider={<StackDivider/>} spacing='4'>
                                <Box>
                                    <LocationSelect
                                        placeholder='Pickup location'
                                        locations={locations.Locations}
                                        onSelect={locationID => setSelectedLocations(prev => ({
                                            ...prev,
                                            pickupId: locationID
                                        }))}
                                        selectedLocationID={selectedLocations.pickupId}
                                    />
                                </Box>
                                <Box>
                                    <LocationSelect
                                        placeholder='Dropoff location'
                                        locations={locations.Locations}
                                        onSelect={locationID => setSelectedLocations(prev => ({
                                            ...prev,
                                            dropoffId: locationID
                                        }))}
                                        selectedLocationID={selectedLocations.dropoffId}
                                    />
                                </Box>
                                <Box>
                                    <Button
                                        variant='solid'
                                        colorScheme='cyan'
                                        onClick={handleRequestDrive}
                                        isDisabled={
                                            selectedLocations.pickupId === -1 ||
                                            isNaN(selectedLocations.pickupId) ||
                                            selectedLocations.dropoffId === -1 ||
                                            isNaN(selectedLocations.dropoffId)
                                        }
                                    >
                                        Request Ride
                                    </Button>
                                </Box>
                            </Stack>
                        </CardBody>
                    </Card>

                    {lastRequestedDrive.driverArrival &&
                        <HStack divider={<StackDivider/>} spacing='4' justifyContent="center" alignItems="center"
                                bg="teal.200" p="1rem" maxW={600} borderRadius={4}>
                            <Countdown
                                date={new Date(Date.now()).setSeconds(lastRequestedDrive.driverArrival)}
                                renderer={(props) => countdownRenderer({
                                    ...props,
                                    driverName: lastRequestedDrive.driverDetails.name,
                                    driverPlate: lastRequestedDrive.driverDetails.plate,
                                })}
                            />
                        </HStack>}
                </Stack>
                <Stack flexGrow={1} justifyContent='space-between' w='50%' h='100%' maxH={'900px'}>
                    <div className={`${styles.drawer} ${logsModal.isOpen ? styles.open : ''}`}>
                        <div className={styles.drawerHeader}>
                            <Heading size="md">Notifications Logs</Heading>
                        </div>
                        <div className={styles.drawerBody}>
                            <Logs logs={logs}/>
                        </div>
                        <div className={styles.drawerFooter}>
                            <Button variant="outline" onClick={logsModal.onClose}>
                                Close
                            </Button>
                        </div>
                    </div>

                    <Stack flexGrow={1} flexShrink={1}>
                        <Map dropoffLocationID={selectedLocations.dropoffId}
                             pickupLocationID={selectedLocations.pickupId}/>
                    </Stack>
                    <Button variant='outline' size="sm" onClick={logsModal.onToggle}>Show Logs</Button>
                </Stack>
            </HStack>
        </MainLayout>
    )
}