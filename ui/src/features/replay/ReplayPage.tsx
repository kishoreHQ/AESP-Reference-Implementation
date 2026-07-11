import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useReplay } from '@/shared/api/hooks'
import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Card'
import { Mono } from '@/shared/ui/Mono'

export function ReplayPage() {
  const { runId = '' } = useParams()
  const { data } = useReplay(runId)
  const [idx, setIdx] = useState(0)
  const events = data ?? []
  const current = events[idx]

  const banner = useMemo(() => 'REPLAY — read-only reconstruction', [])

  return (
    <div className="page anim-in max-w-4xl">
      <div className="mb-4 rounded-[6px] border border-accent/40 bg-accent/10 px-3 py-2 text-center text-[13px] font-semibold text-accent">
        {banner}
      </div>
      <h1 className="page-title">Replay</h1>
      <Mono className="mb-4 block">{runId}</Mono>

      <Card className="mb-4">
        <div className="mb-3 flex flex-wrap gap-2">
          <Button variant="secondary" onClick={() => setIdx((i) => Math.max(0, i - 1))} disabled={idx === 0}>
            Prev
          </Button>
          <Button onClick={() => setIdx((i) => Math.min(events.length - 1, i + 1))} disabled={idx >= events.length - 1}>
            Next
          </Button>
          <input
            type="range"
            min={0}
            max={Math.max(0, events.length - 1)}
            value={idx}
            onChange={(e) => setIdx(Number(e.target.value))}
            className="min-w-[200px] flex-1"
            aria-label="Seek"
          />
          <Mono>
            {idx + 1}/{events.length}
          </Mono>
        </div>
        {current && (
          <div>
            <div className="font-mono text-[13px] text-accent">{current.type}</div>
            <Mono>seq {current.seq}</Mono>
            <pre className="mt-2 overflow-auto rounded-[6px] bg-bg-0 p-3 font-mono text-[12px]">
              {JSON.stringify(current, null, 2)}
            </pre>
          </div>
        )}
      </Card>
    </div>
  )
}
