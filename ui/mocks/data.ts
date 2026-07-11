import type {
  Approval,
  Artifact,
  Budget,
  EvalRun,
  ExecutionTree,
  LogLine,
  MemoryRecord,
  Mission,
  PolicyDoc,
  RegistryItem,
  ReplayEvent,
  TreeNode,
  KgNode,
  KgEdge,
} from '../src/shared/types/host'

export const missions: Mission[] = [
  {
    id: 'mis_demo_running',
    name: 'Ship auth service',
    goal: 'Implement and validate authentication service',
    state: 'running',
    requiredCapabilities: ['coding', 'tools', 'reasoning'],
    agentsActive: 3,
    elapsedMs: 184000,
    costUsd: 0.42,
    createdAt: new Date(Date.now() - 184000).toISOString(),
    updatedAt: new Date().toISOString(),
    progressPhase: 'execute',
    workflowType: 'build-ship',
  },
  {
    id: 'mis_await_hitl',
    name: 'Production deploy gate',
    goal: 'Roll out v1.4.0 to production',
    state: 'awaiting_approval',
    requiredCapabilities: ['tools'],
    agentsActive: 1,
    elapsedMs: 92000,
    costUsd: 0.11,
    createdAt: new Date(Date.now() - 92000).toISOString(),
    updatedAt: new Date().toISOString(),
    progressPhase: 'hitl',
    workflowType: 'deploy',
  },
  {
    id: 'mis_queued',
    name: 'Nightly eval suite',
    goal: 'Run regression campaign',
    state: 'queued',
    requiredCapabilities: ['reasoning'],
    agentsActive: 0,
    elapsedMs: 0,
    costUsd: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    progressPhase: 'queued',
    workflowType: 'eval',
  },
  {
    id: 'mis_recent_ok',
    name: 'Memory curation pass',
    goal: 'Relabel and quarantine suspect memories',
    state: 'succeeded',
    requiredCapabilities: ['reasoning', 'tools'],
    agentsActive: 0,
    elapsedMs: 61000,
    costUsd: 0.07,
    createdAt: new Date(Date.now() - 3600000).toISOString(),
    updatedAt: new Date(Date.now() - 3500000).toISOString(),
    progressPhase: 'done',
    workflowType: 'memory',
  },
  {
    id: 'mis_recent_fail',
    name: 'Provider failover drill',
    goal: 'Prove capability failover path',
    state: 'failed',
    requiredCapabilities: ['coding', 'tools'],
    agentsActive: 0,
    elapsedMs: 22000,
    costUsd: 0.02,
    createdAt: new Date(Date.now() - 7200000).toISOString(),
    updatedAt: new Date(Date.now() - 7100000).toISOString(),
    progressPhase: 'failed',
    workflowType: 'failover',
  },
]

function deepTree(depth: number, breadth: number, prefix = 'n', status: TreeNode['status'] = 'succeeded'): TreeNode {
  const id = prefix
  const node: TreeNode = {
    id,
    label: depth === 0 ? 'orchestrator' : `agent.${prefix}`,
    kind: depth === 0 ? 'agent' : depth % 3 === 0 ? 'tool' : 'agent',
    status: depth < 2 ? 'running' : status,
    runtimeId: 'runtime.generic-loop',
    providerId: depth % 2 ? 'provider.mock-local' : 'provider.mock-remote',
    modelId: depth % 2 ? 'mock-local-small' : 'mock-remote-large',
    durationMs: 1000 + depth * 120,
    costUsd: 0.001 * (depth + 1),
    children: [],
  }
  if (depth >= 4) return node
  for (let i = 0; i < breadth; i++) {
    const child = deepTree(depth + 1, Math.max(1, breadth - 1), `${prefix}.${i}`, i === 0 && depth < 2 ? 'running' : 'succeeded')
    child.parentId = id
    node.children!.push(child)
  }
  return node
}

/** ~500+ nodes for Spine perf gate */
export function bigTree(missionId: string): ExecutionTree {
  const root = deepTree(0, 4, 'root')
  let count = 0
  const walk = (n: TreeNode) => {
    count++
    n.children?.forEach(walk)
  }
  walk(root)
  return { missionId, root, nodeCount: count }
}

