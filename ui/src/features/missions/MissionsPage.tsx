import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { useLaunchMission, useMissions } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Button } from '@/shared/ui/Button'
import { Badge } from '@/shared/ui/Badge'
import { StatusPill } from '@/shared/ui/StatusDot'
import { Mono } from '@/shared/ui/Mono'
import { EmptyState } from '@/shared/ui/EmptyState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Skeleton } from '@/shared/ui/Skeleton'
import { useUi } from '@/shared/store/ui'
import type { Mission, MissionState } from '@/shared/types/host'
import { detectProfile } from '@/profiles/detect'
import { HostApiError } from '@/shared/api/client'
import { IconPlus, IconSearch } from '@/shared/ui/Icons'
import { cn } from '@/shared/lib/cn'

const groups: { key: MissionState | 'recent'; title: string; filter: (m: Mission) => boolean }[] = [
  { key: 'running', title: 'Running', filter: (m) => m.state === 'running' },
  { key: 'awaiting_approval', title: 'Awaiting approval', filter: (m) => m.state === 'awaiting_approval' },
  { key: 'queued', title: 'Queued', filter: (m) => m.state === 'queued' },
  {
    key: 'recent',
    title: 'Recent',
    filter: (m) => m.state === 'succeeded' || m.state === 'failed' || m.state === 'cancelled',
  },
]

function formatElapsed(ms: number) {
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(0)}s`
  return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`
}

export function MissionsPage() {
  const { data, isLoading, error, refetch } = useMissions()
  const launch = useLaunchMission()
  const toast = useUi((s) => s.pushToast)
  const [q, setQ] = useState('')
  const [view, setView] = useState<'board' | 'kanban'>('board')
  const profile = detectProfile()

  const filtered = useMemo(() => {
    const list = data ?? []
    if (!q) return list
    const s = q.toLowerCase()
    return list.filter((m) => m.name.toLowerCase().includes(s) || m.id.includes(s))
  }, [data, q])

  async function onLaunch() {
    try {
      await launch.mutateAsync({
        name: 'Ad-hoc mission',
        goal: 'Explore agent loop from Mission Control',
        requiredCapabilities: ['coding', 'tools'],
      })
      toast('Mission launched', 'ok')
    } catch (e) {
      toast((e as HostApiError).message, 'fail')
    }
  }

  return (
    <div className="page anim-in">
      <header className="mb-7 flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="mb-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-ink-2">Orchestration</p>
          <h1 className="page-title">
            Mission <span className="text-accent">board</span>
          </h1>
          <p className="page-sub">
            {profile.label}
            <span className="text-ink-2"> · </span>
            {profile.useMocks ? 'mock Host Interface' : 'live Host Interface'}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <div className="flex rounded-[10px] border border-[var(--line)] bg-bg-1 p-0.5">
            {(['board', 'kanban'] as const).map((v) => (
              <button
                key={v}
                type="button"
                onClick={() => setView(v)}
                className={cn(
                  'rounded-[8px] px-3 py-1.5 text-[12px] font-semibold capitalize',
                  view === v ? 'bg-[rgba(0,191,255,0.15)] text-accent' : 'text-ink-2 hover:text-ink-0',
                )}
              >
                {v}
              </button>
            ))}
          </div>
          <div className="relative">
            <span className="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-ink-2">
              <IconSearch />
            </span>
            <input
              className="input w-[min(100vw-2rem,16rem)] pl-8"
              placeholder="Search missions…"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              aria-label="Search missions"
            />
          </div>
          <Button onClick={onLaunch} disabled={launch.isPending}>
            <IconPlus />
            Launch mission
          </Button>
        </div>
      </header>

      {isLoading && (
        <div className="grid gap-3 md:grid-cols-2">
          <Skeleton className="h-36" />
          <Skeleton className="h-36" />
          <Skeleton className="h-36" />
        </div>
      )}
      {error && (
        <ErrorState
          message={(error as HostApiError).message}
          remediation={(error as HostApiError).remediation}
          onRetry={() => refetch()}
        />
      )}
      {!isLoading && !error && filtered.length === 0 && (
        <EmptyState
          title="No missions yet"
          body="Launch one from here, or run an example from the Agent OS CLI. The Mission Spine will light up as work starts."
          actionLabel="Launch mission"
          onAction={onLaunch}
        />
      )}

      {!isLoading && !error && view === 'kanban' && (
        <div className="grid gap-3 md:grid-cols-3">
          {[
            { title: 'Pending', filter: (m: Mission) => m.state === 'queued' || m.state === 'awaiting_approval', tone: 'warn' },
            { title: 'In progress', filter: (m: Mission) => m.state === 'running', tone: 'live' },
            { title: 'Done', filter: (m: Mission) => m.state === 'succeeded' || m.state === 'failed' || m.state === 'cancelled', tone: 'ok' },
          ].map((col) => {
            const items = filtered.filter(col.filter)
            return (
              <div key={col.title} className="panel-glass min-h-[20rem] p-3">
                <div className="section-label !mb-3">
                  {col.title}
                  <span className="font-mono normal-case tracking-normal text-ink-2">{items.length}</span>
                </div>
                <div className="space-y-2">
                  {items.map((m) => (
                    <MissionCard key={m.id} m={m} compact />
                  ))}
                </div>
              </div>
            )
          })}
        </div>
      )}

      {!isLoading &&
        !error &&
        view === 'board' &&
        groups.map((g) => {
          const items = filtered.filter(g.filter)
          if (!items.length) return null
          return (
            <section key={g.key} className="mb-9" aria-labelledby={`g-${g.key}`}>
              <h2 id={`g-${g.key}`} className="section-label">
                <span
                  className={cn(
                    'h-1.5 w-1.5 rounded-full',
                    g.key === 'running' && 'bg-live spine-pulse',
                    g.key === 'awaiting_approval' && 'bg-warn',
                    g.key === 'queued' && 'bg-ink-2',
                    g.key === 'recent' && 'bg-ink-1',
                  )}
                />
                {g.title}
                <span className="font-mono text-[11px] font-normal normal-case tracking-normal text-ink-2">
                  {items.length}
                </span>
              </h2>
              <div className="grid gap-3 md:grid-cols-2">
                {items.map((m) => (
                  <MissionCard key={m.id} m={m} />
                ))}
              </div>
            </section>
          )
        })}
    </div>
  )
}

