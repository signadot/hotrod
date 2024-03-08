import {MainLayout} from "../components/layouts";
import {
    Box, Button,
    Card,
    CardBody,
    CardHeader,
    Heading,
    HStack,
    Stack,
    StackDivider,
    Text
} from "@chakra-ui/react";

import {Logs} from "../components/features/logs/logs.tsx";
import {Map} from "../components/features/map/map.tsx";
import {useSession} from "../context/sessionContext/context.tsx";
import {useEffect, useRef, useState} from "react";
import {apiGet, apiPost} from "../services/http.ts";
import {Locations} from "../types/location.ts";
import {LocationSelect} from "../components/common/locationSelect/locationSelect.tsx";
import {useLogs} from "../hooks/useLogs.tsx";
import {NotificationResponse} from "../types/notifications.ts";



export const HomePage = () => {
    const session = useSession();
    const notificationCursorRef = useRef(-1);
    const [locations, setLocations] = useState<Locations | undefined>();

    const [selectedLocations, setSelectedLocations] = useState({ pickupId: -1, dropoffId: -1});
    const { logs, addNewLog, addErrorEntry, addInformationEntry } = useLogs();

    useEffect(() => {
        const fetchLocations = async () => {
            const locations = await apiGet<Locations>('/splash');
            setLocations(locations)
        }

        fetchLocations();

    }, []);

    useEffect(() => {
        console.log('logs changed', logs);
    }, [logs]);

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
            <MainLayout>
                <Heading>Loading</Heading>
            </MainLayout>
        )
    }

    const handleRequestDrive = async () => {
        const { getLastRequestID, sessionID} = session;
        const requestID = getLastRequestID();

        const {pickupId, dropoffId} = selectedLocations;
        const data = {
            sessionID: sessionID,
            requestID: requestID,
            pickupLocationID: pickupId,
            dropoffLocationID: dropoffId
        }

        const pickupLocation = locations?.Locations.find(l => l.ID === pickupId);
        const dropoffLocation = locations?.Locations.find(l => l.ID === dropoffId);

        addNewLog(pickupLocation!, dropoffLocation!, requestID, {
            messageType: "info",
            service: 'browser',
            date: new Date(),
            status: 'Requesting a ride.'

        });

        try {
            await apiPost<{}>('/dispatch', data, { "baggage": `session=${sessionID}, request=${requestID}`});
        } catch (e) {
            addErrorEntry(requestID, {
                status: 'Error requesting a ride to frontend API',
                service: 'api',
                date: new Date(),
            })

            // logError(new Date(), 'Error requesting a ride to frontend API', {
            //     request: {
            //         id: lastRequestID,
            //         pickupLocationID: pickupLocationID,
            //         dropoffLocationID: dropoffLocationID
            //     },
            //     error: error,
            // })
        }
    }

    return (
        <MainLayout>
            <HStack alignItems='flex-start' p={4} gap={8} justifyContent='space-between'>
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
                            <Stack divider={<StackDivider />} spacing='4'>
                                <Box>
                                    <LocationSelect
                                        placeholder='Pickup location'
                                        locations={locations.Locations}
                                        // locations={[{"ID":1,"Name":"My Home","Coordinates":"231,773"},{"ID":123,"Name":"Rachel's Floral Designs","Coordinates":"115,277"},{"ID":392,"Name":"Trom Chocolatier","Coordinates":"577,322"},{"ID":567,"Name":"Amazing Coffee Roasters","Coordinates":"211,653"},{"ID":731,"Name":"Japanese Desserts","Coordinates":"728,326"}]}
                                        onSelect={locationID => setSelectedLocations(prev => ({...prev, pickupId: locationID}))}
                                        selectedLocationID={selectedLocations.pickupId}
                                    />
                                </Box>
                                <Box>
                                    <LocationSelect
                                        placeholder='Dropoff location'
                                        locations={locations.Locations}
                                        // locations={[{"ID":1,"Name":"My Home","Coordinates":"231,773"},{"ID":123,"Name":"Rachel's Floral Designs","Coordinates":"115,277"},{"ID":392,"Name":"Trom Chocolatier","Coordinates":"577,322"},{"ID":567,"Name":"Amazing Coffee Roasters","Coordinates":"211,653"},{"ID":731,"Name":"Japanese Desserts","Coordinates":"728,326"}]}
                                        onSelect={locationID => setSelectedLocations(prev => ({...prev, dropoffId: locationID}))}
                                        selectedLocationID={selectedLocations.dropoffId}
                                    />
                                </Box>
                                <Box>
                                    <Button
                                        variant='solid'
                                        colorScheme='cyan'
                                        onClick={handleRequestDrive}
                                        isDisabled={selectedLocations.pickupId === -1 || selectedLocations.dropoffId === -1}
                                    >
                                        Request Drive
                                    </Button>
                                    <Text color='gray' mt={4}>The distance between the two points is 40 miles</Text>
                                </Box>
                            </Stack>
                        </CardBody>
                    </Card>
                </Stack>
                <Stack flexGrow={1} h={650} position='relative' justifyContent='space-between' w='50%'>
                    <Map />

                    <Stack position='absolute' bottom={0} left={0} w='100%' backgroundColor='white' overflowY='auto' maxH={'60%'}>
                        <Logs logs={logs} />
                    </Stack>
                </Stack>
            </HStack>
        </MainLayout>
    )
}