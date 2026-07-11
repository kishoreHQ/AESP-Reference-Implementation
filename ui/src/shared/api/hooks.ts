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
