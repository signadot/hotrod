import { ServiceNode } from './ServiceNode.tsx';
import { FlowEdge } from './FlowEdge.tsx';
import { RequestFlowState, buildPathEdgeKey } from './useRequestFlowState.ts';
import {
    BASELINE_NODES, SANDBOX_NODES, KAFKA_NODE, MYSQL_NODE,
    INFRA_WIDTH, INFRA_HEIGHT, DB_WIDTH, DB_HEIGHT,
    KAFKA_COLOR, MYSQL_COLOR, IDLE_COLOR,
    VIEWBOX_WIDTH, VIEWBOX_HEIGHT,
    ServiceId, NodeId, NodeVariant, getNodeEdgePoint,
} from './constants.ts';

type FlowCanvasProps = { flowState: RequestFlowState; showTopology: boolean };
type EdgeDef = { id: string; pathD: string; labelX: number; labelY: number; protocol: string; variant: NodeVariant; isTopology?: boolean };

const GAP = 10;

// === EDGE PATH BUILDERS ===
// Horizontal: right→left
function hPath(fId: NodeId, fV: NodeVariant|undefined, tId: NodeId, tV: NodeVariant|undefined) {
    const f = getNodeEdgePoint(fId, fV, 'right'), t = getNodeEdgePoint(tId, tV, 'left');
    return { pathD: `M ${f.x} ${f.y} L ${t.x-GAP} ${t.y}`, labelX: (f.x+t.x)/2, labelY: Math.min(f.y,t.y)-8 };
}

// Frontend → Location: right side of Frontend, then up-right to Location (L-shape)
function frontendToLocation(fV: NodeVariant, tV: NodeVariant) {
    const f = getNodeEdgePoint('frontend', fV, 'right');
    const t = getNodeEdgePoint('location', tV, 'left');
    // Go right from frontend, then up to location row, then right to location
    const midX = f.x + 30;
    return {
        pathD: `M ${f.x} ${f.y} L ${midX} ${f.y} L ${midX} ${t.y} L ${t.x-GAP} ${t.y}`,
        labelX: midX + 30, labelY: (f.y + t.y) / 2 - 4,
    };
}

// Frontend → Kafka: right side of Frontend, stays on same row to Kafka
function frontendToKafka(fV: NodeVariant) {
    const f = getNodeEdgePoint('frontend', fV, 'right');
    const t = getNodeEdgePoint('kafka', undefined, 'left');
    if (Math.abs(f.y - t.y) < 10) {
        // Same row — straight horizontal
        return { pathD: `M ${f.x} ${f.y} L ${t.x-GAP} ${t.y}`, labelX: (f.x+t.x)/2, labelY: f.y - 8 };
    }
    // Different row (sandbox frontend below) — L-shape up to Kafka
    const midX = f.x + 30;
    return {
        pathD: `M ${f.x} ${f.y} L ${midX} ${f.y} L ${midX} ${t.y} L ${t.x-GAP} ${t.y}`,
        labelX: midX - 25, labelY: (f.y + t.y) / 2,
    };
}

// Kafka → Driver: horizontal on bottom branch
function kafkaToDriver(tV: NodeVariant) {
    const f = getNodeEdgePoint('kafka', undefined, 'right');
    const t = getNodeEdgePoint('driver', tV, 'left');
    if (Math.abs(f.y - t.y) < 10) {
        return { pathD: `M ${f.x} ${f.y} L ${t.x-GAP} ${t.y}`, labelX: (f.x+t.x)/2, labelY: f.y - 8 };
    }
    // Sandbox driver below — L-shape
    const midX = (f.x + t.x) / 2;
    return {
        pathD: `M ${f.x} ${f.y} L ${midX} ${f.y} L ${midX} ${t.y} L ${t.x-GAP} ${t.y}`,
        labelX: midX + 14, labelY: (f.y + t.y) / 2,
    };
}

// === TOPOLOGY EDGES ===
function buildTopologyEdges(): EdgeDef[] {
    const fToL = frontendToLocation('baseline', 'baseline');
    const lToM = hPath('location', 'baseline', 'mysql', undefined);
    const fToK = frontendToKafka('baseline');
    const kToD = kafkaToDriver('baseline');
    const dToR = hPath('driver', 'baseline', 'route', 'baseline');
    return [
        { id:'topo-fl', ...fToL, protocol:'HTTP', variant:'baseline', isTopology:true },
        { id:'topo-lm', ...lToM, protocol:'MySQL', variant:'baseline', isTopology:true },
        { id:'topo-fk', ...fToK, protocol:'publish', variant:'baseline', isTopology:true },
        { id:'topo-kd', ...kToD, protocol:'consume', variant:'baseline', isTopology:true },
        { id:'topo-dr', ...dToR, protocol:'gRPC', variant:'baseline', isTopology:true },
    ];
}

