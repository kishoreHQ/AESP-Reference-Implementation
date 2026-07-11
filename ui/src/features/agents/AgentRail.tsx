import { Link } from 'react-router-dom'
import { useConnections, useSessions, useStopSession } from '@/shared/api/hooks'
import { StatusDot } from '@/shared/ui/StatusDot'
import { Mono } from '@/shared/ui/Mono'
import { Button } from '@/shared/ui/Button'
import { cn } from '@/shared/lib/cn'

/** Persistent agent rail — 100% registry + session state (extends UI-FLT-01). */
export function AgentRail({ compact }: { compact?: boolean }) {
  const conns = useConnections()
  const sessions = useSessions()
  const stop = useStopSession()

  const runtimes = (conns.data || []).filter((c: any) => c.kind === 'runtime' && c.status === 'connected')

  return (
    <aside className={cn('flex h-full flex-col border-r border-[var(--line)] bg-bg-1/90', compact && 'border-0')}>
      <div className="border-b border-[var(--line)] px-3 py-3">
        <div className="font-display text-[11px] font-bold uppercase tracking-[0.12em] text-ink-2">Agent rail</div>
        <div className="mt-0.5 font-mono text-[10px] text-ink-2">{runtimes.length} connected</div>
      </div>
      <div className="flex-1 space-y-1 overflow-auto p-2">
        {runtimes.map((c: any) => {
          const sess = (sessions.data || []).find((s: any) => s.runtimeId === c.pluginId && s.status === 'working')
          const status =
            sess?.status === 'working' ? 'running' : c.status === 'connected' ? 'idle' : 'error'
          return (
            <div
              key={c.id}
              className="rounded-[10px] border border-[var(--line)] bg-bg-0/50 p-2.5 hover:border-[rgba(0,191,255,0.35)]"
            >
              <div className="flex items-start justify-between gap-1">
                <Link to={`/agents/${encodeURIComponent(c.pluginId)}`} className="min-w-0 flex-1">
                  <div className="truncate text-[12.5px] font-semibold text-ink-0">{c.name}</div>
                  <Mono className="truncate text-[10px]">{c.pluginId}</Mono>
                </Link>
                <StatusDot status={status === 'idle' ? 'ok' : status} />
              </div>
              {c.unsandboxed && (
                <div className="mt-1 text-[9px] font-bold uppercase tracking-wide text-warn">unsandboxed</div>
              )}
              <div className="mt-1.5 flex flex-wrap gap-1">
                {(c.capabilities || []).slice(0, 3).map((cap: string) => (
                  <span key={cap} className="rounded bg-bg-3 px-1 font-mono text-[9px] text-ink-1">
                    {cap}
                  </span>
                ))}
              </div>
              {sess && (
                <div className="mt-1 truncate text-[10px] text-accent">session · {sess.id.slice(0, 12)}…</div>
              )}
              <div className="mt-2 flex gap-1">
                <Link
                  to={`/sessions/new?runtime=${encodeURIComponent(c.pluginId)}`}
                  className="flex-1 rounded-[6px] bg-[rgba(0,191,255,0.12)] py-1 text-center text-[10px] font-semibold text-accent"
                >
                  Chat
                </Link>
                <Button
                  size="sm"
                  variant="danger"
                  className="!min-h-7 !px-2 !text-[10px]"
                  title="Stop adapter process"
                  onClick={() => sess && stop.mutate(sess.id)}
                >
                  Stop
                </Button>
              </div>
            </div>
          )
        })}
        {!runtimes.length && (
          <p className="p-2 text-[12px] text-ink-2">
            No agents. <Link className="text-accent" to="/connect">Connect</Link>
          </p>
        )}
      </div>
    </aside>
  )
}
