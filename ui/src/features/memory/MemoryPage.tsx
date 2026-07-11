import { useState } from 'react'
import { useKg, useMemorySearch } from '@/shared/api/hooks'
import { apiPost } from '@/shared/api/client'
import { Card } from '@/shared/ui/Card'
import { TrustChip } from '@/shared/ui/TrustChip'
import { Mono } from '@/shared/ui/Mono'
import { Button } from '@/shared/ui/Button'
import { useUi } from '@/shared/store/ui'
import { useQueryClient } from '@tanstack/react-query'
import { qk } from '@/shared/api/hooks'

export function MemoryPage() {
  const [q, setQ] = useState('')
  const mem = useMemorySearch(q)
  const kg = useKg()
  const toast = useUi((s) => s.pushToast)
  const qc = useQueryClient()
  const [view, setView] = useState<'list' | 'graph'>('list')

  async function curate(id: string, action: 'pin' | 'quarantine') {
    await apiPost(`/v1/memory/${id}/${action}`, {})
    toast(action === 'pin' ? 'Memory pinned' : 'Memory quarantined', action === 'pin' ? 'ok' : 'warn')
    qc.invalidateQueries({ queryKey: qk.memory(q) })
  }

  return (
    <div className="page anim-in">
      <h1 className="page-title mb-4">Memory & knowledge</h1>
      <div className="mb-4 flex flex-wrap gap-2">
        <input
          className="min-h-11 min-w-[240px] flex-1 rounded-[6px] border border-line bg-bg-1 px-3"
          placeholder="Search memory…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <Button variant={view === 'list' ? 'primary' : 'secondary'} onClick={() => setView('list')}>
          List
        </Button>
        <Button className="hidden md:inline-flex" variant={view === 'graph' ? 'primary' : 'secondary'} onClick={() => setView('graph')}>
          Graph
        </Button>
      </div>

      {view === 'list' && (
        <div className="space-y-3">
          {mem.data?.map((r) => (
            <Card key={r.id}>
              <div className="mb-2 flex flex-wrap items-center gap-2">
                <TrustChip trust={r.trust} />
                <Mono>{r.id}</Mono>
                <span className="text-[12px] text-ink-1">{r.kind}</span>
                {r.pinned && <span className="text-[12px] text-accent">pinned</span>}
                {r.quarantined && <span className="text-[12px] text-fail">quarantined</span>}
              </div>
              <p className="text-[14px]">{r.text}</p>
              <div className="mt-2 font-mono text-[12px] text-ink-1">
                {r.missionId && <span>mission {r.missionId} · </span>}
                {r.agentId && <span>agent {r.agentId}</span>}
              </div>
              <div className="mt-3 flex gap-2">
                <Button variant="secondary" onClick={() => curate(r.id, 'pin')}>
                  Pin
                </Button>
                <Button variant="danger" onClick={() => curate(r.id, 'quarantine')}>
                  Quarantine
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}

      {view === 'graph' && (
        <Card className="min-h-[360px]">
          <p className="mb-3 text-[13px] text-ink-1">Force-layout style canvas (≤300 nodes). Mobile uses list only.</p>
          <svg viewBox="0 0 600 320" className="h-[320px] w-full rounded-[6px] bg-bg-0">
            {kg.data?.edges.map((e, i) => {
              const from = kg.data!.nodes.findIndex((n) => n.id === e.from)
              const to = kg.data!.nodes.findIndex((n) => n.id === e.to)
              const x1 = 100 + from * 180
              const x2 = 100 + to * 180
              return <line key={i} x1={x1} y1={160} x2={x2} y2={160} stroke="var(--line)" />
            })}
            {kg.data?.nodes.map((n, i) => {
              const x = 100 + i * 180
              return (
                <g key={n.id}>
                  <circle cx={x} cy={160} r={28} fill="var(--bg-2)" stroke="var(--accent)" />
                  <text x={x} y={165} textAnchor="middle" fill="var(--ink-0)" fontSize="11">
                    {n.label}
                  </text>
                </g>
              )
            })}
          </svg>
        </Card>
      )}
    </div>
  )
}
