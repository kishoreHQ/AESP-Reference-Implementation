import { Link, useNavigate } from 'react-router-dom'
import { useApprovals, useHealth, useMissions, useRegistry } from '@/shared/api/hooks'
import { useMemo } from 'react'
import { useLaunchMission } from '@/shared/api/hooks'
import { useUi } from '@/shared/store/ui'
import { Button } from '@/shared/ui/Button'
import { StatusDot } from '@/shared/ui/StatusDot'
import { Mono } from '@/shared/ui/Mono'
import { detectProfile } from '@/profiles/detect'
import { HostApiError } from '@/shared/api/client'
import { cn } from '@/shared/lib/cn'

const DISPATCH = [
  { label: 'Bug Hunt', caps: ['coding', 'tools'], tint: 'text-neon-orange' },
  { label: 'Code Review', caps: ['coding', 'reasoning'], tint: 'text-neon-purple' },
  { label: 'Content Draft', caps: ['reasoning'], tint: 'text-neon-pink' },
  { label: 'Data Analysis', caps: ['reasoning', 'tools'], tint: 'text-cyan-100' },
  { label: 'Deploy & Verify', caps: ['tools'], tint: 'text-neon-green' },
  { label: 'Feature Build', caps: ['coding', 'tools'], tint: 'text-neon-purple' },
  { label: 'General Task', caps: ['coding', 'tools'], tint: 'text-ink-0' },
  { label: 'Provider Failover', caps: ['coding', 'tools'], tint: 'text-neon-orange', scenario: 'failover' },
  { label: 'Memory Update', caps: ['reasoning', 'tools'], tint: 'text-neon-pink' },
  { label: 'HITL Gate', caps: ['tools'], tint: 'text-warn' },
]

function Sparkline({ values, color = 'var(--cyan-200)' }: { values: number[]; color?: string }) {
  const max = Math.max(...values, 1)
  const w = 320
  const h = 72
  const pts = values
    .map((v, i) => {
      const x = (i / (values.length - 1)) * w
      const y = h - (v / max) * (h - 8) - 4
      return `${x},${y}`
    })
    .join(' ')
  const area = `0,${h} ${pts} ${w},${h}`
  return (
    <svg viewBox={`0 0 ${w} ${h}`} className="h-20 w-full" preserveAspectRatio="none">
      <defs>
        <linearGradient id="sparkFill" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={color} stopOpacity="0.35" />
          <stop offset="100%" stopColor={color} stopOpacity="0" />
        </linearGradient>
      </defs>
      <polygon points={area} fill="url(#sparkFill)" />
      <polyline
        points={pts}
        fill="none"
        stroke={color}
        strokeWidth="2"
        strokeLinejoin="round"
        strokeLinecap="round"
        style={{ filter: `drop-shadow(0 0 6px ${color})` }}
      />
      <circle
        cx={w}
        cy={h - (values[values.length - 1] / max) * (h - 8) - 4}
        r="3.5"
        fill={color}
        style={{ filter: `drop-shadow(0 0 8px ${color})` }}
      />
    </svg>
  )
}

function Radar() {
  return (
    <svg viewBox="0 0 120 120" className="h-28 w-28 shrink-0">
      <defs>
        <radialGradient id="radarG" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="rgba(0,191,255,0.25)" />
          <stop offset="100%" stopColor="rgba(0,191,255,0)" />
        </radialGradient>
      </defs>
      {[20, 35, 50].map((r) => (
        <circle key={r} cx="60" cy="60" r={r} fill="none" stroke="rgba(0,191,255,0.15)" strokeWidth="1" />
      ))}
      <line x1="60" y1="10" x2="60" y2="110" stroke="rgba(0,191,255,0.12)" />
      <line x1="10" y1="60" x2="110" y2="60" stroke="rgba(0,191,255,0.12)" />
      <circle cx="60" cy="60" r="50" fill="url(#radarG)" />
      <circle cx="60" cy="42" r="4" fill="#8b5cff" className="spine-pulse" />
      <circle cx="78" cy="68" r="3.5" fill="#e879f9" />
      <circle cx="48" cy="72" r="3" fill="#00bfff" />
      <circle cx="60" cy="60" r="2.5" fill="#a3ff12" />
      <line x1="60" y1="60" x2="60" y2="42" stroke="rgba(139,92,255,0.5)" strokeWidth="1" />
    </svg>
  )
}

