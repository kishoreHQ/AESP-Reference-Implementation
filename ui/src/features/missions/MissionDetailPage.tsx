import { useMemo, useRef, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { useVirtualizer } from '@tanstack/react-virtual'
import { useLogs, useMission, useTree } from '@/shared/api/hooks'
import { MissionSpine } from '@/shared/spine/MissionSpine'
import { StatusDot } from '@/shared/ui/StatusDot'
import { Mono } from '@/shared/ui/Mono'
import { Badge } from '@/shared/ui/Badge'
import { Skeleton } from '@/shared/ui/Skeleton'
import { ErrorState } from '@/shared/ui/ErrorState'
import { HostApiError } from '@/shared/api/client'
import { cn } from '@/shared/lib/cn'

export function MissionDetailPage() {
  const { id = '' } = useParams()
  const mission = useMission(id)
  const tree = useTree(id)
  const [nodeId, setNodeId] = useState<string | undefined>()
  const logs = useLogs(id, nodeId)
  const [tab, setTab] = useState<'tree' | 'focus' | 'meta'>('focus')
  const [follow, setFollow] = useState(true)
  const [level, setLevel] = useState<string>('all')
  const [q, setQ] = useState('')

  const filteredLogs = useMemo(() => {
    let list = logs.data ?? []
    if (level !== 'all') list = list.filter((l) => l.level === level)
    if (q) list = list.filter((l) => l.message.toLowerCase().includes(q.toLowerCase()))
    return list
  }, [logs.data, level, q])

  const parentRef = useRef<HTMLDivElement>(null)
  const virtualizer = useVirtualizer({
    count: filteredLogs.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 28,
    overscan: 20,
  })

  if (mission.isLoading) return <div className="p-6"><Skeleton className="h-40 w-full" /></div>
  if (mission.error)
    return (
      <div className="p-6">
        <ErrorState message={(mission.error as HostApiError).message} remediation={(mission.error as HostApiError).remediation} />
      </div>
    )
  const m = mission.data!

  return (
    <div className="flex h-full min-h-0 flex-col anim-in">
      <header className="flex flex-wrap items-center justify-between gap-2 border-b border-line bg-bg-1/50 px-4 py-3.5">
        <div>
          <div className="text-[12px] text-ink-1">
            <Link to="/" className="text-accent">
              Missions
            </Link>{' '}
            / <Mono className="text-ink-0">{m.id}</Mono>
          </div>
          <h1 className="font-display text-[18px] font-bold tracking-tight">{m.name}</h1>
        </div>
        <StatusDot status={m.state} label={m.state.replaceAll('_', ' ')} />
      </header>

      {/* mobile tabs */}
      <div className="flex border-b border-line md:hidden">
        {(['tree', 'focus', 'meta'] as const).map((t) => (
          <button
            key={t}
            className={cn('min-h-11 flex-1 capitalize', tab === t && 'border-b-2 border-accent text-accent')}
            onClick={() => setTab(t)}
          >
            {t}
          </button>
        ))}
      </div>

      <div className="grid min-h-0 flex-1 md:grid-cols-[280px_1fr_260px]">
        <div className={cn('min-h-0 border-r border-line', tab !== 'tree' && 'hidden md:block')}>
          {tree.data ? (
            <MissionSpine tree={tree.data} onSelectNode={setNodeId} compact />
          ) : (
            <Skeleton className="m-3 h-40" />
          )}
        </div>

        <section className={cn('flex min-h-0 flex-col', tab !== 'focus' && 'hidden md:flex')}>
          <div className="flex flex-wrap items-center gap-2 border-b border-line px-3 py-2">
            <span className="text-[13px] text-ink-1">Focus {nodeId ? <Mono>{nodeId}</Mono> : 'root'}</span>
            <label className="ml-auto flex items-center gap-1 text-[12px] text-ink-1">
              <input type="checkbox" checked={follow} onChange={(e) => setFollow(e.target.checked)} /> follow tail
            </label>
            <select
              className="rounded-[6px] border border-line bg-bg-1 px-2 py-1 text-[12px]"
              value={level}
              onChange={(e) => setLevel(e.target.value)}
            >
              <option value="all">all levels</option>
              <option value="info">info</option>
              <option value="warn">warn</option>
              <option value="error">error</option>
            </select>
            <input
              className="rounded-[6px] border border-line bg-bg-1 px-2 py-1 text-[12px]"
              placeholder="Search in log"
              value={q}
              onChange={(e) => setQ(e.target.value)}
            />
          </div>
          <div ref={parentRef} className="min-h-0 flex-1 overflow-auto bg-bg-0 font-mono text-[12.5px]">
            <div style={{ height: virtualizer.getTotalSize(), position: 'relative' }}>
              {virtualizer.getVirtualItems().map((vi) => {
                const line = filteredLogs[vi.index]
                return (
                  <div
                    key={line.seq}
                    className="absolute left-0 top-0 flex w-full gap-2 border-b border-line/30 px-3"
                    style={{ height: vi.size, transform: `translateY(${vi.start}px)` }}
                  >
                    <span className="text-ink-1">{line.level}</span>
                    <span className="min-w-0 flex-1 truncate text-ink-0">{line.message}</span>
                    {line.tool && (
                      <details className="text-accent">
                        <summary>{line.tool.name}</summary>
                        <pre className="whitespace-pre-wrap text-ink-1">{JSON.stringify(line.tool, null, 2)}</pre>
                      </details>
                    )}
                  </div>
                )
              })}
            </div>
          </div>
        </section>

        <aside className={cn('min-h-0 overflow-auto border-l border-line bg-bg-1 p-3', tab !== 'meta' && 'hidden md:block')}>
          <h2 className="mb-2 font-display text-[14px] font-semibold">Meta</h2>
          <dl className="space-y-2 text-[13px]">
            <div>
              <dt className="text-ink-1">Cost</dt>
              <dd className="font-mono">${m.costUsd.toFixed(4)}</dd>
            </div>
            <div>
              <dt className="text-ink-1">Elapsed</dt>
              <dd className="font-mono">{(m.elapsedMs / 1000).toFixed(1)}s</dd>
            </div>
            <div>
              <dt className="text-ink-1">Capabilities</dt>
              <dd className="mt-1 flex flex-wrap gap-1">
                {m.requiredCapabilities.map((c) => (
                  <Badge key={c}>{c}</Badge>
                ))}
              </dd>
            </div>
            <div>
              <dt className="text-ink-1">Tree nodes</dt>
              <dd className="font-mono">{tree.data?.nodeCount ?? '—'}</dd>
            </div>
            <div>
              <dt className="text-ink-1">Goal</dt>
              <dd className="text-ink-0">{m.goal}</dd>
            </div>
          </dl>
          <Link className="mt-4 inline-block text-[13px] text-accent" to={`/replay/${m.id}`}>
            Open replay →
          </Link>
        </aside>
      </div>
    </div>
  )
}
