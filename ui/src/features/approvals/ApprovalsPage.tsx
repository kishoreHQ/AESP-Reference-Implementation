import { useDecision, useApprovals } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Button } from '@/shared/ui/Button'
import { Mono } from '@/shared/ui/Mono'
import { EmptyState } from '@/shared/ui/EmptyState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Skeleton } from '@/shared/ui/Skeleton'
import { useUi } from '@/shared/store/ui'
import { HostApiError } from '@/shared/api/client'
import { Link } from 'react-router-dom'
import { useState } from 'react'

export function ApprovalsPage() {
  const { data, isLoading, error, refetch } = useApprovals()
  const decision = useDecision()
  const toast = useUi((s) => s.pushToast)
  const [commentFor, setCommentFor] = useState<string | null>(null)
  const [comment, setComment] = useState('')

  async function act(id: string, d: 'approve' | 'reject', c?: string) {
    try {
      await decision.mutateAsync({ id, decision: d, comment: c })
      toast(d === 'approve' ? 'Approval recorded' : 'Rejection recorded', d === 'approve' ? 'ok' : 'warn')
      setCommentFor(null)
      setComment('')
    } catch (e) {
      toast((e as HostApiError).message, 'fail')
    }
  }

  return (
    <div className="page anim-in max-w-3xl">
      <header className="mb-6">
        <p className="mb-1 text-[11px] font-semibold uppercase tracking-[0.1em] text-warn">Human attention</p>
        <h1 className="page-title">Approvals</h1>
        <p className="page-sub">
          HITL gates. Amber is reserved for you — preview before you decide.
        </p>
      </header>

      {isLoading && <Skeleton className="h-40 w-full" />}
      {error && (
        <ErrorState
          message={(error as HostApiError).message}
          remediation={(error as HostApiError).remediation}
          onRetry={() => refetch()}
        />
      )}
      {!isLoading && !error && (data?.length ?? 0) === 0 && (
        <EmptyState
          title="Inbox clear"
          body="No pending approvals. When agents hit a policy gate, requests appear here with a full preview."
        />
      )}

      <div className="space-y-4">
        {data?.map((a) => (
          <Card key={a.id} className="border-warn/35 bg-gradient-to-br from-[var(--warn-dim)] to-bg-1 p-0 overflow-hidden">
            <div className="border-b border-warn/20 px-4 py-3">
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <div className="font-display text-[15px] font-semibold tracking-tight">{a.missionName}</div>
                  <Mono className="mt-0.5 block">
                    {a.id} ·{' '}
                    <Link className="text-accent hover:underline" to={`/missions/${a.missionId}`}>
                      {a.missionId}
                    </Link>
                  </Mono>
                </div>
                <span className="inline-flex items-center gap-1.5 rounded-full border border-warn/40 bg-[var(--warn-dim)] px-2.5 py-1 text-[11px] font-semibold text-warn">
                  <span className="h-1.5 w-1.5 rounded-full bg-warn" />
                  Awaiting you
                </span>
              </div>
            </div>

            <div className="space-y-3 p-4">
              <p className="text-[14px] leading-relaxed text-ink-0">{a.reason}</p>

              <dl className="grid gap-2 rounded-[8px] border border-line bg-bg-0/50 p-3 text-[12px] sm:grid-cols-2">
                <div>
                  <dt className="text-ink-2">Agent</dt>
                  <dd className="font-mono text-ink-0">{a.agentId}</dd>
                </div>
                <div>
                  <dt className="text-ink-2">Policy</dt>
                  <dd className="font-mono text-ink-0">{a.policy}</dd>
                </div>
                <div className="sm:col-span-2">
                  <dt className="text-ink-2">Blast radius</dt>
                  <dd className="text-ink-0">{a.blastRadius}</dd>
                </div>
                <div>
                  <dt className="text-ink-2">Age</dt>
                  <dd className="font-mono text-ink-0">{Math.round(a.ageMs / 1000)}s</dd>
                </div>
              </dl>

              <div className="overflow-hidden rounded-[8px] border border-line bg-bg-0">
                <div className="flex items-center justify-between border-b border-line px-3 py-1.5 text-[11px] text-ink-2">
                  <span>Preview · {a.preview.title}</span>
                  <span className="uppercase tracking-wide">{a.preview.kind}</span>
                </div>
                <pre className="max-h-52 overflow-auto p-3 font-mono text-[12.5px] leading-relaxed text-ink-0 whitespace-pre-wrap">
                  {a.preview.content}
                </pre>
              </div>

              {commentFor === a.id && (
                <textarea
                  className="input w-full min-h-[4.5rem] resize-y"
                  placeholder="Comment (optional)"
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                />
              )}

              <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                <Button className="w-full" variant="warn" onClick={() => act(a.id, 'approve')} disabled={decision.isPending}>
                  Approve
                </Button>
                <Button className="w-full" variant="danger" onClick={() => act(a.id, 'reject')} disabled={decision.isPending}>
                  Reject
                </Button>
                <Button
                  className="w-full"
                  variant="secondary"
                  onClick={() => {
                    if (commentFor === a.id) act(a.id, 'approve', comment)
                    else setCommentFor(a.id)
                  }}
                >
                  Approve with comment
                </Button>
                <Button className="w-full" variant="ghost" onClick={() => toast('Escalation recorded (demo)', 'info')}>
                  Escalate
                </Button>
              </div>
            </div>
          </Card>
        ))}
      </div>
    </div>
  )
}
