import { useUi } from '@/shared/store/ui'
import { cn } from '@/shared/lib/cn'

export function Toasts() {
  const { toasts, dismissToast } = useUi()
  return (
    <div
      className="pointer-events-none fixed bottom-20 right-4 z-50 flex w-[min(100vw-2rem,20rem)] flex-col gap-2 md:bottom-5"
      aria-live="polite"
    >
      {toasts.map((t) => (
        <div
          key={t.id}
          className={cn(
            'anim-in pointer-events-auto rounded-[10px] border bg-bg-1/95 px-3.5 py-2.5 text-[13px] backdrop-blur-sm',
            t.tone === 'ok' && 'border-ok/40',
            t.tone === 'warn' && 'border-warn/40',
            t.tone === 'fail' && 'border-fail/40',
            t.tone === 'info' && 'border-line',
          )}
        >
          <div className="flex items-start justify-between gap-3">
            <span className="leading-snug text-ink-0">{t.message}</span>
            <button
              className="shrink-0 text-ink-2 hover:text-ink-0"
              onClick={() => dismissToast(t.id)}
              aria-label="Dismiss"
            >
              ×
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
