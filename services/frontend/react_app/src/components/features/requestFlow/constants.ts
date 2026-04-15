export type ServiceId = 'frontend' | 'location' | 'driver' | 'route';
export type InfraId = 'kafka' | 'mysql';
export type NodeId = ServiceId | InfraId;
export type NodeVariant = 'baseline' | 'sandbox';

export type NodeConfig = {
    id: ServiceId;
    label: string;
    x: number;
    y: number;
    color: string;
};

export type StageConfig = {
    serviceId: ServiceId;
    bodyPattern: RegExp;
    isFinal: boolean;
};

export const SERVICE_COLORS: Record<ServiceId, string> = {
    frontend: '#76E4F7',
    location: '#48BB78',
    driver: '#4299E1',
    route: '#ECC94B',
};

export const BASELINE_COLOR = '#00B5D8';
export const SANDBOX_COLOR = '#9F7AEA';
export const KAFKA_COLOR = '#F6AD55';
export const MYSQL_COLOR = '#63B3ED';
export const IDLE_COLOR = '#4A5568';

export const NODE_WIDTH = 130;
export const NODE_HEIGHT = 52;
export const INFRA_WIDTH = 110;
export const INFRA_HEIGHT = 44;
export const DB_WIDTH = 70;
export const DB_HEIGHT = 50;

export const VIEWBOX_WIDTH = 780;
export const VIEWBOX_HEIGHT = 340;

// === HORIZONTAL BRANCHING LAYOUT ===
// Inspired by the reference diagram (rotated horizontal):
//
//                   ┌─── Location ─── MySQL        (top branch)
//   Frontend ───────┤
//                   └─── Kafka ─── Driver ─── Route  (bottom branch)
//
// Frontend is on the left. It branches into two horizontal paths:
//   - Top path: Location → MySQL (HTTP + database)
//   - Bottom path: Kafka → Driver → Route (messaging + gRPC)
//
// Sandbox versions appear BELOW their baseline counterpart, offset down.

const TOP_Y = 30;       // top branch (Location, MySQL)
const CENTER_Y = 130;   // Frontend sits between the two branches
const BOTTOM_Y = 130;   // bottom branch (Kafka, Driver, Route) — same as center
const SANDBOX_OFFSET = 120; // sandbox nodes appear this far below baseline

const COL0_X = 20;      // Frontend
const COL1_X = 200;     // Location / Kafka
const COL2_X = 380;     // MySQL / Driver
const COL3_X = 560;     // Route

export const BASELINE_NODES: NodeConfig[] = [
    { id: 'frontend', label: 'Frontend', x: COL0_X, y: CENTER_Y - NODE_HEIGHT / 2, color: SERVICE_COLORS.frontend },
    { id: 'location', label: 'Location', x: COL1_X, y: TOP_Y, color: SERVICE_COLORS.location },
    { id: 'driver',   label: 'Driver',   x: COL2_X, y: BOTTOM_Y, color: SERVICE_COLORS.driver },
    { id: 'route',    label: 'Route',    x: COL3_X, y: BOTTOM_Y, color: SERVICE_COLORS.route },
];

export const SANDBOX_NODES: NodeConfig[] = [
    { id: 'frontend', label: 'Frontend', x: COL0_X, y: CENTER_Y - NODE_HEIGHT / 2 + SANDBOX_OFFSET, color: SERVICE_COLORS.frontend },
    { id: 'location', label: 'Location', x: COL1_X, y: TOP_Y + SANDBOX_OFFSET, color: SERVICE_COLORS.location },
    { id: 'driver',   label: 'Driver',   x: COL2_X, y: BOTTOM_Y + SANDBOX_OFFSET, color: SERVICE_COLORS.driver },
    { id: 'route',    label: 'Route',    x: COL3_X, y: BOTTOM_Y + SANDBOX_OFFSET, color: SERVICE_COLORS.route },
];

// Infrastructure: Kafka on the bottom branch at col1, MySQL on top branch at col2
export const KAFKA_NODE = { x: COL1_X + 10, y: BOTTOM_Y + 4 };
export const MYSQL_NODE = { x: COL2_X + 25, y: TOP_Y - 2 };

export const STAGE_MATCHERS: StageConfig[] = [
    { serviceId: 'frontend', bodyPattern: /Processing dispatch/i, isFinal: false },
    { serviceId: 'location', bodyPattern: /Resolving locations/i, isFinal: false },
    { serviceId: 'driver',   bodyPattern: /Finding an available driver/i, isFinal: false },
    { serviceId: 'route',    bodyPattern: /Resolving routes/i, isFinal: false },
    { serviceId: 'driver',   bodyPattern: /Driver .* arriving/i, isFinal: true },
];

export function getNodeCenter(nodeId: NodeId, variant?: NodeVariant): { x: number; y: number } {
    if (nodeId === 'kafka') return { x: KAFKA_NODE.x + INFRA_WIDTH / 2, y: KAFKA_NODE.y + INFRA_HEIGHT / 2 };
    if (nodeId === 'mysql') return { x: MYSQL_NODE.x + DB_WIDTH / 2, y: MYSQL_NODE.y + DB_HEIGHT / 2 };
    const nodes = variant === 'sandbox' ? SANDBOX_NODES : BASELINE_NODES;
    const node = nodes.find(n => n.id === nodeId)!;
    return { x: node.x + NODE_WIDTH / 2, y: node.y + NODE_HEIGHT / 2 };
}

export function getNodeEdgePoint(
    nodeId: NodeId, variant: NodeVariant | undefined,
    side: 'left' | 'right' | 'top' | 'bottom'
): { x: number; y: number } {
    const center = getNodeCenter(nodeId, variant);
    const isKafka = nodeId === 'kafka';
    const isDb = nodeId === 'mysql';
    const hw = isDb ? DB_WIDTH / 2 : isKafka ? INFRA_WIDTH / 2 : NODE_WIDTH / 2;
    const hh = isDb ? DB_HEIGHT / 2 : isKafka ? INFRA_HEIGHT / 2 : NODE_HEIGHT / 2;
    switch (side) {
        case 'left':   return { x: center.x - hw, y: center.y };
        case 'right':  return { x: center.x + hw, y: center.y };
        case 'top':    return { x: center.x, y: center.y - hh };
        case 'bottom': return { x: center.x, y: center.y + hh };
    }
}
