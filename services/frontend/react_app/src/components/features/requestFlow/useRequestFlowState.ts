import { useMemo, useRef } from 'react';
import { Log } from '../../../hooks/useLogs.tsx';
import { ServiceId, NodeId, NodeVariant, STAGE_MATCHERS } from './constants.ts';

export type FlowNodeState = {
    status: 'idle' | 'active' | 'completed';
    sandboxName: string;
    statusText: string;
};

export type FlowEdgeState = {
    status: 'idle' | 'animating' | 'completed';
};

export type RequestFlowState = {
    nodes: Record<ServiceId, {
        baseline: FlowNodeState;
        sandbox: FlowNodeState;
    }>;
    kafka: { status: 'idle' | 'active' | 'completed' };
    edges: Record<string, FlowEdgeState>;
    activePath: Record<ServiceId, NodeVariant>;
    isComplete: boolean;
    hasSandbox: boolean;
    requestId: number | null;
};

const defaultNodeState = (): FlowNodeState => ({
    status: 'idle',
    sandboxName: '',
    statusText: '',
});

const defaultFlowState = (): RequestFlowState => ({
    nodes: {
        frontend: { baseline: defaultNodeState(), sandbox: defaultNodeState() },
        location: { baseline: defaultNodeState(), sandbox: defaultNodeState() },
        driver: { baseline: defaultNodeState(), sandbox: defaultNodeState() },
        route: { baseline: defaultNodeState(), sandbox: defaultNodeState() },
    },
    kafka: { status: 'idle' },
    edges: {},
    activePath: {
        frontend: 'baseline',
        location: 'baseline',
        driver: 'baseline',
        route: 'baseline',
    },
    isComplete: false,
    hasSandbox: false,
    requestId: null,
});

// Build an edge key that encodes the actual from-variant and to-variant.
// This allows cross-lane edges (e.g., sandbox-frontend → baseline-location).
// Format: "path-{fromVariant}-{fromService}-to-{toVariant}-{toService}"
export function buildPathEdgeKey(
    fromService: NodeId,
    fromVariant: NodeVariant | 'shared',
    toService: NodeId,
    toVariant: NodeVariant | 'shared',
): string {
    return `path-${fromVariant}-${fromService}-to-${toVariant}-${toService}`;
}

export function useRequestFlowState(currentLog: Log | undefined): RequestFlowState {
    const prevEntryCountRef = useRef(0);
    const prevRequestIdRef = useRef<number | null>(null);

    return useMemo(() => {
        if (!currentLog || currentLog.entries.length === 0) {
            return defaultFlowState();
        }

        const state = defaultFlowState();
        state.requestId = currentLog.requestID;

        // Track whether this is a new request
        const isNewRequest = currentLog.requestID !== prevRequestIdRef.current;
        if (isNewRequest) {
            prevEntryCountRef.current = 0;
            prevRequestIdRef.current = currentLog.requestID;
        }

        let hasSandbox = false;
        const activatedServices = new Set<ServiceId>();
        let kafkaActivated = false;

        for (const entry of currentLog.entries) {
            // Skip browser entries
            if (entry.service === 'browser' || entry.service === 'api') continue;

            const variant: NodeVariant = entry.sandboxName && entry.sandboxName.length > 0
                ? 'sandbox'
                : 'baseline';

            if (variant === 'sandbox') hasSandbox = true;

            // Match against stage patterns
            for (const matcher of STAGE_MATCHERS) {
                if (!matcher.bodyPattern.test(entry.status)) continue;

                const entryService = entry.service === '' ? 'frontend' : entry.service;
                if (entryService !== matcher.serviceId) continue;

                const nodeState = state.nodes[matcher.serviceId][variant];
                nodeState.status = matcher.isFinal ? 'completed' : 'active';
                nodeState.sandboxName = entry.sandboxName || '';
                nodeState.statusText = entry.status;

                state.activePath[matcher.serviceId] = variant;
                activatedServices.add(matcher.serviceId);

                // Activate path edges based on actual from→to variant pairs
                if (matcher.serviceId === 'location') {
                    // frontend → location (HTTP): use actual variants
                    const fromV = state.activePath.frontend;
                    const toV = variant;
                    const key = buildPathEdgeKey('frontend', fromV, 'location', toV);
                    state.edges[key] = { status: 'completed' };
                    // location → mysql (always baseline, shared infra)
                    const mysqlKey = buildPathEdgeKey('location', toV, 'mysql', 'shared');
                    state.edges[mysqlKey] = { status: 'completed' };
                }

                if (matcher.serviceId === 'driver' && !matcher.isFinal) {
                    kafkaActivated = true;
                    // frontend → kafka: use frontend's variant
                    const frontendV = state.activePath.frontend;
                    const fkKey = buildPathEdgeKey('frontend', frontendV, 'kafka', 'shared');
                    state.edges[fkKey] = { status: 'completed' };
                    // kafka → driver: use driver's variant
                    const kdKey = buildPathEdgeKey('kafka', 'shared', 'driver', variant);
                    state.edges[kdKey] = { status: 'completed' };
                }

                if (matcher.serviceId === 'route') {
                    // driver → route (gRPC): use actual variants
                    const fromV = state.activePath.driver;
                    const toV = variant;
                    const key = buildPathEdgeKey('driver', fromV, 'route', toV);
                    state.edges[key] = { status: 'completed' };
                }

                if (matcher.isFinal) {
                    state.isComplete = true;
                    for (const svc of activatedServices) {
                        const v = state.activePath[svc];
                        if (state.nodes[svc][v].status === 'active') {
                            state.nodes[svc][v].status = 'completed';
                        }
                    }
                }

                break;
            }
        }

        // Mark the most recently activated node as 'active' (others as 'completed')
        if (!state.isComplete && activatedServices.size > 0) {
            const serviceOrder: ServiceId[] = ['frontend', 'location', 'driver', 'route'];
            let lastActivated: ServiceId | null = null;
            for (const svc of serviceOrder) {
                if (activatedServices.has(svc)) {
                    lastActivated = svc;
                }
            }
            for (const svc of serviceOrder) {
                if (!activatedServices.has(svc)) continue;
                const v = state.activePath[svc];
                if (svc === lastActivated) {
                    state.nodes[svc][v].status = 'active';
                } else {
                    state.nodes[svc][v].status = 'completed';
                }
            }
        }

        state.kafka.status = kafkaActivated
            ? (state.isComplete ? 'completed' : 'active')
            : 'idle';
        state.hasSandbox = hasSandbox;

        prevEntryCountRef.current = currentLog.entries.length;

        return state;
    }, [currentLog?.requestID, currentLog?.entries.length]);
}
