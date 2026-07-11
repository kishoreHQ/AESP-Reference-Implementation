import { Button } from './Button'

export function ErrorState({
  message,
  remediation,
  onRetry,
}: {
  message: string
  remediation?: string
  onRetry?: () => void
}) {
  return (
    <div
      className="anim-in rounded-[10px] border border-fail/35 bg-[var(--fail-dim)] p-4"
      role="alert"
    >
      <p className="font-medium text-fail">{message}</p>
      {remediation && <p className="mt-1.5 text-[13px] leading-relaxed text-ink-1">{remediation}</p>}
      {onRetry && (
        <Button className="mt-3" variant="secondary" size="sm" onClick={onRetry}>
          Retry
        </Button>
      )}
    </div>
  )
}
