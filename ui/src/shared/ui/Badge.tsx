import type { ReactNode } from 'react'
import { cn } from '@/shared/lib/cn'

export function Badge({
  children,
  tone = 'default',
}: {
  children: ReactNode
  tone?: 'default' | 'accent' | 'live' | 'warn' | 'fail' | 'ok'
}) {
  const tones = {
    default: 'bg-bg-2 text-ink-1 border-line',
    accent: 'bg-[var(--accent-dim)] text-accent border-accent/25',
    live: 'bg-[var(--live-dim)] text-live border-live/25',
    warn: 'bg-[var(--warn-dim)] text-warn border-warn/30',
    fail: 'bg-[var(--fail-dim)] text-fail border-fail/30',
    ok: 'bg-[var(--ok-dim)] text-ok border-ok/30',
  }
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-[5px] border px-1.5 py-0.5 font-mono text-[11px] font-medium tracking-wide',
        tones[tone],
      )}
    >
      {children}
    </span>
  )
}