function MissionCard({ m, compact }: { m: Mission; compact?: boolean }) {
  const phasePct =
    m.state === 'succeeded' ? 100 : m.state === 'failed' ? 100 : m.state === 'queued' ? 8 : m.state === 'awaiting_approval' ? 65 : 45

  return (
    <Link to={`/missions/${m.id}`} className="group block focus:outline-none">
      <Card
        className={cn(
          'card-state-' + m.state,
          'h-full p-4 hover:border-[rgba(0,191,255,0.35)] hover:glow-cyan group-focus-visible:ring-2 group-focus-visible:ring-accent',
          compact && 'p-3',
        )}
      >
        <div className="mb-3 flex items-start justify-between gap-3">
          <div className="min-w-0">
            <div className="truncate font-display text-[15px] font-semibold tracking-tight text-ink-0 group-hover:text-accent">
              {m.name}
            </div>
            <Mono className="mt-0.5 block truncate">{m.id}</Mono>
          </div>
          <StatusPill status={m.state} />
        </div>

        <p className="mb-3 line-clamp-2 min-h-[2.5rem] text-[13px] leading-relaxed text-ink-1">{m.goal}</p>

        <div className="mb-3 flex flex-wrap gap-1">
          {m.requiredCapabilities.map((c) => (
            <Badge key={c}>{c}</Badge>
          ))}
        </div>

        {/* progress */}
        <div className="mb-3 h-1 overflow-hidden rounded-full bg-bg-3">
          <div
            className={cn(
              'h-full rounded-full transition-all',
              m.state === 'failed' ? 'bg-fail' : m.state === 'awaiting_approval' ? 'bg-warn' : m.state === 'succeeded' ? 'bg-ok' : 'bg-live',
            )}
            style={{ width: `${phasePct}%` }}
          />
        </div>

        <div className="flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-[11px] text-ink-2">
          <span className="text-ink-1">{m.agentsActive} agents</span>
          <span>{formatElapsed(m.elapsedMs)}</span>
          <span className="text-ink-0">${m.costUsd.toFixed(2)}</span>
          {m.progressPhase && <span className="ml-auto uppercase tracking-wide">{m.progressPhase}</span>}
        </div>
      </Card>
    </Link>
  )
}
