import { useParams, Link } from 'react-router-dom'
import { useState } from 'react'
import { useAgentAnalytics, useConnections, useSessions, useTasks } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Badge } from '@/shared/ui/Badge'
import { Mono } from '@/shared/ui/Mono'
import { StatusDot } from '@/shared/ui/StatusDot'
import { cn } from '@/shared/lib/cn'

const tabs = ['Overview', 'Sessions', 'Credentials', 'Skills', 'Tasks', 'Analytics'] as const

export function ControlRoomPage() {
  const { id = '' } = useParams()
  const agentId = decodeURIComponent(id)
  const [tab, setTab] = useState<(typeof tabs)[number]>('Overview')
  const conns = useConnections()
  const sessions = useSessions()
  const tasks = useTasks()
  const analytics = useAgentAnalytics(agentId)
  const conn = (conns.data || []).find((c: any) => c.pluginId === agentId || c.id === agentId)

  return (
    <div className="page anim-in">
      <div className="mb-1 text-[12px] text-ink-2">
        <Link to="/connect" className="text-accent">
          Connections
        </Link>{' '}
        / Control room
      </div>
      <h1 className="page-title">{conn?.name || agentId}</h1>
      <Mono className="mb-4 block">{agentId}</Mono>
      {conn?.unsandboxed && <Badge tone="warn">unsandboxed PTY</Badge>}

      <div className="mb-4 flex flex-wrap gap-1 border-b border-[var(--line)]">
        {tabs.map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => setTab(t)}
            className={cn(
              'min-h-11 px-3 text-[13px] font-medium',
              tab === t ? 'border-b-2 border-accent text-accent' : 'text-ink-1',
            )}
          >
            {t}
          </button>
        ))}
      </div>

      {tab === 'Overview' && (
        <Card>
          <StatusDot status={conn?.status === 'connected' ? 'ok' : 'failed'} label={conn?.status || 'unknown'} />
          <div className="mt-3 flex flex-wrap gap-1">
            {(conn?.capabilities || []).map((c: string) => (
              <Badge key={c}>{c}</Badge>
            ))}
          </div>
          <p className="mt-3 text-[13px] text-ink-1">
            Default capability routing — no model-name hardcoding (INV-03).
          </p>
          <Link className="mt-3 inline-block text-accent" to={`/sessions/new?runtime=${encodeURIComponent(agentId)}`}>
            Open live session →
          </Link>
        </Card>
      )}

      {tab === 'Sessions' && (
        <div className="space-y-2">
          {(sessions.data || [])
            .filter((s: any) => s.runtimeId === agentId || s.agentId === agentId)
            .map((s: any) => (
              <Card key={s.id} className="flex justify-between">
                <div>
                  <Mono>{s.id}</Mono>
                  <div className="text-[12px] text-ink-1">{s.status}</div>
                </div>
                <Link className="text-[12px] text-accent" to={`/replay/${s.id}`}>
                  Replay
                </Link>
              </Card>
            ))}
        </div>
      )}

      {tab === 'Credentials' && (
        <Card>
          <p className="text-[13px] text-ink-1">
            Bindings via unified credential manager (INV-07). Secrets write-only — rotate in Settings.
          </p>
          <Mono className="mt-2 block">credential ref: {conn?.credentialId || 'none'}</Mono>
        </Card>
      )}

      {tab === 'Skills' && (
        <Card>
          <p className="text-[13px] text-ink-1">Tools/skills from unified tool registry (INV-06). Enable/disable is policy-gated.</p>
        </Card>
      )}

      {tab === 'Tasks' && (
        <div className="space-y-2">
          {(tasks.data || [])
            .filter((t: any) => !t.assignee || t.assignee === agentId)
            .map((t: any) => (
              <Card key={t.id}>
                <div className="font-medium">{t.title}</div>
                <Badge>{t.column}</Badge>
              </Card>
            ))}
        </div>
      )}

      {tab === 'Analytics' && analytics.data && (
        <div className="grid gap-3 md:grid-cols-2">
          <Card>
            <div className="kpi-label">Sessions</div>
            <div className="font-display text-2xl text-accent">{analytics.data.sessions}</div>
          </Card>
          <Card>
            <div className="kpi-label">Tool calls</div>
            <div className="font-display text-2xl text-neon-purple">{analytics.data.toolCalls}</div>
          </Card>
          <Card>
            <div className="kpi-label">Tokens</div>
            <div className="font-display text-2xl text-neon-orange">{analytics.data.tokens}</div>
          </Card>
          <Card>
            <div className="kpi-label">Error rate</div>
            <div className="font-display text-2xl text-fail">
              {(analytics.data.errorRate * 100).toFixed(1)}%
            </div>
          </Card>
          <Card className="md:col-span-2">
            <div className="kpi-label mb-2">Hour-of-day activity</div>
            <div className="flex h-16 items-end gap-0.5">
              {(analytics.data.hourHistogram || Array(24).fill(0)).map((v: number, i: number) => (
                <div
                  key={i}
                  className="flex-1 rounded-t bg-accent/70"
                  style={{ height: `${Math.max(4, (v / Math.max(1, ...analytics.data.hourHistogram)) * 100)}%` }}
                  title={`${i}:00 → ${v}`}
                />
              ))}
            </div>
          </Card>
          <Card className="md:col-span-2">
            <div className="kpi-label">Models used (from journal)</div>
            <ul className="mt-2 space-y-1 font-mono text-[12px]">
              {Object.entries(analytics.data.modelsUsed || {}).map(([m, n]) => (
                <li key={m} className="flex justify-between">
                  <span>{m}</span>
                  <span>{n as number}</span>
                </li>
              ))}
              {!Object.keys(analytics.data.modelsUsed || {}).length && (
                <li className="text-ink-2">No model_switch events yet</li>
              )}
            </ul>
          </Card>
        </div>
      )}
    </div>
  )
}
