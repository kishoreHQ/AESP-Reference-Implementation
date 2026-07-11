import { Button } from './Button'

export function EmptyState({
  title,
  body,
  actionLabel,
  onAction,
}: {
  title: string
  body: string
  actionLabel?: string
  onAction?: () => void
}) {
  return (
    <div className="anim-in flex flex-col items-start gap-3 rounded-[14px] border border-dashed border-line bg-bg-1/80 p-8">
      <div className="flex h-10 w-10 items-center justify-center rounded-[8px] border border-line bg-bg-2 text-ink-1">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75">
          <rect x="3" y="4" width="18" height="16" rx="2" />
          <path d="M3 10h18" />
        </svg>
      </div>
      <h3 className="font-display text-[16px] font-semibold text-ink-0">{title}</h3>
      <p className="max-w-md text-[13.5px] leading-relaxed text-ink-1">{body}</p>
      {actionLabel && onAction && <Button onClick={onAction}>{actionLabel}</Button>}
    </div>
  )
}
