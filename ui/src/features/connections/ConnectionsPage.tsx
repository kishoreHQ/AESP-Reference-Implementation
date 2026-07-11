import { useState } from 'react'
import { useProbe, useRegisterConnection, useConnections } from '@/shared/api/hooks'
import { Button } from '@/shared/ui/Button'
import { Card } from '@/shared/ui/Card'
import { Badge } from '@/shared/ui/Badge'
import { Mono } from '@/shared/ui/Mono'
import { StatusDot } from '@/shared/ui/StatusDot'
import { Skeleton } from '@/shared/ui/Skeleton'
import { ErrorState } from '@/shared/ui/ErrorState'
import { useUi } from '@/shared/store/ui'
import { HostApiError } from '@/shared/api/client'
import { cn } from '@/shared/lib/cn'

export function ConnectionsPage() {
  const probe = useProbe()
  const conns = useConnections()
  const register = useRegisterConnection()
  const toast = useUi((s) => s.pushToast)
  const [step, setStep] = useState(0)
  const [kind, setKind] = useState<'runtime' | 'provider'>('runtime')
  const [selected, setSelected] = useState<any>(null)
  const [secret, setSecret] = useState('')
  const [result, setResult] = useState<any>(null)

  async function doConnect() {
    try {
      const c = await register.mutateAsync({
        kind,
        pluginId: selected?.id,
        name: selected?.name || selected?.id,
        credential: secret || undefined,
        capabilities: selected?.capabilities,
      })
      setResult(c)
      setStep(2)
      toast(c.status === 'connected' ? 'Connected' : 'Handshake failed', c.status === 'connected' ? 'ok' : 'fail')
    } catch (e) {
      const err = e as HostApiError
      toast(err.message, 'fail')
      setResult({ status: 'error', lastError: err.remediation || err.message })
      setStep(2)
    }
  }

  return (
    <div className="page anim-in max-w-3xl">
      <p className="mb-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-ink-2">Fleet</p>
      <h1 className="page-title">
        Connect <span className="text-accent">agents & providers</span>
      </h1>
      <p className="page-sub mb-6">
        Runtime plugins via registry (INV-09). Providers separate (INV-01). Keys write-only (INV-07).
      </p>

      <div className="mb-4 flex gap-2">
        {(['runtime', 'provider'] as const).map((k) => (
          <button
            key={k}
            type="button"
            onClick={() => { setKind(k); setStep(0); setSelected(null); setResult(null) }}
            className={cn(
              'rounded-full border px-3 py-1.5 text-[12px] font-semibold capitalize',
              kind === k ? 'border-accent bg-[rgba(0,191,255,0.12)] text-accent' : 'border-[var(--line)] text-ink-1',
            )}
          >
            {k === 'runtime' ? 'Connect runtime / agent' : 'Connect provider'}
          </button>
        ))}
      </div>

      {step === 0 && (
        <Card className="space-y-3">
          <div className="flex items-center justify-between">
            <div className="section-label !mb-0">1 · Detect / pick plugin</div>
            <Button size="sm" variant="secondary" onClick={() => probe.refetch()}>
              Scan this machine
            </Button>
          </div>
          {probe.isLoading && <Skeleton className="h-24" />}
          {probe.error && <ErrorState message={(probe.error as HostApiError).message} />}
          <div className="space-y-2">
            {(probe.data || [])
              .filter((c: any) => c.kind === kind)
              .map((c: any) => (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => setSelected(c)}
                  className={cn(
                    'flex w-full items-center justify-between rounded-[10px] border px-3 py-2.5 text-left transition-colors',
                    selected?.id === c.id
                      ? 'border-accent bg-[rgba(0,191,255,0.1)]'
                      : 'border-[var(--line)] hover:border-[rgba(0,191,255,0.3)]',
                  )}
                >
                  <div>
                    <div className="font-medium text-ink-0">{c.name || c.id}</div>
                    <Mono>{c.id}</Mono>
                    <div className="mt-1 text-[11px] text-ink-2">{c.detail}</div>
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <StatusDot status={c.detected ? 'ok' : 'down'} label={c.detected ? 'detected' : 'missing'} />
                    {c.version && <Mono className="text-[10px]">{c.version}</Mono>}
                    {c.unsandboxed && <Badge tone="warn">unsandboxed</Badge>}
                  </div>
                </button>
              ))}
          </div>
          <Button disabled={!selected} onClick={() => setStep(1)}>
            Continue
          </Button>
        </Card>
      )}

      {step === 1 && (
        <Card className="space-y-3">
          <div className="section-label">2 · Credentials & handshake</div>
          <p className="text-[13px] text-ink-1">
            Selected <Mono className="text-ink-0">{selected?.id}</Mono>
            {selected?.needsCredential ? ' — paste API key (never re-shown).' : ' — no key required for local probe.'}
          </p>
          {selected?.needsCredential && (
            <input
              type="password"
              className="input w-full"
              placeholder="Paste secret once"
              value={secret}
              onChange={(e) => setSecret(e.target.value)}
              autoComplete="off"
            />
          )}
          <div className="flex gap-2">
            <Button variant="secondary" onClick={() => setStep(0)}>
              Back
            </Button>
            <Button onClick={doConnect} disabled={register.isPending}>
              Run handshake
            </Button>
          </div>
        </Card>
      )}

      {step === 2 && (
        <Card className="space-y-3">
          <div className="section-label">3 · Result</div>
          <StatusDot
            status={result?.status === 'connected' ? 'ok' : 'failed'}
            label={result?.status || 'unknown'}
          />
          {result?.lastError && (
            <pre className="max-h-48 overflow-auto rounded-[8px] border border-fail/40 bg-[var(--fail-dim)] p-3 font-mono text-[12px] text-fail whitespace-pre-wrap">
              {result.lastError}
            </pre>
          )}
          {result?.status === 'connected' && (
            <p className="text-[13px] text-neon-green">
              Connected. Agent appears on the rail · stop control always available.
            </p>
          )}
          <Button onClick={() => { setStep(0); setSelected(null); setSecret(''); setResult(null) }}>
            Connect another
          </Button>
        </Card>
      )}

      <div className="mt-8">
        <div className="section-label">Active connections</div>
        <div className="space-y-2">
          {(conns.data || []).map((c: any) => (
            <Card key={c.id} className="flex items-center justify-between py-3">
              <div>
                <div className="font-medium">{c.name}</div>
                <Mono>{c.pluginId}</Mono>
              </div>
              <div className="flex items-center gap-2">
                {c.unsandboxed && <Badge tone="warn">unsandboxed</Badge>}
                <StatusDot status={c.status === 'connected' ? 'ok' : 'failed'} label={c.status} />
              </div>
            </Card>
          ))}
          {!conns.data?.length && <p className="text-[13px] text-ink-2">No connections yet.</p>}
        </div>
      </div>
    </div>
  )
}
