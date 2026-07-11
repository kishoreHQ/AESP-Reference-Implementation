import { cn } from '@/shared/lib/cn'

const map: Record<string, string> = {
  running: 'bg-accent spine-pulse',
  live: 'bg-accent spine-pulse',
  awaiting_approval: 'bg-warn',
  blocked: 'bg-warn',
  pending: 'bg-ink-2',
  queued: 'bg-ink-2',
  succeeded: 'bg-ok',
  ok: 'bg-ok',
  passed: 'bg-ok',
  failed: 'bg-fail',
  fail: 'bg-fail',
  cancelled: 'bg-ink-2',
  healthy: 'bg-ok',
  degraded: 'bg-warn',
  down: 'bg-fail',
}

export function StatusDot({ status, label }: { status: string; label?: string }) {
  const pretty = (label || status).replaceAll('_', ' ')
  return (
    <span className="inline-flex items-center gap-1.5" title={pretty}>
      <span
        className={cn(
          'inline-block h-2 w-2 shrink-0 rounded-full ring-2 ring-bg-1',
          map[status] || 'bg-ink-2',
        )}
        aria-hidden
      />
      {label != null && <span className="text-[12px] capitalize text-ink-1">{pretty}</span>}
      <span className="sr-only">{pretty}</span>
    </span>
  )
}

export function StatusPill({ status }: { status: string }) {
  const styles: Record<string, string> = {
    running: 'bg-[var(--live-dim)] text-accent border-[rgba(0,191,255,0.35)]',
    awaiting_approval: 'bg-[var(--warn-dim)] text-warn border-[rgba(255,159,28,0.4)]',
    blocked: 'bg-[var(--warn-dim)] text-warn border-[rgba(255,159,28,0.4)]',
    succeeded: 'bg-[var(--ok-dim)] text-ok border-[rgba(163,255,18,0.35)]',
    ok: 'bg-[var(--ok-dim)] text-ok border-[rgba(163,255,18,0.35)]',
    failed: 'bg-[var(--fail-dim)] text-fail border-[rgba(248,113,113,0.35)]',
    queued: 'bg-bg-2 text-ink-1 border-[var(--line)]',
    cancelled: 'bg-bg-2 text-ink-1 border-[var(--line)]',
    pending: 'bg-bg-2 text-ink-1 border-[var(--line)]',
  }
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-[11px] font-semibold capitalize',
        styles[status] || styles.pending,
      )}
    >
      <StatusDot status={status} />
      {status.replaceAll('_', ' ')}
    </span>
  )
}
