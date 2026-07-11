import { Link } from 'react-router-dom'
import { useEvaluations } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Badge } from '@/shared/ui/Badge'
import { Mono } from '@/shared/ui/Mono'

export function EvaluationsPage() {
  const { data } = useEvaluations()
  return (
    <div className="page anim-in">
      <h1 className="page-title mb-4">Evaluations</h1>
      <div className="space-y-3">
        {data?.map((e) => (
          <Card key={e.id}>
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div>
                <div className="font-display font-semibold">{e.suite}</div>
                <Mono>{e.id}</Mono>
              </div>
              <div className="flex items-center gap-2">
                {e.regression && <Badge tone="fail">regression</Badge>}
                <Badge tone={e.passed ? 'ok' : 'fail'}>{e.passed ? 'pass' : 'fail'}</Badge>
                <Mono className="text-ink-0">score {e.score.toFixed(2)} ({e.baselineDelta >= 0 ? '+' : ''}{e.baselineDelta.toFixed(2)})</Mono>
              </div>
            </div>
            <ul className="mt-3 space-y-1 text-[13px]">
              {e.metrics.map((m) => (
                <li key={m.name} className="flex justify-between border-b border-line/40 py-1">
                  <span>{m.name}</span>
                  <Mono className={m.pass ? 'text-ok' : 'text-fail'}>
                    {m.value} {m.pass ? 'pass' : 'fail'}
                  </Mono>
                </li>
              ))}
            </ul>
            {e.traceMissionId && (
              <Link className="mt-2 inline-block text-[13px] text-accent" to={`/missions/${e.traceMissionId}`}>
                Open failing trace →
              </Link>
            )}
          </Card>
        ))}
      </div>
    </div>
  )
}
