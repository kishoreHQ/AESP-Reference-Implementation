import { useClaimTask, useMoveTask, useTasks } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Badge } from '@/shared/ui/Badge'
import { Button } from '@/shared/ui/Button'
import { Mono } from '@/shared/ui/Mono'
import { useUi } from '@/shared/store/ui'
import { cn } from '@/shared/lib/cn'

const COLS = [
  { id: 'backlog', title: 'Backlog' },
  { id: 'queued', title: 'Queued' },
  { id: 'in_progress', title: 'In progress' },
  { id: 'review', title: 'Review' },
  { id: 'done', title: 'Done' },
]

export function BoardPage() {
  const tasks = useTasks()
  const move = useMoveTask()
  const claim = useClaimTask()
  const toast = useUi((s) => s.pushToast)

  return (
    <div className="page anim-in !max-w-[96rem]">
      <div className="mb-4 flex flex-wrap items-end justify-between gap-3">
        <div>
          <p className="mb-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-ink-2">Command deck</p>
          <h1 className="page-title">
            Task <span className="text-accent">board</span>
          </h1>
          <p className="page-sub">Agents pull from queued · review integrates HITL</p>
        </div>
        <Button
          size="sm"
          onClick={async () => {
            try {
              const t = await claim.mutateAsync({ agentId: 'agent.default' })
              toast(`Claimed: ${t.title}`, 'ok')
            } catch (e: any) {
              toast(e.message || 'No claimable task', 'warn')
            }
          }}
        >
          Agent claim next
        </Button>
      </div>

      <div className="grid gap-3 md:grid-cols-3 xl:grid-cols-5">
        {COLS.map((col) => {
          const items = (tasks.data || []).filter((t: any) => t.column === col.id)
          return (
            <div key={col.id} className="panel-glass min-h-[18rem] p-2.5">
              <div className="section-label !mb-2">
                {col.title}
                <span className="font-mono normal-case tracking-normal text-ink-2">{items.length}</span>
              </div>
              <div className="space-y-2">
                {items.map((t: any) => (
                  <Card key={t.id} className="!p-2.5">
                    <div className="text-[13px] font-medium text-ink-0">{t.title}</div>
                    <Mono className="text-[10px]">{t.id}</Mono>
                    <div className="mt-1 flex flex-wrap gap-1">
                      {(t.capabilities || []).map((c: string) => (
                        <Badge key={c}>{c}</Badge>
                      ))}
                    </div>
                    <div className="mt-1 text-[11px] text-ink-2">
                      {t.assignee || 'any capable'}
                      {t.missionId && ` · ${t.missionId}`}
                    </div>
                    <div className="mt-2 flex flex-wrap gap-1">
                      {COLS.filter((c) => c.id !== col.id).map((c) => (
                        <button
                          key={c.id}
                          type="button"
                          className={cn(
                            'rounded px-1.5 py-0.5 text-[9px] font-semibold uppercase tracking-wide',
                            'bg-bg-3 text-ink-1 hover:bg-[rgba(0,191,255,0.15)] hover:text-accent',
                          )}
                          onClick={() => move.mutate({ id: t.id, column: c.id })}
                        >
                          → {c.title}
                        </button>
                      ))}
                    </div>
                  </Card>
                ))}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
