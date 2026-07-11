import { useEffect, useState } from 'react'
import { useSearchParams, useParams, Link } from 'react-router-dom'
import {
  useCreateSession,
  useSessionMessage,
  useSessions,
  useStopSession,
} from '@/shared/api/hooks'
import { EventBridge } from '@/shared/events/EventBridge'
import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Card'
import { Mono } from '@/shared/ui/Mono'
import { Badge } from '@/shared/ui/Badge'
import { useUi } from '@/shared/store/ui'

export function LiveSessionPage() {
  const { id: paramId } = useParams()
  const [sp] = useSearchParams()
  const runtime = sp.get('runtime') || 'runtime.generic-pty'
  const create = useCreateSession()
  const message = useSessionMessage()
  const stop = useStopSession()
  const sessions = useSessions()
  const toast = useUi((s) => s.pushToast)
  const [sessionId, setSessionId] = useState(paramId || '')
  const [text, setText] = useState('')
  const [lines, setLines] = useState<{ t: string; kind: string }[]>([])

  useEffect(() => {
    if (paramId === 'new' || (!paramId && !sessionId)) {
      create
        .mutateAsync({ runtime, agent: runtime })
        .then((s) => {
          setSessionId(s.id)
          toast('Session opened', 'ok')
        })
        .catch((e) => toast(String(e.message || e), 'fail'))
    } else if (paramId && paramId !== 'new') {
      setSessionId(paramId)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (!sessionId) return
    EventBridge.connect(sessionId)
    return EventBridge.subscribe((ev) => {
      if (ev.missionId !== sessionId && (ev.data as any)?.sessionId !== sessionId) return
      const kind = ev.type
      const payload = (ev.data as any)?.payload || ev.data || {}
      const t =
        payload.text ||
        payload.rawType ||
        (kind === 'session.model_switch' ? `model → ${payload.model || payload.payload?.model}` : kind)
      setLines((L) => [...L.slice(-200), { t: String(t), kind }])
    })
  }, [sessionId])

  const sess = (sessions.data || []).find((s: any) => s.id === sessionId)

  async function send() {
    if (!text.trim() || !sessionId) return
    const msg = text
    setText('')
    setLines((L) => [...L, { t: msg, kind: 'user' }])
    try {
      await message.mutateAsync({ id: sessionId, text: msg })
    } catch (e: any) {
      toast(e.message, 'fail')
    }
  }

  return (
    <div className="flex h-full min-h-0 flex-col anim-in">
      <header className="flex flex-wrap items-center justify-between gap-2 border-b border-[var(--line)] bg-bg-1/60 px-4 py-3">
        <div>
          <div className="text-[11px] text-ink-2">
            <Link to="/" className="text-accent">
              Home
            </Link>{' '}
            / Live session
          </div>
          <h1 className="font-display text-[16px] font-bold">
            {runtime} {sess?.unsandboxed && <Badge tone="warn">unsandboxed</Badge>}
          </h1>
          <Mono>{sessionId || 'starting…'}</Mono>
        </div>
        <div className="flex gap-2">
          <Button size="sm" variant="danger" onClick={() => sessionId && stop.mutate(sessionId)}>
            Stop
          </Button>
        </div>
      </header>

      <div className="grid min-h-0 flex-1 lg:grid-cols-[1fr_260px]">
        <section className="flex min-h-0 flex-col">
          <div className="min-h-0 flex-1 space-y-1 overflow-auto bg-bg-0 p-3 font-mono text-[12.5px]">
            {lines.map((l, i) => (
              <div
                key={i}
                className={
                  l.kind.includes('model')
                    ? 'rounded border border-neon-purple/40 bg-[rgba(139,92,255,0.1)] px-2 py-1 text-neon-purple'
                    : l.kind.includes('tool')
                      ? 'rounded border border-[var(--line)] bg-bg-2/50 px-2 py-1'
                      : l.kind === 'user'
                        ? 'text-accent'
                        : 'text-ink-0'
                }
              >
                {l.kind.includes('model') && <span className="mr-2 text-[10px] uppercase">model switch</span>}
                {l.t}
              </div>
            ))}
            {!lines.length && <p className="text-ink-2">Stream will appear here…</p>}
          </div>
          <div className="flex gap-2 border-t border-[var(--line)] p-3">
            <input
              className="input flex-1"
              placeholder="Message agent…"
              value={text}
              onChange={(e) => setText(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && send()}
            />
            <Button onClick={send}>Send</Button>
          </div>
        </section>
        <aside className="hidden border-l border-[var(--line)] bg-bg-1 p-3 lg:block">
          <div className="section-label">Meta</div>
          <dl className="space-y-2 text-[12px]">
            <div>
              <dt className="text-ink-2">Status</dt>
              <dd>{sess?.status || '—'}</dd>
            </div>
            <div>
              <dt className="text-ink-2">Model</dt>
              <dd className="font-mono">{sess?.model || 'routed'}</dd>
            </div>
            <div>
              <dt className="text-ink-2">Tokens / cost</dt>
              <dd className="font-mono">
                {sess?.tokens ?? 0} · ${(sess?.costUsd ?? 0).toFixed(4)}
              </dd>
            </div>
            <div>
              <dt className="text-ink-2">Tool calls</dt>
              <dd className="font-mono">{sess?.toolCalls ?? 0}</dd>
            </div>
          </dl>
          <p className="mt-4 text-[11px] text-ink-2">
            Memory writes from this session default to trust <Badge>agent</Badge> (never verified).
          </p>
        </aside>
      </div>
    </div>
  )
}
