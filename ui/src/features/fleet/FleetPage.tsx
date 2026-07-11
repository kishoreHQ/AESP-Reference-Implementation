import { useState } from 'react'
import { useRegistry } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Badge } from '@/shared/ui/Badge'
import { StatusDot } from '@/shared/ui/StatusDot'
import { Mono } from '@/shared/ui/Mono'
import { Skeleton } from '@/shared/ui/Skeleton'
import type { RegistryItem } from '@/shared/types/host'

const tabs = ['agents', 'runtimes', 'providers', 'tools'] as const

export function FleetPage() {
  const [tab, setTab] = useState<(typeof tabs)[number]>('providers')
  const [selected, setSelected] = useState<RegistryItem | null>(null)
  const q = useRegistry(tab)

  return (
    <div className="page anim-in">
      <h1 className="page-title mb-4">Fleet</h1>
      <div className="mb-4 flex gap-1 border-b border-line">
        {tabs.map((t) => (
          <button
            key={t}
            className={`min-h-11 px-3 capitalize ${tab === t ? 'border-b-2 border-accent text-accent' : 'text-ink-1'}`}
            onClick={() => {
              setTab(t)
              setSelected(null)
            }}
          >
            {t}
          </button>
        ))}
      </div>
      {q.isLoading && <Skeleton className="h-24 w-full" />}
      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        {q.data?.map((item) => (
          <Card key={item.id} className="cursor-pointer hover:bg-bg-2" onClick={() => setSelected(item)} role="button" tabIndex={0}>
            <div className="mb-2 flex items-center justify-between">
              <div className="font-display font-semibold">{item.name}</div>
              {item.health && <StatusDot status={item.health} label={item.health} />}
            </div>
            <Mono>{item.id}</Mono>
            <div className="mt-2 flex flex-wrap gap-1">
              {item.capabilities?.map((c) => (
                <Badge key={c}>{c}</Badge>
              ))}
            </div>
            {item.budgetCap != null && item.budgetCap > 0 && (
              <div className="mt-3">
                <div className="mb-1 flex justify-between font-mono text-[11px] text-ink-1">
                  <span>budget</span>
                  <span>
                    ${item.budgetUsed?.toFixed(1)} / ${item.budgetCap}
                  </span>
                </div>
                <div className="h-1.5 overflow-hidden rounded bg-bg-2">
                  <div
                    className="h-full bg-accent"
                    style={{ width: `${Math.min(100, ((item.budgetUsed ?? 0) / item.budgetCap) * 100)}%` }}
                  />
                </div>
              </div>
            )}
            {item.latencyMs && (
              <div className="mt-2 flex h-6 items-end gap-0.5">
                {item.latencyMs.map((v, i) => (
                  <div key={i} className="w-1.5 bg-live/70" style={{ height: `${Math.min(100, v)}%` }} title={`${v}ms`} />
                ))}
              </div>
            )}
          </Card>
        ))}
      </div>

      {selected && (
        <div className="fixed inset-y-0 right-0 z-40 w-full max-w-md border-l border-line bg-bg-1 p-4 shadow-none">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="font-display text-[18px] font-semibold">{selected.name}</h2>
            <button onClick={() => setSelected(null)} aria-label="Close">
              ×
            </button>
          </div>
          <Mono>{selected.id}</Mono>
          <p className="mt-2 text-[13px] text-ink-1">{selected.description}</p>
          <h3 className="mt-4 text-[12px] uppercase tracking-wide text-ink-1">Manifest</h3>
          <pre className="mt-1 max-h-64 overflow-auto rounded-[6px] border border-line bg-bg-0 p-2 font-mono text-[12px]">
            {selected.manifest || '—'}
          </pre>
        </div>
      )}
    </div>
  )
}