// === ACTIVE EDGES ===
function buildActiveEdges(flowState: RequestFlowState): EdgeDef[] {
    const edges: EdgeDef[] = [];
    const path = flowState.activePath;

    type PE = { key:string; build:()=>{pathD:string;labelX:number;labelY:number}; protocol:string; variant:NodeVariant };
    const possible: PE[] = [];

    // Baseline
    possible.push({ key:buildPathEdgeKey('frontend','baseline','location','baseline'), build:()=>frontendToLocation('baseline','baseline'), protocol:'HTTP', variant:'baseline' });
    possible.push({ key:buildPathEdgeKey('location','baseline','mysql','shared'), build:()=>hPath('location','baseline','mysql',undefined), protocol:'MySQL', variant:'baseline' });
    possible.push({ key:buildPathEdgeKey('frontend','baseline','kafka','shared'), build:()=>frontendToKafka('baseline'), protocol:'publish', variant:'baseline' });
    possible.push({ key:buildPathEdgeKey('kafka','shared','driver','baseline'), build:()=>kafkaToDriver('baseline'), protocol:'consume', variant:'baseline' });
    possible.push({ key:buildPathEdgeKey('driver','baseline','route','baseline'), build:()=>hPath('driver','baseline','route','baseline'), protocol:'gRPC', variant:'baseline' });

    // Sandbox
    possible.push({ key:buildPathEdgeKey('frontend','sandbox','location','sandbox'), build:()=>frontendToLocation('sandbox','sandbox'), protocol:'HTTP', variant:'sandbox' });
    possible.push({ key:buildPathEdgeKey('location','sandbox','mysql','shared'), build:()=>hPath('location','sandbox','mysql',undefined), protocol:'MySQL', variant:'sandbox' });
    possible.push({ key:buildPathEdgeKey('frontend','sandbox','kafka','shared'), build:()=>frontendToKafka('sandbox'), protocol:'publish', variant:'sandbox' });
    possible.push({ key:buildPathEdgeKey('kafka','shared','driver','sandbox'), build:()=>kafkaToDriver('sandbox'), protocol:'consume', variant:'sandbox' });
    possible.push({ key:buildPathEdgeKey('driver','sandbox','route','sandbox'), build:()=>hPath('driver','sandbox','route','sandbox'), protocol:'gRPC', variant:'sandbox' });

    // Cross-lane: sandbox frontend → baseline location
    if (path.frontend !== path.location) {
        const v = path.frontend === 'sandbox' ? 'sandbox' as const : 'baseline' as const;
        possible.push({ key:buildPathEdgeKey('frontend',path.frontend,'location',path.location),
            build:()=>frontendToLocation(path.frontend, path.location), protocol:'HTTP', variant:v });
    }
    if (path.driver !== path.route) {
        const v = path.driver === 'sandbox' ? 'sandbox' as const : 'baseline' as const;
        possible.push({ key:buildPathEdgeKey('driver',path.driver,'route',path.route),
            build:()=>hPath('driver',path.driver,'route',path.route), protocol:'gRPC', variant:v });
    }

    for (const pe of possible) {
        if (flowState.edges[pe.key]) edges.push({ id:pe.key, ...pe.build(), protocol:pe.protocol, variant:pe.variant });
    }
    return edges;
}

// === NODE SHAPES ===
function DatabaseNode({ x, y }: { x:number; y:number; color:string; isIdle:boolean }) {
    const w=DB_WIDTH, h=DB_HEIGHT, ry=7;
    const cx=x+w/2, topY=y+ry, bodyH=h-ry*2, bottomY=y+h-ry;
    return (
        <g>
            <rect x={x} y={topY} width={w} height={bodyH} fill="#1A202C" stroke={IDLE_COLOR} strokeWidth={1.5}/>
            <ellipse cx={cx} cy={bottomY} rx={w/2} ry={ry} fill="#1A202C" stroke={IDLE_COLOR} strokeWidth={1.5}/>
            <ellipse cx={cx} cy={topY} rx={w/2} ry={ry} fill="#1A202C" stroke={IDLE_COLOR} strokeWidth={1.5}/>
            <text x={cx} y={topY+bodyH/2+4} textAnchor="middle" fontSize={10} fontWeight={700} fill="#A0AEC0" fontFamily="system-ui, sans-serif">MySQL</text>
        </g>
    );
}

