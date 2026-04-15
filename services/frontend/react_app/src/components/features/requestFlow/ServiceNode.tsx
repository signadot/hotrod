import { motion } from 'framer-motion';
import { FlowNodeState } from './useRequestFlowState.ts';
import { NodeVariant, NODE_WIDTH, NODE_HEIGHT, SERVICE_COLORS, ServiceId, SANDBOX_COLOR, IDLE_COLOR } from './constants.ts';

type ServiceNodeProps = {
    serviceId: ServiceId;
    variant: NodeVariant;
    state: FlowNodeState;
    x: number;
    y: number;
    label: string;
    visible: boolean;
};

export const ServiceNode = ({ serviceId, variant, state, x, y, label, visible }: ServiceNodeProps) => {
    if (!visible) return null;

    const serviceColor = SERVICE_COLORS[serviceId];
    const isSandbox = variant === 'sandbox';
    const isActive = state.status === 'active';
    const isCompleted = state.status === 'completed';
    // baseline always looks the same regardless of idle/active/completed

    // BASELINE: Always gray/muted — same appearance idle or active.
    // Completed baseline gets a subtle checkmark but stays gray.
    // SANDBOX: Vibrant and colorful — purple gradient bg, bright border, glow.
    let borderColor: string;
    let bgColor: string;
    let borderStyle: string;
    let labelColor: string;
    let sublabelColor: string;
    let glowShadow = 'none';

    if (isSandbox) {
        // Sandbox — colorful and vibrant
        borderColor = SANDBOX_COLOR;
        bgColor = (isActive || isCompleted)
            ? 'linear-gradient(135deg, rgba(159, 122, 234, 0.25) 0%, rgba(128, 90, 213, 0.15) 100%)'
            : 'rgba(159, 122, 234, 0.10)';
        borderStyle = '3px dashed';
        labelColor = '#E9D8FD';
        sublabelColor = '#B794F4';
        if (isActive) glowShadow = `0 0 24px ${SANDBOX_COLOR}50, 0 0 8px ${SANDBOX_COLOR}30`;
        else if (isCompleted) glowShadow = `0 0 12px ${SANDBOX_COLOR}35`;
    } else {
        // Baseline — always gray/muted, consistent appearance
        borderColor = IDLE_COLOR;
        bgColor = '#1A202C';
        borderStyle = '2px solid';
        labelColor = '#A0AEC0';
        sublabelColor = '#718096';
        // Completed baseline gets a very subtle glow
        if (isCompleted) glowShadow = `0 0 6px ${serviceColor}15`;
    }

    return (
        <g>
            <foreignObject x={x} y={y} width={NODE_WIDTH} height={NODE_HEIGHT}>
                <motion.div
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{
                        opacity: 1,
                        scale: (isSandbox && isActive) ? [1, 1.03, 1] : 1,
                        boxShadow: glowShadow,
                    }}
                    transition={{
                        opacity: { duration: 0.3 },
                        scale: (isSandbox && isActive)
                            ? { duration: 0.8, repeat: Infinity, ease: 'easeInOut' }
                            : { duration: 0.3 },
                        boxShadow: { duration: 0.4 },
                    }}
                    style={{
                        width: '100%',
                        height: '100%',
                        background: bgColor,
                        border: `${borderStyle} ${borderColor}`,
                        borderRadius: isSandbox ? '12px' : '10px',
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: 'center',
                        padding: '4px 8px',
                        position: 'relative',
                    }}
                >
                    <div style={{
                        fontSize: '13px',
                        fontWeight: 600,
                        color: labelColor,
                        lineHeight: 1.2,
                    }}>
                        {label}
                    </div>
                    <div style={{
                        fontSize: '10px',
                        fontWeight: 500,
                        color: sublabelColor,
                        marginTop: '2px',
                        maxWidth: '125px',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                    }}>
                        {isSandbox ? (state.sandboxName || 'sandbox') : 'baseline'}
                    </div>
                    {isCompleted && (
                        <motion.div
                            initial={{ scale: 0 }}
                            animate={{ scale: 1 }}
                            style={{
                                position: 'absolute',
                                top: 4,
                                right: 6,
                                fontSize: '10px',
                                color: isSandbox ? SANDBOX_COLOR : serviceColor,
                            }}
                        >
                            &#10003;
                        </motion.div>
                    )}
                </motion.div>
            </foreignObject>
        </g>
    );
};