export const logs: LogLine[] = Array.from({ length: 200 }, (_, i) => ({
  seq: i + 1,
  ts: new Date(Date.now() - (200 - i) * 1000).toISOString(),
  level: i % 17 === 0 ? 'warn' : i % 41 === 0 ? 'error' : 'info',
  nodeId: i % 3 === 0 ? 'root.0' : 'root',
  message:
    i % 5 === 0
      ? `tool invocation completed step=${i}`
      : `mission progress phase=execute line=${i}`,
  tool:
    i % 9 === 0
      ? {
          name: 'workspace.write',
          input: { path: `src/file_${i}.ts`, content: 'export const x = 1' },
          output: { ok: true },
          durationMs: 40 + (i % 30),
        }
      : undefined,
}))

export let approvals: Approval[] = [
  {
    id: 'hitl_1',
    missionId: 'mis_await_hitl',
    missionName: 'Production deploy gate',
    agentId: 'agent.deployer',
    reason: 'Destructive rollout to production requires human approval',
    policy: 'policy.destructive.requires_hitl',
    blastRadius: 'production / auth-service / 12 replicas',
    preview: {
      kind: 'diff',
      title: 'deploy.yaml',
      content: `- image: auth:1.3.9\n+ image: auth:1.4.0\n  replicas: 12\n  strategy: rolling`,
    },
    createdAt: new Date(Date.now() - 45000).toISOString(),
    ageMs: 45000,
    state: 'pending',
  },
  {
    id: 'hitl_2',
    missionId: 'mis_demo_running',
    missionName: 'Ship auth service',
    agentId: 'agent.security',
    reason: 'Network egress to external IdP',
    policy: 'policy.egress.review',
    blastRadius: 'outbound HTTPS to idp.example',
    preview: {
      kind: 'text',
      title: 'egress request',
      content: 'Allow tool net.fetch to https://idp.example/oauth/token for mission mis_demo_running',
    },
    createdAt: new Date(Date.now() - 120000).toISOString(),
    ageMs: 120000,
    state: 'pending',
  },
]

export const registry: Record<string, RegistryItem[]> = {
  agents: [
    {
      id: 'agent.default',
      name: 'Builder',
      kind: 'agent',
      capabilities: ['coding', 'tools'],
      evalScore: 0.91,
      enabled: true,
      description: 'Default builder agent',
      manifest: 'role: builder\ncapabilities: [coding, tools]',
    },
    {
      id: 'agent.reviewer',
      name: 'Reviewer',
      kind: 'agent',
      capabilities: ['reasoning', 'coding'],
      evalScore: 0.88,
      enabled: true,
    },
  ],
  runtimes: [
    {
      id: 'runtime.generic-loop',
      name: 'Generic loop',
      kind: 'runtime',
      capabilities: ['tools', 'streaming', 'coding'],
      sandbox: 'process',
      sessionsActive: 2,
      enabled: true,
      manifest: 'apiVersion: aesp.runtime/v1\nkind: RuntimePlugin\nmetadata:\n  id: runtime.generic-loop',
    },
  ],
  providers: [
    {
      id: 'provider.mock-remote',
      name: 'Mock remote',
      kind: 'provider',
      capabilities: ['coding', 'tools', 'vision', 'reasoning', 'streaming', 'planning'],
      health: 'healthy',
      budgetUsed: 12.4,
      budgetCap: 100,
      latencyMs: [40, 42, 38, 55, 44, 41, 39],
      enabled: true,
    },
    {
      id: 'provider.mock-local',
      name: 'Mock local',
      kind: 'provider',
      capabilities: ['coding', 'tools', 'local', 'reasoning', 'planning'],
      health: 'healthy',
      budgetUsed: 0.2,
      budgetCap: 0,
      latencyMs: [12, 11, 13, 12, 10],
      enabled: true,
    },
  ],
  tools: [
    {
      id: 'echo',
      name: 'echo',
      kind: 'tool',
      description: 'Echo message',
      usageCount: 120,
      lastInvocation: new Date().toISOString(),
      enabled: true,
    },
    {
      id: 'memory.write',
      name: 'memory.write',
      kind: 'tool',
      description: 'Write unified memory with trust label',
      usageCount: 44,
      enabled: true,
    },
    {
      id: 'kg.upsert',
      name: 'kg.upsert',
      kind: 'tool',
      usageCount: 18,
      enabled: true,
    },
  ],
}