function KafkaNode({ x, y }: { x:number; y:number; color:string; isIdle:boolean }) {
    const w=INFRA_WIDTH, h=INFRA_HEIGHT, indent=12;
    const pts = `${x+indent},${y} ${x+w-indent},${y} ${x+w},${y+h/2} ${x+w-indent},${y+h} ${x+indent},${y+h} ${x},${y+h/2}`;
    return (
        <g>
            <polygon points={pts} fill="#1A202C" stroke={IDLE_COLOR} strokeWidth={1.5} strokeLinejoin="round"/>
            <text x={x+w/2} y={y+h/2-2} textAnchor="middle" fontSize={11} fontWeight={700} fill="#A0AEC0" fontFamily="system-ui, sans-serif">Kafka</text>
            <text x={x+w/2} y={y+h/2+10} textAnchor="middle" fontSize={8} fill="#718096" fontFamily="system-ui, sans-serif">message bus</text>
        </g>
    );
}

// === MAIN CANVAS ===
export const FlowCanvas = ({ flowState, showTopology }: FlowCanvasProps) => {
    const activeEdges = buildActiveEdges(flowState);
    const hasActive = activeEdges.length > 0;
    // Infra nodes are always shown at full visibility when topology is visible
    const showInfra = showTopology || hasActive;

    let edges: EdgeDef[] = [];
    if (hasActive && showTopology) edges = [...buildTopologyEdges(), ...activeEdges];
    else if (hasActive) edges = activeEdges;
    else if (showTopology) edges = buildTopologyEdges();

    const sbxServices = (['frontend','location','driver','route'] as ServiceId[]).filter(s => flowState.nodes[s].sandbox.status!=='idle');

    return (
        <svg viewBox={`0 0 ${VIEWBOX_WIDTH} ${VIEWBOX_HEIGHT}`} preserveAspectRatio="xMidYMid meet" width="100%" height="100%">
            <defs>
                {/* Topology arrows — visible gray */}
                <marker id="arrow-topo" markerWidth="8" markerHeight="6" refX="7" refY="3" orient="auto" markerUnits="userSpaceOnUse">
                    <path d="M 0 0.5 L 7 3 L 0 5.5" fill="none" stroke="#718096" strokeWidth={1.2} strokeLinejoin="round"/>
                </marker>
                {/* Active flow arrows — cyan */}
                <marker id="arrow-active" markerWidth="10" markerHeight="8" refX="9" refY="4" orient="auto" markerUnits="userSpaceOnUse">
                    <path d="M 0 1 L 9 4 L 0 7" fill="none" stroke="#00B5D8" strokeWidth={1.5} strokeLinejoin="round"/>
                </marker>
            </defs>

            {/* Separator line when sandbox nodes are present */}
            {sbxServices.length > 0 && (
                <line x1={0} y1={208} x2={VIEWBOX_WIDTH} y2={208} stroke="#2D3748" strokeWidth={1} strokeDasharray="4 6" opacity={0.25}/>
            )}

            {/* Edges */}
            {edges.map(e => (
                <FlowEdge key={e.id} pathD={e.pathD} protocol={e.protocol}
                    state={e.isTopology ? {status:'idle'} : (flowState.edges[e.id]||{status:'idle'})}
                    variant={e.variant} labelX={e.labelX} labelY={e.labelY}/>
            ))}

            {/* Infrastructure — always shown at full opacity when visible */}
            {showInfra && <DatabaseNode x={MYSQL_NODE.x} y={MYSQL_NODE.y} color={MYSQL_COLOR} isIdle={true}/>}
            {showInfra && <KafkaNode x={KAFKA_NODE.x} y={KAFKA_NODE.y} color={KAFKA_COLOR} isIdle={true}/>}

            {/* Baseline nodes */}
            {BASELINE_NODES.map(n => (
                <ServiceNode key={`b-${n.id}`} serviceId={n.id} variant="baseline" state={flowState.nodes[n.id].baseline} x={n.x} y={n.y} label={n.label} visible={true}/>
            ))}

            {/* Sandbox nodes — only active ones, no label */}
            {sbxServices.map(s => {
                const n = SANDBOX_NODES.find(nd=>nd.id===s)!;
                return <ServiceNode key={`s-${s}`} serviceId={s} variant="sandbox" state={flowState.nodes[s].sandbox} x={n.x} y={n.y} label={n.label} visible={true}/>;
            })}
        </svg>
    );
};
