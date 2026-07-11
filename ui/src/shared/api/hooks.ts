import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost, apiPut } from './client'
import type {
  Approval,
  Artifact,
  Budget,
  Envelope,
  EvalRun,
  ExecutionTree,
  Health,
  LogLine,
  MemoryRecord,
  Mission,
  PolicyDoc,
  RegistryItem,
  ReplayEvent,
  KgEdge,
  KgNode,
} from '@/shared/types/host'

export const qk = {
  missions: (state?: string) => ['missions', state ?? 'all'] as const,
  mission: (id: string) => ['mission', id] as const,
  tree: (id: string) => ['tree', id] as const,
  logs: (id: string, node?: string) => ['logs', id, node ?? ''] as const,
  approvals: ['approvals'] as const,
  registry: (k: string) => ['registry', k] as const,
  memory: (q: string) => ['memory', q] as const,
  kg: ['kg'] as const,
  artifacts: (m?: string) => ['artifacts', m ?? ''] as const,
  evals: ['evals'] as const,
  replay: (id: string) => ['replay', id] as const,
  budgets: ['budgets'] as const,
  policies: ['policies'] as const,
  health: ['health'] as const,
}

export function useHealth() {
  return useQuery({ queryKey: qk.health, queryFn: () => apiGet<Health>('/v1/health'), refetchInterval: 10000 })
}

export function useMissions(state?: string) {
  return useQuery({
    queryKey: qk.missions(state),
    queryFn: () => apiGet<Mission[]>('/v1/missions', { state }),
  })
}

export function useMission(id: string) {
  return useQuery({
    queryKey: qk.mission(id),
    queryFn: () => apiGet<Mission>(`/v1/missions/${id}`),
    enabled: !!id,
  })
}

export function useTree(id: string) {
  return useQuery({
    queryKey: qk.tree(id),
    queryFn: () => apiGet<ExecutionTree>(`/v1/missions/${id}/tree`),
    enabled: !!id,
    refetchInterval: 3000,
  })
}

export function useLogs(id: string, node?: string) {
  return useQuery({
    queryKey: qk.logs(id, node),
    queryFn: () => apiGet<LogLine[]>(`/v1/missions/${id}/logs`, { node }),
    enabled: !!id,
    refetchInterval: 2000,
  })
}

export function useApprovals() {
  return useQuery({
    queryKey: qk.approvals,
    queryFn: () => apiGet<Approval[]>('/v1/approvals', { state: 'pending' }),
    refetchInterval: 2000,
  })
}

export function useDecision() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (p: { id: string; decision: 'approve' | 'reject'; comment?: string }) =>
      apiPost<Approval>(`/v1/approvals/${p.id}/decision`, {
        decision: p.decision,
        comment: p.comment,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: qk.approvals })
      qc.invalidateQueries({ queryKey: ['missions'] })
      qc.invalidateQueries({ queryKey: ['tree'] })
    },
  })
}

export function useLaunchMission() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: { name: string; goal: string; requiredCapabilities: string[] }) =>
      apiPost<Mission>('/v1/missions', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['missions'] }),
  })
}

export function useRegistry(kind: 'agents' | 'runtimes' | 'providers' | 'tools') {
  return useQuery({
    queryKey: qk.registry(kind),
    queryFn: () => apiGet<RegistryItem[]>(`/v1/registry/${kind}`),
  })
}

export function useMemorySearch(q: string) {
  return useQuery({
    queryKey: qk.memory(q),
    queryFn: () => apiGet<MemoryRecord[]>('/v1/memory/search', { q }),
  })
}

export function useKg() {
  return useQuery({
    queryKey: qk.kg,
    queryFn: () => apiGet<{ nodes: KgNode[]; edges: KgEdge[] }>('/v1/memory/kg'),
  })
}

export function useArtifacts(mission?: string) {
  return useQuery({
    queryKey: qk.artifacts(mission),
    queryFn: () => apiGet<Artifact[]>('/v1/artifacts', { mission }),
  })
}

