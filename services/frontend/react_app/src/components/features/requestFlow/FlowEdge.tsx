import { FlowEdgeState } from './useRequestFlowState.ts';
import { NodeVariant } from './constants.ts';

// Unified active color for ALL request flow lines (both baseline and sandbox)
const ACTIVE_COLOR = '#00B5D8'; // cyan.400
// Topology (architecture view) color — visible but not animated
const TOPO_COLOR = '#718096'; // gray.500

type FlowEdgeProps = {
    pathD: string;
    protocol: string;
    state: FlowEdgeState;
    variant: NodeVariant;
    labelX: number;
    labelY: number;
};

export const FlowEdge = ({ pathD, protocol, state, labelX, labelY }: FlowEdgeProps) => {
    const isIdle = state.status === 'idle';

    if (isIdle) {
        // Topology/architecture view: visible, static, with arrows
        return (
            <g>
                <path d={pathD} fill="none" stroke={TOPO_COLOR} strokeWidth={2}
                    strokeDasharray="8 5" strokeLinecap="round" strokeLinejoin="round"
                    markerEnd="url(#arrow-topo)" opacity={0.6} />
                <text x={labelX} y={labelY} textAnchor="middle" fontSize={16} fontWeight={600}
                    fontFamily="system-ui, sans-serif" fill={TOPO_COLOR} opacity={0.65}>
                    {protocol}
                </text>
            </g>
        );
    }

    // Active request flow: all lines same cyan color, animated dashes
    return (
        <g>
            <path d={pathD} fill="none" stroke={ACTIVE_COLOR} strokeWidth={3.5}
                strokeDasharray="14 8" strokeLinecap="round" strokeLinejoin="round"
                markerEnd="url(#arrow-active)" opacity={0.9}>
                <animate attributeName="stroke-dashoffset" from="44" to="0" dur="1.2s" repeatCount="indefinite" />
            </path>
            <text x={labelX} y={labelY} textAnchor="middle" fontSize={16} fontWeight={700}
                fontFamily="system-ui, sans-serif" fill={ACTIVE_COLOR} opacity={0.85}>
                {protocol}
            </text>
        </g>
    );
};