export const memory: MemoryRecord[] = [
  {
    id: 'mem_1',
    text: 'Auth service uses RS256 for JWT verification',
    kind: 'session',
    trust: 'agent',
    missionId: 'mis_demo_running',
    agentId: 'agent.default',
    createdAt: new Date().toISOString(),
  },
  {
    id: 'mem_2',
    text: 'Retrieved RFC excerpt on OAuth — treat as untrusted until verified',
    kind: 'semantic',
    trust: 'retrieved',
    missionId: 'mis_demo_running',
    createdAt: new Date().toISOString(),
  },
  {
    id: 'mem_3',
    text: 'Human-approved production cutover checklist',
    kind: 'session',
    trust: 'verified',
    missionId: 'mis_await_hitl',
    agentId: 'operator',
    createdAt: new Date().toISOString(),
    pinned: true,
  },
  {
    id: 'mem_4',
    text: 'Suspicious tool output claiming admin credentials',
    kind: 'working',
    trust: 'poison-suspect',
    missionId: 'mis_recent_fail',
    createdAt: new Date().toISOString(),
    quarantined: true,
  },
]

export const kgNodes: KgNode[] = [
  { id: 'svc.auth', label: 'auth-service', type: 'service' },
  { id: 'svc.db', label: 'postgres', type: 'datastore' },
  { id: 'svc.idp', label: 'idp', type: 'external' },
]
export const kgEdges: KgEdge[] = [
  { from: 'svc.auth', to: 'svc.db', predicate: 'depends_on' },
  { from: 'svc.auth', to: 'svc.idp', predicate: 'federates_with' },
]

export const artifacts: Artifact[] = [
  {
    id: 'art_1',
    name: 'plan.json',
    mediaType: 'application/json',
    missionId: 'mis_demo_running',
    digest: 'sha256:aaa111',
    sizeBytes: 2048,
    createdAt: new Date().toISOString(),
    version: 1,
    provenance: ['planner', 'mis_demo_running'],
    contentPreview: '{\n  "goal": "Implement auth service",\n  "revision": 1\n}',
  },
  {
    id: 'art_2',
    name: 'mission-report.md',
    mediaType: 'text/markdown',
    missionId: 'mis_recent_ok',
    digest: 'sha256:bbb222',
    sizeBytes: 4096,
    createdAt: new Date().toISOString(),
    version: 2,
    provenance: ['docgen', 'mis_recent_ok'],
    contentPreview: '# Mission Report\n\nAll memory labels verified.\n',
  },
]

export const evals: EvalRun[] = [
  {
    id: 'eval_1',
    suite: 'core-agent-loop',
    score: 0.94,
    baselineDelta: 0.02,
    passed: true,
    metrics: [
      { name: 'task_success', pass: true, value: 0.96 },
      { name: 'tool_precision', pass: true, value: 0.91 },
    ],
    createdAt: new Date().toISOString(),
  },
  {
    id: 'eval_2',
    suite: 'hitl-safety',
    score: 0.81,
    baselineDelta: -0.05,
    passed: false,
    regression: true,
    metrics: [
      { name: 'no_auto_approve', pass: true, value: 1 },
      { name: 'preview_present', pass: false, value: 0.7 },
    ],
    createdAt: new Date().toISOString(),
    traceMissionId: 'mis_await_hitl',
  },
]

export const replayEvents: ReplayEvent[] = [
  { seq: 1, type: 'aesp.control.mission.accepted', ts: new Date().toISOString() },
  { seq: 2, type: 'aesp.session.opened', ts: new Date().toISOString() },
  { seq: 3, type: 'aesp.control.route.selected', ts: new Date().toISOString(), data: { providerId: 'provider.mock-local' } },
  { seq: 4, type: 'aesp.runtime.completed', ts: new Date().toISOString(), nodeId: 'root' },
]

export const budgets: Budget[] = [
  { id: 'b1', scope: 'provider.mock-remote', usedUsd: 12.4, capUsd: 100 },
  { id: 'b2', scope: 'tenant.default', usedUsd: 18.1, capUsd: 250 },
]

export const policies: PolicyDoc[] = [
  {
    id: 'pol_hitl',
    name: 'Destructive requires HITL',
    body: 'effect: require_approval\nwhen: sideEffect in [destructive, admin]',
    version: 3,
  },
]
