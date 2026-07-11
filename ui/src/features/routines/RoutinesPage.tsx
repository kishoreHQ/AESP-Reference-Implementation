import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useCreateRoutine, useRoutines } from '@/shared/api/hooks'
import { apiPost } from '@/shared/api/client'
import { Card } from '@/shared/ui/Card'
import { Button } from '@/shared/ui/Button'
import { Mono } from '@/shared/ui/Mono'
import { Badge } from '@/shared/ui/Badge'
import { useUi } from '@/shared/store/ui'
import { useQueryClient } from '@tanstack/react-query'

export function RoutinesPage() {
  const { data, refetch } = useRoutines()
  void refetch
  const create = useCreateRoutine()
  const toast = useUi((s) => s.pushToast)
  const qc = useQueryClient()
  const [name, setName] = useState('')
  const [schedule, setSchedule] = useState('@every 30m')
  const [prompt, setPrompt] = useState('')

  async function add() {
    try {
      await create.mutateAsync({ name, schedule, prompt, capabilities: ['tools'] })
      toast('Routine created', 'ok')
      setName('')
    } catch (e: any) {
      toast(e.message, 'fail')
    }
  }

  async function action(id: string, act: string) {
    await apiPost(`/v1/routines/${id}/${act}`, {})
    qc.invalidateQueries({ queryKey: ['routines'] })
    toast(`${act} ok`, 'ok')
  }

  return (
    <div className="page anim-in max-w-4xl">
      <p className="mb-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-ink-2">Command deck</p>
      <h1 className="page-title">
        Routines <span className="text-accent">cron</span>
      </h1>
      <p className="page-sub mb-6">Scheduled missions with next-fire countdown and run history</p>

      <Card className="mb-6 space-y-2">
        <div className="section-label">Create routine</div>
        <input className="input w-full" placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} />
        <input
          className="input w-full font-mono"
          placeholder="Cron e.g. @every 1h or 0 * * * *"
          value={schedule}
          onChange={(e) => setSchedule(e.target.value)}
        />
        <input className="input w-full" placeholder="Prompt / workflow" value={prompt} onChange={(e) => setPrompt(e.target.value)} />
        <Button onClick={add} disabled={!name || !schedule}>
          Save routine
        </Button>
      </Card>

      <div className="space-y-3">
        {(data || []).map((r: any) => {
          const next = r.nextFireAt ? new Date(r.nextFireAt) : null
          const countdown = next ? Math.max(0, Math.floor((next.getTime() - Date.now()) / 1000)) : null
          return (
            <Card key={r.id}>
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <div className="font-display text-[15px] font-semibold">{r.name}</div>
                  <Mono>{r.schedule}</Mono>
                  {r.paused && <Badge tone="warn">paused</Badge>}
                </div>
                <div className="text-right">
                  <div className="text-[11px] text-ink-2">Next fire</div>
                  <div className="font-mono text-[13px] text-accent">
                    {countdown != null ? `${Math.floor(countdown / 60)}m ${countdown % 60}s` : '—'}
                  </div>
                </div>
              </div>
              <p className="mt-2 text-[13px] text-ink-1">{r.prompt}</p>
              <div className="mt-2 flex flex-wrap gap-1">
                {(r.capabilities || []).map((c: string) => (
                  <Badge key={c}>{c}</Badge>
                ))}
              </div>
              {r.lastMissionId && (
                <div className="mt-2 text-[12px]">
                  Last:{' '}
                  <Link className="text-accent" to={`/missions/${r.lastMissionId}`}>
                    {r.lastMissionId}
                  </Link>{' '}
                  · {r.lastStatus}
                </div>
              )}
              <div className="mt-3 flex flex-wrap gap-2">
                <Button size="sm" variant="secondary" onClick={() => action(r.id, r.paused ? 'resume' : 'pause')}>
                  {r.paused ? 'Resume' : 'Pause'}
                </Button>
                <Button size="sm" onClick={() => action(r.id, 'fire')}>
                  Fire now
                </Button>
                <Button size="sm" variant="ghost" onClick={() => refetch()}>
                  Refresh
                </Button>
              </div>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
