import {
    Accordion,
    AccordionButton,
    AccordionIcon,
    AccordionItem,
    AccordionPanel,
    Box,
    Highlight, HStack, Stack, Text,
} from "@chakra-ui/react";
import {Log} from "../../../hooks/useLogs.tsx";
import {useMemo} from "react";

const getTime = (dt: Date) => {
    return  String(dt.getHours()).padStart(2, '0') + ":" +
        String(dt.getMinutes()).padStart(2, '0') + ":" +
        String(dt.getSeconds()).padStart(2, '0') + "." +
        String(dt.getMilliseconds()).padStart(3, '0');
}


type LocationHighlightProps = {
    type: 'pickup' | 'dropoff',
    value: string
}

const LocationHighlight = ({ type, value }: LocationHighlightProps) => {
    return (
        <Highlight
            query={value}
            styles={{
                px: '1',
                py: '1',
                bg: type === 'pickup' ? 'orange.900' : 'teal.900',
                color: type === 'pickup' ? 'orange.200' : 'teal.200',
                mx: '2',
                borderRadius: '2px',
        }}>
            {value}
        </Highlight>
    )
}

type BaseLogProps = {
    log: Log
};

const BaseLog = ({ log }: BaseLogProps) => {
    const { requestID, dropoffLocation, pickupLocation, entries} = log;

    const servicesColor: Record<string, string> = {
        'route': '#ECC94B',
        'driver': '#4299E1',
        'location': '#48BB78',
        'frontend': '#76E4F7',
        'browser': '#B794F4',
    };

    const entriesMemo = useMemo(() => {
        return entries.map(e => {
                const serviceColor = e.service.length > 0 ? servicesColor[e.service] : "black";

                return (
                    <HStack fontWeight='bold'>
                        <Text >{getTime(e.date)}</Text>
                        <Text color={serviceColor}>{e.service}</Text>
                        <Text color={serviceColor}>({e.sandboxName && e.sandboxName.length > 0 ? e.sandboxName : 'baseline'})</Text>
                        <Text color='green.300'>{e.status}</Text>
                    </HStack>
                )
            })
    }, [entries])


    return (
        <AccordionItem key={requestID}>
            <h2>
                <AccordionButton>
                    <Box as="span" flex='1' textAlign='left'>
                        Request ID: #{requestID} from <LocationHighlight value={pickupLocation.name} type='pickup'/> to <LocationHighlight value={dropoffLocation.name} type='dropoff'/>
                    </Box>
                    <AccordionIcon />
                </AccordionButton>
            </h2>
            <AccordionPanel pb={4}>
                <Stack>
                    { entriesMemo }
                </Stack>
            </AccordionPanel>
        </AccordionItem>
    )
}


type LogsProps = {
    logs: Log[]
}



export const Logs = ({ logs }: LogsProps) => {
    const logsMemo = useMemo(() => {
        return logs.map(l => <BaseLog log={l}/>)
    }, [logs])

    return (
        <Accordion allowMultiple defaultIndex={[0]} pb={8} mb={8}>
            {logsMemo}
        </Accordion>
    )
}