import { Box, Flex, Text } from '@chakra-ui/react';
import { FlowCanvas } from './FlowCanvas.tsx';
import { useRequestFlowState } from './useRequestFlowState.ts';
import { Log } from '../../../hooks/useLogs.tsx';

type RequestFlowDiagramProps = {
    currentLog: Log | undefined;
    showArchitecture: boolean;
};

export const RequestFlowDiagram = ({ currentLog, showArchitecture }: RequestFlowDiagramProps) => {
    const flowState = useRequestFlowState(currentLog);

    if (!showArchitecture && !currentLog) {
        return null;
    }

    return (
        <Box
            bg="gray.800"
            borderRadius="lg"
            border="1px solid"
            borderColor="gray.700"
            flex={1}
            display="flex"
            flexDirection="column"
            overflow="hidden"
            minH={0}
        >
            {/* Header */}
            <Flex
                px={6}
                py={4}
                borderBottom="1px solid"
                borderColor="gray.700"
                alignItems="center"
                justifyContent="space-between"
                minH="52px"
            >
                <Text fontSize="md" fontWeight={700} color="whiteAlpha.800" letterSpacing="1px">
                    {currentLog ? 'REQUEST FLOW' : 'ARCHITECTURE'}
                </Text>
                {flowState.requestId && (
                    <Text fontSize="sm" color="whiteAlpha.600" fontFamily="mono">
                        Request #{flowState.requestId}
                        {flowState.isComplete && ' — completed'}
                    </Text>
                )}
            </Flex>

            {/* Canvas */}
            <Box flex={1} p={4} minH={0}>
                <FlowCanvas flowState={flowState} showTopology={showArchitecture} />
            </Box>
        </Box>
    );
};
