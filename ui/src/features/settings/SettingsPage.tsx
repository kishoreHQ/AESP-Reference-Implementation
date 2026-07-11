import { useBudgets, usePolicies, useSavePolicy } from '@/shared/api/hooks'
import { apiPost } from '@/shared/api/client'
import { Card } from '@/shared/ui/Card'
import { Button } from '@/shared/ui/Button'
import { Mono } from '@/shared/ui/Mono'
import { useUi } from '@/shared/store/ui'
import { detectProfile } from '@/profiles/detect'
import { useState } from 'react'

export function SettingsPage() {
  const budgets = useBudgets()
  const policies = usePolicies()
  const save = useSavePolicy()
  const { theme, setTheme, pushToast } = useUi()
  const profile = detectProfile()
  const [secret, setSecret] = useState('')

  async function storeCred() {
    await apiPost('/v1/credentials', { name: 'provider.default', value: secret })
    setSecret('')
    pushToast('Credential stored (write-only)', 'ok')
  }

  return (
    <div className="page anim-in max-w-3xl space-y-6">
      <h1 className="page-title">Settings</h1>

      <Card>
        <h2 className="mb-2 font-display font-semibold">Profile connection</h2>
        <p className="text-[13px] text-ink-1">
          Active: <strong className="text-ink-0">{profile.label}</strong> · base <Mono>{profile.baseUrl}</Mono>
          {profile.useMocks ? ' · mocks' : ' · live'}
        </p>
      </Card>

      <Card>
        <h2 className="mb-2 font-display font-semibold">Credentials</h2>
        <p className="mb-2 text-[13px] text-ink-1">Secrets are write-only and never re-displayed.</p>
        <input
          type="password"
          className="mb-2 w-full rounded-[6px] border border-line bg-bg-0 px-3 py-2"
          placeholder="Paste secret once"
          value={secret}
          onChange={(e) => setSecret(e.target.value)}
          autoComplete="off"
        />
        <Button onClick={storeCred} disabled={!secret}>
          Store credential
        </Button>
      </Card>

      <Card>
        <h2 className="mb-2 font-display font-semibold">Budgets</h2>
        <ul className="space-y-2">
          {budgets.data?.map((b) => (
            <li key={b.id} className="flex justify-between text-[13px]">
              <Mono>{b.scope}</Mono>
              <span>
                ${b.usedUsd.toFixed(2)} / ${b.capUsd.toFixed(2)}
              </span>
            </li>
          ))}
        </ul>
      </Card>

      <Card>
        <h2 className="mb-2 font-display font-semibold">Policies</h2>
        {policies.data?.map((p) => (
          <div key={p.id} className="mb-3">
            <div className="font-medium">
              {p.name} <Mono>v{p.version}</Mono>
            </div>
            <textarea
              className="mt-1 w-full rounded-[6px] border border-line bg-bg-0 p-2 font-mono text-[12px]"
              rows={4}
              defaultValue={p.body}
              id={`pol-${p.id}`}
            />
            <Button
              className="mt-2"
              variant="secondary"
              onClick={() => {
                const body = (document.getElementById(`pol-${p.id}`) as HTMLTextAreaElement).value
                save.mutate({ ...p, body })
                pushToast('Policy saved', 'ok')
              }}
            >
              Save policy
            </Button>
          </div>
        ))}
      </Card>

      <Card>
        <h2 className="mb-2 font-display font-semibold">Appearance</h2>
        <div className="flex gap-2">
          <Button variant={theme === 'dark' ? 'primary' : 'secondary'} onClick={() => setTheme('dark')}>
            Dark
          </Button>
          <Button variant={theme === 'light' ? 'primary' : 'secondary'} onClick={() => setTheme('light')}>
            Light
          </Button>
        </div>
      </Card>
    </div>
  )
}