export function useEvaluations() {
  return useQuery({ queryKey: qk.evals, queryFn: () => apiGet<EvalRun[]>('/v1/evaluations') })
}

export function useReplay(runId: string) {
  return useQuery({
    queryKey: qk.replay(runId),
    queryFn: () => apiGet<ReplayEvent[]>(`/v1/replay/${runId}/events`),
    enabled: !!runId,
  })
}

export function useBudgets() {
  return useQuery({ queryKey: qk.budgets, queryFn: () => apiGet<Budget[]>('/v1/budgets') })
}

export function usePolicies() {
  return useQuery({ queryKey: qk.policies, queryFn: () => apiGet<PolicyDoc[]>('/v1/policies') })
}

export function useSavePolicy() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (p: PolicyDoc) => apiPut<PolicyDoc>(`/v1/policies/${p.id}`, p),
    onSuccess: () => qc.invalidateQueries({ queryKey: qk.policies }),
  })
}

// silence unused Envelope import for some TS configs
export type { Envelope }


// ——— Feature pack: connections & command deck ———

export function useProbe() {
  return useQuery({
    queryKey: ['connections', 'probe'],
    queryFn: () => apiGet<any[]>('/v1/connections/probe'),
  })
}

export function useConnections() {
  return useQuery({
    queryKey: ['connections'],
    queryFn: () => apiGet<any[]>('/v1/connections'),
    refetchInterval: 4000,
  })
}

export function useRegisterConnection() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: any) => apiPost<any>('/v1/connections', body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['connections'] })
      qc.invalidateQueries({ queryKey: ['registry'] })
    },
  })
}

export function useSessions() {
  return useQuery({
    queryKey: ['sessions'],
    queryFn: () => apiGet<any[]>('/v1/sessions'),
    refetchInterval: 2000,
  })
}

export function useCreateSession() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: any) => apiPost<any>('/v1/sessions', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['sessions'] }),
  })
}

export function useSessionMessage() {
  return useMutation({
    mutationFn: (p: { id: string; text: string }) =>
      apiPost<any>(`/v1/sessions/${p.id}/message`, { text: p.text }),
  })
}

export function useStopSession() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => apiPost<any>(`/v1/sessions/${id}/stop`, {}),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['sessions'] }),
  })
}

export function useBoards() {
  return useQuery({ queryKey: ['boards'], queryFn: () => apiGet<any[]>('/v1/boards') })
}

export function useTasks(board = 'board_default') {
  return useQuery({
    queryKey: ['tasks', board],
    queryFn: () => apiGet<any[]>('/v1/tasks', { board }),
    refetchInterval: 2000,
  })
}

export function useMoveTask() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (p: { id: string; column: string; assignee?: string }) =>
      apiPost<any>(`/v1/tasks/${p.id}`, { column: p.column, assignee: p.assignee }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tasks'] }),
  })
}

export function useClaimTask() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (p: { agentId: string }) =>
      apiPost<any>('/v1/tasks/_/claim', { agentId: p.agentId }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tasks'] }),
  })
}

export function useRoutines() {
  return useQuery({
    queryKey: ['routines'],
    queryFn: () => apiGet<any[]>('/v1/routines'),
    refetchInterval: 5000,
  })
}

export function useCreateRoutine() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (body: any) => apiPost<any>('/v1/routines', body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['routines'] }),
  })
}

export function useGoals() {
  return useQuery({ queryKey: ['goals'], queryFn: () => apiGet<any[]>('/v1/goals') })
}

export function useJournal() {
  return useQuery({ queryKey: ['journal'], queryFn: () => apiGet<any[]>('/v1/journal') })
}

export function useAddJournal() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (text: string) => apiPost<any>('/v1/journal', { text }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['journal'] }),
  })
}

export function useAgentAnalytics(id: string) {
  return useQuery({
    queryKey: ['analytics', id],
    queryFn: () => apiGet<any>(`/v1/analytics/agents/${id}`),
    enabled: !!id,
  })
}