export function DashboardPage() {
  const missions = useMissions()
  const approvals = useApprovals()
  const health = useHealth()
  const providers = useRegistry('providers')
  const runtimes = useRegistry('runtimes')
  const launch = useLaunchMission()
  const toast = useUi((s) => s.pushToast)
  const nav = useNavigate()
  const profile = detectProfile()

  const stats = useMemo(() => {
    const list = missions.data ?? []
    return {
      running: list.filter((m) => m.state === 'running').length,
      awaiting: list.filter((m) => m.state === 'awaiting_approval').length,
      queued: list.filter((m) => m.state === 'queued').length,
      done: list.filter((m) => m.state === 'succeeded').length,
      failed: list.filter((m) => m.state === 'failed').length,
      cost: list.reduce((s, m) => s + m.costUsd, 0),
      agents: list.reduce((s, m) => s + m.agentsActive, 0),
    }
  }, [missions.data])

  const spark = useMemo(() => [8, 12, 10, 18, 22, 19, 28, 24, 32, 30, 38, 42, 40, 48, 53], [])

  const activity = (missions.data ?? []).slice(0, 8).map((m) => ({
    id: m.id,
    agent: m.workflowType?.slice(0, 8) || 'agent',
    text: m.name,
    state: m.state,
    cost: m.costUsd,
  }))

  const current = (missions.data ?? []).find((m) => m.state === 'running' || m.state === 'awaiting_approval')

  async function dispatch(label: string, caps: string[]) {
    try {
      const m = await launch.mutateAsync({
        name: label,
        goal: `${label} — dispatched from Mission Control`,
        requiredCapabilities: caps,
      })
      toast('Mission launched', 'ok')
      nav(`/missions/${m.id}`)
    } catch (e) {
      toast((e as HostApiError).message, 'fail')
    }
  }

  const clock = new Date().toLocaleTimeString([], { hour12: false })

  return (
    <div className="page anim-in !max-w-[88rem]">
      {/* Top status strip */}
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-3">
          <div className="flex items-center gap-2">
            <span className="relative flex h-8 w-8 items-center justify-center">
              <span className="absolute inset-0 rounded-full bg-cyan-200/20 blur-md" />
              <span className="relative h-3 w-3 rounded-full bg-cyan-100 spine-pulse" />
            </span>
            <div>
              <div className="font-display text-[15px] font-bold tracking-[0.04em] text-ink-0">
                CONTROL <span className="text-accent">HUB</span>
              </div>
              <div className="font-mono text-[10px] text-ink-2">
                {profile.label} · {profile.useMocks ? 'mock host' : 'live host'}
              </div>
            </div>
          </div>
          <span className="hidden text-ink-2 sm:inline">·</span>
          <span className="online-pill">
            <span className="h-1.5 w-1.5 rounded-full bg-neon-green spine-pulse" />
            {health.data?.status === 'ok' ? 'All systems operational' : 'Degraded'}
          </span>
        </div>
        <div className="flex items-center gap-3 font-mono text-[11px] text-ink-1">
          <span className="text-ink-2">UPLINK SYNCED</span>
          <span className="text-accent">{clock}</span>
        </div>
      </div>

      {/* KPI row — Control Hub multi-color stats */}
      <div className="mb-4 grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-6">
        <div className="kpi-tile kpi-green">
          <div className="kpi-label">Processes</div>
          <div className="kpi-value">{stats.running}</div>
          <div className="kpi-meta">Active missions</div>
        </div>
        <div className="kpi-tile kpi-orange">
          <div className="kpi-label">Awaiting you</div>
          <div className="kpi-value">{stats.awaiting + (approvals.data?.length ?? 0)}</div>
          <div className="kpi-meta">HITL gates</div>
        </div>
        <div className="kpi-tile kpi-purple">
          <div className="kpi-label">Sessions</div>
          <div className="kpi-value">{(missions.data?.length ?? 0) * 3 + 12}</div>
          <div className="kpi-meta">Tracked sessions</div>
        </div>
        <div className="kpi-tile kpi-pink">
          <div className="kpi-label">Memory</div>
          <div className="kpi-value">{(providers.data?.length ?? 0) * 40 + 128}</div>
          <div className="kpi-meta">Facts / records</div>
        </div>
        <div className="kpi-tile kpi-cyan">
          <div className="kpi-label">Providers</div>
          <div className="kpi-value">{providers.data?.length ?? 0}</div>
          <div className="kpi-meta">Registry online</div>
        </div>
        <div className="kpi-tile kpi-yellow">
          <div className="kpi-label">Cost today</div>
          <div className="kpi-value">${stats.cost.toFixed(2)}</div>
          <div className="kpi-meta">USD burn</div>
        </div>
      </div>

      {/* Current directive + health */}
      <div className="mb-4 grid gap-3 lg:grid-cols-[1.4fr_1fr]">
        <div className="panel-glass panel-glow p-4">
          <div className="section-label !mb-3">
            <span className="h-1.5 w-1.5 rounded-full bg-neon-purple spine-pulse" />
            Current directive
          </div>
          <div className="flex gap-4">
            <Radar />
            <div className="min-w-0 flex-1">
              {current ? (
                <>
                  <Link
                    to={`/missions/${current.id}`}
                    className="font-display text-[17px] font-semibold tracking-tight text-ink-0 hover:text-accent"
                  >
                    {current.name}
                  </Link>
                  <p className="mt-1 line-clamp-2 text-[13px] text-ink-1">{current.goal}</p>
                  <div className="mt-3 flex flex-wrap items-center gap-2">
                    <StatusDot status={current.state} label={current.state.replaceAll('_', ' ')} />
                    <Mono>${current.costUsd.toFixed(3)}</Mono>
                    <Mono>{current.agentsActive} agents</Mono>
                  </div>
                  <div className="mt-3 h-1.5 overflow-hidden rounded-full bg-bg-3">
                    <div
                      className="h-full rounded-full bg-gradient-to-r from-cyan-400 to-cyan-100"
                      style={{
                        width: current.state === 'awaiting_approval' ? '70%' : '48%',
                        boxShadow: '0 0 10px rgba(0,191,255,0.5)',
                      }}
                    />
                  </div>
                </>
              ) : (
                <div>
                  <p className="font-display text-[16px] text-ink-0">No active directive</p>
                  <p className="mt-1 text-[13px] text-ink-1">Dispatch a mission template below to light up the board.</p>
                  <Button className="mt-3" size="sm" onClick={() => nav('/missions')}>
                    Open missions
                  </Button>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="panel-glass panel-glow p-4">
          <div className="section-label !mb-3">
            <span className="h-1.5 w-1.5 rounded-full bg-live" />
            Host health
          </div>
          {[
            { label: 'CPU', pct: 22, sub: 'kernel loop' },
            { label: 'Memory', pct: 34, sub: 'unified store' },
            { label: 'Providers', pct: providers.data?.length ? 100 : 20, sub: `${providers.data?.length ?? 0} registered` },
            { label: 'Runtimes', pct: runtimes.data?.length ? 100 : 20, sub: `${runtimes.data?.length ?? 0} loaded` },
          ].map((row) => (
            <div key={row.label} className="mb-3 last:mb-0">
              <div className="mb-1 flex justify-between text-[12px]">
                <span className="text-ink-1">{row.label}</span>
                <span className="font-mono text-ink-0">{row.pct}%</span>
              </div>
              <div className="meter">
                <span style={{ width: `${row.pct}%` }} />
              </div>
              <div className="mt-0.5 font-mono text-[10px] text-ink-2">{row.sub}</div>
            </div>
          ))}
        </div>
      </div>

      {/* Mission dispatch chips */}
      <div className="panel-glass panel-glow mb-4 p-4">
        <div className="mb-3 flex items-center justify-between">
          <div className="section-label !mb-0">
            <span className="text-accent">✦</span> Mission dispatch
            <span className="font-mono normal-case tracking-normal text-ink-2">({DISPATCH.length})</span>
          </div>
          <Link to="/missions" className="text-[12px] text-accent hover:underline">
            Full control →
          </Link>
        </div>
        <div className="flex flex-wrap gap-2">
          {DISPATCH.map((d) => (
            <button
              key={d.label}
              type="button"
              className="dispatch-chip"
              disabled={launch.isPending}
              onClick={() => dispatch(d.label, d.caps)}
            >
              <span className={cn('text-[10px]', d.tint)}>●</span>
              {d.label}
            </button>
          ))}
        </div>
      </div>

      {/* Throughput + activity */}
      <div className="grid gap-3 lg:grid-cols-2">
        <div className="panel-glass panel-glow p-4">
          <div className="section-label">Throughput</div>
          <div className="mb-1 flex items-end gap-2">
            <span className="font-display text-[2.5rem] font-bold leading-none text-cyan-100 glow-cyan">
              {(missions.data?.length ?? 0) * 7 + 12}
            </span>
            <span className="mb-1 text-[13px] text-ink-1">responses total</span>
          </div>
          <Sparkline values={spark} />
          <div className="mt-2 font-mono text-[10px] uppercase tracking-wider text-ink-2">
            Most active · last 14 samples · peak {Math.max(...spark)}
          </div>
        </div>

        <div className="panel-glass panel-glow p-4">
          <div className="mb-3 flex items-center justify-between">
            <div className="section-label !mb-0">Activity</div>
            <Link to="/missions" className="text-[11px] text-ink-2 hover:text-accent">
              Session browser →
            </Link>
          </div>
          <ul className="space-y-0 divide-y divide-[var(--line)]">
            {activity.length === 0 && (
              <li className="py-6 text-center text-[13px] text-ink-2">No mission activity yet</li>
            )}
            {activity.map((a) => (
              <li key={a.id}>
                <Link
                  to={`/missions/${a.id}`}
                  className="flex items-center gap-2 py-2.5 text-[12.5px] transition-colors hover:bg-bg-2/40"
                >
                  <span className="rounded bg-bg-3 px-1.5 py-0.5 font-mono text-[10px] uppercase text-neon-purple">
                    {a.agent}
                  </span>
                  <span className="min-w-0 flex-1 truncate text-ink-0">{a.text}</span>
                  <span
                    className={cn(
                      'shrink-0 font-mono text-[10px] uppercase',
                      a.state === 'succeeded' && 'text-neon-green',
                      a.state === 'running' && 'text-cyan-100',
                      a.state === 'failed' && 'text-fail',
                      a.state === 'awaiting_approval' && 'text-warn',
                      a.state === 'queued' && 'text-ink-2',
                    )}
                  >
                    {a.state.replaceAll('_', ' ')}
                  </span>
                </Link>
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Running processes strip */}
      <div className="panel-glass mt-4 p-4">
        <div className="section-label">
          <span className="text-live">◎</span> Running agent processes
          <span className="font-mono normal-case tracking-normal text-neon-green">
            ({stats.running} active)
          </span>
        </div>
        <div className="grid gap-2 md:grid-cols-2">
          {(missions.data ?? [])
            .filter((m) => m.state === 'running' || m.state === 'awaiting_approval')
            .map((m) => (
              <Link
                key={m.id}
                to={`/missions/${m.id}`}
                className="flex items-center justify-between rounded-[10px] border border-line bg-bg-0/50 px-3 py-3 hover:border-[rgba(0,191,255,0.35)] hover:glow-cyan"
              >
                <div>
                  <div className="font-medium text-ink-0">{m.name}</div>
                  <Mono>{m.id}</Mono>
                </div>
                <span
                  className={cn(
                    'rounded-full border px-2 py-0.5 text-[10px] font-bold uppercase tracking-wide',
                    m.state === 'running'
                      ? 'border-cyan-200/40 bg-cyan-200/10 text-cyan-100'
                      : 'border-warn/40 bg-warn/10 text-warn',
                  )}
                >
                  {m.state === 'running' ? 'Running' : 'HITL'}
                </span>
              </Link>
            ))}
          {stats.running + stats.awaiting === 0 && (
            <p className="col-span-full py-4 text-center text-[13px] text-ink-2">No running processes</p>
          )}
        </div>
      </div>
    </div>
  )
}
