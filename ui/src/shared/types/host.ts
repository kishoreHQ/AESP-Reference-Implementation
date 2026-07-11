/** Host Interface types — aligned with AESP / UI §6 contract (UI-TEC-07). */

export type MissionState =
  | 'running'
  | 'awaiting_approval'
  | 'queued'
  | 'succeeded'
  | 'failed'
  | 'cancelled'

export type NodeStatus = 'pending' | 'running' | 'blocked' | 'succeeded' | 'failed' | 'cancelled'

export type TrustLabel =
  | 'system'
  | 'verified'
  | 'agent'
  | 'retrieved'
  | 'untrusted'
  | 'poison-suspect'

export type ProfileId = 'platform' | 'local-first' | 'embedded'

export interface ApiError {
  code: string
  message: string
  remediation?: string
}

export interface Envelope<T> {
  data: T
  cursor?: string
  error?: ApiError
}

export interface Mission {
  id: string
  name: string
  goal: string
  state: MissionState
  requiredCapabilities: string[]
  agentsActive: number
  elapsedMs: number
  costUsd: number
  createdAt: string
  updatedAt: string
  progressPhase?: string
  workflowType?: string
}

export interface TreeNode {
  id: string
  parentId?: string
  label: string
  kind: 'agent' | 'step' | 'tool' | 'provider'
  status: NodeStatus
  runtimeId?: string
  providerId?: string
  modelId?: string
  durationMs?: number
  costUsd?: number
  children?: TreeNode[]
}

export interface ExecutionTree {
  missionId: string
  root: TreeNode
  nodeCount: number
}

export interface LogLine {
  seq: number
  ts: string
  level: 'debug' | 'info' | 'warn' | 'error'
  nodeId?: string
  message: string
  tool?: {
    name: string
    input?: unknown
    output?: unknown
    durationMs?: number
  }
}

export interface Approval {
  id: string
  missionId: string
  missionName: string
  agentId: string
  reason: string
  policy: string
  blastRadius: string
  preview: { kind: 'diff' | 'text' | 'file'; title: string; content: string }
  createdAt: string
  ageMs: number
  state: 'pending' | 'approved' | 'rejected'
}

export interface RegistryItem {
  id: string
  name: string
  kind: 'agent' | 'runtime' | 'provider' | 'tool'
  description?: string
  capabilities?: string[]
  health?: 'healthy' | 'degraded' | 'down'
  budgetUsed?: number
  budgetCap?: number
  latencyMs?: number[]
  sandbox?: string
  sessionsActive?: number
  evalScore?: number
  usageCount?: number
  lastInvocation?: string
  enabled?: boolean
  manifest?: string
}

export interface MemoryRecord {
  id: string
  text: string
  kind: 'working' | 'session' | 'semantic' | 'artifact'
  trust: TrustLabel
  missionId?: string
  agentId?: string
  createdAt: string
  pinned?: boolean
  quarantined?: boolean
}

export interface KgNode {
  id: string
  label: string
  type: string
}

export interface KgEdge {
  from: string
  to: string
  predicate: string
}

export interface Artifact {
  id: string
  name: string
  mediaType: string
  missionId: string
  digest: string
  sizeBytes: number
  createdAt: string
  version: number
  provenance: string[]
  contentPreview?: string
}

export interface EvalRun {
  id: string
  suite: string
  score: number
  baselineDelta: number
  passed: boolean
  metrics: { name: string; pass: boolean; value: number }[]
  createdAt: string
  regression?: boolean
  traceMissionId?: string
}

export interface ReplayEvent {
  seq: number
  type: string
  ts: string
  data?: Record<string, unknown>
  nodeId?: string
}

export interface Budget {
  id: string
  scope: string
  usedUsd: number
  capUsd: number
}

export interface PolicyDoc {
  id: string
  name: string
  body: string
  version: number
}

export interface Health {
  status: 'ok' | 'degraded' | 'down'
  profile: ProfileId
  version?: string
}

export type HostEventType =
  | 'mission.updated'
  | 'node.updated'
  | 'log.append'
  | 'approval.created'
  | 'approval.resolved'
  | 'artifact.created'
  | 'memory.written'
  | 'budget.threshold'
  | 'eval.completed'

export interface HostEvent {
  seq: number
  type: HostEventType | string
  ts: string
  missionId?: string
  data?: Record<string, unknown>
}
