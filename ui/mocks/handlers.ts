import { http, HttpResponse } from 'msw'
import {
  approvals,
  artifacts,
  bigTree,
  budgets,
  evals,
  kgEdges,
  kgNodes,
  logs,
  memory,
  missions,
  policies,
  registry,
  replayEvents,
} from './data'

const ok = <T,>(data: T) => HttpResponse.json({ data })

export const handlers = [
  http.get('/api/v1/health', () =>
    ok({ status: 'ok', profile: 'local-first', version: 'ui-mock-1.0.0' }),
  ),
  http.get('/api/v1/missions', ({ request }) => {
    const url = new URL(request.url)
    const state = url.searchParams.get('state')
    const list = state ? missions.filter((m) => m.state === state) : missions
    return ok(list)
  }),
  http.get('/api/v1/missions/:id', ({ params }) => {
    const m = missions.find((x) => x.id === params.id)
    if (!m) return HttpResponse.json({ data: null, error: { code: 'not_found', message: 'Mission not found', remediation: 'Check mission id.' } }, { status: 404 })
    return ok(m)
  }),
  http.post('/api/v1/missions', async ({ request }) => {
    const body = (await request.json()) as { name?: string; goal?: string; requiredCapabilities?: string[] }
    const m = {
      id: `mis_${Date.now()}`,
      name: body.name || 'New mission',
      goal: body.goal || '',
      state: 'running' as const,
      requiredCapabilities: body.requiredCapabilities || ['coding'],
      agentsActive: 1,
      elapsedMs: 0,
      costUsd: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      progressPhase: 'plan',
    }
    missions.unshift(m)
    return ok(m)
  }),
  http.post('/api/v1/missions/:id/cancel', ({ params }) => {
    const m = missions.find((x) => x.id === params.id)
    if (m) m.state = 'cancelled'
    return ok(m)
  }),
  http.get('/api/v1/missions/:id/tree', ({ params }) => ok(bigTree(String(params.id)))),
  http.get('/api/v1/missions/:id/logs', () => ok(logs)),
  http.get('/api/v1/approvals', () => ok(approvals.filter((a) => a.state === 'pending'))),
  http.post('/api/v1/approvals/:id/decision', async ({ params, request }) => {
    const body = (await request.json()) as { decision: string; comment?: string }
    const a = approvals.find((x) => x.id === params.id)
    if (!a)
      return HttpResponse.json(
        { data: null, error: { code: 'not_found', message: 'Approval not found', remediation: 'Refresh inbox.' } },
        { status: 404 },
      )
    a.state = body.decision === 'approve' ? 'approved' : 'rejected'
    const m = missions.find((x) => x.id === a.missionId)
    if (m && body.decision === 'approve') m.state = 'running'
    if (m && body.decision === 'reject') m.state = 'failed'
    return ok(a)
  }),
  http.get('/api/v1/registry/:kind', ({ params }) => ok(registry[String(params.kind)] || [])),
  http.get('/api/v1/memory/search', ({ request }) => {
    const q = new URL(request.url).searchParams.get('q')?.toLowerCase() || ''
    return ok(memory.filter((m) => !q || m.text.toLowerCase().includes(q)))
  }),
  http.get('/api/v1/memory/kg', () => ok({ nodes: kgNodes, edges: kgEdges })),
  http.post('/api/v1/memory/:id/:action', ({ params }) => {
    const rec = memory.find((m) => m.id === params.id)
    if (!rec) return HttpResponse.json({ data: null, error: { code: 'not_found', message: 'Memory not found' } }, { status: 404 })
    if (params.action === 'pin') rec.pinned = true
    if (params.action === 'quarantine') {
      rec.quarantined = true
      rec.trust = 'poison-suspect'
    }
    return ok(rec)
  }),
  http.get('/api/v1/artifacts', ({ request }) => {
    const mission = new URL(request.url).searchParams.get('mission')
    return ok(mission ? artifacts.filter((a) => a.missionId === mission) : artifacts)
  }),
  http.get('/api/v1/artifacts/:id/versions', ({ params }) =>
    ok(artifacts.filter((a) => a.id === params.id || a.name === params.id)),
  ),
  http.get('/api/v1/evaluations', () => ok(evals)),
  http.get('/api/v1/replay/:runId/events', () => ok(replayEvents)),
  http.get('/api/v1/budgets', () => ok(budgets)),
  http.get('/api/v1/policies', () => ok(policies)),
  http.put('/api/v1/policies/:id', async ({ params, request }) => {
    const body = (await request.json()) as (typeof policies)[0]
    const idx = policies.findIndex((p) => p.id === params.id)
    if (idx >= 0) policies[idx] = { ...policies[idx], ...body, version: policies[idx].version + 1 }
    return ok(policies[idx])
  }),
  http.post('/api/v1/credentials', async () =>
    ok({ id: 'cred_new', status: 'stored' }),
  ),
]
