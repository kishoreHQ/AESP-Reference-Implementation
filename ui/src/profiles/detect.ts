import type { ProfileId } from '@/shared/types/host'

export interface RuntimeProfile {
  id: ProfileId
  baseUrl: string
  token?: string
  useMocks: boolean
  label: string
}

declare global {
  interface Window {
    __HERMES_UI__?: {
      baseUrl?: string
      token?: string
      profile?: ProfileId
      useMocks?: boolean
    }
  }
}

/** Runtime-detected profile config (UI-AUTH-01, UI-TEC-06). */
export function detectProfile(): RuntimeProfile {
  const emb = typeof window !== 'undefined' ? window.__HERMES_UI__ : undefined
  if (emb?.baseUrl || emb?.profile === 'embedded') {
    return {
      id: 'embedded',
      baseUrl: emb.baseUrl ?? '/api',
      token: emb.token,
      useMocks: emb.useMocks ?? false,
      label: 'P3 Embedded',
    }
  }
  const host = typeof window !== 'undefined' ? window.location.hostname : 'localhost'
  const qs = typeof window !== 'undefined' ? new URLSearchParams(window.location.search) : null
  const forceMocks = qs?.get('mocks') === '1' || import.meta.env.VITE_USE_MOCKS === '1'
  if (host === 'localhost' || host === '127.0.0.1') {
    return {
      id: 'local-first',
      baseUrl: import.meta.env.VITE_API_BASE ?? '/api',
      useMocks: forceMocks || import.meta.env.VITE_USE_MOCKS !== '0',
      label: 'P2 Local-first',
    }
  }
  return {
    id: 'platform',
    baseUrl: import.meta.env.VITE_API_BASE ?? '/api',
    useMocks: forceMocks,
    label: 'P1 Platform',
  }
}
