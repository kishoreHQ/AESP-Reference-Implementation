import { createRoot, type Root } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router-dom'
import App from './app/App'
import type { ProfileId } from './shared/types/host'
import './index.css'

export interface MountOptions {
  baseUrl?: string
  token?: string
  profile?: ProfileId
  useMocks?: boolean
}

let root: Root | null = null

/**
 * P3 Embedded profile entry (UI-TEC-06).
 * mountHermesUI(el, { baseUrl, token })
 */
export function mountHermesUI(el: HTMLElement, opts: MountOptions = {}) {
  window.__HERMES_UI__ = {
    baseUrl: opts.baseUrl ?? '/api',
    token: opts.token,
    profile: opts.profile ?? 'embedded',
    useMocks: opts.useMocks ?? false,
  }
  const qc = new QueryClient()
  root = createRoot(el)
  root.render(
    <QueryClientProvider client={qc}>
      <MemoryRouter>
        <App />
      </MemoryRouter>
    </QueryClientProvider>,
  )
  return {
    unmount() {
      root?.unmount()
      root = null
    },
  }
}
